package v1

import (
	"github.com/google/uuid"
	"github.com/shft1/grpc-notes/internal/domain/notes"
	pb "github.com/shft1/grpc-notes/pkg/api/notes/v1"
)

func toDTORequest(in *notes.NoteCreate) *pb.NoteCreateRequest {
	return &pb.NoteCreateRequest{
		Title: in.Title,
		Desc:  in.Desc,
	}
}

func toDomainResponse(in *pb.Note) *notes.Note {
	id, _ := uuid.Parse(in.Id)
	return &notes.Note{
		ID:        id,
		Title:     in.Title,
		Desc:      in.Desc,
		IsDel:     in.IsDel,
		CreatedAt: in.CreatedAt.AsTime(),
		UpdatedAt: in.UpdatedAt.AsTime(),
	}
}

func toDomainListResponse(in []*pb.Note) []*notes.Note {
	items := make([]*notes.Note, 0, len(in))
	for _, dto := range in {
		items = append(items, toDomainResponse(dto))
	}
	return items
}
