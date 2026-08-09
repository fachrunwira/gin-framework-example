package appvalidation

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ErrorMessages map[string]map[string]string

func CustomFirstError(err error, message ErrorMessages) string {
	var errors string

	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return errors
	}

	e := errs[0]

	errors = formatError(e, message)

	return errors
}

func formatError(fl validator.FieldError, messages ErrorMessages) string {
	field := convertToSnake(fl.Field())
	namespace := fl.Namespace()
	tag := fl.Tag()
	var errMsg string

	regex := regexp.MustCompile(`\[\d+\]`)
	if regex.MatchString(field) {
		newField := regex.ReplaceAllString(field, "[n]")
		fieldParam := regexp.MustCompile(`\D+`).ReplaceAllString(field, "")
		fieldIndex := regex.ReplaceAllString(field, "")

		if fieldMsg, ok := messages[newField]; ok {
			if msg, ok := fieldMsg[tag]; ok {
				errMsg = formatArrayMessage(msg, fieldParam, fl.Param(), fl.Value(), fieldIndex)
			}
		}

		return errMsg
	}

	if regex.MatchString(namespace) {
		fieldParam := regexp.MustCompile(`\D+`).ReplaceAllString(namespace, "")

		if fieldMsg, ok := messages[field]; ok {
			if msg, ok := fieldMsg[tag]; ok {
				errMsg = formatArrayMessage(msg, field, fl.Param(), fl.Value(), fieldParam)
			}
		}

		return errMsg
	}

	if fieldMsg, ok := messages[field]; ok {
		if msg, ok := fieldMsg[tag]; ok {
			errMsg = formatMessage(msg, field, fl.Param(), fl.Value())
		}
	}

	return errMsg
}

func formatMessage(msg, field, param string, value any) string {
	args := make([]any, 0, 2)

	if strings.Contains(msg, ":field") {
		msg = strings.ReplaceAll(msg, ":field", "%s")
		args = append(args, field)
	}

	if strings.Contains(msg, ":param") {
		msg = strings.ReplaceAll(msg, ":param", "%v")
		args = append(args, param)
	}

	if strings.Contains(msg, ":value") {
		msg = strings.ReplaceAll(msg, ":value", "%v")
		args = append(args, value)
	}

	if len(args) == 0 {
		return msg
	}

	return fmt.Sprintf(msg, args...)
}

func formatArrayMessage(msg, field, param string, value, index any) string {
	args := make([]any, 0, 2)

	if strings.Contains(msg, ":field") {
		msg = strings.ReplaceAll(msg, ":field", "%s")
		args = append(args, field)
	}

	if strings.Contains(msg, ":param") {
		msg = strings.ReplaceAll(msg, ":param", "%v")
		args = append(args, param)
	}

	if strings.Contains(msg, ":value") {
		msg = strings.ReplaceAll(msg, ":value", "%v")
		args = append(args, value)
	}

	if strings.Contains(msg, ":row") {
		msg = strings.ReplaceAll(msg, ":row", "%v")
		args = append(args, index)
	}

	if len(args) == 0 {
		return msg
	}

	return fmt.Sprintf(msg, args...)
}

func convertToSnake(txt string) string {
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

	txt = matchFirstCap.ReplaceAllString(txt, "${1}_${2}")
	txt = matchAllCap.ReplaceAllString(txt, "${1}_${2}")
	return strings.ToLower(txt)
}
