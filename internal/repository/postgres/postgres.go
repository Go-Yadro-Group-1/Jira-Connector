package postgres

type PostgresRepository struct {
}

func New() (*PostgresRepository, error) {
	return &PostgresRepository{}, nil
}
