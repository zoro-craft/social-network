package models 

type Groupe struct {
	Id int `json:"id"`
	Admin int  `json:"admin"`
	Titel string `json:"titel"` 
	Description string `json:"description"`
}