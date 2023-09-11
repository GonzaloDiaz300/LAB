package main

import (
	"bufio"
	"context"
	"fmt"
	"math"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"

	pb "github.com/GonzaloDiaz300/LAB/america/proto"
	"google.golang.org/grpc"
)

var interesados int
var interesados_iniciales int

type europa struct {
	pb.UnimplementedNotificacionServer
}

func (a *europa) Notificar(ctx context.Context, in *pb.NotiReq) (*pb.NotiResp, error) {
	fmt.Printf("Se envia el 3 de vuelta a la central\n")
	crearInteresados(0)
	return &pb.NotiResp{Respuesta: 3}, nil
}

func crearInteresados(no_registrados int) {

	if interesados_iniciales == 0 {
		fileName := "america/parametros_de_inicio.txt"

		// Intenta abrir el archivo
		file, err := os.Open(fileName)
		if err != nil {
			fmt.Println("Error al abrir el archivo:", err)
			return
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
				return
			}
			interesados_iniciales = intValue
			// Almacenar el valor entero globalmente
		}
		interesados = interesados_iniciales / 2
	} else {
		interesados = (interesados_iniciales - (interesados - no_registrados)) / 2
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
}

func main() {
	listner, err := net.Listen("tcp", ":50053")

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
