package teacher

import (
	"context"
	"errors"
	"myapp/core/student"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type TeacherService interface {
	ImportStudentsFromCSV(ctx context.Context, teacherID int, csvData string) error
	GetAllStudents(ctx context.Context, teacherID int) ([]student.Student, error)
	GetAllClasses(ctx context.Context, teacherID int) ([]Class, error)
	Login(ctx context.Context, username, password string) error
	Register(ctx context.Context, username, password string) error
}

type teacherApp struct {
	teacherRepo TeacherRepo
}

type TeacherAppOption func(*teacherApp)

func WithTeacherRepo(repo TeacherRepo) TeacherAppOption {
	return func(o *teacherApp) {
		o.teacherRepo = repo
	}
}
func NewTeacherApp(opts ...TeacherAppOption) *teacherApp {
	o := &teacherApp{}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

func (a *teacherApp) GetAllStudents(ctx context.Context, teacherID int) ([]student.Student, error) {
	return a.teacherRepo.GetAllStudentsByTeacherID(ctx, teacherID)
}

func (a *teacherApp) ImportStudentsFromCSV(ctx context.Context, teacherID int, csvData string) error {
	// Implementation goes here
	return nil
}
func (a *teacherApp) Login(ctx context.Context, username, password string) error {
	teacher, err := a.teacherRepo.GetUsernameAndPassword(ctx, username)
	if err != nil {
		return err
	}
	if teacher.Username != username {
		return errors.New("invalid username")
	}
	err = bcrypt.CompareHashAndPassword([]byte(teacher.PasswordHash), []byte(password))
	if err != nil {
		return errors.New("invalid password")
	}
	return nil
}
func (a *teacherApp) Register(ctx context.Context, username, password string) error {
	teacher, err := a.teacherRepo.GetUsernameAndPassword(ctx, username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if len(teacher.Username) != 0 {
		return errors.New("username already exists")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	err = a.teacherRepo.CreateTeacher(ctx, username, string(hashedPassword))
	if err != nil {
		return err
	}
	return nil
}

func (a *teacherApp) GetAllClasses(ctx context.Context, teacherID int) ([]Class, error) {
	// Implementation goes here
	return nil, nil
}
