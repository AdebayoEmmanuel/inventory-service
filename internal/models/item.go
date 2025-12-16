package models

type Item struct {
    ID       string `json:"id"`
    Name     string `json:"name"`
    Quantity int    `json:"quantity"`
}