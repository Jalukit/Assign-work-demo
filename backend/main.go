package main

import (
    "log"
    "os"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"

    "assignwork/src/db"
    "assignwork/src/routes"
)

func main() {
    _ = godotenv.Load()

    r := gin.Default()

    mongoURI := getenv("MONGO_URI", "mongodb://localhost:27017")
    dbName := getenv("DB_NAME", "assignwork")

    if err := db.Connect(mongoURI, dbName); err != nil {
        log.Fatal(err)
    }

    routes.Register(r)

    port := getenv("PORT", "8080")
    log.Println("Backend listening on :" + port)
    if err := r.Run(":" + port); err != nil {
        log.Fatal(err)
    }
}

func getenv(k, fallback string) string {
    if v := os.Getenv(k); v != "" {
        return v
    }
    return fallback
}
