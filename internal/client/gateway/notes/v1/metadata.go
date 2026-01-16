package v1

import (
	"context"

	"google.golang.org/grpc/metadata"
)

func (gw *noteGateway) createMetadata(token string) metadata.MD {
	md := metadata.New(map[string]string{
		"authorization": token,
	})
	return md
}

func (gw *noteGateway) putMetadataToContext(ctx context.Context, md metadata.MD) context.Context {
	return metadata.NewOutgoingContext(ctx, md)
}
