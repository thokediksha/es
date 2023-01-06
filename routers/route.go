package routers

import (
    "fmt"
    controllers "es/controllers"

    "github.com/gin-gonic/gin"
)


func SetupRouter() *gin.Engine {
    r := gin.Default()
    r.Use(gin.Logger())
    r.Use(func(c *gin.Context) {
        //allow all
        c.Writer.Header().Set("Content-Type", "application/json")
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Max-Age", "86400")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "*")
        c.Writer.Header().Set("Access-Control-Allow-Headers",  "*")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(200)
        } else {
            platform := c.GetHeader("User-agent")
            fmt.Println("platofor", platform)
        }
        c.Next()
    })
    r.POST("/users", controllers.CreateUser)
    r.PUT("users/:index/:id", controllers.UpdateUser)
    r.GET("users/:index/:id", controllers.GetUser)
    r.GET("users/:index", controllers.GetAllUser)
    r.DELETE("users/:index/:id", controllers.DeleteUser)
    return r
}