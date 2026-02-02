package handlers

import (
    "context"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"

    "assignwork/src/db"
    "assignwork/src/models"
    "assignwork/src/utils"
)

type registerReq struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
    Role     string `json:"role"` // optional: "admin" or "user"
}

func Register(c *gin.Context) {
    var req registerReq
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    role := models.RoleUser
    if req.Role == string(models.RoleAdmin) {
        role = models.RoleAdmin
    }

    hash, err := utils.HashPassword(req.Password)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "hash password failed"})
        return
    }

    u := models.User{
        ID:           primitive.NewObjectID(),
        Email:        req.Email,
        PasswordHash: hash,
        Role:         role,
        CreatedAt:    time.Now().Unix(),
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    _, err = db.DB.Collection("users").InsertOne(ctx, u)
    if err != nil {
        // likely duplicate email
        c.JSON(http.StatusBadRequest, gin.H{"error": "email already exists"})
        return
    }

    token, err := utils.SignToken(u.ID, u.Role)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "sign token failed"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "token": token,
        "user":  u.Public(),
    })
}

type loginReq struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
    var req loginReq
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var u models.User
    err := db.DB.Collection("users").FindOne(ctx, bson.M{"email": req.Email}).Decode(&u)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
        return
    }

    if !utils.CheckPassword(u.PasswordHash, req.Password) {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
        return
    }

    token, err := utils.SignToken(u.ID, u.Role)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "sign token failed"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "token": token,
        "user":  u.Public(),
    })
}

func Me(c *gin.Context) {
    uidAny, _ := c.Get("userId")
    uid := uidAny.(primitive.ObjectID)

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var u models.User
    err := db.DB.Collection("users").FindOne(ctx, bson.M{"_id": uid}).Decode(&u)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"user": u.Public()})
}
