package cfgs

import (
	"log"
	"os"
	"reflect"

	"github.com/joho/godotenv"
	"github.com/pion/webrtc/v3"
)

const (
	RECORDER_WEBSOCKET_CHANNEL_SIZE = 256 // RECORDER channel buffer size for websocket

	RECORDER_ENVKEY_SESSIONID      = "RECORDER_SESSIONID" // name of env var to keep sessionid value
	RECORDER_WEBRTC_DATA_CHANNEL   = "RECORDER"           // lable name of webrtc data channel to exchange byte data
	RECORDER_WEBRTC_CONFIG_CHANNEL = "config"             // lable name of webrtc config channel to exchange config
	RECORDER_WEBSOCKET_HOST_ID     = "host"               // ID of message sent by the host

	RECORDER_VERSION  = "0.0.4"
	SUPPORTED_VERSION = "0.0.4" // the oldest TRANSFER version of client that the host could support
)

var TRANSFER_ICE_SERVER_STUNS = []webrtc.ICEServer{{URLs: []string{"stun:stun.l.google.com:19302", "stun:stun1.l.google.com:19302"}}}

// var TRANSFER_ICE_SERVER_TURNS = []webrtc.ICEServer{{
// 	URLs:       []string{"turn:104.237.1.191:3478"},
// 	Username:   "TRANSFER",
// 	Credential: "termishareisfun"}}

type Configs struct {
	DB_CONNECTION_URI    string
	REDIS_CONNECTION_URI string
	LLM_URI              string
}

func LoadConfigs() Configs {
	defaultConfig := Configs{}
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	configValue := reflect.ValueOf(&defaultConfig).Elem()
	configType := configValue.Type()

	for i := 0; i < configType.NumField(); i++ {
		field := configType.Field(i)
		envKey := "GJ_" + field.Name
		if value := os.Getenv(envKey); value != "" {
			configValue.Field(i).SetString(value)
		} else {
			log.Fatalf("Error: Required environment variable %s not set", envKey)
		}
	}
	return defaultConfig
}
