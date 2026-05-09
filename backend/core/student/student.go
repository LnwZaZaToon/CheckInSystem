package student

type Student struct {
	ID        int
	Firstname string
	Lastname  string
	School    string
}

type Class struct {
	ID          int
	SubjectName string
}

type StudentClass struct {
	StudentID int
	ClassID   int
}
