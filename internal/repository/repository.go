package repository

type Deps struct {
}

type Repository struct {
}

func NewService(deps *Deps) *Repository {
	return &Repository{}
}
