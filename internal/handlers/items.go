package handlers

import (
    "encoding/json"
    "net/http"

    "github.com/AdebayoEmmanuel/inventory-service/internal/models"
)

func ItemsHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    
    items := []models.Item{
        {ID: "1", Name: "Laptop", Quantity: 10},
        {ID: "2", Name: "Keyboard", Quantity: 25},
        {ID: "3", Name: "Phone", Quantity: 15},
    }
    
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(items)
}