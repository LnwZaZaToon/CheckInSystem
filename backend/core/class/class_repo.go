package class

import "context"

type ClassRepo interface {
	GetAllClassesByTeacherID(ctx context.Context, teacherID int) ([]Class, error)
	CreateClass(ctx context.Context, teacherID int, subjectName string) (int, error)
}
