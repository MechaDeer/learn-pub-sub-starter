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
	const connString = "amqp://guest:guest@localhost:5672/"
	fmt.Println("Starting Peril client...")
	conn, err := amqp.Dial(connString)
	if err != nil {
		log.Fatalf("failed to connect to rabbitMQ: %v", err)
	}
	defer conn.Close()

	uname, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("can't start without username: %v", err)
	}
	_, queue, err := pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilDirect,
		fmt.Sprintf("%s.%s", routing.PauseKey, uname),
		routing.PauseKey,
		pubsub.SimpleQueueTransient,
	)

	if err != nil {
		log.Fatalf("could not subscribe to pause: %v", err)
	}
	fmt.Printf("Queue %v declared and bound!\n", queue.Name)

	gs := gamelogic.NewGameState(uname)

	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}
		switch words[0] {
		case "move":
			_, err := gs.CommandMove(words)
			if err != nil {
				fmt.Println(err)
				continue
			}

			// TODO: publish the move
		case "spawn":
			err = gs.CommandSpawn(words)
			if err != nil {
				fmt.Println(err)
				continue
			}
		case "status":
			gs.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			// TODO: publish n malicious logs
			fmt.Println("Spamming not allowed yet!")
		case "quit":
			gamelogic.PrintQuit()
			return
		default:
			fmt.Println("unknown command")
		}
	}
}
