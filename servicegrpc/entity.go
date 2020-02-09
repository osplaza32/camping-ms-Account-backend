package servicegrpc

import (
	entityv1 "awesomeProject/gen/bussine"
	healthv1 "awesomeProject/gen/core"
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	//"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"net"
	"net/http"
	"os"
)


type Server struct {
	addr             string
	// mejorar log
	loguber          *zap.Logger

}
func (s *Server) Check(ctx context.Context,req *healthv1.CheckRequest) (*healthv1.CheckResponse, error) {
	return &healthv1.CheckResponse{Status:healthv1.CheckResponse_SERVING},nil
}

func (s *Server) Watch(cctx *healthv1.CheckRequest,req healthv1.HealthAPI_WatchServer) error {
	panic("implement me")
}

func (s *Server)Theotherfn(ctx context.Context,req *http.Request) metadata.MD{
	s.GetLogUber().Info("GET entro GRPC")

	as := metadata.MD{}
	as.Set("apikey",req.Header.Get("apikey"))
	return as
}

func (s *Server) Entity(ctx context.Context, req *entityv1.EntityRequest) (*entityv1.EntityResponse, error) {
	panic("implement me")
}

func (s *Server) GetEntity(ctx context.Context,req *entityv1.GetEntityRequest) (*entityv1.GetEntityResponse, error) {
	return &entityv1.GetEntityResponse{MyBool:true,Code:"200"},nil

}

type ServerOption func(*Server) error

func NewServer(opts ...ServerOption) (*Server, error) {
	logger, _ := zap.NewProduction()

	s := &Server{
		addr: os.Getenv("PORT_GRPC"),
		loguber:logger,
	}
	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}
	return s, nil
}

func SetAddr(addr string) ServerOption {
	return func(s *Server) error {
		s.addr = addr
		return nil
	}
}

func (s *Server) Start() (*grpc.Server, net.Listener, error) {
	s.GetLogUber().Info("GET STARED GRPC")

	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return nil, nil, err
	}
	gs := grpc.NewServer(grpc.UnaryInterceptor(s.unaryInterceptor))
	return gs, lis,nil
}
func (s *Server)unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error)  {
	//meta, ok := metadata.FromIncomingContext(ctx)
	/*
	if !ok {
		return nil, grpc.Errorf(codes.Unauthenticated, "missing context metadata")
	}
	s.loguber.Sugar().Info(meta)

	if len(meta["apikey"]) != 1 {
		return nil, grpc.Errorf(codes.Unauthenticated, "invalid token")
	}
	if meta["apikey"][0] != "123456" {
		return nil, grpc.Errorf(codes.Unauthenticated, "invalid token")
	}*/

	return handler(ctx, req)
}
func (s *Server) GetPort() string{
	return s.addr
}
func (s *Server) GetLogUber() *zap.Logger {
	return s.loguber
}
