package repository

/*This package Will Responsibel For Manipulating the database and Genereating an Instance of User to be used by the Service */
import (
	"fmt"

	entity "github.com/Projects/Inovide/models"
	"github.com/jinzhu/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(sqlite *gorm.DB) *UserRepo {
	return &UserRepo{db: sqlite}
}

func (users *UserRepo) CreateUser(enti *entity.Person) (error, int64) {

	rowsaffected := users.db.Table("users").Create(enti).RowsAffected
	if rowsaffected <= 0 {
		// handle error
		fmt.Println("its ok ladies man ")
		fmt.Println(rowsaffected)

	}
	defer recover()
	// if err.Error() == "pq: duplicate key value violates unique constraint \"users_username_key\"" {
	// 	return err, 1
	// }
	// if err.Error() == "pq: duplicate key value violates unique constraint \"users_username_key\"" {
	// 	return err, 2
	// }
	return nil, rowsaffected
}

func (users *UserRepo) CheckUser(enti *entity.Person) bool {

	err := users.db.Table("users").Where("Username=?", enti.Username).Find(enti /*enti.Username, enti.Password*/).Error
	if err == gorm.ErrRecordNotFound {
		return false
	}
	fmt.Println(enti.ID, "_______-------<< User Repo")
	fmt.Println(enti.Username, enti.Password, enti.ID, enti.Firstname, enti.Lastname)
	return true
}
func (users *UserRepo) GetUser(enti *entity.Person) int64 {

	geterr := users.db.Debug().Table("users").Where("id=?", enti.ID).Find(&enti).RowsAffected
	//updateerr := users.db.Debug().Table("users").Model(&entity.Person{}).Set("Firstname = ?Firstname").Where("id = ?id").Update(enti)
	// if geterr <= {
	// 	return false
	// }
	defer recover()
	return geterr

}
func (users *UserRepo) GetUserById(enti *entity.Person) bool {
	geterr := users.db.Debug().Table("users").Model(&entity.Person{}).Where(&entity.Person{}, enti.ID).Find(&enti).Error
	if geterr != nil {
		return false
	}
	return true
}
func (userrepo *UserRepo) FollowUser(followingid, followerid int) error {
	person := &entity.Person{}
	err := userrepo.db.Debug().Table("users").Where(entity.Person{}, followingid).Find(person).Error
	fmt.Println(person.Followers)
	err = userrepo.db.Debug().Table("users").Where(" id=?", followingid).Update(&entity.Person{ID: person.ID, Firstname: person.Firstname, Lastname: person.Lastname, Username: person.Username, Password: person.Password, Email: person.Email, Biography: person.Biography, Followers: person.Followers + 1, Ideas: person.Ideas, Paid: person.Paid}).Error
	if err != nil {
		return err
	}
	err = userrepo.db.Debug().Table("following").Create(&entity.Following{FollowerId: followerid, FollowingId: followingid}).Error
	if err != nil {
		return err
	}
	return nil
}
func (Userrepo *UserRepo) UpdateUser(person *entity.Person) (error, int) {

	perso := &entity.Person{Username: person.Username,
		// Email:     person.Email,
		Firstname: person.Firstname,
		Lastname:  person.Lastname,
		Imagedir:  person.Imagedir,
		Paid:      person.Paid,
		Biography: person.Biography}

	erro := Userrepo.db.Table("users").Debug().Where("id=?", person.ID).Begin().Update(perso).Error

	defer recover()

	if person.ID > 0 {
		return nil, 0
	}
	if erro.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"" {
		return erro, 1
	}
	if erro.Error() == "pq: duplicate key value violates unique constraint \"users_username_key\"" {
		return erro, 2
	}
	return nil, 0
}

func (UserRepo *UserRepo) DeleteUser(person *entity.Person) []error {
	erro := UserRepo.db.Table("users").Model(&entity.Person{}).Debug().Delete(person).GetErrors()
	return erro
}
func (userrepo *UserRepo) UnFollowUser(followingid, followerid int) []error {
	person := &entity.Person{}
	err := userrepo.db.Debug().Table("users").Where(entity.Person{}, followingid).Find(person).GetErrors()
	fmt.Println(person.Followers)
	person.Followers -= 1
	err = userrepo.db.Debug().Table("users").Save(person).GetErrors()
	if err != nil {
		return err
	}
	err = userrepo.db.Debug().Table("following").Delete(&entity.Following{FollowerId: followerid, FollowingId: followingid}).GetErrors()
	if err != nil {
		return err
	}
	return nil
}

func (userrepo *UserRepo) ListOfFolowingId(person *entity.Person) ([]int, []error) {

	var IdList []int
	holder := userrepo.db.Model(&[]entity.Following{}).Table("following").Select("followingid").Where(&entity.Following{FollowerId: int(person.ID)}).Find(IdList).GetErrors()
	return IdList, holder
}
func (UserRepo *UserRepo) ListOfIdeasById(idis []int) (*[]entity.Idea, []error) {
	var ListOfIdeas []entity.Idea = []entity.Idea{}
	for index, val := range idis {
		var Idea entity.Idea = entity.Idea{}
		Idea.Ideaownerid = val
		theError := UserRepo.db.Table("idea").Model(&entity.Idea{}).Debug().Find(Idea).GetErrors()
		if theError != nil {
			return &ListOfIdeas, theError
		}
		ListOfIdeas[index] = Idea
	}
	return &ListOfIdeas, nil
}
func (userrepo *UserRepo) NumberOfFollowers(person *entity.Person) []error {
	return nil
}
func (userrepo *UserRepo) UploadProfilePicture(person *entity.Person) []error {
	theerror := userrepo.db.Debug().Table("users").Save(person).Find(person).GetErrors()
	if theerror != nil {
		return theerror
	}
	return nil
}
