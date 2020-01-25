package repository

/*This package Will Responsibel For Manipulating the database and Genereating an Instance of User to be used by the Service */
import (
	"fmt"

	entity "github.com/Projects/Inovide/models"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(sqlite *gorm.DB) *UserRepo {
	return &UserRepo{db: sqlite}
}

func (users *UserRepo) CreateUser(enti *entity.Person) error {

	err := users.db.Debug().Table("users").Model(&entity.Person{}).Create(enti).Error
	if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
		// handle error
	}
	fmt.Println("-----------------------")
	if err != nil {
		return err
	}
	return nil
}

func (users *UserRepo) CheckUser(enti *entity.Person) bool {

	person := entity.Person{}
	err := users.db.Debug().Table("users").Where("Username=?", enti.Username).Find(person /*enti.Username, enti.Password*/).Error
	//users.db.Table("users").Select("ID").Debug().Model(&entity.Person{}).Where("UserName=$1 AND Password=$2", enti.Username, enti.Password).Find(&person) //Select([]string{"UserName", "Email", "Password"}).Find(person  , )

	if err != nil {
		return false
	}
	fmt.Println(person.ID, "_______-------<< User Repo")
	// fmt.Println(peoples.Username, peoples.Password, peoples.Email)
	if person.Username == "" || person.Password == "" || person.Email == "" {
		return false
	}
	enti = &person
	return true
}
func (users *UserRepo) GetUser(enti *entity.Person) bool {

	geterr := users.db.Debug().Table("users").Model(&entity.Person{}).Where("UserName=? and Password=?", enti.Username, enti.Password).Find(&enti).Error
	//updateerr := users.db.Debug().Table("users").Model(&entity.Person{}).Set("Firstname = ?Firstname").Where("id = ?id").Update(enti)
	if geterr != nil {
		return false
	}
	return true

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
func (Userrepo *UserRepo) UpdateUser(person *entity.Person) []error {
	erro := Userrepo.db.Model(&entity.Person{}).Table("users").Save(person).GetErrors()
	return erro
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
