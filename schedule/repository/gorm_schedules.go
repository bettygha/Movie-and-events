package repository

import (
	"github.com/hannasamuel20/Movie-and-events/model"
	"github.com/hannasamuel20/Movie-and-events/schedule"
	"github.com/jinzhu/gorm"
)

type ScheduleGormRepo struct {
	conn *gorm.DB
}

// NewCommentGormRepo returns new object of CommentGormRepo
func NewScheduleGormRepo(db *gorm.DB) schedule.ScheduleRepository {
	return &ScheduleGormRepo{conn: db}
}

// Comments returns all customer comments stored in the database
func (scheduleRepo *ScheduleGormRepo) Schedules() ([]model.Schedule, []error) {
	schdls := []model.Schedule{}
	errs := scheduleRepo.conn.Find(&schdls).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return schdls, errs
}

// StoreComment stores a given customer comment in the database
func (scheduleRepo *ScheduleGormRepo) StoreSchedule(schedule *model.Schedule) (*model.Schedule, []error) {
	schdl := schedule
	errs := scheduleRepo.conn.Create(schdl).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return schdl, errs
}

// Comment retrieves a customer comment from the database by its id
// func (cmntRepo *CommentGormRepo) Comment(id uint) (*entity.Comment, []error) {
// 	cmnt := entity.Comment{}
// 	errs := cmntRepo.conn.First(&cmnt, id).GetErrors()
// 	if len(errs) > 0 {
// 		return nil, errs
// 	}
// 	return &cmnt, errs
// }

// UpdateComment updates a given customer comment in the database
// func (cmntRepo *CommentGormRepo) UpdateComment(comment *entity.Comment) (*entity.Comment, []error) {
// 	cmnt := comment
// 	errs := cmntRepo.conn.Save(cmnt).GetErrors()
// 	if len(errs) > 0 {
// 		return nil, errs
// 	}
// 	return cmnt, errs
// }

// DeleteComment deletes a given customer comment from the database
// func (cmntRepo *CommentGormRepo) DeleteComment(id uint) (*entity.Comment, []error) {
// 	cmnt, errs := cmntRepo.Comment(id)

// 	if len(errs) > 0 {
// 		return nil, errs
// 	}

// 	errs = cmntRepo.conn.Delete(cmnt, id).GetErrors()
// 	if len(errs) > 0 {
// 		return nil, errs
// 	}
// 	return cmnt, errs
// }
