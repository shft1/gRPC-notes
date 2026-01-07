package v1

import (
	"github.com/shft1/grpc-notes/observability/logger"
	pb "github.com/shft1/grpc-notes/pkg/api/notes/v1"
)

type noteGateway struct {
	log    logger.Logger
	client pb.NoteAPIClient
}

func NewNoteGateway(log logger.Logger, client pb.NoteAPIClient) *noteGateway {
	return &noteGateway{
		log:    log,
		client: client,
	}
}
