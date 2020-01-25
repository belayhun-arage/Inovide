package service

import (
	"fmt"

	repository "github.com/Projects/Inovide/User/Repository"
	entity "github.com/Projects/Inovide/models"
)

var er error

type UserService struct {
	userrepo *repository.UserRepo
}

func NewUserService(userep *repository.UserRepo) *UserService {
	return &UserService{userrepo: userep}
}

func (userService *UserService) RegisterUser(person *entity.Person) *entity.SystemMessage {
	var message = entity.SystemMessage{}

	if person.Username == "" || person.Email == "" || person.Password == "" {

		message.Message = "The Input Is Not Fully FIlled Please Submitt Again Filling The Data Appropriately"
		message.Succesful = false

	} else {
		er = userService.userrepo.CreateUser(person)
		if er != nil {
			panic(er)
		}
	}
	message.Message = "Succesfully Inserted "
	message.Succesful = true
	return &message

}

//This Method Checks Whether the User Presents or Not and Populate the Person struct Pointer with the data in the database
// if and only if the ser exists
// this method send and Recieve a message to and from UserHandler and UserRepository Respectively
//Returning a SystemMessage
//            For More Info About The System Messagelook for The System Message Struct in the models Folder entity Package
func (userService *UserService) CheckUser(person *entity.Person) *entity.SystemMessage {
	var message = entity.SystemMessage{}

	theBool := userService.userrepo.CheckUser(person)

	fmt.Println(person.ID, "The Final Time Of test")
	if theBool {
		message.Message = "Ok The User Exists "
		message.Succesful = true
	} else {
		message.Message = "No Can't Get The User "
		message.Succesful = false
	}
	return &message
}

func (userService *UserService) GetUser(person *entity.Person) *entity.SystemMessage {

	message := &entity.SystemMessage{}
	bools := userService.userrepo.GetUser(person)

	if bools {
		message.Message = "Succesfully Fetched "
		message.Succesful = true
		return message
	}
	message.Message = "Noooo "
	message.Succesful = false
	return message
}

func (UserService *UserService) AdminDeleteuser(person *entity.Person) *entity.SystemMessage {
	systemmessage := &entity.SystemMessage{}
	errors := UserService.userrepo.DeleteUser(person)
	if errors != nil {

		systemmessage.Message = "Can't Delete The User "
		systemmessage.Succesful = false

	} else {
		systemmessage.Message = "Succesfully Deleted The user"
		systemmessage.Succesful = true
	}
	return systemmessage
}

func (userService *UserService) GetUserById(person *entity.Person) *entity.SystemMessage {
	message := &entity.SystemMessage{}
	bools := userService.userrepo.GetUserById(person)

	if bools {
		message.Message = "Succesfully Fetched "
		message.Succesful = true
		return message
	}
	message.Message = "Noooo "
	message.Succesful = false
	return message

}

func (userservice *UserService) FollowUser(followerid, followingid int) *entity.SystemMessage {
	systemmessage := &entity.SystemMessage{}

	erro := userservice.userrepo.FollowUser(followingid, followerid)
	if erro != nil {
		return systemmessage
	}
	return systemmessage
}
func (userservice *UserService) UpdateUser(person *entity.Person) *entity.SystemMessage {
	systemmessage := &entity.SystemMessage{}

	erro := userservice.userrepo.UpdateUser(person)

	if erro != nil {

		systemmessage.Message = "Can't Update The Profile Of The User "
		systemmessage.Succesful = true
	} else {
		systemmessage.Message = " Succesfull yupdated the user"
		systemmessage.Succesful = true
	}
	return systemmessage
}
func (userservice *UserService) UnFollowUser(followingid, followerid int) *entity.SystemMessage {
	systemmessage := &entity.SystemMessage{}

	erro := userservice.userrepo.UnFollowUser(followingid, followerid)
	if erro != nil {
		systemmessage.Message = "Operation Is Not Succesfull"
		systemmessage.Succesful = false
	} else {
		systemmessage.Message = "Succesfully Removed Followship "
		systemmessage.Succesful = true
	}
	return systemmessage
}

func (userservice *UserService) NewIdeas(person *entity.Person) (*entity.SystemMessage, *[]entity.Idea) {
	idArrayOFPersonsIAmFollowing, err := userservice.userrepo.ListOfFolowingId(person) // Returning []int and error
	systemmessage := &entity.SystemMessage{}
	if err != nil {
		systemmessage.Message = "Can't Get The Id OF Users I Am Following "
		systemmessage.Succesful = false
		return systemmessage, nil
	} else {
		systemmessage.Succesful = true
		systemmessage.Message = "The List of Users Id I AM Followisng Found "
	}

	listofIdeas, err := userservice.userrepo.ListOfIdeasById(idArrayOFPersonsIAmFollowing) // returning the list of ideas that for each User Id Passing the List of following

	if err != nil {
		systemmessage.Message = "Can't Get The Id OF Users I Am Following "
		systemmessage.Succesful = false
		return systemmessage, nil
	}
	return systemmessage, listofIdeas
}
