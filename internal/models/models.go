package models

type Item struct {
	Name  string `json:"name"`
	ID    string `json:"id"`
	Price int    `json:"price"`
}

type Items []Item

type User struct {
	ID      string `json:"id"`
	Balance uint   `json:"balance"`
}

type Users []User
