package consumer

import "github.com/danielMensah/bullhorn-sync-poc/internal/repository"

type Consumer struct {
	repo repository.Repository
}

func New(repo repository.Repository) *Consumer {
	return &Consumer{
		repo: repo,
	}
}
