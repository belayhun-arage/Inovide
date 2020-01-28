package Admin

type AdminRepo interface {
	CountUsers() int64
	CountIdeas() int64
	CountAdmins() int64
	CountMessages() int64
	CountActiveUsers() int64
}
