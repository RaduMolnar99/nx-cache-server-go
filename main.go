package main

import (
	"fmt"
	"io"
	"os"

	"github.com/gin-gonic/gin"
)

var Secret = os.Getenv("SECRET")

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")

		if auth == Secret {
			c.Next()
		} else {
			c.Data(401, "text/plain", []byte("Missing or invalid authentication token."))
		}
	}
}

func uploadCache(c *gin.Context) {
	hash := c.Param("hash")

	file, err := io.ReadAll(c.Request.Body)

	filePath := fmt.Sprintf("cache/%s", hash)

	if err != nil {
		c.Data(400, "application/octet-stream", []byte("Invalid request body"))
	}

	if _, err := os.Stat(filePath); err == nil {

		fmt.Printf("HERERE: %s\n", err)
		c.Data(409, "text/plain", []byte("Cannot override an existing record"))
	} else {
		fmt.Printf("HERERE: %s", filePath)
		os.WriteFile(filePath, file, 0644)

		c.Data(200, "text/plain", []byte("Successfully uploaded the output"))
	}
}

func getCache(c *gin.Context) {
	hash := c.Param("hash")
	byteFile, err := os.ReadFile(fmt.Sprintf("./cache/%s", hash))

	if err != nil {
		c.Data(404, "text/plain", []byte("Record was not found"))
	}

	c.Data(200, "application/octet-stream", byteFile)
}

func main() {
	router := gin.Default()

	router.Use(authMiddleware())

	router.GET("/v1/cache/:hash", getCache)
	router.PUT("/v1/cache/:hash", uploadCache)

	router.Run(":" + os.Getenv("PORT"))
}
