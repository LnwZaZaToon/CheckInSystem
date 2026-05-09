package class

import (
	"context"
	coreClass "myapp/core/class"

	"gorm.io/gorm"
)

type classRow struct {
	ID          int    `gorm:"column:id;primaryKey;autoIncrement"`
	SubjectName string `gorm:"column:subject_name"`
}

func (classRow) TableName() string { return "class" }

type teacherClassRow struct {
	TeacherID int `gorm:"column:teacher_id"`
	ClassID   int `gorm:"column:class_id"`
}

func (teacherClassRow) TableName() string { return "teacher_class" }

type GormClassRepo struct {
	db *gorm.DB
}

func NewGormClassRepo(db *gorm.DB) *GormClassRepo {
	return &GormClassRepo{db: db}
}

func (r *GormClassRepo) GetAllClassesByTeacherID(ctx context.Context, teacherID int) ([]coreClass.Class, error) {
	var rows []classRow
	err := r.db.WithContext(ctx).
		Joins("JOIN teacher_class ON teacher_class.class_id = class.id").
		Where("teacher_class.teacher_id = ?", teacherID).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	classes := make([]coreClass.Class, len(rows))
	for i, row := range rows {
		classes[i] = coreClass.Class{
			ID:          row.ID,
			SubjectName: row.SubjectName,
		}
	}
	return classes, nil
}

func (r *GormClassRepo) CreateClass(ctx context.Context, teacherID int, subjectName string) (int, error) {
	row := classRow{SubjectName: subjectName}
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&row).Error; err != nil {
			return err
		}
		return tx.Create(&teacherClassRow{TeacherID: teacherID, ClassID: row.ID}).Error
	})
	if err != nil {
		return 0, err
	}
	return row.ID, nil
}
