package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/log"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

func SendToPubSub(msg WebSocketMsg, rdb PubSub) {
	if msg.Type == "trade" {
		for _, data := range msg.Data {
			bytes, err := json.Marshal(data)
			if err != nil {
				log.Errorf("Error marshaling message: %v", err)
			} else {
				log.Infof("%s", bytes)

				if err := rdb.Publish("trades", bytes); err != nil {
					log.Errorf("Error publishing trade data to Redis: %v", err.Error())
				}
			}
		}
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Warnf("Error loading .env file: %v", err)
	}

	for _, envVar := range REQUIRED_ENV_VARS {
		key, isKeySet := os.LookupEnv(envVar)
		if !isKeySet || key == "" {
			log.Fatalf("%s is not set in the environment variables", envVar)
		}
	}

	pubsub := NewRedisPubSub()

	w, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("wss://ws.finnhub.io?token=%s", os.Getenv("FINNHUB_API_KEY")), nil)
	if err != nil {
		log.Fatalf("Error creating websocket connection: %v", err)
	}
	defer w.Close()

	symbols := strings.Split(os.Getenv("FINNHUB_SYMBOLS"), ",")
	for _, s := range symbols {
		msg, err := json.Marshal(map[string]interface{}{"type": "subscribe", "symbol": s})
		if err != nil {
			log.Errorf("Error marshaling message: %v", err)
		} else {
			w.WriteMessage(websocket.TextMessage, msg)
		}
	}

	var msg WebSocketMsg

	for {
		err := w.ReadJSON(&msg)
		if err != nil {
			log.Fatal(err)
		} else {
			go SendToPubSub(msg, pubsub)
		}
	}
}
