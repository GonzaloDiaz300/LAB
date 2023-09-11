package main

import (
	"context"
	"net"
	"log"
	"time"
	"fmt"
	"google.golang.org/grpc"
	pb "github.com/GonzaloDiaz300/LAB/europa/proto"
	amqp "github.com/rabbitmq/amqp091-go"
)

type europa struct{
	pb.UnimplementedNotificacionServer
}

func failOnError(err error, msg string) {
	if err != nil {
	  log.Panicf("%s: %s", msg, err)
	}
}

func (a *europa) Notificar(ctx context.Context, in *pb.NotiReq) (*pb.NotiResp, error) {
	fmt.Printf("Se envia el 3 de vuelta a la central como confirmación\n")
	go encolarse(in.Solicitud)
	return &pb.NotiResp{Respuesta: 3}, nil
}

func encolarse(postulantes int32){
	//then connect to RabbitMQ server
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	//The connection abstracts the socket connection, and takes care of protocol version negotiation 
	//and authentication and so on for us. Next we create a channel, which is where most of the API 
	//for getting things done resides:
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()
	//To send, we must declare a queue for us to send to; then we can publish a message to the queue:
	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	  )
	  failOnError(err, "Failed to declare a queue")
	  
	  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	  defer cancel()
	  
	  body := []byte("89")
	  err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing {
		  ContentType: "text/plain",
		  Body:        body,
		})
	  failOnError(err, "Failed to publish a message")
	  log.Printf(" [x] Sent %s\n", body)
}

func main(){
	listner, err := net.Listen("tcp", ":50056")

	if err != nil {
		panic("cannot create tcp connection" + err.Error())
	}

	serv:= grpc.NewServer()
	fmt.Printf("Servidor Europa Activo\n")
	pb.RegisterNotificacionServer(serv, &europa{})
	if err = serv.Serve(listner); err != nil {
		panic("cannot initialize the server" + err.Error())
	}
}