package consumer

import "github.com/danielMensah/bullhorn-sync-poc/internal/repository"

type Client struct {
	repo repository.Repository
}

func New(repo repository.Repository) *Client {
	return &Client{
		repo: repo,
	}
}
