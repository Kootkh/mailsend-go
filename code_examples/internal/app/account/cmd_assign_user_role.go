package account

import (
	"context"
	"errors"
	"fmt"

	"github.com/samverrall/gopherjobs/internal/app"
	"github.com/samverrall/gopherjobs/internal/repository"
)

func (s *Service) AssignUserRole(ctx context.Context, userID string, roleID int) error {

	role, err := RoleFormID(roleID)

	if err != nil {
		return fmt.Errorf("%w: %w", app.ErrInvalidInput, err)
	}

	_, err = s.repo.GetUserByUUID(ctx, userID)

	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return (err)
		}

		return fmt.Errorf("failed to get user: %w", err)
	}

	return nil
}
