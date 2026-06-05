package auth

import (
	"regexp"
	"strings"
)

var (
	emailPattern    = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
	usernamePattern = regexp.MustCompile(`^[a-z0-9_]+$`)
	letterPattern   = regexp.MustCompile(`[A-Za-z]`)
	numberPattern   = regexp.MustCompile(`[0-9]`)
)

type ValidationError struct {
	Fields map[string]string
}

func (e ValidationError) Error() string {
	return "validation failed"
}

func NormalizeRegisterRequest(req RegisterRequest) RegisterRequest {
	req.Name = strings.TrimSpace(req.Name)
	req.Email = normalizeEmail(req.Email)
	req.Username = normalizeOptional(req.Username)

	return req
}

func NormalizeLoginRequest(req LoginRequest) LoginRequest {
	req.Email = normalizeEmail(req.Email)

	return req
}

func ValidateRegisterRequest(req RegisterRequest) (RegisterRequest, error) {
	req = NormalizeRegisterRequest(req)
	fields := make(map[string]string)

	if req.Name == "" {
		fields["name"] = "Name is required"
	} else if len(req.Name) < 2 {
		fields["name"] = "Name must be at least 2 characters"
	} else if len(req.Name) > 80 {
		fields["name"] = "Name must be at most 80 characters"
	}

	validateEmail(req.Email, fields)

	if req.Username != nil {
		username := *req.Username
		if len(username) < 3 {
			fields["username"] = "Username must be at least 3 characters"
		} else if len(username) > 30 {
			fields["username"] = "Username must be at most 30 characters"
		} else if !usernamePattern.MatchString(username) {
			fields["username"] = "Username can only use lowercase letters, numbers, and underscore"
		}
	}

	if req.TargetBand == nil {
		fields["targetBand"] = "Target band is required"
	} else if !isAllowedTargetBand(*req.TargetBand) {
		fields["targetBand"] = "Choose a valid IELTS target band"
	}

	validateRegisterPassword(req.Password, fields)

	if len(fields) > 0 {
		return req, ValidationError{Fields: fields}
	}

	return req, nil
}

func ValidateLoginRequest(req LoginRequest) (LoginRequest, error) {
	req = NormalizeLoginRequest(req)
	fields := make(map[string]string)

	validateEmail(req.Email, fields)
	validateLoginPassword(req.Password, fields)

	if len(fields) > 0 {
		return req, ValidationError{Fields: fields}
	}

	return req, nil
}

func validateEmail(email string, fields map[string]string) {
	if email == "" {
		fields["email"] = "Email is required"
		return
	}

	if len(email) > 255 {
		fields["email"] = "Email must be at most 255 characters"
		return
	}

	if !emailPattern.MatchString(email) {
		fields["email"] = "Please enter a valid email address"
	}
}

func validateLoginPassword(password string, fields map[string]string) {
	if password == "" {
		fields["password"] = "Password is required"
		return
	}

	if len(password) < 8 {
		fields["password"] = "Password must be at least 8 characters"
		return
	}

	if len(password) > 72 {
		fields["password"] = "Password must be at most 72 characters"
	}
}

func validateRegisterPassword(password string, fields map[string]string) {
	validateLoginPassword(password, fields)
	if _, exists := fields["password"]; exists {
		return
	}

	if !letterPattern.MatchString(password) {
		fields["password"] = "Password must contain at least one letter"
		return
	}

	if !numberPattern.MatchString(password) {
		fields["password"] = "Password must contain at least one number"
	}
}

func isAllowedTargetBand(value float64) bool {
	switch value {
	case 5.0, 5.5, 6.0, 6.5, 7.0, 7.5, 8.0, 8.5:
		return true
	default:
		return false
	}
}
