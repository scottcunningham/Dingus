package main

// Slack outgoing webhooks are handled here. Requests come in and are run through
// the markov chain to generate a response, which is sent back to Slack.
//
// Create an outgoing webhook in your Slack here:
// https://my.slack.com/services/new/outgoing-webhook

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
//	"fmt"
	"io/ioutil"
)



type whr struct {
	Token		 string `json:"token"`
	Team_ID	 string `json:"team_id"`
	ChannelID  string `json:"channel_id"`
	ChannelName string `json:"channel_name"`
	TS				string `json:"timestamp"`
	UserID		string `json:"user_id"`
	Username string `json:"user_name"`
	Text     string `json:"text"`
	Trigger  string `json:"trigger_word"`
}

func init() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//	incomingText := r.PostFormValue("text")
			

			err := r.ParseForm()
			if err != nil {
				panic(err)
			}
			//log.Println("r.Form", r.Form)
		
			txt := r.Form["text"]
			usr := r.Form["user_name"]
			incomingText := txt[0]
			rUsername := usr[0]
		if incomingText != "" && rUsername != "" && rUsername != "slackbot" {
			text := parseText(incomingText)
			log.Printf("Handling incoming request: %s from %s", text,rUsername )

			if text != "" {
				markovChain.Write(text)
			}


			go func() {
				markovChain.Save(stateFile)
			}()

			if rand.Intn(100) <= responseChance || strings.HasPrefix(text, botUsername) || strings.Contains(text,"Dingus") {
				var response whr
				response.Username = botUsername

				response.Text = markovChain.Generate(numWords)
				if strings.HasPrefix(text,"Dingus: compliment me"){
					url := "http://127.0.0.1:56735"
					comp, err := http.Get(url)
					if err != nil {
						log.Fatal(err)
					}
					defer comp.Body.Close()
					compData, err := ioutil.ReadAll(comp.Body)
					if err != nil {
						log.Fatal(err)
					}
					compString := string(compData)
					response.Text = rUsername + ": " + compString
				}
				log.Printf("Sending response: %s", response.Text)

				b, err := json.Marshal(response)
				if err != nil {
					log.Fatal(err)
				}

				

				time.Sleep(5 * time.Second)
				w.Write(b)
				r.Body.Close()
			}
		}
	})
}

func StartServer(port int) {
	log.Printf("Starting HTTP server on %d", port)
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
