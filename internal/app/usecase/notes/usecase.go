package notes

type useCase struct {
	repo repository
}

func NewUseCase(repo repository) *useCase {
	return &useCase{repo: repo}
}
