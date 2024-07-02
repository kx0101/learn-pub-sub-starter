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
	fmt.Println("Starting Peril client...")
	url := "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(url)
	if err != nil {
		fmt.Println("error: ", err)
		return
	}

	defer conn.Close()

	fmt.Println("Connection was successful")

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		fmt.Println("error: ", err)
		return
	}

	_, queue, err := pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilDirect,
		routing.PauseKey+"."+username,
		routing.PauseKey,
		pubsub.TransientQueue,
	)
	if err != nil {
		log.Fatalf("could not subscribe to pause: %v", err)
	}
	fmt.Printf("Queue %v declared and bound!\n", queue.Name)

	gameState := gamelogic.NewGameState(username)

	for {
		words := gamelogic.GetInput()

		if len(words) == 0 {
			continue
		}

		switch words[0] {
		case "spawn":
			{
				gameState.CommandSpawn(words)
				continue
			}
		case "move":
			{
				gameState.CommandMove(words)
				continue
			}
		case "status":
			{
				gameState.CommandStatus()
				continue
			}
		case "help":
			{
				gamelogic.PrintClientHelp()
				continue
			}
		case "spam":
			{
				log.Printf("Spamming not allowed yet!")
				continue
			}
		case "quit":
			{
				gamelogic.PrintQuit()
				continue
			}
		}

		log.Print("I don't understand the command")
	}
}
