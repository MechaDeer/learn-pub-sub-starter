package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	const connString = "amqp://guest:guest@localhost:5672/"

	fmt.Println("Starting Peril server...")
	conn, err := amqp.Dial(connString)
	if err != nil {
		log.Fatalf("failed to connect to rabbitMQ: %v", err)
	}
	defer conn.Close()
	fmt.Println("Connection to Peril server was successful")

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to create channel: %v", err)
	}

	pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
		IsPaused: true,
	})

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	<-sigc
	fmt.Println("Shutting down Peril server...")
}
