package data

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/vj-2303/voting-api-go/internal/validator"
)

var (
	ErrDuplicateVote = errors.New("user have already voted this poll")
)

type Vote struct {
	ID           int64     `json:"id"`
	PollID       int64     `json:"poll_id"`
	UserID       int64     `json:"user_id"`
	ChosenOption string    `json:"chosen_option"`
	CreatedAt    time.Time `json:"created_at"`
}

func ValidateVote(v *validator.Validator, chosenOption string, pollOptions []string) {
	v.Check(chosenOption != "", "option", "must be provided")
	v.Check(validator.In(chosenOption, pollOptions...), "option", "is not a valid poll option")
}

type VotesModel struct {
	DB *sql.DB
}

func (m VotesModel) Insert(vote *Vote) error {
	query := `
		INSERT INTO votes(poll_id,user_id,chosen_option)
		VALUES($1,$2,$3)
		RETURNING id, created_at
			 `
	args := []any{vote.PollID, vote.UserID, vote.ChosenOption}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&vote.ID,
		&vote.CreatedAt,
	)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return ErrDuplicateVote
		}
		return err
	}
	return nil
}
