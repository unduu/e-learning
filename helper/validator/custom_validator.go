package validator

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"reflect"
	"sync"

	"github.com/dongri/phonenumber"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
)

type CustomValidator struct {
	Validator   *validator.Validate
	Translation ut.Translator
}

func NewCustomValidator(db *sqlx.DB) *CustomValidator {
	translator := en.New()
	uni := ut.New(translator, translator)
	trans, found := uni.GetTranslator("en")
	if !found {
		fmt.Println("translator not found")
	}

	binding.Validator = new(defaultValidator)
	v, ok := binding.Validator.Engine().(*validator.Validate)

	if ok {
		v.RegisterValidation("isphonenumber", func(fl validator.FieldLevel) bool {
			fieldValue := fl.Field().String()
			normalize := phonenumber.Parse(fieldValue, "ID")
			if normalize == "" {
				return false
			}
			return true
		})

		v.RegisterValidation("emailExists", func(fl validator.FieldLevel) bool {
			row := db.QueryRow("SELECT count(*) as total FROM users WHERE email=?", fl.Field().String())
			var exists bool
			err := row.Scan(&exists)
			if err != nil {
				fmt.Println("ERROR emailExists ", err)
			}
			return !exists
		})

		v.RegisterValidation("phoneExists", func(fl validator.FieldLevel) bool {
			phone := phonenumber.Parse(fl.Field().String(), "ID")
			row := db.QueryRow("SELECT count(*) as total FROM users WHERE phone=?", phone)
			var exists bool
			err := row.Scan(&exists)
			if err != nil {
				fmt.Println("ERROR phoneExists ", err)
			}
			return !exists
		})

		v.RegisterValidation("usernameExists", func(fl validator.FieldLevel) bool {
			row := db.QueryRow("SELECT count(*) as total FROM users WHERE username=?", fl.Field().String())
			var exists bool
			err := row.Scan(&exists)
			if err != nil {
				fmt.Println("ERROR usernameExists ", err)
			}
			return !exists
		})
	}

	if err := en_translations.RegisterDefaultTranslations(v, trans); err != nil {
		fmt.Println(err)
	}

	setTranslations(v, trans)

	return &CustomValidator{Validator: v, Translation: trans}
}

// V8 to V9
type defaultValidator struct {
	once     sync.Once
	validate *validator.Validate
}

var _ binding.StructValidator = &defaultValidator{}

func (v *defaultValidator) ValidateStruct(obj interface{}) error {

	if kindOfData(obj) == reflect.Struct {

		v.lazyinit()

		if err := v.validate.Struct(obj); err != nil {
			return error(err)
		}
	}

	return nil
}

func (v *defaultValidator) Engine() interface{} {
	v.lazyinit()
	return v.validate
}

func (v *defaultValidator) lazyinit() {
	v.once.Do(func() {
		v.validate = validator.New()
		v.validate.SetTagName("binding")

		// add any custom validations etc. here
	})
}

func kindOfData(data interface{}) reflect.Kind {

	value := reflect.ValueOf(data)
	valueType := value.Kind()

	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}
