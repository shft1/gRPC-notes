package v1

import (
	"github.com/shft1/grpc-notes/internal/domain/notes"
	pb "github.com/shft1/grpc-notes/pkg/api/notes/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func toDomainCreate(in *pb.NoteCreateRequest) *notes.NoteCreate {
	return &notes.NoteCreate{
		Title: in.Title,
		Desc:  in.Desc,
	}
}

func toDTOResponse(item *notes.Note) *pb.Note {
	return &pb.Note{
		Uuid:      item.UUID.String(),
		Title:     item.Title,
		Desc:      item.Desc,
		IsDel:     item.IsDel,
		CreatedAt: timestamppb.New(item.CreatedAt),
		UpdatedAt: timestamppb.New(item.UpdatedAt),
	}
}

func toDTOListResponse(items []*notes.Note) *pb.NoteList {
	dtoList := make([]*pb.Note, 0, len(items))
	for _, item := range items {
		dtoList = append(dtoList, toDTOResponse(item))
	}
	return &pb.NoteList{Notes: dtoList}
}
