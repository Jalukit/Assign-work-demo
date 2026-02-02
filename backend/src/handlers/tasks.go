package handlers

import (
    "context"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo/options"

    "assignwork/src/db"
    "assignwork/src/middleware"
    "assignwork/src/models"
)

type createTaskReq struct {
    Title       string `json:"title" binding:"required,min=2"`
    Description string `json:"description"`
    Capacity    int    `json:"capacity" binding:"required,min=1"`
}

// Admin create task
func CreateTask(c *gin.Context) {
    var req createTaskReq
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    uidAny, _ := c.Get(middleware.CtxUserID)
    uid := uidAny.(primitive.ObjectID)

    t := models.Task{
        ID:                primitive.NewObjectID(),
        Title:             req.Title,
        Description:       req.Description,
        Capacity:          req.Capacity,
        RegisteredUserIds: []primitive.ObjectID{},
        CreatedByUserId:   uid,
        CreatedAt:         time.Now().Unix(),
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    _, err := db.DB.Collection("tasks").InsertOne(ctx, t)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "create task failed"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"task": t})
}

// list tasks (both admin + user)
func ListTasks(c *gin.Context) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    cur, err := db.DB.Collection("tasks").Find(ctx, bson.M{}, options.Find().SetSort(bson.M{"createdAt": -1}))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
        return
    }
    defer cur.Close(ctx)

    tasks := []models.Task{}
    if err := cur.All(ctx, &tasks); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "decode tasks error"})
        return
    }

    // fetch users for each task (simple; for demo)
    userCol := db.DB.Collection("users")
    views := []models.TaskView{}
    for _, t := range tasks {
        registrants := []models.UserPublic{}
        if len(t.RegisteredUserIds) > 0 {
            ucur, _ := userCol.Find(ctx, bson.M{"_id": bson.M{"$in": t.RegisteredUserIds}})
            var us []models.User
            _ = ucur.All(ctx, &us)
            for _, u := range us {
                registrants = append(registrants, u.Public())
            }
        }
        registeredCount := len(t.RegisteredUserIds)
        remaining := t.Capacity - registeredCount
        if remaining < 0 {
            remaining = 0
        }
        views = append(views, models.TaskView{
            ID:              t.ID,
            Title:           t.Title,
            Description:     t.Description,
            Capacity:        t.Capacity,
            RegisteredCount: registeredCount,
            Remaining:       remaining,
            Registrants:     registrants,
            CreatedAt:       t.CreatedAt,
        })
    }

    c.JSON(http.StatusOK, gin.H{"tasks": views})
}

// user register for task
func RegisterForTask(c *gin.Context) {
    taskIDHex := c.Param("id")
    taskID, err := primitive.ObjectIDFromHex(taskIDHex)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
        return
    }

    uidAny, _ := c.Get(middleware.CtxUserID)
    uid := uidAny.(primitive.ObjectID)

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Ensure capacity not exceeded AND user not already registered.
    // We'll do it in two steps for simplicity; in production you'd do a more robust atomic update.
    var t models.Task
    if err := db.DB.Collection("tasks").FindOne(ctx, bson.M{"_id": taskID}).Decode(&t); err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
        return
    }
    for _, rid := range t.RegisteredUserIds {
        if rid == uid {
            c.JSON(http.StatusBadRequest, gin.H{"error": "already registered"})
            return
        }
    }
    if len(t.RegisteredUserIds) >= t.Capacity {
        c.JSON(http.StatusBadRequest, gin.H{"error": "task is full"})
        return
    }

    _, err = db.DB.Collection("tasks").UpdateOne(ctx,
        bson.M{"_id": taskID},
        bson.M{"$addToSet": bson.M{"registeredUserIds": uid}},
    )
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "register failed"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"ok": true})
}

// user cancel register
func UnregisterFromTask(c *gin.Context) {
    taskIDHex := c.Param("id")
    taskID, err := primitive.ObjectIDFromHex(taskIDHex)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
        return
    }

    uidAny, _ := c.Get(middleware.CtxUserID)
    uid := uidAny.(primitive.ObjectID)

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    _, err = db.DB.Collection("tasks").UpdateOne(ctx,
        bson.M{"_id": taskID},
        bson.M{"$pull": bson.M{"registeredUserIds": uid}},
    )
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "unregister failed"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"ok": true})
}
