package main

import (
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"os"
)

func main() {

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
