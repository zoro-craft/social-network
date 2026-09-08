package users

import (
    "net/http"
    "zone/internal/database"
    "zone/internal/utils"
)


type User struct {
    Full_name string `json:"full_name"`
    Avatar    string `json:"avatar,omitempty"`
}

func GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
        return
    }

    users, err := GetAllUsers()
    if err != nil {
        utils.HandleError(w, http.StatusInternalServerError, "Could not load users")
        return
    }

    utils.RespondJSON(w, http.StatusOK, users)
}

func GetAllUsers() ([]User, error) {
    rows, err := database.Database.Query(`SELECT full_name, COALESCE(avatar,'') FROM users ORDER BY id`)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []User
    for rows.Next() {
        var u User
        if err := rows.Scan(&u.Full_name, &u.Avatar); err != nil {
            return nil, err
        }
        users = append(users, u)
    }

    return users, rows.Err()
}
