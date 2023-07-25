package panic

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func PanicExample() {
	router := gin.Default()
	defer func() {
		_ = router.Run(":8080")
	}()

	router.Use(panicHandler)
	router.GET(
		"/panic", func(context *gin.Context) {
			panic("Error Panic")
		},
	)
}

func panicHandler(context *gin.Context) {
	if err := recover(); err != nil {
		context.JSON(http.StatusAccepted, gin.H{"message": err})
		context.Abort()
		return
	}
	context.Next()
}
