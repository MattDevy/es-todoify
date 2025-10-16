package todo

import (
	"fmt"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

// validate is a global validator instance used throughout the package.
// It is initialized with required struct enabled to enforce non-empty fields.
// This is a best practice to ensure that all fields are validated.
// We want to reuse the validator to allow caching of struct tags and avoid reflection whenever possible.
var validate = validator.New(validator.WithRequiredStructEnabled())

// trans is the universal translator for English locale.
// It provides human-readable error messages for validation failures.
var trans ut.Translator

func init() {
	// Initialize the English locale
	english := en.New()
	uni := ut.New(english, english)

	// Get the translator for English
	var found bool
	trans, found = uni.GetTranslator("en")
	if !found {
		panic("translator not found")
	}

	// Register default English translations
	if err := en_translations.RegisterDefaultTranslations(validate, trans); err != nil {
		panic(fmt.Sprintf("failed to register translations: %v", err))
	}
}

// TranslateError converts validator errors into human-readable messages.
// It returns a map of field names to their translated error messages.
func TranslateError(err error) map[string]string {
	if err == nil {
		return nil
	}

	validatorErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return map[string]string{"error": err.Error()}
	}

	errs := make(map[string]string)
	for _, e := range validatorErrs {
		errs[e.Field()] = e.Translate(trans)
	}

	return errs
}
