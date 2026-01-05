package v1

import pb "github.com/shft1/grpc-notes/pkg/api/notes/v1"

type NoteHandler struct {
	pb.UnimplementedNoteAPIServer
	noteUsecase noteUsecase
}

func NewNoteHandler(noteUsecase noteUsecase) *NoteHandler {
	return &NoteHandler{noteUsecase: noteUsecase}
}
