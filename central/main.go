package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	pb "github.com/GonzaloDiaz300/LAB/proto"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

var flag bool

func enviarMensaje(servidor string, mensaje int, wg *sync.WaitGroup) {
	defer wg.Done()
	// Aquí puedes implementar la lógica para enviar el mensaje al servidor especificado
	conn, err := grpc.Dial(servidor, grpc.WithInsecure())
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
	if res != nil && !flag {
		flag = true
	}
}

func enviarInscripcion(mensaje int, servidor string) {
	conn, err := grpc.Dial("localhost:"+servidor, grpc.WithInsecure())
	if err != nil {
		fmt.Printf("Error al conectar con %s: %v\n", servidor, err)
		return
	}
	defer conn.Close()
	serviceClient := pb.NewNotificacionClient(conn)
	res, err := serviceClient.Inscribir(context.Background(), &pb.InscritosReq{Solicitud_2: int32(mensaje)})
	if err != nil {
		panic("No se llego el mensaje " + err.Error())
	}
	if res != nil && flag {
		flag = false
	}
}

func crearLlaves(limiteInferior int, limiteSuperior int) int {
	// Inicializa el generador de números aleatorios con una semilla única
	rand.Seed(time.Now().UnixNano())
	// Genera un número aleatorio dentro del rango
	numeroAleatorio := rand.Intn(limiteSuperior-limiteInferior+1) + limiteInferior
	return numeroAleatorio
}

func isIn(slice []int, item int) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

var numero_llaves int
var limiteInferior int
var limiteSuperior int

func main() {
	inicio := time.Now()
	// El código de configuración inicial (lectura de archivo, inicialización de variables, etc.) permanece igual
	totalIteraciones := 0

	// Abre el archivo para lectura
	archivo, err := os.Open("central/parametros_de_inicio.txt")
	if err != nil {
		fmt.Println("Error al abrir el archivo:", err)
		return
	}
	defer archivo.Close()

	// Crea un lector para el archivo
	lector := bufio.NewScanner(archivo)

	// Itera sobre las líneas del archivo
	if lector.Scan() {
		// Lee la línea actual
		linea := lector.Text()

		// Divide la línea en dos partes utilizando "-"
		partes := strings.Split(linea, "-")
		// Convierte las partes en enteros
		int1, err1 := strconv.Atoi(partes[0])
		int2, err2 := strconv.Atoi(partes[1])

		if err1 != nil || err2 != nil {
			fmt.Println("Error al convertir a entero en la línea:", linea)
		}

		// Realiza la operación deseada con int1 e int2

		limiteInferior = int1
		limiteSuperior = int2
	}

	if lector.Scan() {
		linea := lector.Text()
		iteraciones, err := strconv.Atoi(linea)
		if err != nil {
			fmt.Println("Error al convertir a entero en la segunda línea:", linea)
			return
		}
		totalIteraciones = iteraciones

	}
	numero_llaves = crearLlaves(limiteInferior, limiteSuperior)

	// Abre la cola rabbit para permitir una comunicación asíncrona
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

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

	if totalIteraciones == -1 {
		var servidores_listos []int
		iteracion := 1
		for {
			if len(servidores_listos) == 4 {
				break
			}
			fmt.Printf("\nGeneración %d/infinito\n", iteracion)
			// Lógica de cada iteración aquí
			/*log.Printf("Received a message: %d", iteracion)*/
			// Avisa a los servidores que tiene cupo mediante comunicación asíncrona

			servidores := []string{"10.6.46.90:50051", "10.6.46.87:50052", "10.6.46.88:50056", "10.6.46.89:50054"}
			var wg sync.WaitGroup
			for _, servidor := range servidores {
				wg.Add(1)
				go enviarMensaje(servidor, numero_llaves, &wg)
			}

			wg.Wait()
			var count = 0
			desiredMessageCount := 4
			for d := range msgs {
				/*log.Printf("Received a message: %s", d.Body)*/
				var nombre_servidor string
				partes_mensaje := strings.Split(string(d.Body), ",")
				puerto, err1 := strconv.Atoi(partes_mensaje[0])
				interesados, err2 := strconv.Atoi(partes_mensaje[1])
				if puerto == 50051 {
					nombre_servidor = "América"
					fmt.Printf("Mensaje asíncrono del servidor %s leído\n", nombre_servidor)
				} else if puerto == 50052 {
					nombre_servidor = "Oceanía"
					fmt.Printf("Mensaje asíncrono del servidor %s leído\n", nombre_servidor)
				} else if puerto == 50054 {
					nombre_servidor = "Europa"
					fmt.Printf("Mensaje asíncrono del servidor %s leído\n", nombre_servidor)
				} else {
					nombre_servidor = "Asia"
					fmt.Printf("Mensaje asíncrono del servidor %s leído\n", nombre_servidor)
				}

				if err1 != nil {
					fmt.Printf("Error al splitear mensaje")
					return
				}
				if err2 != nil {
					fmt.Printf("Error al splitear mensaje")
					return
				} // comprueba si el server ya está listo
				if interesados == -1 && !isIn(servidores_listos, puerto) {
					servidores_listos = append(servidores_listos, puerto)
					count++
					if count >= desiredMessageCount {
						break
					}
				} else {
					if numero_llaves > 0 {
						resultado := numero_llaves - interesados
						if resultado < 0 {
							//Se ocuparon todas las llaves, entonces se envian todas las llaves restantes a ese servidor y se dejan 0 en la centralv
							enviarInscripcion(-resultado, strconv.Itoa(puerto))
							if interesados == -1 {
								fmt.Printf("\nServidor %s ya abastecido\n", nombre_servidor)
							} else {
								fmt.Printf("Se inscribieron %d cupos de servidor %s\n", interesados, nombre_servidor)
							}
							numero_llaves = 0
							count++
						} else if resultado > 0 {
							//Se ocuparon llaves pero no todas
							enviarInscripcion(0, strconv.Itoa(puerto))
							if interesados == -1 {
								fmt.Printf("\nServidor %s ya abastecido\n", nombre_servidor)
							} else {
								fmt.Printf("Se inscribieron %d cupos de servidor %s\n", interesados, nombre_servidor)
							}
							numero_llaves = resultado
							count++
						} else {
							//cantidad de llaves = interesados
							enviarInscripcion(0, strconv.Itoa(puerto))
							if interesados == -1 {
								fmt.Printf("\nServidor %s ya abastecido\n", nombre_servidor)
							} else {
								fmt.Printf("Se inscribieron %d cupos de servidor %s\n", interesados, nombre_servidor)
							}
							numero_llaves = 0
							count++
						}
						if count >= desiredMessageCount {
							break
						}
					} else {
						enviarInscripcion(interesados, strconv.Itoa(puerto))
						count++
						if count >= desiredMessageCount {
							break
						}
					}
				}

			}
			numero_llaves = crearLlaves(limiteInferior, limiteSuperior)
			iteracion++
			// Realizar cualquier limpieza necesaria antes de la siguiente iteración
		}
	} else {
		for iteracion := 0; iteracion < totalIteraciones; iteracion++ {
			fmt.Printf("\nGeneración %d/%d\n", iteracion, totalIteraciones)
			// Lógica de cada iteración aquí
			/*log.Printf("Received a message: %d", iteracion)*/
			// Avisa a los servidores que tiene cupo mediante comunicación asíncrona

			servidores := []string{"10.6.46.90:50051", "10.6.46.87:50052", "10.6.46.88:50056", "10.6.46.89:50054"}
			var wg sync.WaitGroup
			for _, servidor := range servidores {
				wg.Add(1)
				go enviarMensaje(servidor, numero_llaves, &wg)
			}

			wg.Wait()
			var count = 0
			desiredMessageCount := 4
			for d := range msgs {
				/*log.Printf("Received a message: %s", d.Body)*/
				var nombre_servidor string
				partes_mensaje := strings.Split(string(d.Body), ",")
				puerto, err1 := strconv.Atoi(partes_mensaje[0])
				interesados, err2 := strconv.Atoi(partes_mensaje[1])
				if puerto == 50051 {
					nombre_servidor = "América"
					fmt.Printf("Mensaje asíncrono del servidor %s leído\n", nombre_servidor)
				} else if puerto == 50052 {
					nombre_servidor = "Oceanía"
					fmt.Printf("Mensaje asíncrono del servidor %s leído\n", nombre_servidor)
				} else if puerto == 50054 {
					nombre_servidor = "Europa"
					fmt.Printf("Mensaje asíncrono del servidor %s leído\n", nombre_servidor)
				} else {
					nombre_servidor = "Asia"
					fmt.Printf("Mensaje asíncrono del servidor %s leído\n", nombre_servidor)
				}

				if err1 != nil {
					fmt.Printf("Error al splitear mensaje")
					return
				}
				if err2 != nil {
					fmt.Printf("Error al splitear mensaje")
					return
				}
				if numero_llaves > 0 {
					resultado := numero_llaves - interesados
					if resultado < 0 {
						//Se ocuparon todas las llaves, entonces se envian todas las llaves restantes a ese servidor y se dejan 0 en la centralv
						enviarInscripcion(-resultado, strconv.Itoa(puerto))
						if interesados == -1 {
							fmt.Printf("\nServidor %s ya abastecido\n", nombre_servidor)
						} else {
							fmt.Printf("Se inscribieron %d cupos de servidor %s\n", interesados, nombre_servidor)
						}
						numero_llaves = 0
						count++
					} else if resultado > 0 {
						//Se ocuparon llaves pero no todas
						enviarInscripcion(0, strconv.Itoa(puerto))
						if interesados == -1 {
							fmt.Printf("\nServidor %s ya abastecido\n", nombre_servidor)
						} else {
							fmt.Printf("Se inscribieron %d cupos de servidor %s\n", interesados, nombre_servidor)
						}
						numero_llaves = resultado
						count++
					} else {
						//cantidad de llaves = interesados
						enviarInscripcion(0, strconv.Itoa(puerto))
						if interesados == -1 {
							fmt.Printf("\nServidor %s ya abastecido\n", nombre_servidor)
						} else {
							fmt.Printf("Se inscribieron %d cupos de servidor %s\n", interesados, nombre_servidor)
						}
						numero_llaves = 0
						count++
					}
					if count >= desiredMessageCount {
						break
					}
				} else {
					enviarInscripcion(interesados, strconv.Itoa(puerto))
					count++
					if count >= desiredMessageCount {
						break
					}
				}
			}
			numero_llaves = crearLlaves(limiteInferior, limiteSuperior)
			// Realizar cualquier limpieza necesaria antes de la siguiente iteración
		}
	}
	fin := time.Now()
	duracion := fin.Sub(inicio)

	fmt.Printf("El programa tardó %s en ejecutarse.\n", duracion)
}
