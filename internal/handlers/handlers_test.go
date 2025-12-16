package handlers

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestStatusHandler(t *testing.T) {
    req, err := http.NewRequest("GET", "/status", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(StatusHandler)

    handler.ServeHTTP(rr, req)

    assert.Equal(t, http.StatusOK, rr.Code)
    assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
    
    var response map[string]string
    err = json.Unmarshal(rr.Body.Bytes(), &response)
    if err != nil {
        t.Fatal(err)
    }
    
    assert.Equal(t, "healthy", response["status"])
}

func TestItemsHandler(t *testing.T) {
    req, err := http.NewRequest("GET", "/items", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(ItemsHandler)

    handler.ServeHTTP(rr, req)

    assert.Equal(t, http.StatusOK, rr.Code)
    assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
    
    var response []map[string]interface{}
    err = json.Unmarshal(rr.Body.Bytes(), &response)
    if err != nil {
        t.Fatal(err)
    }
    
    assert.Equal(t, 3, len(response))
    assert.Equal(t, "Laptop", response[0]["name"])
}