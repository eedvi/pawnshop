package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	v := New()
	assert.NotNil(t, v)
	assert.NotNil(t, v.validate)
}

func TestGet_ReturnsSingleton(t *testing.T) {
	instance = nil
	v1 := Get()
	v2 := Get()
	assert.NotNil(t, v1)
	assert.Equal(t, v1, v2)
}

func TestValidate_GlobalFunction(t *testing.T) {
	instance = nil
	type testStruct struct {
		Name string `json:"name" validate:"required"`
	}
	errs := Validate(testStruct{Name: ""})
	assert.NotNil(t, errs)
	assert.Len(t, errs, 1)
	assert.Equal(t, "name", errs[0].Field)
}

func TestValidate_RequiredField(t *testing.T) {
	v := New()

	type testStruct struct {
		Name  string `json:"name" validate:"required"`
		Email string `json:"email" validate:"required,email"`
	}

	errs := v.Validate(testStruct{})
	assert.Len(t, errs, 2)
}

func TestValidate_ValidStruct(t *testing.T) {
	v := New()

	type testStruct struct {
		Name  string `json:"name" validate:"required"`
		Email string `json:"email" validate:"required,email"`
	}

	errs := v.Validate(testStruct{Name: "John", Email: "john@example.com"})
	assert.Nil(t, errs)
}

func TestValidate_EmailField(t *testing.T) {
	v := New()

	type testStruct struct {
		Email string `json:"email" validate:"required,email"`
	}

	errs := v.Validate(testStruct{Email: "not-an-email"})
	assert.Len(t, errs, 1)
	assert.Equal(t, "email", errs[0].Field)
	assert.Equal(t, "Invalid email address", errs[0].Message)
}

func TestValidate_MinLength(t *testing.T) {
	v := New()

	type testStruct struct {
		Name string `json:"name" validate:"required,min=3"`
	}

	errs := v.Validate(testStruct{Name: "ab"})
	assert.Len(t, errs, 1)
	assert.Equal(t, "Must be at least 3 characters", errs[0].Message)
}

func TestValidate_MaxLength(t *testing.T) {
	v := New()

	type testStruct struct {
		Name string `json:"name" validate:"max=5"`
	}

	errs := v.Validate(testStruct{Name: "toolongname"})
	assert.Len(t, errs, 1)
	assert.Equal(t, "Must be at most 5 characters", errs[0].Message)
}

func TestValidate_MinNumeric(t *testing.T) {
	v := New()

	type testStruct struct {
		Age int `json:"age" validate:"min=18"`
	}

	errs := v.Validate(testStruct{Age: 10})
	assert.Len(t, errs, 1)
	assert.Equal(t, "Must be at least 18", errs[0].Message)
}

func TestValidate_Gte(t *testing.T) {
	v := New()

	type testStruct struct {
		Amount float64 `json:"amount" validate:"gte=0"`
	}

	errs := v.Validate(testStruct{Amount: -1})
	assert.Len(t, errs, 1)
	assert.Equal(t, "Must be greater than or equal to 0", errs[0].Message)
}

func TestValidate_Gt(t *testing.T) {
	v := New()

	type testStruct struct {
		Amount float64 `json:"amount" validate:"gt=0"`
	}

	errs := v.Validate(testStruct{Amount: 0})
	assert.Len(t, errs, 1)
	assert.Equal(t, "Must be greater than 0", errs[0].Message)
}

func TestValidate_Lte(t *testing.T) {
	v := New()

	type testStruct struct {
		Rate float64 `json:"rate" validate:"lte=100"`
	}

	errs := v.Validate(testStruct{Rate: 150})
	assert.Len(t, errs, 1)
	assert.Equal(t, "Must be less than or equal to 100", errs[0].Message)
}

func TestValidate_Lt(t *testing.T) {
	v := New()

	type testStruct struct {
		Rate float64 `json:"rate" validate:"lt=100"`
	}

	errs := v.Validate(testStruct{Rate: 100})
	assert.Len(t, errs, 1)
	assert.Equal(t, "Must be less than 100", errs[0].Message)
}

func TestValidate_OneOf(t *testing.T) {
	v := New()

	type testStruct struct {
		Status string `json:"status" validate:"required,oneof=active inactive"`
	}

	errs := v.Validate(testStruct{Status: "unknown"})
	assert.Len(t, errs, 1)
	assert.Equal(t, "Must be one of: active inactive", errs[0].Message)
}

func TestValidate_Len(t *testing.T) {
	v := New()

	type testStruct struct {
		Code string `json:"code" validate:"len=5"`
	}

	errs := v.Validate(testStruct{Code: "abc"})
	assert.Len(t, errs, 1)
	assert.Equal(t, "Must be exactly 5 characters", errs[0].Message)
}

// DPI validation tests

func TestValidate_DPI_Valid(t *testing.T) {
	v := New()

	type testStruct struct {
		DPI string `json:"dpi" validate:"dpi"`
	}

	errs := v.Validate(testStruct{DPI: "1234567890123"})
	assert.Nil(t, errs)
}

func TestValidate_DPI_TooShort(t *testing.T) {
	v := New()

	type testStruct struct {
		DPI string `json:"dpi" validate:"dpi"`
	}

	errs := v.Validate(testStruct{DPI: "123456"})
	assert.Len(t, errs, 1)
	assert.Equal(t, "Invalid DPI format", errs[0].Message)
}

func TestValidate_DPI_NonDigits(t *testing.T) {
	v := New()

	type testStruct struct {
		DPI string `json:"dpi" validate:"dpi"`
	}

	errs := v.Validate(testStruct{DPI: "12345678901ab"})
	assert.Len(t, errs, 1)
	assert.Equal(t, "Invalid DPI format", errs[0].Message)
}

func TestValidate_DPI_TooLong(t *testing.T) {
	v := New()

	type testStruct struct {
		DPI string `json:"dpi" validate:"dpi"`
	}

	errs := v.Validate(testStruct{DPI: "12345678901234"})
	assert.Len(t, errs, 1)
}

// Phone GT validation tests

func TestValidate_PhoneGT_Valid(t *testing.T) {
	v := New()

	type testStruct struct {
		Phone string `json:"phone" validate:"phone_gt"`
	}

	tests := []string{"21234567", "31234567", "41234567", "51234567", "61234567", "71234567"}
	for _, phone := range tests {
		errs := v.Validate(testStruct{Phone: phone})
		assert.Nil(t, errs, "phone %s should be valid", phone)
	}
}

func TestValidate_PhoneGT_WithSeparators(t *testing.T) {
	v := New()

	type testStruct struct {
		Phone string `json:"phone" validate:"phone_gt"`
	}

	errs := v.Validate(testStruct{Phone: "2123-4567"})
	assert.Nil(t, errs)

	errs = v.Validate(testStruct{Phone: "2123 4567"})
	assert.Nil(t, errs)
}

func TestValidate_PhoneGT_InvalidFirstDigit(t *testing.T) {
	v := New()

	type testStruct struct {
		Phone string `json:"phone" validate:"phone_gt"`
	}

	errs := v.Validate(testStruct{Phone: "01234567"})
	assert.Len(t, errs, 1)
	assert.Equal(t, "Invalid phone number format", errs[0].Message)

	errs = v.Validate(testStruct{Phone: "11234567"})
	assert.Len(t, errs, 1)

	errs = v.Validate(testStruct{Phone: "81234567"})
	assert.Len(t, errs, 1)

	errs = v.Validate(testStruct{Phone: "91234567"})
	assert.Len(t, errs, 1)
}

func TestValidate_PhoneGT_WrongLength(t *testing.T) {
	v := New()

	type testStruct struct {
		Phone string `json:"phone" validate:"phone_gt"`
	}

	errs := v.Validate(testStruct{Phone: "2123456"})
	assert.Len(t, errs, 1)

	errs = v.Validate(testStruct{Phone: "212345678"})
	assert.Len(t, errs, 1)
}

func TestValidate_PhoneGT_NonDigits(t *testing.T) {
	v := New()

	type testStruct struct {
		Phone string `json:"phone" validate:"phone_gt"`
	}

	errs := v.Validate(testStruct{Phone: "2123abcd"})
	assert.Len(t, errs, 1)
}

// Password validation tests

func TestValidate_Password_Valid(t *testing.T) {
	v := New()

	type testStruct struct {
		Password string `json:"password" validate:"password"`
	}

	errs := v.Validate(testStruct{Password: "Abcdefg1"})
	assert.Nil(t, errs)
}

func TestValidate_Password_TooShort(t *testing.T) {
	v := New()

	type testStruct struct {
		Password string `json:"password" validate:"password"`
	}

	errs := v.Validate(testStruct{Password: "Ab1"})
	assert.Len(t, errs, 1)
}

func TestValidate_Password_NoUppercase(t *testing.T) {
	v := New()

	type testStruct struct {
		Password string `json:"password" validate:"password"`
	}

	errs := v.Validate(testStruct{Password: "abcdefg1"})
	assert.Len(t, errs, 1)
}

func TestValidate_Password_NoLowercase(t *testing.T) {
	v := New()

	type testStruct struct {
		Password string `json:"password" validate:"password"`
	}

	errs := v.Validate(testStruct{Password: "ABCDEFG1"})
	assert.Len(t, errs, 1)
}

func TestValidate_Password_NoDigit(t *testing.T) {
	v := New()

	type testStruct struct {
		Password string `json:"password" validate:"password"`
	}

	errs := v.Validate(testStruct{Password: "Abcdefgh"})
	assert.Len(t, errs, 1)
}

// JSON tag field name test

func TestValidate_UsesJSONTagNames(t *testing.T) {
	v := New()

	type testStruct struct {
		FirstName string `json:"first_name" validate:"required"`
	}

	errs := v.Validate(testStruct{})
	assert.Len(t, errs, 1)
	assert.Equal(t, "first_name", errs[0].Field)
}

func TestValidate_IgnoredJSONField(t *testing.T) {
	v := New()

	type testStruct struct {
		Internal string `json:"-" validate:"required"`
		Name     string `json:"name" validate:"required"`
	}

	errs := v.Validate(testStruct{})
	// The "-" field may still require validation but the field name should be empty or skipped
	assert.NotNil(t, errs)
}

func TestValidate_URL(t *testing.T) {
	v := New()

	type testStruct struct {
		Website string `json:"website" validate:"url"`
	}

	errs := v.Validate(testStruct{Website: "not-a-url"})
	assert.Len(t, errs, 1)
	assert.Equal(t, "Invalid URL", errs[0].Message)

	errs = v.Validate(testStruct{Website: "https://example.com"})
	assert.Nil(t, errs)
}

func TestValidate_UUID(t *testing.T) {
	v := New()

	type testStruct struct {
		ID string `json:"id" validate:"uuid"`
	}

	errs := v.Validate(testStruct{ID: "not-a-uuid"})
	assert.Len(t, errs, 1)
	assert.Equal(t, "Invalid UUID", errs[0].Message)

	errs = v.Validate(testStruct{ID: "550e8400-e29b-41d4-a716-446655440000"})
	assert.Nil(t, errs)
}

func TestValidate_Numeric(t *testing.T) {
	v := New()

	type testStruct struct {
		Code string `json:"code" validate:"numeric"`
	}

	errs := v.Validate(testStruct{Code: "abc"})
	assert.Len(t, errs, 1)
	assert.Equal(t, "Must be numeric", errs[0].Message)

	errs = v.Validate(testStruct{Code: "12345"})
	assert.Nil(t, errs)
}

func TestValidate_Alpha(t *testing.T) {
	v := New()

	type testStruct struct {
		Name string `json:"name" validate:"alpha"`
	}

	errs := v.Validate(testStruct{Name: "abc123"})
	assert.Len(t, errs, 1)
	assert.Equal(t, "Must contain only letters", errs[0].Message)

	errs = v.Validate(testStruct{Name: "abc"})
	assert.Nil(t, errs)
}

func TestValidate_Alphanum(t *testing.T) {
	v := New()

	type testStruct struct {
		Code string `json:"code" validate:"alphanum"`
	}

	errs := v.Validate(testStruct{Code: "abc-123"})
	assert.Len(t, errs, 1)
	assert.Equal(t, "Must contain only letters and numbers", errs[0].Message)

	errs = v.Validate(testStruct{Code: "abc123"})
	assert.Nil(t, errs)
}

func TestValidate_EqField(t *testing.T) {
	v := New()

	type testStruct struct {
		Password        string `json:"password" validate:"required"`
		ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
	}

	errs := v.Validate(testStruct{Password: "abc123", ConfirmPassword: "different"})
	assert.Len(t, errs, 1)
	assert.Equal(t, "Must match Password", errs[0].Message)

	errs = v.Validate(testStruct{Password: "abc123", ConfirmPassword: "abc123"})
	assert.Nil(t, errs)
}

func TestValidate_NeField(t *testing.T) {
	v := New()

	type testStruct struct {
		OldPassword string `json:"old_password" validate:"required"`
		NewPassword string `json:"new_password" validate:"required,nefield=OldPassword"`
	}

	errs := v.Validate(testStruct{OldPassword: "same", NewPassword: "same"})
	assert.Len(t, errs, 1)
	assert.Equal(t, "Must not match OldPassword", errs[0].Message)

	errs = v.Validate(testStruct{OldPassword: "old", NewPassword: "new"})
	assert.Nil(t, errs)
}

func TestValidate_MaxNumeric(t *testing.T) {
	v := New()

	type testStruct struct {
		Count int `json:"count" validate:"max=10"`
	}

	errs := v.Validate(testStruct{Count: 20})
	assert.Len(t, errs, 1)
	assert.Equal(t, "Must be at most 10", errs[0].Message)
}

func TestValidate_DefaultErrorMessage(t *testing.T) {
	v := New()

	type testStruct struct {
		IP string `json:"ip" validate:"ip"`
	}

	errs := v.Validate(testStruct{IP: "not-an-ip"})
	assert.Len(t, errs, 1)
	assert.Equal(t, "Invalid value", errs[0].Message)
}
