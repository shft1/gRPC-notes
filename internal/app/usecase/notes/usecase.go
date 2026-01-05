package notes

type noteUsecase struct {
	repo noteRepository
}

func NewNotesUseCase(repo noteRepository) *noteUsecase {
	return &noteUsecase{repo: repo}
}
