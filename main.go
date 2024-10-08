package main

import (
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"os"
)

func main() {
	if os.Getenv("SUPA_URL") == "" {
		log.Println("SUPA_URL is required")
		return
	}
	if os.Getenv("SUPA_SECRET_KEY") == "" {
		log.Println("SUPA_SECRET_KEY is required")
		return
	}
	port := func() string {
		if os.Getenv("PORT") == "" {
			return "8080"
		}
		return os.Getenv("PORT")
	}()

	router := gin.New()

	router.POST("webhook", func(c *gin.Context) {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Println(err)
			}
		}(c.Request.Body)

		bytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Println(err)
			return
		}
		WebhookHandler(bytes)
	})

	err := router.Run(":" + port)
	if err != nil {
		log.Println(err)
	}
}
