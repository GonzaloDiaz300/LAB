package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"

	pb "github.com/GonzaloDiaz300/LAB/proto"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
)

var interesados int // Variable global para modificar los interesados en obtener la key
var interesados_actuales = 0

type europa struct {
	pb.UnimplementedNotificacionServer
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func (a *europa) Inscribir(ctx context.Context, in *pb.InscritosReq) (*pb.InscritosResp, error) {
	fmt.Printf("Se recibe inscripciones que no lograron pasar la cola\n")
	interesados_actuales = interesados_actuales - (interesados - int(in.Solicitud_2)) //700 = 700-(290-190)=600
	return &pb.InscritosResp{Respuesta_2: 7}, nil
}

func (a *europa) Notificar(ctx context.Context, in *pb.NotiReq) (*pb.NotiResp, error) {
	fmt.Printf("Se envia el 1 de vuelta a la central para confirmar llegada de mensaje\n")
	//aqui deberia procesarse la request
	go encolarse(int(in.Solicitud))
	return &pb.NotiResp{Respuesta: 3}, nil
}

// funcion para generar el numero de interesados en cada iteración, se llama cuando llaman la funcion de notificar
func crearInteresados(no_registrados int) int {
	if interesados_actuales == 0 {
		fileName := "europa/parametros_de_inicio.txt"

		// Intenta abrir el archivo
		file, err := os.Open(fileName)
		if err != nil {
			fmt.Println("Error al abrir el archivo:", err)
		}
		defer file.Close()

		// Leer el contenido de la primera línea
		scanner := bufio.NewScanner(file)
		if scanner.Scan() {
			// Obtener el texto de la primera línea
			line := scanner.Text()

			// Convertir el texto a un entero
			intValue, err := strconv.Atoi(line)
			if err != nil {
				fmt.Println("Error al convertir a entero:", err)
			}
			interesados_actuales = intValue
			// Almacenar el valor entero globalmente
		}
		interesados = interesados_actuales / 2
	} else {
		interesados_actuales = interesados_actuales - (interesados - no_registrados)
		interesados = interesados_actuales / 2
	}
	limiteInferior_interesados := math.Round(float64(interesados) - (float64(interesados) * 0.2))
	limiteSuperior_interesados := math.Round(float64(interesados) + (float64(interesados) * 0.2))

	// Inicializa el generador de números aleatorios con una semilla única
	rand.Seed(time.Now().UnixNano())
	// Genera un número aleatorio dentro del rango
	numeroAleatorio := rand.Intn(int(limiteSuperior_interesados)-int(limiteInferior_interesados)+1) + int(limiteInferior_interesados)
	fmt.Printf("Número aleatorio dentro del rango [%d, %d]: %d\n", int(limiteInferior_interesados), int(limiteSuperior_interesados), numeroAleatorio)

	interesados = numeroAleatorio
	fmt.Println("Valor entero global:", interesados)
	return interesados
}

func encolarse(cupos int) {
	postulantes_finales := crearInteresados(cupos)
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
	message := fmt.Sprintf("%d,%d", 50054, postulantes_finales)
	body := []byte(message)
	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", body)
}

func main() {
	listner, err := net.Listen("tcp", ":50056")

	if err != nil {
		panic("cannot create tcp connection" + err.Error())
	}

	serv := grpc.NewServer()
	fmt.Printf("Servidor Europa Activo\n")
	pb.RegisterNotificacionServer(serv, &europa{})
	if err = serv.Serve(listner); err != nil {
		panic("cannot initialize the server" + err.Error())
	}
}
