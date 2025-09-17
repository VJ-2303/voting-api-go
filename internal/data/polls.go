package data

import (
	"context"
	"database/sql"
	"errors"
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

type PollWithResults struct {
	*Poll
	Results map[string]int `json:"results"`
}

func ValidatePoll(v *validator.Validator, poll *Poll) {
	v.Check(poll.Title != "", "title", "must be provided")
	v.Check(len(poll.Title) <= 500, "title", "must be less than 500 chars")

	v.Check(poll.Options != nil, "options", "must be provided")
	v.Check(len(poll.Options) >= 2, "options", "must contain at least 2 options")
	v.Check(len(poll.Options) <= 20, "options", "must be less than 20 options")
	v.Check(validator.Unique(poll.Options), "options", "must not contain duplicate values")
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

func (m PollsModel) GetByID(id int64) (*Poll, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `
		SELECT id, created_at, title, description,options,created_by, version
		FROM polls
		WHERE id = $1
			 `
	var poll Poll

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&poll.ID,
		&poll.CreatedAt,
		&poll.Title,
		&poll.Description,
		pq.Array(&poll.Options), // Use pq.Array to scan the text array
		&poll.CreatedBy,
		&poll.Version,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &poll, nil
}

func (m PollsModel) GetWithResults(id int64) (*PollWithResults, error) {

	poll, err := m.GetByID(id)
	if err != nil {
		return nil, err
	}
	query := `
		SELECT chosen_option, count(*)
		FROM votes
		WHERE poll_id = $1
		GROUP BY chosen_option
			 `
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make(map[string]int)

	for rows.Next() {
		var option string
		var count int
		if err := rows.Scan(&option, &count); err != nil {
			return nil, err
		}
		results[option] = count
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	PollWithResults := &PollWithResults{
		Poll:    poll,
		Results: results,
	}
	return PollWithResults, nil
}
