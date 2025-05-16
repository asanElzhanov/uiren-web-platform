package users

import (
	"net/mail"
	"regexp"
	"strings"
	"unicode"

	"github.com/nyaruka/phonenumbers"
)

func (dto *UserDTO) normalize() {
	dto.Username = strings.TrimSpace(dto.Username)
	dto.Firstname = strings.TrimSpace(dto.Firstname)
	dto.Lastname = strings.TrimSpace(dto.Lastname)
	dto.Phone = strings.TrimSpace(dto.Phone)
	dto.Email = strings.TrimSpace(dto.Email)
}

func (dto *CreateUserDTO) normalize() {
	dto.Username = strings.TrimSpace(dto.Username)
	dto.Email = strings.TrimSpace(dto.Email)
}

func (dto *UpdateUserDTO) normalize() {
	dto.Firstname = strings.TrimSpace(dto.Firstname)
	dto.Lastname = strings.TrimSpace(dto.Lastname)
	dto.Phone = strings.TrimSpace(dto.Phone)
	dto.PhoneRegion = strings.TrimSpace(dto.PhoneRegion)
}

func (dto *CreateUserDTO) Validate() error {
	dto.normalize()

	if !isValidEmail(dto.Email) {
		return ErrIncorrectEmail
	}
	if !isValidPassword(dto.Password) {
		return ErrIncorrectPassword
	}

	return nil
}

func (dto *UpdateUserDTO) Validate() error {
	dto.normalize()

	if !isValidPhone(dto.Phone, dto.PhoneRegion) {
		return ErrIncorrectPhone
	}

	return nil
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func isValidPhone(phone, region string) bool {
	num, err := phonenumbers.Parse(phone, region)
	if err != nil {
		return false
	}

	if !phonenumbers.IsPossibleNumber(num) {
		return false
	}

	if !phonenumbers.IsValidNumber(num) {
		return false
	}

	if !phonenumbers.IsValidNumberForRegion(num, region) {
		return false
	}

	return true
}

func isValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	var hasDigit, hasSpecial, hasLowerCase, hasUpperCase bool
	for _, ch := range password {
		switch {
		case unicode.IsLetter(ch):
			if unicode.IsLower(ch) {
				hasLowerCase = true
			} else if unicode.IsUpper(ch) {
				hasUpperCase = true
			}
		case unicode.IsDigit(ch):
			hasDigit = true
		case isSpecialChar(ch):
			hasSpecial = true
		}
	}

	return hasLowerCase && hasUpperCase && hasDigit && hasSpecial
}

func isSpecialChar(ch rune) bool {
	matched, _ := regexp.MatchString(`[@$!%*?&]`, string(ch))
	return matched
}
