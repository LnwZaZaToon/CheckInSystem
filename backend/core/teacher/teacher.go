package teacher

type Teacher struct {
	ID           int    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Username     string `gorm:"column:username" json:"username"`
	PasswordHash string `gorm:"column:password_hash" json:"password_hash"`
}

func (Teacher) TableName() string {
	return "teacher"
}

type Class struct {
	ID          int    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	SubjectName string `gorm:"column:subject_name" json:"subject_name"`
}

func (Class) TableName() string {
	return "class"
}

type TeacherClass struct {
	TeacherID int `gorm:"column:teacher_id;primaryKey" json:"teacher_id"`
	ClassID   int `gorm:"column:class_id;primaryKey" json:"class_id"`
}

func (TeacherClass) TableName() string {
	return "teacher_class"
}
