package main

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
	"strconv"

	pb "github.com/GonzaloDiaz300/LAB/europa/proto"
	"google.golang.org/grpc"
)

type europa struct {
	pb.UnimplementedNotificacionServer
}

func (a *europa) Notificar(ctx context.Context, in *pb.NotiReq) (*pb.NotiResp, error) {
	fmt.Printf("Se envia el 1 de vuelta a la central")
	return &pb.NotiResp{Respuesta: 3}, nil
}

func main() {
	var interesados int

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

		// Almacenar el valor entero globalmente
		interesados = intValue
	}

	// Imprimir el valor almacenado globalmente
	fmt.Println("Valor entero global:", interesados)

	/////////////////////////////////////////////////////////////////////////////////////////////////////////////
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
