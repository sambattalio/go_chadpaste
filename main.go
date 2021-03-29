package main

import (
	"fmt"
	"context"
	"time"
	"strings"
	"net/http"
	"encoding/json"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	r := buildRouter()
	http.ListenAndServe(":8000", r)
}

func buildRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/post", createPostHandler).Methods("POST")
	
	return r
}

