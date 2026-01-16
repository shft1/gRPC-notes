package server

import pb "github.com/shft1/grpc-notes/pkg/api/notes/v1"

func (gs *grpcServer) RegisterRPC(noteRPC pb.NoteAPIServer) {
	pb.RegisterNoteAPIServer(gs.srv, noteRPC)
}
