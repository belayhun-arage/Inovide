package SessionRepo

import (
	entity "github.com/Projects/Inovide/models"
	"github.com/jinzhu/gorm"
)

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepo() *SessionRepository {
	return &SessionRepository{}
}

func (sessionrepo *SessionRepository) CreateSession(session *entity.Session) []error {
	errors := sessionrepo.db.Table("session").Debug().Create(session).GetErrors()
	return errors
}
func (sessionrepo *SessionRepository) DeleteSession(session *entity.Session) []error {
	errors := sessionrepo.db.Debug().Table("session").Delete(session).GetErrors()
	return errors
}
