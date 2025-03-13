package utils

import (
	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func ValidationError(err error) []ErrorResponse {
	var errors []ErrorResponse
	validatorErrs := err.(validator.ValidationErrors)

	for _, e := range validatorErrs {
		var element ErrorResponse
		element.Field = e.Field()
		element.Message = getErrorMsg(e)
		errors = append(errors, element)
	}
	return errors
}

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "Ce champ est requis"
	case "email":
		return "Format d'email invalide"
	case "min":
		return "La longueur minimale n'est pas respectée"
	case "max":
		return "La longueur maximale est dépassée"
	}
	return "Erreur de validation inconnue"
}
