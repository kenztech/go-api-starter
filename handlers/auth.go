package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/kenztech/go-api-starter/models"
	"github.com/kenztech/go-api-starter/models/api"
	"github.com/kenztech/go-api-starter/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AuthHandler struct {
	db *mongo.Database
}

func NewAuthHandler(db *mongo.Database) *AuthHandler {
	return &AuthHandler{db}
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	// Retrieve user data from context (assumes user data is stored in the context)
	data, ok := utils.GetUserDataFromContext(r.Context())
	if !ok {
		utils.SendError(w, "Unable to retrieve user from context", http.StatusInternalServerError)
		return
	}

	// Log email to debug
	log.Printf("User email retrieved from context: %s", data.Email)

	// Get user from MongoDB by email
	var user models.User
	collection := h.db.Collection("users")
	filter := bson.M{"email": data.Email}

	err := collection.FindOne(r.Context(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.SendError(w, "User not found", http.StatusNotFound)
		} else {
			utils.SendError(w, "Error retrieving user", http.StatusInternalServerError)
		}
		return
	}

	// Prepare the response object
	response := api.UserResponse{
		Success: true,
		User: api.UserData{
			ID:     uint(user.ID[0]),
			Role:   user.Role,
			Email:  user.Email,
			Status: user.Status,
		},
	}

	// Send the JSON response with user details
	utils.SendJSON(w, http.StatusOK, response)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var request api.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.SendError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if !utils.ValidateStruct(w, request) {
		return
	}

	// Get user from MongoDB
	var user models.User
	filter := bson.M{}

	if request.Email != "" {
		filter["email"] = request.Email
	} else {
		filter["username"] = request.Username
	}

	log.Printf("User to find: Email: %s, Username: %s, Password: %s", request.Email, request.Username, request.Password)
	log.Printf("Constructed filter: %+v", filter)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := h.db.Collection("users").FindOne(ctx, filter).Decode(&user)
	if err != nil {
		log.Printf("Error finding user: %+v", err)
		utils.SendError(w, "User not found", http.StatusUnauthorized)
		return
	}

	if !utils.ComparePassword(user.Password, request.Password) {
		utils.SendError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateRefreshToken(user.Email, user.Role)
	if err != nil {
		utils.SendError(w, "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	response := api.UserResponse{
		Success: true,
		User: api.UserData{
			ID:       uint(user.ID[0]),
			Role:     user.Role,
			Email:    user.Email,
			Status:   user.Status,
			Username: user.Username,
		},
	}

	http.SetCookie(w, &http.Cookie{
		Path:     "/",
		Name:     "token",
		Value:    token,
		Secure:   true,
		SameSite: http.SameSiteDefaultMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	utils.SendJSON(w, http.StatusOK, response)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Path:     "/",
	})

	res := api.SuccessResponse{
		Success: true,
		Message: "Logged out successfully.",
	}

	utils.SendJSON(w, http.StatusOK, res)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {}
