package postgres

import (
	"api/internal/domain"
	"api/internal/domain/errdomain"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUserRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresUserRepository(pool *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{pool: pool}
}

// Get возвращает пользователя по его ID.
func (r *PostgresUserRepository) Get(ctx context.Context, id string) (*domain.User, error) {
	const sql = `
		SELECT id, name 
		FROM users WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, sql, id)

	var user domain.User
	if err := row.Scan(&user.ID, &user.Name); err != nil {
		if err == pgx.ErrNoRows {
			return nil, errdomain.ErrObjectNotFound
		}
		return nil, err
	}
	return &user, nil
}

// Search возвращает список пользователей.
func (r *PostgresUserRepository) Search(ctx context.Context, queryOpts domain.QueryOptions) ([]*domain.User, error) {
	const sql = `
		SELECT id, name 
		FROM users
		LIMIT $1
		OFFSET $2
	`

	rows, err := r.pool.Query(ctx, sql, queryOpts.Limit, queryOpts.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Name); err != nil {
			return nil, err
		}

		users = append(users, &user)
	}
	return users, nil
}

// Create добавляет нового пользователя, используя транзакцию для обеспечения атомарности операции.
func (r *PostgresUserRepository) Create(ctx context.Context, user domain.User) (*string, error) {
	var createdUserID string

	err := WithTransaction(ctx, r.pool, func(tx pgx.Tx) error {
		const sql = `
			INSERT INTO users (name, email) 
			VALUES ($1, $2) RETURNING id
		`
		return tx.QueryRow(ctx, sql, user.Name).Scan(&createdUserID)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return &createdUserID, nil
}

// Delete удаляет пользователя по его ID.
func (r *PostgresUserRepository) Delete(ctx context.Context, id string) error {
	const sql = `
		DELETE FROM users 
		WHERE id = $1
	`

	cmdTag, err := r.pool.Exec(ctx, sql, id)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no user found with id %s", id)
	}

	return nil
}
