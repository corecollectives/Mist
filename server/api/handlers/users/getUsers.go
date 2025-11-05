package users

import (
	"net/http"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/models"
)

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {

	users, err := models.GetAllUsers()
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to retrieve users", err.Error())
		return
	}
	handlers.SendResponse(w, http.StatusOK, true, users, "Users retrieved successfully", "")
}
