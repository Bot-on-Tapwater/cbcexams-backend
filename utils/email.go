package utils

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

// SendPasswordResetEmail sends a password reset email to the specified recipient.
// It generates a reset link using the provided token and the FRONTEND_URL environment variable.
// The email is sent using SMTP credentials and settings defined in environment variables.
//
// Parameters:
//   - email: The recipient's email address.
//   - token: The password reset token to include in the reset link.
//
// Returns:
//   - error: An error if the email could not be sent or if there is an issue with the SMTP configuration.
//
// Environment Variables:
//   - FRONTEND_URL: The base URL of the frontend application.
//   - SMTP_FROM: The sender's email address.
//   - SMTP_HOST: The SMTP server host.
//   - SMTP_PORT: The SMTP server port (must be a valid integer).
//   - SMTP_USER: The SMTP username.
//   - SMTP_PASS: The SMTP password.
func SendPasswordResetEmail(email, token string) error {
	resetLink := fmt.Sprintf("%s/reset-password?token=%s", os.Getenv("FRONTEND_URL"), token)
	body := fmt.Sprintf("Click <a href='%s'>here</a> to reset your password. Link expires in 1 hour.", resetLink)

	/* Run the email-sending process in a Goroutine */
	go func() {
		if err := SendEmail(email, "Password Reset Request", body); err != nil {
			fmt.Printf("Failed to send password reset email: %s\n", err)
		}
	}()

	/* Return immediately without waiting for the email to be sent */
	return nil
}

func SendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("SMTP_FROM"))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	/* Convert SMTP_PORT to an integer */
	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return fmt.Errorf("invalid SMTP_PORT: %v", err)
	}

	d := gomail.NewDialer(
		os.Getenv("SMTP_HOST"),
		port,
		os.Getenv("SMTP_USER"),
		os.Getenv("SMTP_PASS"),
	)

	/* Attempt to send the email and log any errors */
	if err := d.DialAndSend(m); err != nil {
		fmt.Printf("Failed to send email: %s\n", err) // Log the error
		return fmt.Errorf("failed to send email: %v", err)
	}
	return nil
}
