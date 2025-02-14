package models

import (
	"context"
	"log"
	"time"

	"github.com/kenztech/go-api-starter/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var DB *mongo.Database

func ConnectDB() (*mongo.Database, error) {
	utils.LoadEnv()

	db := utils.GetEnv("MONGO_DB", "test")
	uri := utils.GetEnv("MONGO_URI", "mongodb://localhost:27017")

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	// Ensure the connection is successful
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("Error pinging MongoDB:", err)
		return nil, err
	}

	log.Println("Connected to MongoDB successfully!")
	DB := client.Database(db)

	// Ensure an admin user exists
	intitAdminUser(DB)

	return DB, nil
}

// Ensure an admin user exists in the database
func intitAdminUser(db *mongo.Database) {
	collection := db.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if an admin user already exists
	var existingUser User
	err := collection.FindOne(ctx, bson.M{"role": "admin"}).Decode(&existingUser)
	if err == mongo.ErrNoDocuments {
		// Create an admin user if none exists
		hashedPassword, err := utils.HashPassword("Admin@kenz.io")
		if err != nil {
			log.Fatal("Error hashing admin password:", err)
		}

		admin := User{
			ID:       primitive.NewObjectID(),
			Name:     "Admin",
			Role:     "admin",
			Username: "admin",
			Status:   "active",
			Email:    "admin@kenz.io",
			Password: string(hashedPassword),
		}

		_, err = collection.InsertOne(ctx, admin)
		if err != nil {
			log.Println("Error creating admin user:", err)
		} else {
			log.Println("Admin user created successfully!")
		}
	} else if err != nil {
		log.Println("Error checking for admin user:", err)
	}
}
