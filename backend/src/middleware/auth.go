package middleware

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson/primitive"

    "assignwork/src/models"
    "assignwork/src/utils"
)

const (
    CtxUserID = "userId"
    CtxRole   = "role"
)

func AuthRequired() gin.HandlerFunc {
    return func(c *gin.Context) {
        h := c.GetHeader("Authorization")
        if h == "" || !strings.HasPrefix(h, "Bearer ") {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
            return
        }
        tokenStr := strings.TrimPrefix(h, "Bearer ")
        claims, err := utils.ParseToken(tokenStr)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
            return
        }
        uid, err := primitive.ObjectIDFromHex(claims.UserID)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
            return
        }
        c.Set(CtxUserID, uid)
        c.Set(CtxRole, claims.Role)
        c.Next()
    }
}

func AdminOnly() gin.HandlerFunc {
    return func(c *gin.Context) {
        roleAny, ok := c.Get(CtxRole)
        if !ok {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing auth context"})
            return
        }
        role := roleAny.(models.UserRole)
        if role != models.RoleAdmin {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin only"})
            return
        }
        c.Next()
    }
}
