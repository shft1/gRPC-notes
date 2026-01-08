package v1

import "github.com/shft1/grpc-notes/internal/domain/notes"

func toDomainCreate(in *noteCreateRequest) *notes.NoteCreate {
	return &notes.NoteCreate{
		Title: in.Title,
		Desc:  in.Desc,
	}
}

func toDTOResponse(item *notes.Note) *noteResponse {
	return &noteResponse{
		UUID:      item.UUID.String(),
		Title:     item.Title,
		Desc:      item.Desc,
		IsDel:     item.IsDel,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
}

func toDTOResponseList(items []*notes.Note) []*noteResponse {
	out := make([]*noteResponse, 0, len(items))
	for _, item := range items {
		out = append(out, toDTOResponse(item))
	}
	return out
}
