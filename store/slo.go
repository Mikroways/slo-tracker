package store

import (
	"slo-tracker/pkg/errors"
	"slo-tracker/schema"

	// appstore "slo-tracker/store"

	"gorm.io/gorm"
)

// SLOStore implements the SLO interface
type SLOStore struct {
	*Conn
}

// NewSLOStore ...
func NewSLOStore(st *Conn) *SLOStore {
	cs := &SLOStore{st}
	return cs
}

// All returns all the SLOs
func (cs *SLOStore) All() ([]*schema.SLO, *errors.AppError) {
	var SLOs []*schema.SLO
	if err := cs.DB.Find(&SLOs).Error; err != nil { // For displaying all the columns
		// if err := cs.DB.Select("SliName, Alertsource, State, CreatedAt, ErrorBudgetSpent, MarkFalsePositive").Find(&SLOs).Error; err != nil {
		return nil, errors.InternalServerStd().AddDebug(err)
	}

	return SLOs, nil
}

// GetByID returns the matched record for the given id
func (cs *SLOStore) GetByID(SLOID uint) (*schema.SLO, *errors.AppError) {
	var SLO schema.SLO
	if err := cs.DB.First(&SLO, "id=?", SLOID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.BadRequest("invalid SLO id").AddDebug(err)
		}
		return nil, errors.InternalServerStd().AddDebug(err)
	}

	return &SLO, nil
}

// GetByName returns the matched record for the given slo_name
func (cs *SLOStore) GetByName(SLOName string) (*schema.SLO, *errors.AppError) {
	var SLO schema.SLO
	if err := cs.DB.First(&SLO, "slo_name=?", SLOName).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.BadRequest("invalid SLO name").AddDebug(err)
		}
		return nil, errors.InternalServerStd().AddDebug(err)
	}

	return &SLO, nil
}

// Create a new SLO
func (cs *SLOStore) Create(req *schema.SLO) (*schema.SLO, *errors.AppError) {

	slo := &schema.SLO{
		SLOName:   req.SLOName,
		TargetSLO: req.TargetSLO,
		OpenHour:  req.OpenHour,
		CloseHour: req.CloseHour,
	}
	if err := cs.DB.Save(slo).Error; err != nil {
		return nil, errors.InternalServerStd().AddDebug(err)
	}

	return slo, nil
}

// Update the SLO record..
func (cs *SLOStore) Update(SLO *schema.SLO, update *schema.SLO) (*schema.SLO, *errors.AppError) {
	if err := cs.DB.Model(SLO).Updates(update).Error; err != nil {
		return nil, errors.InternalServerStd().AddDebug(err)
	}
	return SLO, nil
}

// Delete the SLO record..
func (cs *SLOStore) Delete(SLO *schema.SLO) *errors.AppError {
	if err := cs.DB.Delete(&SLO).Error; err != nil {
		return errors.InternalServerStd().AddDebug(err)
	}
	return nil
}
