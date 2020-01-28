package Admin

import (
	entity "github.com/Projects/Inovide/models"
)

type AdminService interface {
	CreateAdmin(admin *entity.Person) *entity.SystemMessage
}
