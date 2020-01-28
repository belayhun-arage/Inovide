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

var message *entity.SystemMessage

func (userService *UserService) RegisterUser(person *entity.Person) *entity.SystemMessage {
	message = &entity.SystemMessage{}
	_, val := userService.userrepo.CreateUser(person)
	defer recover()
	message.Message = "Can't Save The User "
	message.Succesful = false
	// if er != nil {
	if val == 0 {
		message.Message = "Invalid Username Or Email "
		return message
	}

	message.Message = "Succesfully Inserted "
	message.Succesful = true
	return message
}

//This Method Checks Whether the User Presents or Not and Populate the Person struct Pointer with the data in the database
// if and only if the ser exists
// this method send and Recieve a message to and from UserHandler and UserRepository Respectively
//Returning a SystemMessage
//            For More Info About The System Messagelook for The System Message Struct in the models Folder entity Package
var systemmessagechu *entity.SystemMessage

func (userService *UserService) CheckUser(person *entity.Person) *entity.SystemMessage {
	systemmessagechu = &entity.SystemMessage{}

	theBool := userService.userrepo.CheckUser(person)

	fmt.Print(person.ID, " ## ChechUser\n\n")
	if theBool {
		systemmessagechu.Message = "Ok The User Exists "
		systemmessagechu.Succesful = true
	} else {
		systemmessagechu.Message = "No Can't Get The User "
		systemmessagechu.Succesful = false
	}
	return systemmessagechu
}

func (userService *UserService) GetUser(person *entity.Person) *entity.SystemMessage {

	message := &entity.SystemMessage{}
	bools := userService.userrepo.GetUser(person)

	if bools == 1 {

		message.Message = "Succesfully Fetched "
		message.Succesful = true
		return message
	}
	fmt.Println("failure ... ")

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

	erro, val := userservice.userrepo.UpdateUser(person)

	if erro != nil {

		if val == 1 {
			systemmessage.Message = "You Can't Use This Email "

		} else if val == 2 {
			systemmessage.Message = "There Is One By This Username! Please Make A change "

		}
		systemmessage.Succesful = false

	} else {
		systemmessage.Message = " Succesfull updated the user"
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
