package format

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func FormatExample() {
	router := gin.Default()
	defer func() {
		_ = router.Run(":8080")
	}()
	router.GET(
		"/json", func(context *gin.Context) {
			context.JSON(http.StatusOK, gin.H{"message": "Json message"})
		},
	)
	router.GET(
		"/xml", func(context *gin.Context) {
			context.XML(http.StatusOK, gin.H{"message": "XML message"})
		},
	)
	router.GET(
		"/yaml", func(context *gin.Context) {
			context.YAML(http.StatusOK, gin.H{"message": "Yaml message"})
		},
	)
	router.GET(
		"/toml", func(context *gin.Context) {
			context.TOML(http.StatusOK, gin.H{"message": "Toml message"})
		},
	)
}
