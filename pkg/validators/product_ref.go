package validators

import (
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func ProductReferenceValidator(fl validator.FieldLevel) bool {
	re := regexp.MustCompile(`^[A-Z0-9-]+$`) // Only uppercase letters, numbers, and hyphens
	return re.MatchString(fl.Field().String())
}

func RegisterCustomValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// Register all custom validators in here //
		_ = v.RegisterValidation("product_ref", ProductReferenceValidator)
	}
}
