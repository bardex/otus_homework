package hw09structvalidator

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5|regexp:^\\d\\.\\d\\.\\d$"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	validID := strings.Repeat("1", 36)
	validAge := 18
	validRole := UserRole("admin")
	validEmail := "test@example.com"
	validPhones := []string{"01234567890", "00987654321"}

	tests := []struct {
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			name: "valid_user",
			in: User{
				ID:     validID,
				Name:   "Alex",
				Age:    validAge,
				Email:  validEmail,
				Role:   validRole,
				Phones: validPhones,
				meta:   nil,
			},
			expectedErr: nil,
		},
		{
			name:        "valid_app",
			in:          App{Version: "3.2.1"},
			expectedErr: nil,
		},
		{
			name: "valid_token",
			in: Token{
				Header:    bytes.Repeat([]byte("0"), 10),
				Payload:   []byte{},
				Signature: bytes.Repeat([]byte("1"), 20),
			},
			expectedErr: nil,
		},
		{
			name:        "valid_response",
			in:          Response{Code: 200},
			expectedErr: nil,
		},
		{
			name: "invalid_user_1",
			in: User{
				ID:     "10",
				Name:   "",
				Age:    17,
				Email:  "example.com",
				Role:   "user",
				Phones: []string{strings.Repeat("1", 11), strings.Repeat("2", 3), strings.Repeat("2", 3)},
				meta:   nil,
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: "ID", Err: LenValidationError{Expected: 36, Actual: 2}},
				ValidationError{Field: "Age", Err: MinValidationError{Min: "18", Actual: "17"}},
				ValidationError{Field: "Email", Err: RegexpValidationError{Regexp: "^\\w+@\\w+\\.\\w+$", Actual: "example.com"}},
				ValidationError{Field: "Role", Err: InValidationError{In: []string{"admin", "stuff"}, Actual: "user"}},
				ValidationError{Field: "Phones", Err: LenValidationError{Expected: 11, Actual: 3}},
			},
		},
		{
			name: "validator_error_type",
			in: struct {
				ID int `validate:"len:10"`
			}{ID: 10},
			expectedErr: ValidatorError{},
		},
		{
			name: "validator_error_syntax",
			in: struct {
				ID int `validate:"min:xyz"`
			}{ID: 10},
			expectedErr: ValidatorError{},
		},
		{
			name: "validator_error_required_param",
			in: struct {
				ID int `validate:"min:"`
			}{ID: 10},
			expectedErr: ValidatorError{},
		},
		{
			name:        "validator_error_no_struct",
			in:          10,
			expectedErr: ValidatorError{},
		},
		{
			name: "validator_error_unknown_validator",
			in: struct {
				Email string `validate:"email"`
			}{Email: "example.com"},
			expectedErr: ValidatorError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)

			require.IsType(t, tt.expectedErr, err)
			if errors.As(tt.expectedErr, &ValidationErrors{}) {
				require.Equal(t, tt.expectedErr.Error(), err.Error())
			}
		})
	}
}
