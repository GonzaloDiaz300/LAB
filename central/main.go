package main

import (
	"bufio"
	"context"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	pb "github.com/GonzaloDiaz300/LAB/central/proto"
	"google.golang.org/grpc"
)

var contador = 0
var numero_iteraciones int
var numero_llaves int
var limiteInferior int
var limiteSuperior int

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

func crearLlaves(limiteInferior int, limiteSuperior int) int {
	// Inicializa el generador de números aleatorios con una semilla única
	rand.Seed(time.Now().UnixNano())
	// Genera un número aleatorio dentro del rango
	numeroAleatorio := rand.Intn(limiteSuperior-limiteInferior+1) + limiteInferior

	fmt.Printf("Número aleatorio dentro del rango [%d, %d]: %d\n", limiteInferior, limiteSuperior, numeroAleatorio)
	return numeroAleatorio
}

func main() {

	// Esto solo se hace la primera vez, por eso la condicional de contador < 1, asigna variables como el limite superior e inferior, el numero de iteraciones
	// y el numero de llaves de la primera iteración
	if contador < 1 {

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
			fmt.Printf("Límite inferior: %d, Límite superior: %d\n", int1, int2)
		}

		if lector.Scan() {
			linea := lector.Text()
			iteraciones, err := strconv.Atoi(linea)
			if err != nil {
				fmt.Println("Error al convertir a entero en la segunda línea:", linea)
				return
			}
			numero_iteraciones = iteraciones

			fmt.Printf("Número interaciones: %d\n", iteraciones)
		}
		numero_llaves = crearLlaves(limiteInferior, limiteSuperior)
	}

	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	servidores := []string{"50051", "50052", "50053", "50054"}
	var wg sync.WaitGroup

	// Si el número de iteraciones que dice el documento es -1, entonces se repetirá para siempre, a lo mejor hay que ponerle alguna condición si es que
	// todos los servidores ya estan abastecidos
	if numero_iteraciones == -1 {
		for {
			for _, servidor := range servidores {
				wg.Add(1)
				go enviarMensaje(servidor, numero_llaves, &wg)
			}
			wg.Wait()
			contador += 1
		}
	} else {
		// se envían llaves las veces que pida el archivo, quizá cambie un poco al hacer la comunicacion asincrona
		for contador < numero_iteraciones {
			fmt.Printf("Iteración número: %d\n", contador+1)
			numero_llaves = crearLlaves(limiteInferior, limiteSuperior)
			for _, servidor := range servidores {
				wg.Add(1)
				go enviarMensaje(servidor, numero_llaves, &wg)
			}
			wg.Wait()
			contador += 1
		}
	}
}
