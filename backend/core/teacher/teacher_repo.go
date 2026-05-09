package teacher

import (
	"context"
	"myapp/core/student"
)

type TeacherRepo interface {
	GetAllStudentsByTeacherID(ctx context.Context, teacherID int) ([]student.Student, error)
	FindByUsername(ctx context.Context, username string) (Teacher, error)
	CreateTeacher(ctx context.Context, username, passwordHash string) error
}
