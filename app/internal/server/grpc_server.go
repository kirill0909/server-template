package server

import (
	"auth-svc/config"
	"auth-svc/pkg/logger"
	"net"

	pb "gitlab.axarea.ru/main/aiforypay/package/auth-svc-proto"
	"google.golang.org/grpc"
)

type Deps struct {
	AuthHandlers pb.AuthServer
}

type GRPCServer struct {
	srv  *grpc.Server
	deps Deps
	cfg  *config.Config
	log  logger.Logger
}

func NewGRPCServer(srv *grpc.Server, deps Deps, cfg *config.Config, logg logger.Logger) *GRPCServer {
	return &GRPCServer{
		srv:  srv,
		deps: deps,
		cfg:  cfg,
		log:  logg,
	}
}

func (g *GRPCServer) Run() error {
	pb.RegisterAuthServer(g.srv, g.deps.AuthHandlers)

	l, err := net.Listen("tcp", g.cfg.Server.Host+":"+g.cfg.Server.GRPCPort)
	if err != nil {
		return err
	}

	if err := g.srv.Serve(l); err != nil {
		return err
	}

	return nil
}

func (g *GRPCServer) GracefulShutdown() {
	g.srv.GracefulStop()
}
