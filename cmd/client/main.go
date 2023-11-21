package main

import (
	"context"
	"log"
	"net/http"

	greetv1 "github.com/elliotmjackson/pv-demo/gen/greet/v1"
	"github.com/elliotmjackson/pv-demo/gen/greet/v1/greetv1connect"

	"connectrpc.com/connect"
)

func main() {
	client := greetv1connect.NewGreetServiceClient(
		http.DefaultClient,
		"http://localhost:8080",
	)
	res, err := client.Greet(
		context.Background(),
		connect.NewRequest(&greetv1.GreetRequest{Name: "J@ne"}),
	)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(res.Msg.Greeting)
}