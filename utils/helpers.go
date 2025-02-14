package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"

	"github.com/jordan-wright/email"
	"github.com/kenztech/go-api-starter/models/api"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
)

type EmailSender struct {
	SMTPPort   int
	SMTPServer string
	Username   string
	Password   string
}

func SendMail(to, subject, message string) error {
	log.Println("Creating email message...")

	// Create a new email message
	msg := gomail.NewMessage()
	msg.SetHeader("From", "Zemen Innovation<gxforce430@gmail.com>")
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/plain", message)

	// Set up the SMTP dialer
	dialer := gomail.NewDialer("smtp.gmail.com", 587, "gxforce430@gmail.com", "odte rbqr swps ktpm")

	// Attempt to send the email
	log.Println("Attempting to send email...")
	err := dialer.DialAndSend(msg)
	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}

	log.Println("Email sent successfully!")
	return nil
}

func SendError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := api.ErrorResponse{
		Success: false,
		Message: message,
	}
	json.NewEncoder(w).Encode(response)
}

func SendJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func SendEmail(to, subject, body string) error {
	log.Printf("Simulating sending email...\nTo: %s\nSubject: %s\nBody: %s\n", to, subject, body)
	fmt.Println("Email simulated successfully")
	return nil
}

func SendHTMLEmail(to, subject, htmlContent string, cc []string, attachments ...string) error {
	e := email.NewEmail()
	e.From = "Zemen Innovation <gxforce430@gmail.com>"
	e.To = []string{to}
	e.Subject = subject
	e.HTML = []byte(htmlContent)

	if cc != nil {
		e.Cc = cc
	}

	err := e.Send("smtp.gmail.com:587", smtp.PlainAuth("", "gxforce430@gmail.com", "fuvw ylpy xqgn iish", "smtp.gmail.com"))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}

// HashPassword hashes a plain text password using bcrypt.
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// ComparePassword checks if the provided password matches the hashed password.
func ComparePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
