package main

import (
	"context"
	"log"
	"net/http"
	"strings"

	pb "github.com/ahmelsayed/ws-grpc-sample/hello"
	"github.com/gorilla/websocket"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

func wsGrpcHandlerFunc(grpcServer *grpc.Server, wsHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			wsHandler.ServeHTTP(w, r)
		}
	})
}

type helloService struct {
	pb.HelloServiceServer
}

func (s *helloService) Hello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Grpc Response >> " + req.Name}, nil
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		if err := conn.WriteMessage(messageType, append([]byte("Websockets Response >> "), p...)); err != nil {
			log.Println(err)
			return
		}
	}
}

func main() {
	// setup grpc
	grpcServer := grpc.NewServer()
	pb.RegisterHelloServiceServer(grpcServer, &helloService{})

	// setup websocket
	mux := http.NewServeMux()
	mux.HandleFunc("/", wsHandler)

	// setup server
	server := &http.Server{
		Addr: ":8080",
		// handle grpc and websocket based on protocol and content-type
		Handler: h2c.NewHandler(wsGrpcHandlerFunc(grpcServer, mux), &http2.Server{}),
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
