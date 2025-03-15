package repository

import (
	"context"

	"github.com/quydmfl/niveau-test/internal/model"
)

type DocumentsRepository interface {
	Create(ctx context.Context, document *model.Documents) error
}

func NewDocumentsRepository(
	repository *Repository,
) DocumentsRepository {
	return &documentsRepository{
		Repository: repository,
	}
}

type documentsRepository struct {
	*Repository
}

func (r *documentsRepository) Create(ctx context.Context, document *model.Documents) error {
	if err := r.DB(ctx).Create(document).Error; err != nil {
		return err
	}
	return nil
}
