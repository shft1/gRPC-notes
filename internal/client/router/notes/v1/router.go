package v1

import (
	"github.com/go-chi/chi"
	v1 "github.com/shft1/grpc-notes/internal/client/handler/notes/v1"
)

type NoteRouter struct {
	main *chi.Mux
	hand *v1.NoteHandler
}

func NewNoteRouter(mr *chi.Mux, hand *v1.NoteHandler) *NoteRouter {
	return &NoteRouter{
		main: mr,
		hand: hand,
	}
}

func (nr *NoteRouter) SetupRoutesV1() {
	nr.main.Route("/notes/v1", func(r chi.Router) {
		r.Get("/", nr.hand.GetMulti)
		r.Get("/{uuid}", nr.hand.GetByID)
		r.Post("/", nr.hand.Create)
		r.Delete("/{uuid}", nr.hand.DeleteByID)
		r.Get("/subscribe", nr.hand.SubscribeToEvents)
		r.Post("/metrics/upload", nr.hand.UploadMetrics)
	})
}
