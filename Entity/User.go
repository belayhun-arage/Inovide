package entity

type User struct {
	userName   string
	email      string
	password   string
	createdDay string
	profilePic string
	firstName  string
	lastName   string
	biography  string
	id         int32
}

func (user *User) setUserName(name string) {
	user.userName = name
}

func (user *User) generateImageDir() {

	user.profilePic = "sjfksjdkfj"
}
