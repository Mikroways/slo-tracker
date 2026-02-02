package store

import (
	"slo-tracker/pkg/errors"
	"slo-tracker/schema"
	"slo-tracker/utils"

	"gorm.io/gorm"
)

// IncidentStore implements the incident interface
type IncidentStore struct {
	*Conn
}

// NewIncidentStore ...
func NewIncidentStore(st *Conn) *IncidentStore {
	cs := &IncidentStore{st}
	return cs
}

// All returns all the Incidents
func (cs *IncidentStore) All(SLOID uint) ([]*schema.Incident, *errors.AppError) {
	var Incidents []*schema.Incident
	if err := cs.DB.Order("created_at desc").Find(&Incidents, "slo_id=?", SLOID).Error; err != nil { // For displaying all the columns
		// if err := cs.DB.Select("SliName, Alertsource, State, CreatedAt, ErrorBudgetSpent, MarkFalsePositive").Find(&Incidents).Error; err != nil {
		return nil, errors.InternalServerStd().AddDebug(err)
	}

	return Incidents, nil
}

func (cs *IncidentStore) GetByYearMonth(SLOID uint, yearMonthStr string) ([]*schema.Incident, *errors.AppError) {
	var Incidents []*schema.Incident

	year, month, err := utils.ParseYearMonth(yearMonthStr)

	if err != nil {
		return nil, errors.BadRequest("Year and Month are not valid")
	}

	if err = cs.DB.Order("created_at desc").Where("EXTRACT(YEAR FROM created_at) = ? AND EXTRACT(MONTH FROM created_at) = ?", year, month).Find(&Incidents, "slo_id=?", SLOID).Error; err != nil {
		return nil, errors.InternalServerStd().AddDebug(err)
	}

	return Incidents, nil
}

// GetByID returns the matched record for the given id
func (cs *IncidentStore) GetByID(incidentID uint) (*schema.Incident, *errors.AppError) {
	var incident schema.Incident
	if err := cs.DB.First(&incident, "id=?", incidentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.BadRequest("invalid incident id").AddDebug(err)
		}
		return nil, errors.InternalServerStd().AddDebug(err)
	}

	return &incident, nil
}

// GetBySLIName returns the matched record for the given SLI
func (cs *IncidentStore) GetBySLIName(SLOID uint, sliName string) (*schema.Incident, *errors.AppError) {
	var incident schema.Incident
	if err := cs.DB.First(&incident, "state=? AND sli_name=? AND slo_id=?", "open", sliName, SLOID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.InternalServerStd().AddDebug(err)
		}
		return nil, errors.InternalServerStd().AddDebug(err)
	}

	return &incident, nil
}

// GetBySLINameV2 returns the matched record for the given SLI
func (cs *IncidentStore) GetBySLINameV2(sliName string) (*schema.Incident, *errors.AppError) {
	var incident schema.Incident
	if err := cs.DB.First(&incident, "sli_name=?", sliName).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.InternalServerStd().AddDebug(err)
		}
		return nil, errors.InternalServerStd().AddDebug(err)
	}

	return &incident, nil
}

// GetBySLINameAndOpenState returns the matched record for the given SLI
func (cs *IncidentStore) GetBySLINameAndOpenState(sliName string) (*schema.Incident, *errors.AppError) {
	var incident schema.Incident
	if err := cs.DB.First(&incident, "state=? AND sli_name=?", "open", sliName).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.InternalServerStd().AddDebug(err)
		}
		return nil, errors.InternalServerStd().AddDebug(err)
	}

	return &incident, nil
}

// Create a new Incident
func (cs *IncidentStore) Create(req *schema.IncidentReq) (*schema.Incident, *errors.AppError) {

	incident := &schema.Incident{
		SliName:          req.SliName,
		SLOID:            req.SLOID,
		Alertsource:      req.Alertsource,
		State:            req.State,
		CreatedAt:        req.CreatedAt,
		ErrorBudgetSpent: req.ErrorBudgetSpent,
		RealErrorBudget:  req.RealErrorBudget,
		Observations:     req.Observations,
	}
	if err := cs.DB.Save(incident).Error; err != nil {
		return nil, errors.InternalServerStd().AddDebug(err)
	}

	return incident, nil
}

// Update the incident record..
func (cs *IncidentStore) Update(incident *schema.Incident, update *schema.Incident) (*schema.Incident, *errors.AppError) {

	var err *errors.AppError

	if !incident.MarkFalsePositive && update.MarkFalsePositive {
		// Close the open incident if it's being marked as false positive
		if update.State == "open" {
			update.State = "closed"
		}
	}

	if err != nil {
		return nil, errors.BadRequest(err.Error()).AddDebug(err)
	}

	if err := cs.DB.Model(incident).Updates(map[string]interface{}{
		"State":             update.State,
		"ErrorBudgetSpent":  update.ErrorBudgetSpent,
		"MarkFalsePositive": update.MarkFalsePositive,
		"RealErrorBudget":   update.RealErrorBudget,
		"Observations":      update.Observations,
	}).Error; err != nil {
		return nil, errors.InternalServerStd().AddDebug(err)
	}

	return incident, nil
}

// Deletes all the incidents matching SLOID field
func (cs *IncidentStore) Delete(SLOID uint) *errors.AppError {
	if err := cs.DB.Delete(schema.Incident{}, "slo_id = ?", SLOID).Error; err != nil {
		return errors.InternalServerStd().AddDebug(err)
	}
	return nil
}

// GetLatestBySLOID returns the latest incident of a SLO
func (cs *IncidentStore) GetLatestBySLOID(SLOID uint) (*schema.Incident, *errors.AppError) {
	var incident schema.Incident
	if err := cs.DB.Order("created_at DESC").First(&incident, "slo_id=?", SLOID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.InternalServerStd().AddDebug(err)
		}
		return nil, errors.InternalServerStd().AddDebug(err)
	}

	return &incident, nil
}
