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
	router.GET("/download/:file", downloadFile)

	secureStaticFs := router.Group("/", authMiddleware)
	secureStaticFs.StaticFS("preview-file", http.Dir("public"))
	/*
				/download			/preview-file

		docx	download file		download file
		txt		download file		preview file in tab
		png		download file		preview file in tab
		pdf		download file		preview file in tab
		json	download file		preview file in tab
		xslx	download file		download file
		zip		download file		download file
		yaml	download file		preview file in tab
	*/
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
	filepath := "http://localhost:8080/preview-file/" + filename
	c.JSON(http.StatusOK, gin.H{"filepath": filepath})
}

func downloadFile(c *gin.Context) {
	file := c.Param("file")
	response, err := http.Get("http://localhost:8080/preview-file/" + file)
	if err != nil || response.StatusCode != http.StatusOK {
		c.Status(http.StatusServiceUnavailable)
		return
	}

	reader := response.Body
	contentLength := response.ContentLength
	contentType := response.Header.Get("Content-Type")
	extraHeaders := map[string]string{
		// "Content-Disposition": `attachment; filename=wt.png`,
		"Content-Disposition": `attachment; filename=` + file,
	}

	c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
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
