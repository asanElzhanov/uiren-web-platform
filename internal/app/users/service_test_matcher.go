package users

import (
	"fmt"
)

type createUserDTOMatcher struct {
	dto CreateUserDTO
}

func (m createUserDTOMatcher) Matches(x interface{}) bool {
	dto, ok := x.(CreateUserDTO)
	if !ok {
		return false
	}
	return dto.Username == m.dto.Username &&
		dto.Email == m.dto.Email &&
		len(dto.Password) == 60
}

func (m createUserDTOMatcher) String() string {
	return fmt.Sprintf("matches CreateUserDTO %v", m.dto)
}
