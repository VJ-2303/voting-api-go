package data

import "database/sql"

type Models struct {
	Users UserModel
	Polls PollsModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users: UserModel{
			DB: db,
		},
		Polls: PollsModel{
			DB: db,
		},
	}
}
