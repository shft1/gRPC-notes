package notes

import "github.com/shft1/grpc-notes/observability/logger"

type noteUsecase struct {
	log  logger.Logger
	repo noteRepository
}

func NewNotesUseCase(log logger.Logger, repo noteRepository) *noteUsecase {
	return &noteUsecase{
		log:  log,
		repo: repo,
	}
}
