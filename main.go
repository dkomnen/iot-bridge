package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var counter int

func main() {
	host := flag.String("host", "", "host")
	port := flag.String("port", "", "port")
	flag.Parse()

	opts := mqtt.NewClientOptions().
		AddBroker("tcp://localhost:1883").
		SetUsername("user_bridge").
		SetPassword("pass_bridge").
		SetClientID("bridge")

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	if tok := c.Subscribe(
		"topic",
		0,
		func(cl mqtt.Client, msg mqtt.Message) {
			var toSend struct {
				Payload string `json:"payload"`
			}
			toSend.Payload = string(msg.Payload())

			var buff bytes.Buffer
			json.NewEncoder(&buff).Encode(toSend)
			fmt.Printf("req: %s\n", buff.Bytes())

			c := http.DefaultClient
			url := fmt.Sprintf("http://%s:%s/test", *host, *port)
			if req, err := http.NewRequest(http.MethodPost, url, &buff); err != nil {
				log.Fatal(err)
			} else {
				req.Header.Set("Content-Type", "application/json")
				if rsp, err := c.Do(req); err != nil {
					log.Fatal(err)
				} else {
					if _, err := io.Copy(os.Stdout, rsp.Body); err != nil {
						log.Fatal(err)
					}
				}
			}
		},
	); tok.Wait() && tok.Error() != nil {
		log.Fatal(tok.Error())
	}

	block := make(chan struct{})
	<-block
}
