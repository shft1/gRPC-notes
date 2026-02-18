package v1

import (
	"embed"
	"net/http"

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
		r.Post("/chat", nr.hand.Chat)
	})
}

func (nr *NoteRouter) SetupGenRoutesV1(gwMux http.Handler, swaggerUI, specs embed.FS) {
	nr.main.Route("/gen/v1", func(r chi.Router) {
		r.Route("/swagger", func(r chi.Router) {
			r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
				http.FileServerFS(swaggerUI).ServeHTTP(w, r)
			})
			r.Get("/specs/*", func(w http.ResponseWriter, r *http.Request) {
				http.StripPrefix("/specs", http.FileServerFS(specs)).ServeHTTP(w, r)
			})
		})
		r.Mount("/", gwMux)
	})
}
