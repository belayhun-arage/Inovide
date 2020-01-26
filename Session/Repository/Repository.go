package SessionRepo

import (
	"fmt"

	entity "github.com/Projects/Inovide/models"
	"github.com/jinzhu/gorm"
)

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepo(dbs *gorm.DB) *SessionRepository {
	return &SessionRepository{db: dbs}
}

func (sessionrepo *SessionRepository) CreateSession(session *entity.Session) []error {
	errors := sessionrepo.db.Table("session").Create(session).GetErrors()
	return errors
}
func (sessionrepo *SessionRepository) DeleteSession(session *entity.Session) int {
	fmt.Println("Delete Session Service ", session.Userid, session.Username)

	errors := sessionrepo.db.Debug().Table("session").Where("username=? and userid=?", session.Username, session.Userid).Delete(session).Error
	if errors == gorm.ErrRecordNotFound {
		return 1
	}
	return 0
}
