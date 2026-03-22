package repository

import (
	"context"

	"github.com/Murolando/ftem_wallet/pkg/config"
)

var _ Repository = (*RepositoryImpl)(nil)

type RepositoryImpl struct {
}

func NewRepositoryImpl(ctx context.Context, cfg *config.DBConfig) (*RepositoryImpl, error) {

	return &RepositoryImpl{}, nil
}
