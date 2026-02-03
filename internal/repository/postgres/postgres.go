package postgres

import (
	"log"
)

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) SaveIssue(issue interface{}) error {
	log.Printf("POSTGRES: Saving issue: %+v", issue)
	return nil
}
