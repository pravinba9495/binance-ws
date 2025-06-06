package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/log"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

func SendToPubSub(msg WebSocketMsg, rdb PubSub) {
	price, err := strconv.ParseFloat(msg.Data.Price, 64)
	if err != nil {
		log.Errorf("Error parsing price: %v", err)
		return
	}
	data := PublishTradeData{
		Symbol:    msg.Data.Symbol,
		Price:     price,
		Timestamp: msg.Data.Timestamp,
	}
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

func restartService(interval time.Duration) {
	log.Warnf("Service will quit in interval: %v hr(s)", interval.Hours())
	time.Sleep(interval)
	log.Info("Restarting service ...")
	os.Exit(1)
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

	w, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("wss://stream.binance.com:443/stream?streams=%s@trade", strings.ToLower(os.Getenv("SYMBOL"))), nil)
	if err != nil {
		log.Fatalf("Error creating websocket connection: %v", err)
	}
	defer w.Close()

	go restartService(time.Hour * 23)

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
