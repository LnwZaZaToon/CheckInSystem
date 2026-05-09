package teacher

type Teacher struct {
	ID           int
	Username     string
	PasswordHash string
}

type TeacherClass struct {
	TeacherID int
	ClassID   int
}
