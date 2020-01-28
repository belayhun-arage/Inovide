package AdminService

import (
	"github.com/Projects/Inovide/Admin"
	entity "github.com/Projects/Inovide/models"
)

type AdminService struct {
	Adminrepo *Admin.AdminRepo
}

func NewAdminService(adminrepo *Admin.AdminRepo) *AdminService {
	return &AdminService{Adminrepo: adminrepo}
}

func (adminservice *AdminService) CreateAdmin(admin *entity.Person) *entity.SystemMessage {

	systemmessage := &entity.SystemMessage{}
	return systemmessage
}
