package main

import (
	"awesomeProject/Proxy"
	"awesomeProject/external/sbuslistener"
	entityv1 "awesomeProject/gen/bussine"
	healthv1 "awesomeProject/gen/core"
	"awesomeProject/internal/utils"
	"awesomeProject/servicegrpc"
	"context"
	servicebus "github.com/Azure/azure-service-bus-go"
	"google.golang.org/grpc/reflection"
	"os"

	"fmt"
)
type worker struct {
	source chan interface{}
	quit chan struct{}
}
func main() {
	utils.LoadENV()
	c := make(chan *servicebus.Message)
	go sbuslistener.MakeListner(os.Getenv("CONN_STRING"),c)
	ser,err := servicegrpc.NewServer()
	if err != nil {
		fmt.Println(err)
	}
	grsp,let,erro:=ser.Start()
	if erro != nil {
		fmt.Println(err)
	}
	go Proxy.Run(ser)

	entityv1.RegisterEntityserviceAPIServer(grsp,ser)
	healthv1.RegisterHealthAPIServer(grsp,ser)
	reflection.Register(grsp)
	ser.GetLogUber().Info("GET STARED SERVER")
	go func() {
		for {
			msg := <- c
			fmt.Println(msg)
			msg.Complete(context.Background())
		}
	}()

	if err := grsp.Serve(let); err != nil {

	}
}