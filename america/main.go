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

var interesados int // Variable global para modificar los interesados en obtener la key
var interesados_actuales = 0

type america struct {
	pb.UnimplementedNotificacionServer
}

func (a *america) Notificar(ctx context.Context, in *pb.NotiReq) (*pb.NotiResp, error) {
	fmt.Printf("Se envia el 1 de vuelta a la central\n")
	crearInteresados(0) // Acá en vez de 0 tendría que ser el valor de usuarios sin registrar que devuelve la central, lo que se dá con la asincrona
	return &pb.NotiResp{Respuesta: 1}, nil
}

// funcion para generar el numero de interesados en cada iteración, se llama cuando llaman la funcion de notificar
func crearInteresados(no_registrados int) {
	if interesados_actuales == 0 {
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
}

func main() {

	listner, err := net.Listen("tcp", ":50051")

	if err != nil {
		panic("cannot create tcp connection" + err.Error())
	}

	serv := grpc.NewServer()
	fmt.Printf("Servidor America Activo\n")
	pb.RegisterNotificacionServer(serv, &america{})
	if err = serv.Serve(listner); err != nil {
		panic("cannot initialize the server" + err.Error())
	}
}
