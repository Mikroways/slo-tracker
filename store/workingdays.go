package store

import (
	"slo-tracker/pkg/errors"
	"slo-tracker/schema"

	"gorm.io/gorm"
)

func (cs *SLOStore) CreateWorkingSchedule(req *schema.StoreWorkingSchedule) (*schema.StoreWorkingSchedule, *errors.AppError) {

	if err := cs.DB.Save(req).Error; err != nil {
		return nil, errors.InternalServerStd().AddDebug(err)
	}

	return req, nil
}

func (cs *SLOStore) GetWorkingSchedule(SLOID uint) (*[]schema.StoreWorkingSchedule, *errors.AppError) {

	var ws []schema.StoreWorkingSchedule
	if err := cs.DB.Where("slo_id=?", SLOID).Find(&ws).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.BadRequest("invalid SLO id").AddDebug(err)
		}
		return nil, errors.InternalServerStd().AddDebug(err)
	}

	return &ws, nil
}

// Delete the SLO record..
func (cs *SLOStore) DeleteWorkingSchedule(SLOID uint) *errors.AppError {
	if err := cs.DB.Where("slo_id=?", SLOID).Delete(&schema.StoreWorkingSchedule{}).Error; err != nil {
		return errors.InternalServerStd().AddDebug(err)
	}
	return nil
}
