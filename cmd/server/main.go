package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"connectrpc.com/connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/bufbuild/protovalidate-go"
	"github.com/bufbuild/protovalidate-go/legacy"
	greetv1 "github.com/elliotmjackson/pv-demo/gen/greet/v1"
	"github.com/elliotmjackson/pv-demo/gen/greet/v1/greetv1connect"
)

type GreetServer struct {
	greetv1connect.UnimplementedGreetServiceHandler
	*protovalidate.Validator
}

func (s *GreetServer) Greet(
	ctx context.Context,
	req *connect.Request[greetv1.GreetRequest],
) (*connect.Response[greetv1.GreetResponse], error) {
	log.Println("Request headers: ", req.Header())
	if err := s.Validator.Validate(req.Msg); err != nil {
		log.Println(err.Error())
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	res := connect.NewResponse(&greetv1.GreetResponse{
		Greeting: fmt.Sprintf("Hello, %s!", req.Msg.Name),
	})
	res.Header().Set("Greet-Version", "v1")
	return res, nil
}

func main() {
	v, err := protovalidate.New(legacy.WithLegacySupport(legacy.ModeMerge))
	if err != nil {
		log.Fatal(err)
	}
	greeter := &GreetServer{
		Validator: v,
	}
	mux := http.NewServeMux()
	path, handler := greetv1connect.NewGreetServiceHandler(greeter)
	mux.Handle(path, handler)
	if err := http.ListenAndServe(
		"localhost:8080",
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	); err != nil {
		log.Fatal(err)
	}
}
