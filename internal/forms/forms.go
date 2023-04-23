package forms

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

// Form creates a custom form struct, and it embeds a url.Value object
type Form struct {
	url.Values
	Errors errors
}

// New initializes a form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		// empty errors
		errors(map[string][]string{}),
	}
}

// Defined certain check functions to validate the received form data

// 1. Valida returns true if there is no errors, otherwise false
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// 2. Required checks for required fields
// Required allows multiple, unfixed number of parameters with type string
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		// if value does not contains characters
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field can not be blank")
		}
	}
}

// 3. Has checks if certain fields are included in the form
func (f *Form) Has(field string, r *http.Request) bool {
	x := f.Values.Get(field)
	if x == "" {
		return false
	}
	return true
}

// 4. MinLength checks for string minimum length
func (f *Form) MinLength(field string, length int) bool {
	x := f.Values.Get(field)
	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", length))
		return false
	}
	return true
}

// IsEmail checks for a valid email address
func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid email address")
	}
}
