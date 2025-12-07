package student

type Student struct {
	ID        int    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Firstname string `gorm:"column:firstname" json:"firstname"`
	Lastname  string `gorm:"column:lastname" json:"lastname"`
	School    string `gorm:"column:school" json:"school"` // ปรับชื่อ column ให้ตรง DB
}

func (Student) TableName() string {
	return "student"
}

type Class struct {
	ID          int    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	SubjectName string `gorm:"column:subject_name" json:"subject_name"`
}

func (Class) TableName() string {
	return "class"
}

type StudentClass struct {
	StudentID int `gorm:"column:student_id;primaryKey" json:"student_id"`
	ClassID   int `gorm:"column:class_id;primaryKey" json:"class_id"`
}

func (StudentClass) TableName() string {
	return "student_class"
}
