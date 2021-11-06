package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/leksss/hw-test/hw09_struct_validator/validators"
	"github.com/stretchr/testify/require"
)

type UserRole string

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
		Version string `validate:"len:5"`
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

	Custom struct {
		Num    int    `validate:"min:0|max:10"`
		StrNum string `validate:"regexp:\\d+|len:20"`
	}

	IntSlice struct {
		MinMaxNum []int `validate:"min:0|max:10"`
		InNum     []int `validate:"in:333,444,555"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:     "123456789012345678901234567890123456", // string `validate:"len:36"`
				Age:    42,                                     // int `validate:"min:18|max:50"`
				Email:  "match@email.re",                       // string `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
				Role:   "admin",                                // UserRole `validate:"in:admin,stuff"`
				Phones: []string{"12345678901", "12345678901"}, // []string `validate:"len:11"`
				meta:   nil,
			},
			expectedErr: validators.ValidationErrors{},
		},
		{
			in: User{
				ID:    "1234567890123456789012345678901234567",
				Age:   17,
				Email: "name@host",
				Role:  "user",
				Phones: []string{
					"123456789012",
					"1234567890",
					"",
				},
			},
			expectedErr: validators.ValidationErrors{
				{Field: "ID", Err: validators.ErrLenValidatorMustBeExact},
				{Field: "Age", Err: validators.ErrMinValidatorShouldMore},
				{Field: "Email", Err: validators.ErrRegexpValidatorNotMatch},
				{Field: "Role", Err: validators.ErrInValidatorShouldBeInList},
				{Field: "Phones", Err: validators.ErrLenValidatorMustBeExact},
			},
		},
		{
			in: App{"123456"},
			expectedErr: validators.ValidationErrors{
				{Field: "Version", Err: validators.ErrLenValidatorMustBeExact},
			},
		},
		{
			in:          App{"12345"},
			expectedErr: validators.ValidationErrors{},
		},
		{
			in: Token{
				Header:    nil,
				Payload:   nil,
				Signature: nil,
			},
			expectedErr: validators.ValidationErrors{},
		},
		{
			in: Response{
				Code: 404,
				Body: "",
			},
			expectedErr: validators.ValidationErrors{},
		},
		{
			in: Response{
				Code: 777,
				Body: "",
			},
			expectedErr: validators.ValidationErrors{
				{Field: "Code", Err: validators.ErrInValidatorShouldBeInList},
			},
		},
		{
			in: Custom{
				Num:    7,
				StrNum: "12345678901234567890",
			},
			expectedErr: validators.ValidationErrors{},
		},
		{
			in: Custom{
				Num:    22,                    // int `validate:"min:0|max:10"`
				StrNum: "1234567890123456789", // string `validate:"regexp:\\d+|len:20"`
			},
			expectedErr: validators.ValidationErrors{
				{Field: "Num", Err: validators.ErrMaxValidatorShouldBeLess},
				{Field: "StrNum", Err: validators.ErrLenValidatorMustBeExact},
			},
		},
		{
			in: IntSlice{
				MinMaxNum: []int{4, 5, 7, 0, 10}, // []int `validate:"min:0|max:10"`
				InNum:     []int{555, 333},       // []int `validate:"in:333,444,555"`
			},
			expectedErr: validators.ValidationErrors{},
		},
		{
			in: IntSlice{
				MinMaxNum: []int{-1, 8},    // []int `validate:"min:0|max:10"`
				InNum:     []int{333, 111}, // []int `validate:"in:333,444,555"`
			},
			expectedErr: validators.ValidationErrors{
				{Field: "MinMaxNum", Err: validators.ErrMinValidatorShouldMore},
				{Field: "InNum", Err: validators.ErrInValidatorShouldBeInList},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			validationErrors, err := Validate(tt.in)
			require.NoError(t, err)
			require.Equal(t, validationErrors, tt.expectedErr)
		})
	}
}
