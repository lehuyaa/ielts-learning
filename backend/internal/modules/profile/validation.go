package profile

import (
	"regexp"
	"strings"
	"time"
)

var (
	profileUsernamePattern = regexp.MustCompile(`^[a-z0-9_]+$`)
	localePattern          = regexp.MustCompile(`^[a-z]{2}(-[a-z]{2})?$`)
)

type ValidationError struct {
	Fields map[string]string
}

func (e ValidationError) Error() string {
	return "validation failed"
}

func NormalizeUpdateProfileRequest(req UpdateProfileRequest) UpdateProfileRequest {
	if req.Name != nil {
		value := strings.TrimSpace(*req.Name)
		req.Name = &value
	}
	if req.Username != nil {
		value := strings.ToLower(strings.TrimSpace(*req.Username))
		req.Username = &value
	}
	if req.Timezone != nil {
		value := strings.TrimSpace(*req.Timezone)
		req.Timezone = &value
	}
	if req.Locale != nil {
		value := strings.ToLower(strings.TrimSpace(*req.Locale))
		req.Locale = &value
	}

	return req
}

func ValidateUpdateProfileRequest(req UpdateProfileRequest) (UpdateProfileRequest, error) {
	req = NormalizeUpdateProfileRequest(req)
	fields := make(map[string]string)

	if req.Name == nil && req.Username == nil && req.TargetBand == nil && req.Timezone == nil && req.Locale == nil {
		fields["request"] = "At least one profile field must be provided"
	}

	if req.Name != nil {
		switch {
		case *req.Name == "":
			fields["name"] = "Name is required"
		case len(*req.Name) < 2:
			fields["name"] = "Name must be at least 2 characters"
		case len(*req.Name) > 80:
			fields["name"] = "Name must be at most 80 characters"
		}
	}

	if req.Username != nil {
		switch {
		case *req.Username == "":
			fields["username"] = "Username is required"
		case len(*req.Username) < 3:
			fields["username"] = "Username must be at least 3 characters"
		case len(*req.Username) > 30:
			fields["username"] = "Username must be at most 30 characters"
		case !profileUsernamePattern.MatchString(*req.Username):
			fields["username"] = "Username can only use lowercase letters, numbers, and underscore"
		}
	}

	if req.TargetBand != nil && !isAllowedTargetBand(*req.TargetBand) {
		fields["targetBand"] = "Choose a valid IELTS target band"
	}

	if req.Timezone != nil {
		switch {
		case *req.Timezone == "":
			fields["timezone"] = "Timezone is required"
		case len(*req.Timezone) > 80:
			fields["timezone"] = "Timezone must be at most 80 characters"
		default:
			if _, err := time.LoadLocation(*req.Timezone); err != nil {
				fields["timezone"] = "Timezone must be a valid IANA timezone"
			}
		}
	}

	if req.Locale != nil {
		switch {
		case *req.Locale == "":
			fields["locale"] = "Locale is required"
		case len(*req.Locale) < 2:
			fields["locale"] = "Locale must be at least 2 characters"
		case len(*req.Locale) > 20:
			fields["locale"] = "Locale must be at most 20 characters"
		case !localePattern.MatchString(*req.Locale):
			fields["locale"] = "Locale must use language or language-region format"
		}
	}

	if len(fields) > 0 {
		return req, ValidationError{Fields: fields}
	}

	return req, nil
}

func isAllowedTargetBand(value float64) bool {
	switch value {
	case 5.0, 5.5, 6.0, 6.5, 7.0, 7.5, 8.0, 8.5:
		return true
	default:
		return false
	}
}
