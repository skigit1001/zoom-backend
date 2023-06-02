package websocket

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/websocket/v2"
)

// Upgraded websocket request
func handleWebsocket(c *websocket.Conn) {
	fmt.Println(c.Locals("Host")) // "Localhost:3000"
	var videoFile *os.File
	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		if string(msg) == "start_recording" {
			file, err := os.Create(fmt.Sprintf("./static/video/%d.webm", time.Now().Unix()))
			if err != nil {
				panic(err)
			}
			videoFile = file
			c.WriteMessage(websocket.TextMessage, []byte("accepted_recording"))
		} else if string(msg) == "stop_recording" {
			videoFile.Close()
		} else {
			decodedBytes, err := base64.StdEncoding.DecodeString(string(msg))
			if err != nil {
				panic(err)
			}
			videoFile.Write(decodedBytes)
		}
	}
}
