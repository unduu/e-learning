package validator

import (
	"fmt"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
)

// Add your own translation
func setTranslations(v *validator.Validate, trans ut.Translator) {

	_ = v.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is a required field", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	_ = v.RegisterTranslation("isphonenumber", trans, func(ut ut.Translator) error {
		return ut.Add("isphonenumber", "{0} is not a valid number", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		fieldValueStr := fmt.Sprintf("%v", fe.Value())
		t, _ := ut.T("isphonenumber", fieldValueStr)
		return t
	})

	_ = v.RegisterTranslation("emailExists", trans, func(ut ut.Translator) error {
		return ut.Add("emailExists", "{0} already exists", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("emailExists", fe.Field())
		return t
	})

	_ = v.RegisterTranslation("phoneExists", trans, func(ut ut.Translator) error {
		return ut.Add("phoneExists", "{0} already exists", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("phoneExists", fe.Field())
		return t
	})

	_ = v.RegisterTranslation("isValidPhone", trans, func(ut ut.Translator) error {
		return ut.Add("isValidPhone", "{0} not registered", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("isValidPhone", fe.Field())
		return t
	})

	_ = v.RegisterTranslation("usernameExists", trans, func(ut ut.Translator) error {
		return ut.Add("usernameExists", "{0} already exists", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("usernameExists", fe.Field())
		return t
	})

	_ = v.RegisterTranslation("isValidPaswordKey", trans, func(ut ut.Translator) error {
		return ut.Add("isValidPaswordKey", "{0} not found", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("isValidPaswordKey", fe.Field())
		return t
	})
}
