package main

import "net/http"

func (app *application) testAuthHandler(w http.ResponseWriter, r *http.Request) {

	user := app.contextGetUser(r)

	err := app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
