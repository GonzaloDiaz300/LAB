package main

import (
	"context"
	"fmt"
	"sync"
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

func main(){
	servidores := []string{"50051", "50052", "50053", "50054"}
	var wg sync.WaitGroup
	mensaje := 300
	for _, servidor := range servidores {
        wg.Add(1)
        go enviarMensaje(servidor, mensaje, &wg)
    }

    wg.Wait()
}