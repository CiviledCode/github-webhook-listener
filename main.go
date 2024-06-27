package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type Config struct {
	IP       string                     `json:"ip"`
	Port     int                        `json:"port"`
	Webhooks map[string]WebhookEndpoint `json:"webhooks"`
}

type WebhookEndpoint struct {
	Type  string `json:"type"`
	Data  any    `json:"data"`
	Token string `json:"secret_token"`
}

var conf Config
var mux *http.ServeMux

func main() {
	loadConfig()
	address := fmt.Sprintf("%s:%d", conf.IP, conf.Port)
	fmt.Printf("Starting Server at http://%s/\n", address)
	mux = http.NewServeMux()
	handler := http.HandlerFunc(endpointHandler)
	for path := range conf.Webhooks {
		mux.Handle(path, handler)
	}
	http.ListenAndServe(address, mux)
}

func loadConfig() {
	confFile, err := os.Open("./config.json")
	if err != nil {
		panic(fmt.Errorf("error opening config file: %w", err))
	}

	confBytes, _ := io.ReadAll(confFile)

	err = json.Unmarshal(confBytes, &conf)
	if err != nil {
		panic(fmt.Errorf("error unmarshaling json: %w", err))
	}
}

// https://docs.github.com/en/webhooks/using-webhooks/validating-webhook-deliveries
// generateSha256Hmac([]byte("Hello, World!"), []byte("It's a Secret to Everybody")) == 757107ea0eb2509fc211221cce984b8a37570b6d7586c22c46f4379c8b043e17
func generateSha256Hmac(data []byte, secret []byte) string {
	hmac := hmac.New(sha256.New, secret)
	hmac.Write(data)
	return hex.EncodeToString(hmac.Sum(nil))
}

func endpointHandler(w http.ResponseWriter, r *http.Request) {
	if endpoint, ok := conf.Webhooks[r.RequestURI]; ok {
		if endpoint.Token != "" {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(400)
				return
			}

			key := generateSha256Hmac(body, []byte(endpoint.Token))
			keyHeader := r.Header.Get("X-Hub-Signature-256")

			if key != keyHeader[7:] {
				w.WriteHeader(401)
				return
			}
		}

		switch endpoint.Type {
		case "command", "cmd":
			if cmd, ok := endpoint.Data.(string); ok {
				// This is naive and doesn't allow more complex commands
				// like cat "some argument with spaces"
				spl := strings.Split(cmd, " ")
				fmt.Println("Running Command:", cmd)
				builtCmd := exec.Command(spl[0], spl[1:]...)
				err := builtCmd.Run()
				if err != nil {
					fmt.Printf("Error running command: %v\n", err)
				}
			}
		}
	}
}
