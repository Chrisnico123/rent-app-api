package auth

import (
	"errors"
	"log"
	"regexp"
)

type OTP struct {
	Email string `json:"email"`
	Code  string `json:"code"`
	Level int    `json:"level"`
}

type OTPReq struct {
	Email string `json:"email"`
	Level int    `json:"level"`
}

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (o OTP) Validate() error {
	if o.Email == "" {
		log.Println("[ERROR] OTP validation: email is empty")
		return errors.New("email is required")
	}
	if o.Code == "" {
		log.Println("[ERROR] OTP validation: code is empty")
		return errors.New("code is required")
	}

	emailRegex := regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
	if !emailRegex.MatchString(o.Email) {
		log.Println("[ERROR] OTP validation: invalid email format")
		return errors.New("invalid email format")
	}

	if len(o.Code) != 4 {
		log.Println("[ERROR] OTP validation: code must be 4 digits")
		return errors.New("code must be 4 digits")
	}
	return nil
}

func (o OTPReq) Validate() error {
	if o.Email == "" {
		log.Println("[ERROR] OTPReq validation: email is empty")
		return errors.New("email is required")
	}
	emailRegex := regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
	if !emailRegex.MatchString(o.Email) {
		log.Println("[ERROR] OTPReq validation: invalid email format")
		return errors.New("invalid email format")
	}
	return nil
}

func (u User) Validate() error {
	if u.Email == "" {
		log.Println("[ERROR] User validation: email is empty")
		return errors.New("email is required")
	}
	if u.Name == "" {
		log.Println("[ERROR] User validation: name is empty")
		return errors.New("name is required")
	}

	emailRegex := regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
	if !emailRegex.MatchString(u.Email) {
		log.Println("[ERROR] User validation: invalid email format")
		return errors.New("invalid email format")
	}
	return nil
}
