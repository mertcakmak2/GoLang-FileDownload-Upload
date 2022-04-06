package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.POST("/upload-file", uploadFile)
	router.GET("/get-file/:file", getFileAsBase64String)

	secureStaticFs := router.Group("/", authMiddleware)
	secureStaticFs.StaticFS("preview-image", http.Dir("public"))
	router.Run(":8080")
}

func uploadFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("file err : %s", err.Error()))
		return
	}
	filename := header.Filename
	out, err := os.Create("public/" + filename)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		log.Fatal(err)
	}
	filepath := "http://localhost:8080/preview-image/" + filename
	c.JSON(http.StatusOK, gin.H{"filepath": filepath})
}

func getFileAsBase64String(c *gin.Context) {
	file := c.Param("file")
	out, err := os.ReadFile("public/" + file)
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, out)
}

func authMiddleware(c *gin.Context) {
	fmt.Println("auth middleware")
	// c.JSON(http.StatusUnauthorized, "unauth")
	// c.Abort()
}
