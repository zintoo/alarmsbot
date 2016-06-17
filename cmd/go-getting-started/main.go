package main

import (
	"log"
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
	"github.com/line"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

        var (
		channelID     int64
		channelSecret = "01cddcc039fb93e6511aa4fd0179b98e"
		channelMID    = "u9e839472b886bb2162391ae3c0f926a8"
		err           error
	)

	// Setup bot client
	channelID, err = strconv.ParseInt("1471157712", 10, 64)
	if err != nil {
		log.Print(err)
		return 1
	}
	bot, err := linebot.NewClient(channelID, channelSecret, channelMID)
	if err != nil {
		log.Print(err)
		return 1
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")
        router.HandleFunc("/callback",  func(w http.ResponseWriter, req *http.Request) {
		received, err := bot.ParseRequest(req)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}
		for _, result := range received.Results {
			content := result.Content()
			if content != nil && content.IsMessage && content.ContentType == linebot.ContentTypeText {
				text, err := content.TextContent()
				_, err = bot.SendText([]string{content.From}, text.Text)
				if err != nil {
					log.Print(err)
				}
			}
		}
	})

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.Run(":" + port)
}
