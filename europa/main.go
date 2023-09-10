package main

import (
	"context"
	"net"
	"fmt"
	"google.golang.org/grpc"
	pb "github.com/GonzaloDiaz300/LAB/europa/proto"
)

type america struct{
	pb.UnimplementedNotificacionServer
}

func (a *america) Notificar(ctx context.Context, in *pb.NotiReq) (*pb.NotiResp, error) {
	fmt.Printf("Se envia el 1 de vuelta a la central")
	return &pb.NotiResp{Respuesta: 3}, nil
}

func main(){
	listner, err := net.Listen("tcp", ":50053")

	if err != nil {
		panic("cannot create tcp connection" + err.Error())
	}

	serv:= grpc.NewServer()
	fmt.Printf("Servidor Europa Activo\n")
	pb.RegisterNotificacionServer(serv, &america{})
	if err = serv.Serve(listner); err != nil {
		panic("cannot initialize the server" + err.Error())
	}
}