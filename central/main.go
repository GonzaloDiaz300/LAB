package main

import (
	"context"
	"fmt"
	"sync"
	"log"
	pb "github.com/GonzaloDiaz300/LAB/central/proto"
	"google.golang.org/grpc"
	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
	  log.Panicf("%s: %s", msg, err)
	}
  }

func enviarMensaje(servidor string, mensaje int, wg *sync.WaitGroup) {
    defer wg.Done()
    // Aquí puedes implementar la lógica para enviar el mensaje al servidor especificado
	conn, err := grpc.Dial("localhost:"+servidor, grpc.WithInsecure())
	if err != nil {
        fmt.Printf("Error al conectar con %s: %v\n", servidor, err)
        return
    }
    defer conn.Close()
	serviceClient := pb.NewNotificacionClient(conn)
	res, err := serviceClient.Notificar(context.Background(), &pb.NotiReq{Solicitud: 100})
	if err != nil {
		panic("No se llego el mensaje " + err.Error())
	}
    fmt.Printf("Enviando mensaje a %s: %d\n", servidor, mensaje)
	fmt.Printf("Se recibió %d\n ", res.Respuesta)
}

func main(){
	servidores := []string{"50051", "50052", "50056", "50054"}
	var wg sync.WaitGroup
	mensaje := 300
	for _, servidor := range servidores {
        wg.Add(1)
        go enviarMensaje(servidor, mensaje, &wg)
    }

    wg.Wait()

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()
	//Setting up is the same as the publisher; we open a connection and a channel,
	//and declare the queue from which we're going to consume. Note this matches up
	//with the queue that send publishes to.
	q, err := ch.QueueDeclare(
	"hello", // name
	false,   // durable
	false,   // delete when unused
	false,   // exclusive
	false,   // no-wait
	nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")
	//We're about to tell the server to deliver us the messages from the queue. Since it will push us messages asynchronously, 
	//we will read the messages from a channel (returned by amqp::Consume) in a goroutine.
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	  )
	  failOnError(err, "Failed to register a consumer")
	  
	  var forever chan struct{}
	  
	  go func() {
		for d := range msgs {
		  log.Printf("Received a message: %s", d.Body)
		}
	  }()
	  
	  log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	  <-forever
}