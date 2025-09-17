package main

import (
	"net/http"

	"github.com/vj-2303/voting-api-go/internal/data"
	"github.com/vj-2303/voting-api-go/internal/validator"
)

func (app *application) createPollHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Options     []string `json:"options"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	user := app.contextGetUser(r)

	poll := &data.Poll{
		Title:       input.Title,
		Description: input.Description,
		Options:     input.Options,
		CreatedBy:   user.ID,
	}
	v := validator.New()
	if data.ValidatePoll(v, poll); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.models.Polls.Insert(poll)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusCreated, envelope{"poll": poll}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
