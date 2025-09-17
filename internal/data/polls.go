package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
	"github.com/vj-2303/voting-api-go/internal/validator"
)

type Poll struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Options     []string  `json:"options"`
	CreatedBy   int64     `json:"created_by"`
	Version     int       `json:"version"`
}

type PollsModel struct {
	DB *sql.DB
}

func (m PollsModel) Insert(poll *Poll) error {
	query := `
		INSERT INTO polls(title, description, options, created_by)
		VALUES ($1,$2,$3,$4)
		RETURNING id, created_at, version
			 `
	args := []any{poll.Title, poll.Description, pq.Array(poll.Options), poll.CreatedBy}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&poll.ID, &poll.CreatedAt, &poll.Version)
}

func ValidatePoll(v *validator.Validator, poll *Poll) {
	v.Check(poll.Title != "", "title", "must be provided")
	v.Check(len(poll.Title) <= 500, "title", "must be less than 500 chars")

	v.Check(poll.Options != nil, "options", "must be provided")
	v.Check(len(poll.Options) >= 2, "options", "must contain at least 2 options")
	v.Check(len(poll.Options) <= 20, "options", "must be less than 20 options")
	v.Check(validator.Unique(poll.Options), "options", "must not contain duplicate values")
}
