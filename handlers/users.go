package handlers

import (
	"net/http"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserHandler struct {
	db *mongo.Database
}

func NewUserHandler(db *mongo.Database) *UserHandler {
	return &UserHandler{db}
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request)    {}
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request)   {}
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {}
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {}
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {}
