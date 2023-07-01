package server_template

import (
	"context"
	"server-template/internal/models/auth"

	"server-template/internal/models/templates"
	pb "server-template/pkg/proto"
)

type UseCase interface {
	GetSession(ctx context.Context, accessToken string, sessionType uint8) (auth.Session, error)
	HTTPPing(ctx context.Context) (result templates.Templates, err error)
	GRPCPing(ctx context.Context) (res *pb.ServerTemplateResponse, err error)
}
