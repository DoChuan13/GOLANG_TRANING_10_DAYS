package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ApiExample() {
	//Initial Http Router
	router := gin.Default()
	//Group by Path
	api := router.Group("/api")
	user := router.Group("/user")
	{
		api.GET(
			"/test", func(ctx *gin.Context) {
				ctx.JSON(
					http.StatusOK, gin.H{
						"message": "test api get successful",
					},
				)
			},
		)
		user.GET(
			"/test", func(ctx *gin.Context) {
				ctx.JSON(
					http.StatusOK, gin.H{
						"message": "test user get successful",
					},
				)
			},
		)
		user.POST(
			"/test", func(ctx *gin.Context) {
				ctx.JSON(
					http.StatusOK, gin.H{
						"message": "test user post successful",
					},
				)
			},
		)
		_ = router.Run(":8080")
	}
}
