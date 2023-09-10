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

func main() {
	var limiteInferior int
	var limiteSuperior int

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
	for lector.Scan() {
		// Lee la línea actual
		linea := lector.Text()

		// Divide la línea en dos partes utilizando "-"
		partes := strings.Split(linea, "-")
		// Convierte las partes en enteros
		int1, err1 := strconv.Atoi(partes[0])
		int2, err2 := strconv.Atoi(partes[1])

		if err1 != nil || err2 != nil {
			fmt.Println("Error al convertir a entero en la línea:", linea)
			continue
		}

		// Realiza la operación deseada con int1 e int2

		limiteInferior = int1
		limiteSuperior = int2
		fmt.Printf("Límite inferior: %d, Límite superior: %d\n", int1, int2)
	}

	///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	// Inicializa el generador de números aleatorios con una semilla única
	rand.Seed(time.Now().UnixNano())
	// Genera un número aleatorio dentro del rango
	numeroAleatorio := rand.Intn(limiteSuperior-limiteInferior+1) + limiteInferior

	fmt.Printf("Número aleatorio dentro del rango [%d, %d]: %d\n", limiteInferior, limiteSuperior, numeroAleatorio)

	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	servidores := []string{"50051", "50052", "50053", "50054"}
	var wg sync.WaitGroup
	for _, servidor := range servidores {
		wg.Add(1)
		go enviarMensaje(servidor, numeroAleatorio, &wg)
	}

	wg.Wait()
}
