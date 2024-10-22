package internal

import (
	"context"
	"raccoonstash/internal/repository"
)

func VerifyToken(context context.Context, token string) bool {
	repo := repository.New(DB)

	_, err := repo.GetToken(context, token)
	return err == nil
}
