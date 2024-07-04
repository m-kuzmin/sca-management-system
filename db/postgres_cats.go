package db

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/m-kuzmin/sca-management-system/db/sqlc"
)

// CreateCat inserts a cat.
func (p *Postgres) CreateCat(
	ctx context.Context,
	name string,
	years_of_experience uint16,
	breed string,
	salary uint,
) (uuid.UUID, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create a uuid: %w", err)
	}

	if err = p.queries.CreateCat(ctx, sqlc.CreateCatParams{
		ID:                id,
		Name:              name,
		YearsOfExperience: int16(years_of_experience),
		Breed:             breed,
		Salary:            int32(salary),
	}); err != nil {
		return uuid.Nil, fmt.Errorf("while creating a cat: %w", err)
	}
	return id, nil
}

// GetCat returns the cat object with a matching id.
func (p *Postgres) GetCatByID(ctx context.Context, id uuid.UUID) (Cat, error) {
	cat, err := p.queries.GetCatByID(ctx, id)
	if err != nil {
		return Cat{}, fmt.Errorf("failed to get a cat with id %s: %w", id, err)
	}

	return Cat{
		Name:              cat.Name,
		Breed:             cat.Breed,
		YearsOfExperience: uint32(cat.YearsOfExperience),
		Salary:            uint64(cat.Salary),
		ID:                cat.ID,
	}, nil
}

// GetCatsPaginated returns a list of cats given the page number and amount per page.
func (p *Postgres) GetCatsPaginated(ctx context.Context, amountPerPage, pageNumber uint32) ([]Cat, error) {
	cats, err := p.queries.ListCatsPaginated(ctx, sqlc.ListCatsPaginatedParams{
		Limit:   int32(amountPerPage),
		Column2: pageNumber,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list cats on page %d, with page size %d", pageNumber, amountPerPage)
	}

	result := make([]Cat, len(cats))
	for i, cat := range cats {
		result[i] = Cat{
			Name:              cat.Name,
			Breed:             cat.Breed,
			YearsOfExperience: uint32(cat.YearsOfExperience),
			Salary:            uint64(cat.Salary),
			ID:                cat.ID,
		}
	}

	return result, nil
}

// UpdateCatSalaryByID sets a new salary value for a cat with a matching id. Returns nil on successful update.
func (p *Postgres) UpdateCatSalaryByID(ctx context.Context, id uuid.UUID, salary uint32) error {
	if err := p.queries.UpdateCatSalaryByID(ctx, sqlc.UpdateCatSalaryByIDParams{
		ID:     id,
		Salary: int32(salary),
	}); err != nil {
		return fmt.Errorf("failed to update the salary for cat %s: %w", id, err)
	}

	return nil
}

// DeleteCat deletes a cat. Returns nil on successful deletion.
func (p *Postgres) DeleteCatByID(ctx context.Context, id uuid.UUID) error {
	if err := p.queries.DeteleCatByID(ctx, id); err != nil {
		return fmt.Errorf("failed to delete a cat with id %s: %w", id, err)
	}

	return nil
}
