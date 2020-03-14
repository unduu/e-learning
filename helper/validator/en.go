package validator

import (
	"fmt"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
)

// Add your own translation
func setTranslations(v *validator.Validate, trans ut.Translator) {
	_ = v.RegisterTranslation("isphonenumber", trans, func(ut ut.Translator) error {
		return ut.Add("isphonenumber", "{0} is not a valid number", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		fieldValueStr := fmt.Sprintf("%v", fe.Value())
		t, _ := ut.T("isphonenumber", fieldValueStr)
		return t
	})
}
