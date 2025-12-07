package teacher

import (
	"context"
	coreStudent "myapp/core/student"
	coreTeacher "myapp/core/teacher"

	"gorm.io/gorm"
)

type GormTeacherRepo struct {
	db *gorm.DB
}

func NewGormTeacherRepo(db *gorm.DB) *GormTeacherRepo {
	return &GormTeacherRepo{db: db}
}

func (r *GormTeacherRepo) GetAllStudentsByTeacherID(ctx context.Context, teacherID int) ([]coreStudent.Student, error) {
	students := []coreStudent.Student{}
	err := r.db.WithContext(ctx).
		Table("student").
		Select("DISTINCT student.*"). // ใช้ DISTINCT ป้องกันนักเรียนคนเดิมซ้ำ (กรณีเรียนหลายวิชากับครูคนเดิม)
		Joins("JOIN student_class ON student_class.student_id = student.id").
		Joins("JOIN class ON class.id = student_class.class_id").
		Joins("JOIN teacher_class ON teacher_class.class_id = class.id").
		Where("teacher_class.teacher_id = ?", teacherID).
		Find(&students).Error
	if err != nil {
		return nil, err
	}
	return students, nil
}

func (r *GormTeacherRepo) GetUsernameAndPassword(ctx context.Context, username string) (coreTeacher.Teacher, error) {
	type TeacherCredential struct {
		Username     string `gorm:"column:username"`
		PasswordHash string `gorm:"column:password_hash"`
	}
	var cread TeacherCredential
	err := r.db.WithContext(ctx).
		Table("teacher").
		Select("username, password_hash").
		Where("username = ?", username).
		First(&cread).Error
	if err != nil {
		return coreTeacher.Teacher{}, err
	}
	return coreTeacher.Teacher{
		Username:     cread.Username,
		PasswordHash: cread.PasswordHash,
	}, nil
}

func (r *GormTeacherRepo) CreateTeacher(ctx context.Context, username, passwordHash string) error {
	teacher := coreTeacher.Teacher{
		Username:     username,
		PasswordHash: passwordHash,
	}
	err := r.db.WithContext(ctx).Create(&teacher).Error
	if err != nil {
		return err
	}
	return nil
}
