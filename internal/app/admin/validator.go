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

// progress

func validateAchievementProgress(rawData map[string]interface{}) error {
	value, exists := rawData["achievements_progress"]
	if !exists {
		return nil
	}

	list, ok := value.([]interface{})
	if !ok {
		return ErrInvalidAchievementProgress
	}

	for _, item := range list {
		progress, ok := item.(map[string]interface{})
		if !ok {
			return ErrInvalidAchievementProgress
		}

		requiredFields := map[string]string{
			"achievement_id":  "float64",
			"earned_progress": "float64",
		}

		for field, expectedType := range requiredFields {
			val, exists := progress[field]
			if !exists || !isValidType(val, expectedType) {
				return ErrInvalidAchievementProgress
			}
		}
	}

	return nil
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

// exercises
func validatePairs(rawData map[string]interface{}) error {
	v, exists := rawData["type"]
	if !exists {
		return ErrInvalidType
	}
	exerciseType, ok := v.(string)
	if !ok {
		return ErrInvalidType
	}
	if exerciseType != "match_pairs" {
		return nil
	}

	v, exists = rawData["pairs"]
	if !exists {
		return ErrInvalidPairs
	}

	list, ok := v.([]interface{})
	if !ok {
		return ErrInvalidPairs
	}

	for _, item := range list {
		pair, ok := item.(map[string]interface{})
		if !ok {
			return ErrInvalidPairs
		}

		if err := validatePair(pair); err != nil {
			return err
		}
	}

	return nil
}

func validatePair(rawData map[string]interface{}) error {
	requiredFields := map[string]string{
		"term":  "string",
		"match": "string",
	}

	for field, expectedType := range requiredFields {
		val, exists := rawData[field]
		if !exists || !isValidType(val, expectedType) {
			return ErrInvalidPairs
		}
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
