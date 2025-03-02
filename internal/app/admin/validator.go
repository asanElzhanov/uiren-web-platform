package admin

import (
	"regexp"
	"strings"
)

var (
	codeRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
)

func validateObjectKeys(data map[string]interface{}, key string, requiredFields map[string]string, errReturn error) error {
	value, exists := data[key]
	if !exists {
		return nil
	}

	obj, ok := value.(map[string]interface{})
	if !ok {
		return errReturn
	}

	for field, expectedType := range requiredFields {
		val, exists := obj[field]
		if !exists {
			return errReturn
		}

		if !isValidType(val, expectedType) {
			return errReturn
		}
	}

	return nil
}

func isValidType(value interface{}, expectedType string) bool {
	switch expectedType {
	case "string":
		_, ok := value.(string)
		return ok
	case "float64":
		_, ok := value.(float64)
		return ok
	default:
		return false
	}
}

// modules
func validateRewardsAndRequirements(rawData map[string]interface{}) error {
	if err := validateObjectKeys(rawData, "reward", map[string]string{
		"xp":    "float64",
		"badge": "string",
	}, ErrInvalidReward); err != nil {
		return err
	}

	if err := validateObjectKeys(rawData, "unlock_requirements", map[string]string{
		"previous_module": "string",
		"min_xp":          "float64",
	}, ErrInvalidRequirements); err != nil {
		return err
	}

	return nil
}

// modules, lessons, exercises
func validateCode(code string) error {
	code = strings.TrimSpace(code)
	if code == "" || !codeRegex.MatchString(code) || strings.Contains(code, " ") {
		return ErrInvalidCode
	}

	return nil
}
