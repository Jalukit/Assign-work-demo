package routes

import (
    "net/http"
    "os"

    "github.com/gin-gonic/gin"

    "assignwork/src/handlers"
    "assignwork/src/middleware"
)

func Register(r *gin.Engine) {
    // simple health check
    r.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })

    // CORS (dev-friendly)
    corsOrigin := os.Getenv("CORS_ORIGIN")
    if corsOrigin == "" {
        corsOrigin = "http://localhost:5173"
    }
    r.Use(func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", corsOrigin)
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        c.Next()
    })

    api := r.Group("/api")

    auth := api.Group("/auth")
    auth.POST("/register", handlers.Register)
    auth.POST("/login", handlers.Login)

    me := api.Group("/me")
    me.Use(middleware.AuthRequired())
    me.GET("", handlers.Me)

    tasks := api.Group("/tasks")
    tasks.Use(middleware.AuthRequired())
    tasks.GET("", handlers.ListTasks)
    tasks.POST("", middleware.AdminOnly(), handlers.CreateTask)
    tasks.POST("/:id/register", handlers.RegisterForTask)
    tasks.POST("/:id/unregister", handlers.UnregisterFromTask)
}
