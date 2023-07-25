package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func MiddlewareExample() {
	router := gin.Default()
	defer func() {
		_ = router.Run(":8080")
	}()
	router.Use(middlewareHandle)

	router.GET(
		"/private", func(context *gin.Context) {
			query := context.Query("value")
			fmt.Println(query)
			context.JSON(http.StatusAccepted, gin.H{"message": "Private Content " + query})
		},
	)
	router.GET(
		"/private/:id", func(context *gin.Context) {
			fmt.Println(context.Params)
			id := context.Param("id")
			context.JSON(http.StatusAccepted, gin.H{"message": "Private Content (Param)==>" + id})
		},
	)
	router.GET(
		"/public", func(context *gin.Context) {
			context.JSON(http.StatusAccepted, gin.H{"message": "Public Content"})
		},
	)
	router.GET(
		"/common", func(context *gin.Context) {
			context.JSON(http.StatusAccepted, gin.H{"message": "Common Content"})
		},
	)
}

func middlewareHandle(context *gin.Context) {
	isAccess := false
	url := context.Request.URL.String()
	param := ""
	value1 := context.Param("id")
	value2 := context.Query("value")
	if value1 != "" {
		param = value1
	} else {
		param = value2
	}

	if param != "" {
		isAccess = true
	}
	fmt.Println("URL==>", url)
	fmt.Println("Params==>", param)
	if strings.Contains(url, "/common") {
		context.Redirect(http.StatusMovedPermanently, "/public")
		return
	}
	if strings.Contains(url, "/private") && param == "" {
		if !isAccess {
			context.AbortWithStatus(http.StatusForbidden)
			return
		}
	}
	context.Next()
}
