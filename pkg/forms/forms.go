package forms

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"
)

// EmailRX is used to make sure that the email is correct and pretty
var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Form is holding the submitted form data and form errors
type Form struct {
	url.Values
	Errors errors
}

// New creates a new Form Struct that will hold
// submitted form data and form errors
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Requiered checks to make sure a requiered field isn't empty
func (f *Form) Requiered(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, fmt.Sprintf("%v field cannot be empty", field))
		}
	}
}

// MaxLength checks if user passes to many characters
func (f *Form) MaxLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) > d {
		f.Errors.Add(field, fmt.Sprintf("Only %d characters allowed", d))
	}
}

// PermittedValues checks that the user provied a valid option
func (f *Form) PermittedValues(field string, opts ...string) {
	value := f.Get(field)
	if value == "" {
		return
	}

	for _, opt := range opts {
		if value == opt {
			return
		}
	}
	f.Errors.Add(field, "Invalid option")
}

// Valid checks of the form had any errors or not
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// MinLength validates that the string has the correct minimum length
func (f *Form) MinLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) < d {
		f.Errors.Add(field, fmt.Sprintf("You need to enter at least %d characters", d))
	}
}

// MatchesPattern makes sure that specific string matches the pattern
func (f *Form) MatchesPattern(field string, pattern *regexp.Regexp) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if !pattern.MatchString(value) {
		f.Errors.Add(field, fmt.Sprintf("Please enter an correct %s", field))
	}
}
