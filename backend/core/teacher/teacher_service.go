package teacher

import (
	"context"
	"errors"
	"myapp/core/student"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type TeacherService interface {
	ImportStudentsFromCSV(ctx context.Context, teacherID int, csvData string) error
	GetAllStudents(ctx context.Context, teacherID int) ([]student.Student, error)
	Login(ctx context.Context, username, password string) (string, error)
	Register(ctx context.Context, username, password string) error
}

type teacherApp struct {
	teacherRepo TeacherRepo
	jwtSecret   []byte
}

type TeacherAppOption func(*teacherApp)

func WithTeacherRepo(repo TeacherRepo) TeacherAppOption {
	return func(o *teacherApp) {
		o.teacherRepo = repo
	}
}

func WithJWTSecret(secret string) TeacherAppOption {
	return func(o *teacherApp) {
		o.jwtSecret = []byte(secret)
	}
}

func NewTeacherApp(opts ...TeacherAppOption) TeacherService {
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
	return nil
}

func (a *teacherApp) Login(ctx context.Context, username, password string) (string, error) {
	teacher, err := a.teacherRepo.FindByUsername(ctx, username)
	if err != nil {
		return "", errors.New("invalid username or password")
	}
	if err = bcrypt.CompareHashAndPassword([]byte(teacher.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid username or password")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"teacher_id": teacher.ID,
		"username":   teacher.Username,
		"exp":        time.Now().Add(24 * time.Hour).Unix(),
	})
	signed, err := token.SignedString(a.jwtSecret)
	if err != nil {
		return "", err
	}
	return signed, nil
}

func (a *teacherApp) Register(ctx context.Context, username, password string) error {
	_, err := a.teacherRepo.FindByUsername(ctx, username)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return err
	}
	if err == nil {
		return ErrUsernameExists
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return a.teacherRepo.CreateTeacher(ctx, username, string(hashedPassword))
}
