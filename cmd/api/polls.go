package main

import (
	"errors"
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

func (app *application) showPollHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
	}
	poll, err := app.models.Polls.GetByID(id)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"poll": poll}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) castVoteHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	poll, err := app.models.Polls.GetByID(id)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	var input struct {
		Option string `json:"option"`
	}
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	if data.ValidateVote(v, input.Option, poll.Options); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	user := app.contextGetUser(r)

	vote := &data.Vote{
		PollID:       poll.ID,
		UserID:       user.ID,
		ChosenOption: input.Option,
	}
	err = app.models.Votes.Insert(vote)
	if err != nil {
		if errors.Is(err, data.ErrDuplicateVote) {
			app.errorResponse(w, r, http.StatusConflict, "you have already voted on this poll")
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusCreated, envelope{"vote": vote}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showPollResultsHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	poll, err := app.models.Polls.GetWithResults(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"poll": poll}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
