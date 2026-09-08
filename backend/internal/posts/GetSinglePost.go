package posts

import (
	"net/http"
	"strconv"
	"zone/internal/auth"
	"zone/internal/database"
	"zone/internal/utils"
)

func GetPostByIDAPI(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodGet {
    utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
    return
  }

  postID, err := strconv.Atoi(r.PathValue("id"))
  if err != nil || postID <= 0 {
    utils.HandleError(w, http.StatusBadRequest, "Invalid post id")
    return
  }

  viewerID, _ := auth.GetUserID(r)
  var p Post
  
  query := postSelect  + fromPosts + "WHERE p.id = ? AND (" + visibleCondition + ")"

  err = database.Database.QueryRow(query, postID, viewerID, viewerID, viewerID).Scan(
    &p.ID, &p.Title, &p.Content, &p.Image, &p.Privacy,
    &p.GroupID, &p.UserID, &p.Full_name, &p.CreatedAt,
  )
  
  if err != nil {
    utils.HandleError(w, http.StatusNotFound, "Post not found or unauthorized")
    return
  }
	
  utils.RespondJSON(w, http.StatusOK, p)
}