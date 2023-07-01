package server

import (
	"net"
	"server-template/config"
	"server-template/pkg/logger"

	"google.golang.org/grpc"
	pb "server-template/pkg/proto"
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
