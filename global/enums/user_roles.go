package enums

import "fmt"

var UserRoles = map[int]string{
	1: "Administrador",
	2: "Usuário",
}

func GetUserRole(id int) (string, error) {
	role, exists := UserRoles[id]
	if !exists {
		return "", fmt.Errorf("role não encontrada para o ID: %d", id)
	}
	return role, nil
}
