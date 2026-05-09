package class

import (
	"context"
	"myapp/core/student"
)

type ClassService interface {
	GetAllClassesByTeacherID(ctx context.Context, teacherID int) ([]Class, error)
	CreateClasses(ctx context.Context, teacherID int, subjectName string) (int, error)
	ImportStudentsFromCSV(ctx context.Context, teacherID int, student student.Student) (string, error)
}

type classApp struct {
	classRepo ClassRepo
}

type ClassAppOption func(*classApp)

func WithClassRepo(repo ClassRepo) ClassAppOption {
	return func(o *classApp) {
		o.classRepo = repo
	}
}

func NewClassApp(opts ...ClassAppOption) ClassService {
	o := &classApp{}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

func (a *classApp) GetAllClassesByTeacherID(ctx context.Context, teacherID int) ([]Class, error) {
	return a.classRepo.GetAllClassesByTeacherID(ctx, teacherID)
}

func (a *classApp) CreateClasses(ctx context.Context, teacherID int, subjectName string) (int, error) {
	if subjectName == "" {
		return 0, ErrSubjectNameRequired
	}
	return a.classRepo.CreateClass(ctx, teacherID, subjectName)
}

func (a *classApp) ImportStudentsFromCSV(ctx context.Context, teacherID int, s student.Student) (string, error) {
	return "", nil
}
