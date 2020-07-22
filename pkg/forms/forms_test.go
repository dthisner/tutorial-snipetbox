package forms

import (
	"net/url"
	"reflect"
	"testing"
)

func TestRequiered(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		formData map[string]string
		want     errors
	}{
		{name: "Passing all correct data",
			formData: map[string]string{
				"email":    "dennis@bob.com",
				"name":     "bob",
				"password": "Minsalfijl",
			},
			want: errors(map[string][]string{}),
		},
		{name: "Missing Email",
			formData: map[string]string{
				"email":    "",
				"name":     "bob",
				"password": "Minsalfijl",
			},
			want: errors(map[string][]string{
				"email": {"email field cannot be empty"},
			}),
		},
		{name: "Missing Name and Password",
			formData: map[string]string{
				"name":     "",
				"password": "",
			},
			want: errors(map[string][]string{
				"name":     {"name field cannot be empty"},
				"password": {"password field cannot be empty"},
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := url.Values{}
			for k, v := range tt.formData {
				params.Add(k, v)
			}

			form := New(params)
			form.Requiered("name", "email", "password")

			if !reflect.DeepEqual(tt.want, form.Errors) {
				t.Errorf("want %q; got %q", tt.want, form.Errors)
			}
		})
	}
}

func TestMaxLength(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		maxLength int
		formData  map[string]string
		want      errors
	}{
		{name: "Passing correct length",
			maxLength: 30,
			formData: map[string]string{
				"email": "test@bob.com",
			},
			want: errors(map[string][]string{}),
		},
		{name: "Passing exact length",
			maxLength: 12,
			formData: map[string]string{
				"email": "test@bob.com",
			},
			want: errors(map[string][]string{}),
		},
		{name: "Passing to long",
			maxLength: 12,
			formData: map[string]string{
				"email": "ReallyLongStringOfTextLikeOMG-DidYouSeeThatAmazing!",
			},
			want: errors(map[string][]string{
				"email": {"Only 12 characters allowed"},
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := url.Values{}
			for k, v := range tt.formData {
				params.Add(k, v)
			}

			form := New(params)
			form.MaxLength("email", tt.maxLength)

			if !reflect.DeepEqual(tt.want, form.Errors) {
				t.Errorf("want %q; got %q", tt.want, form.Errors)
			}
		})
	}
}

func TestMinLength(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		minLength int
		formData  map[string]string
		want      errors
	}{
		{name: "Passing correct length",
			minLength: 2,
			formData: map[string]string{
				"email": "test@bob.com",
			},
			want: errors(map[string][]string{}),
		},
		{name: "Passing exact length",
			minLength: 12,
			formData: map[string]string{
				"email": "test@bob.com",
			},
			want: errors(map[string][]string{}),
		},
		{name: "Passing to short",
			minLength: 12,
			formData: map[string]string{
				"email": "Rea",
			},
			want: errors(map[string][]string{
				"email": {"You need to enter at least 12 characters"},
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := url.Values{}
			for k, v := range tt.formData {
				params.Add(k, v)
			}

			form := New(params)
			form.MinLength("email", tt.minLength)

			if !reflect.DeepEqual(tt.want, form.Errors) {
				t.Errorf("want %q; got %q", tt.want, form.Errors)
			}
		})
	}
}

func TestPermittedValues(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		permittedValues string
		formData        map[string]string
		want            errors
	}{
		{name: "Passing valid string",
			permittedValues: "bob",
			formData: map[string]string{
				"email": "bob",
			},
			want: errors(map[string][]string{}),
		},
		{name: "Passing lowercase into only permitted upercase",
			permittedValues: "BOB",
			formData: map[string]string{
				"email": "bob",
			},
			want: errors(map[string][]string{
				"email": {"Invalid option"},
			})},
		{name: "Passing empty string into as permitted value",
			permittedValues: "",
			formData: map[string]string{
				"email": "bob",
			},
			want: errors(map[string][]string{
				"email": {"Invalid option"},
			})},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := url.Values{}
			for k, v := range tt.formData {
				params.Add(k, v)
			}

			form := New(params)
			form.PermittedValues("email", tt.permittedValues)

			if !reflect.DeepEqual(tt.want, form.Errors) {
				t.Errorf("want %q; got %q", tt.want, form.Errors)
			}
		})
	}
}

func TestMatchPattern(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		formData map[string]string
		want     errors
	}{
		{name: "Passing correct match",
			formData: map[string]string{
				"email": "test@bob.com",
			},
			want: errors(map[string][]string{}),
		},
		{name: "Missing '@'",
			formData: map[string]string{
				"email": "testbob.com",
			},
			want: errors(map[string][]string{
				"email": {"Please enter an correct email"},
			}),
		},
		{name: "Using space in the email",
			formData: map[string]string{
				"email": "dennis @bob.com",
			},
			want: errors(map[string][]string{
				"email": {"Please enter an correct email"},
			}),
		},
		{name: "Missing '.com'",
			formData: map[string]string{
				"email": "Rea",
			},
			want: errors(map[string][]string{
				"email": {"Please enter an correct email"},
			}),
		},
		{name: "Using symbols",
			formData: map[string]string{
				"email": "bob#!@.nils.com",
			},
			want: errors(map[string][]string{
				"email": {"Please enter an correct email"},
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := url.Values{}
			for k, v := range tt.formData {
				params.Add(k, v)
			}

			form := New(params)
			form.MatchesPattern("email", EmailRX)

			if !reflect.DeepEqual(tt.want, form.Errors) {
				t.Errorf("want %q; got %q", tt.want, form.Errors)
			}
		})
	}
}
