package server

import pb "github.com/shft1/grpc-notes/pkg/api/notes/v1"

func (gs *grpcServer) RegisterRPC(noteHand pb.NoteAPIServer) {
	pb.RegisterNoteAPIServer(gs.srv, noteHand)
}
