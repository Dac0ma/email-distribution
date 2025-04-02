package main

import (
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
    "sync"
)

type Item struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

var items = make(map[string]Item)
var mu sync.Mutex

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/items", getItems).Methods("GET")
    r.HandleFunc("/items/{id}", getItem).Methods("GET")
    r.HandleFunc("/items", createItem).Methods("POST")
    r.HandleFunc("/items/{id}", updateItem).Methods("PUT")
    r.HandleFunc("/items/{id}", deleteItem).Methods("DELETE")

    http.ListenAndServe(":8080", r)
}

func getItems(w http.ResponseWriter, r *http.Request) {
    mu.Lock()
    defer mu.Unlock()
    var result []Item
    for _, item := range items {
        result = append(result, item)
    }
    json.NewEncoder(w).Encode(result)
}

func getItem(w http.ResponseWriter, r *http.Request) {
    mu.Lock()
    defer mu.Unlock()
    vars := mux.Vars(r)
    item, exists := items[vars["id"]]
    if !exists {
        http.Error(w, "Item not found", http.StatusNotFound)
        return
    }
    json.NewEncoder(w).Encode(item)
}

func createItem(w http.ResponseWriter, r *http.Request) {
    mu.Lock()
    defer mu.Unlock()
    var item Item
    json.NewDecoder(r.Body).Decode(&item)
    items[item.ID] = item
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(item)
}

func updateItem(w http.ResponseWriter, r *http.Request) {
    mu.Lock()
    defer mu.Unlock()
    vars := mux.Vars(r)
    var item Item
    json.NewDecoder(r.Body).Decode(&item)
    item.ID = vars["id"]
    items[item.ID] = item
    json.NewEncoder(w).Encode(item)
}

func deleteItem(w http.ResponseWriter, r *http.Request) {
    mu.Lock()
    defer mu.Unlock()
    vars := mux.Vars(r)
    delete(items, vars["id"])
    w.WriteHeader(http.StatusNoContent)
}