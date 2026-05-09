package teacher

import (
	"context"
	"errors"
	coreStudent "myapp/core/student"
	coreTeacher "myapp/core/teacher"

	"gorm.io/gorm"
)

type teacherRow struct {
	ID           int    `gorm:"column:id;primaryKey;autoIncrement"`
	Username     string `gorm:"column:username"`
	PasswordHash string `gorm:"column:password_hash"`
}

func (teacherRow) TableName() string { return "teacher" }

type studentRow struct {
	ID        int    `gorm:"column:id;primaryKey;autoIncrement"`
	Firstname string `gorm:"column:firstname"`
	Lastname  string `gorm:"column:lastname"`
	School    string `gorm:"column:school"`
}

func (studentRow) TableName() string { return "student" }

type GormTeacherRepo struct {
	db *gorm.DB
}

func NewGormTeacherRepo(db *gorm.DB) *GormTeacherRepo {
	return &GormTeacherRepo{db: db}
}

func (r *GormTeacherRepo) GetAllStudentsByTeacherID(ctx context.Context, teacherID int) ([]coreStudent.Student, error) {
	var rows []studentRow
	err := r.db.WithContext(ctx).
		Select("DISTINCT student.*").
		Joins("JOIN student_class ON student_class.student_id = student.id").
		Joins("JOIN class ON class.id = student_class.class_id").
		Joins("JOIN teacher_class ON teacher_class.class_id = class.id").
		Where("teacher_class.teacher_id = ?", teacherID).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	students := make([]coreStudent.Student, len(rows))
	for i, row := range rows {
		students[i] = coreStudent.Student{
			ID:        row.ID,
			Firstname: row.Firstname,
			Lastname:  row.Lastname,
			School:    row.School,
		}
	}
	return students, nil
}

func (r *GormTeacherRepo) FindByUsername(ctx context.Context, username string) (coreTeacher.Teacher, error) {
	var row teacherRow
	err := r.db.WithContext(ctx).
		Where("username = ?", username).
		First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return coreTeacher.Teacher{}, coreTeacher.ErrNotFound
		}
		return coreTeacher.Teacher{}, err
	}
	return coreTeacher.Teacher{
		ID:           row.ID,
		Username:     row.Username,
		PasswordHash: row.PasswordHash,
	}, nil
}

func (r *GormTeacherRepo) CreateTeacher(ctx context.Context, username, passwordHash string) error {
	row := teacherRow{
		Username:     username,
		PasswordHash: passwordHash,
	}
	return r.db.WithContext(ctx).Create(&row).Error
}
