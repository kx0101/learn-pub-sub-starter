package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")
	url := "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(url)
	if err != nil {
		fmt.Println("error: ", err)
		return
	}

	defer conn.Close()

	fmt.Println("Connection was successful")

	c, err := conn.Channel()
	if err != nil {
		fmt.Println("error: ", err)
		return
	}

	gamelogic.PrintServerHelp()

	_, _, err = pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilTopic,
		routing.GameLogSlug,
        routing.GameLogSlug+".*",
		pubsub.DurableQueue,
	)
	if err != nil {
		fmt.Println("error: ", err)
		return
	}

	for {
		words := gamelogic.GetInput()

		if len(words) == 0 {
			continue
		}

		if words[0] == "pause" {
			log.Printf("Pausing...")
			pubsub.PublishJSON(c, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
				IsPaused: true,
			})

			continue
		}

		if words[0] == "resume" {
			log.Printf("Resuming...")
			pubsub.PublishJSON(c, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
				IsPaused: false,
			})

			continue
		}

		if words[0] == "quit" {
			log.Printf("Quitting...")
			break
		}

		log.Printf("I don't understand the command")
	}
}
