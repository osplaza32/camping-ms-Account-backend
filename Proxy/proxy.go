package Proxy

import (
	entityv1 "awesomeProject/gen/bussine"
	healthv1 "awesomeProject/gen/core"
	"awesomeProject/servicegrpc"
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
"google.golang.org/grpc"
	"os"

	"net/http"
)

func Run(s *servicegrpc.Server) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var mux = runtime.NewServeMux(runtime.WithMetadata(s.Theotherfn))
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := entityv1.RegisterEntityserviceAPIHandlerFromEndpoint(ctx, mux,"localhost"+s.GetPort(), opts)
	err = healthv1.RegisterHealthAPIHandlerFromEndpoint(ctx,mux,"localhost"+s.GetPort(),opts)
	if err != nil {
		fmt.Println(err)
		return err
	}
	s.GetLogUber().Info("GET STARED PROXY lsxnlsnxlnlx")

	return http.ListenAndServe(os.Getenv("PORT_HTTP"), mux)
}

