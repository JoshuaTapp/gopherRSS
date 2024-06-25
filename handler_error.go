package main

import (
	"net/http"
)

func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
}
