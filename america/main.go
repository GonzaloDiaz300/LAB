package main

import (
	"context"
	"net"
	"fmt"
	"google.golang.org/grpc"
	pb "github.com/GonzaloDiaz300/LAB/america/proto"
)

type america struct{
	pb.UnimplementedNotificacionServer
}

func (a *america) Notificar(ctx context.Context, in *pb.NotiReq) (*pb.NotiResp, error) {
	fmt.Printf("Se envia el 1 de vuelta a la central")
	return &pb.NotiResp{Respuesta: 1}, nil
}

func main(){
	listner, err := net.Listen("tcp", ":50051")

	if err != nil {
		panic("cannot create tcp connection" + err.Error())
	}

	serv:= grpc.NewServer()
	fmt.Printf("Servidor America Activo\n")
	pb.RegisterNotificacionServer(serv, &america{})
	if err = serv.Serve(listner); err != nil {
		panic("cannot initialize the server" + err.Error())
	}
}