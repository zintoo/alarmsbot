package main

import (
	"bytes"
	"log"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/line"
) 

type Expression struct{
      	InputExpr string
      	SymbolExpr string
}

type Message struct{
	From string
	MathExpr string
}

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
		//log.Print(err)
		return //1
	}
	bot, err := linebot.NewClient(channelID, channelSecret, channelMID)
	if err != nil {
		//log.Print(err)
		return //1
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")
        router.POST("/callback",  func(c *gin.Context) {
        	w := c.Writer
        	req := c.Request
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
				url := "http://122.154.148.234/expr"
				m := Message{content.From, text.Text}
				b, err := json.Marshal(m)
				if err != nil {
				   	return
				}
				req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
				if err != nil {
				   	return
				}
				req.Header.Set("Content-Type", "application/json")
				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
				   	return
				}
                                defer resp.Body.Close()
                                // Read the content into a byte array
                                body, err_json := ioutil.ReadAll(resp.Body)
                                if err_json != nil {
                                    	return 
                                }
                                
                                var expr Expression
                                err = json.Unmarshal(body, &expr)
                                if err_json != nil {
                                    	return 
                                }
				//_, err = bot.SendText([]string{content.From}, text.Text)
				_, err = bot.SendText([]string{content.From}, expr.SymbolExpr)
				if err != nil {
					return
				}
				/*
				_, err = bot.SendImage([]string{content.From}, "http://122.154.148.234/static/imgexprs/expr.png", "http://122.154.148.234/static/imgexprs/expr.png")
				if err != nil {
					//log.Print(err)
				}*/
				
			}
		}
	})

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.Run(":" + port)
}
