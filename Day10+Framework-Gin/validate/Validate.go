package validate

import (
    "github.com/gin-gonic/gin"
    "net/http"
    "strings"
)

type User struct {
    Id   int    `json:"id" binding:"required"`
    Name string `json:"name" binding:"required"`
}

func ValidateExample() {
    router := gin.Default()
    defer func() {
        _ = router.Run(":8080")
    }()

    router.POST(
        "/post/user", func(context *gin.Context) {
            var body User
            if err := context.ShouldBind(&body); err != nil {
                context.AbortWithStatusJSON(
                    http.StatusBadRequest,
                    gin.H{"message": err.Error()},
                )
                return
            }
            body.Name = strings.TrimSpace(body.Name)
            if len(body.Name) == 0 {
                context.JSON(
                    http.StatusBadRequest,
                    gin.H{"message": "Name Field is Empty"},
                )
                return
            }
            context.JSON(http.StatusOK, gin.H{"message": "Post Success", "value": body})
        },
    )
}
