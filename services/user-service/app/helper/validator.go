package helper

import (
	"user-service/app/model"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func Validate(v *validator.Validate, i interface{}) []model.ValidationError {
	validationErrors := []model.ValidationError{}

	errs := v.Struct(i)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			param := err.Param()
			value := err.Value()

			validationErrors = append(validationErrors, model.ValidationError{
				FailedField: ConvertToSpaced(err.Field()),
				Tag:         err.Tag(),
				Value:      &value,
				Param: 	 	&param,
			})
		}
	}

	return validationErrors
}

func GetErrorMessages(errors []model.ValidationError) []string {
	messages := []string{}

	for _, err := range errors {
		messages = append(messages, model.TranslateTag(err))
	}

	return messages
}

func ParseAndValidate(ctx *fiber.Ctx, validator *validator.Validate, request interface{}) error {
    if err := ctx.BodyParser(request); err != nil {
        if errors := Validate(validator, request); len(errors) > 0 {
			return model.NewError(model.StatusBadRequest, "The given data was invalid", errors)
		}
    }

	if errors := Validate(validator, request); len(errors) > 0 {
		return model.NewError(model.StatusBadRequest, "The given data was invalid", errors)
	}
    
    return nil
}

func ValidateFileExtension(extensions []string, extension string) bool {
	for _, ext := range extensions {
		if ext == extension {
			return true
		}
	}

	return false
}