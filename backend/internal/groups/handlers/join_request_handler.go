package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"zone/internal/auth"
	"zone/internal/groups/repositories"
	"zone/internal/groups/services"
	"zone/internal/utils"
)

// URL : /groups/{group_id}/join"

// frontend  send { group_id } -> handler user_id -> get Admin -> notifcation  -> accept -> addmember
// 																			   -> reject ->

/* TO DO
- Check the user isn't already a member
*/

func Join_Request_handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.HandleError(w, 405, http.StatusText(405))
		return
	}
	user_id, err := auth.GetUserIDFromSession(r)
	if err != nil {
		utils.HandleError(w, 401, http.StatusText(401))
		return
	}

	parts := strings.Split((r.URL.Path), "/")

	if len(parts) < 4 {
		utils.HandleError(w, 400, http.StatusText(400))
		return
	}

	group_id, err := strconv.Atoi(parts[2])
	if err != nil {
		utils.HandleError(w, 400, http.StatusText(400))
		return
	}

	if !services.JoinRequestValidator(group_id, user_id) {
		utils.HandleError(w, 400, http.StatusText(400))
		return
	}

	err = repositories.SaveJoinRequest(group_id, user_id)
	if err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}

	// Get the admin from the group_id

	creator_id, err := repositories.GetGroupCreator(group_id)
	if err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}

	fullname    := repositories.GetUserFullName(user_id) 
	
	msg :=  "New join request for your group" + fullname
	// creat the notifcation  
	err  = repositories.CreateNotification(creator_id , "join_request" , nil , msg )
	if err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}

	utils.RespondJSON(
		w,
		http.StatusCreated,
		map[string]string{
			"message": "join request  created successfully!",
		},
	)

}
