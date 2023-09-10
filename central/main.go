package main

import (
	"context"
	"fmt"

	pb "github.com/GonzaloDiaz300/LAB/central/proto"
	"google.golang.org/grpc"
)

func main(){
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		panic("cannot connect with server " + err.Error())
	}

	serviceClient := pb.NewNotificacionClient(conn)
	fmt.Printf("Se envia un 100 de parte del cliente\n")
	res, err := serviceClient.Notificar(context.Background(), &pb.NotiReq{Solicitud: 100})
	if err != nil {
		panic("No se llego el mensaje " + err.Error())
	}
	fmt.Printf("Se recibió %d\n", res.Respuesta)
}