package slo

import (
	"net/http"
	"strconv"
	"time"

	"slo-tracker/pkg/errors"
	"slo-tracker/pkg/respond"
	"slo-tracker/schema"
	"slo-tracker/utils"
)

// getAllSLOsHandler fetches and unmarshal the slo data
func getAllSLOsHandler(w http.ResponseWriter, r *http.Request) *errors.AppError {
	slos, err := store.SLO().All()
	if err != nil {
		return err
	}

	// this is for backward compatibility with FE
	slosResponse := make([]schema.SLOResponse, len(slos))
	for i, slo := range slos {
		slosResponse[i] = schema.SLOResponse{
			ID:                 slo.ID,
			SLOName:            slo.SLOName,
			TargetSLO:          slo.TargetSLO,
			CurrentSLO:         0,
			RemainingErrBudget: 0,
		}
	}

	respond.OK(w, slosResponse)
	return nil
}

// creates a new slo
func createSLOHandler(w http.ResponseWriter, r *http.Request) *errors.AppError {
	var input schema.SLOPayload

	if err := utils.Decode(r, &input); err != nil {
		return errors.BadRequest(err.Error()).AddDebug(err)
	}

	slo, err := store.SLO().Create(&schema.SLO{
		SLOName:         input.SLOName,
		TargetSLO:       input.TargetSLO,
		HolidaysEnabled: &input.HolidaysEnabled,
	})

	if err != nil {
		return err
	}

	for _, workingDay := range input.Days {
		_, err := store.SLO().CreateWorkingSchedule(&schema.StoreWorkingSchedule{
			SLOID:     slo.ID,
			Weekday:   workingDay.Weekday,
			OpenHour:  workingDay.OpenHour,
			CloseHour: workingDay.CloseHour,
		})

		if err != nil {
			return err
		}
	}

	respond.Created(w, slo)
	return nil
}

// Get SLO details by ID
func getSLOHandler(w http.ResponseWriter, r *http.Request) *errors.AppError {
	ctx := r.Context()
	SLO, _ := ctx.Value("SLO").(*schema.SLO)

	yearMonthStr := r.URL.Query().Get("yearMonth")

	if yearMonthStr == "" { // parameter not present or is empty, return all incidents
		return errors.BadRequest("parameter 'yearMonth' must be present and must not be empty")
	} else {
		_, err := time.Parse("2006-01", yearMonthStr)
		if err != nil {
			return errors.BadRequest("Error parsing yearMonth, parameter should be like '2025-08'").AddDebug(err)
		}
	}

	incidents, err := store.Incident().GetByYearMonth(SLO.ID, yearMonthStr)
	if err != nil {
		return err
	}

	remaningErrBudget, currentSLO := utils.CalculateMonthlyErrBudget(SLO, incidents, yearMonthStr)

	SloResponse := schema.SLOResponse{
		ID:                 SLO.ID,
		SLOName:            SLO.SLOName,
		TargetSLO:          SLO.TargetSLO,
		CurrentSLO:         currentSLO,
		RemainingErrBudget: remaningErrBudget,
	}

	respond.OK(w, &SloResponse)
	return nil
}

// Updates the slo
func updateSLOHandler(w http.ResponseWriter, r *http.Request) *errors.AppError {
	var input schema.SLOPayload
	var isReset bool = true

	ctx := r.Context()
	slo, _ := ctx.Value("SLO").(*schema.SLO)

	if err := utils.Decode(r, &input); err != nil {
		return errors.BadRequest(err.Error()).AddDebug(err)
	}

	isReset, _ = strconv.ParseBool(r.URL.Query().Get("isReset"))

	// Consider already spent error budget when isReset set to false
	// If not, reset Error budget completely and delete past incidents
	if isReset {
		store.IncidentConn.Delete(slo.ID)
	}

	updated, err := store.SLO().Update(slo, &schema.SLO{
		SLOName:         input.SLOName,
		TargetSLO:       input.TargetSLO,
		HolidaysEnabled: &input.HolidaysEnabled,
	})
	if err != nil {
		return err
	}

	// delete working hours
	store.SLO().DeleteWorkingSchedule(slo.ID)

	for _, workingDay := range input.Days {
		_, err := store.SLO().CreateWorkingSchedule(&schema.StoreWorkingSchedule{
			SLOID:     slo.ID,
			Weekday:   workingDay.Weekday,
			OpenHour:  workingDay.OpenHour,
			CloseHour: workingDay.CloseHour,
		})

		if err != nil {
			return err
		}
	}

	respond.OK(w, updated)
	return nil
}

// Deletes the slo
func deleteSLOHandler(w http.ResponseWriter, r *http.Request) *errors.AppError {
	ctx := r.Context()
	slo, _ := ctx.Value("SLO").(*schema.SLO)

	// if err := utils.Decode(r, &input); err != nil {
	// 	return errors.BadRequest(err.Error()).AddDebug(err)
	// }

	if err := store.SLO().Delete(slo); err != nil {
		return errors.InternalServerStd().AddDebug(err)
	}

	respond.OK(w, slo)
	return nil
}

func getWorkingSchedule(w http.ResponseWriter, r *http.Request) *errors.AppError {

	ctx := r.Context()
	slo, _ := ctx.Value("SLO").(*schema.SLO)

	ws, err := store.SLO().GetWorkingSchedule(slo.ID)

	if err != nil {
		return err
	}

	sloPayload := &schema.SLOPayload{
		SLOName:         slo.SLOName,
		TargetSLO:       slo.TargetSLO,
		Days:            []schema.WorkingDaySchedule{},
		HolidaysEnabled: *slo.HolidaysEnabled,
	}

	for _, w := range *ws {
		sloPayload.Days = append(sloPayload.Days, schema.WorkingDaySchedule{
			Weekday:   w.Weekday,
			OpenHour:  w.OpenHour,
			CloseHour: w.CloseHour,
		})
	}

	respond.OK(w, sloPayload)
	return nil
}

func getOverview(w http.ResponseWriter, r *http.Request) *errors.AppError {

	yearMonthStr := r.URL.Query().Get("yearMonth")

	if yearMonthStr == "" { // parameter not present or is empty, return all incidents
		return errors.BadRequest("parameter 'yearMonth' must be present and must not be empty")
	} else {
		_, err := time.Parse("2006-01", yearMonthStr)
		if err != nil {
			return errors.BadRequest("Error parsing yearMonth, parameter should be like '2025-08'").AddDebug(err)
		}
	}

	query := `SELECT s.id,
				s.slo_name,
				s.target_slo,
				COUNT(i.id) AS num_incidents,
				COUNT(CASE WHEN i.mark_false_positive = 't' THEN 1 END) AS num_incidents_false_positive
		FROM slos s 
		LEFT JOIN incidents i ON s.id = i.slo_id AND TO_CHAR(i.created_at, 'YYYY-MM') = ?
		WHERE s.deleted_at is null
		GROUP BY s.id, s.slo_name, s.target_slo ORDER BY s.id;`

	var overviewRes []schema.OverviewResult
	store.DB.Raw(query, yearMonthStr).Scan(&overviewRes)

	for i, oslo := range overviewRes {

		slo, err := store.SLO().GetByID(oslo.Id)
		if err != nil {
			return err
		}

		incidents, err := store.Incident().GetByYearMonth(oslo.Id, yearMonthStr)
		if err != nil {
			return err
		}

		rem, c := utils.CalculateMonthlyErrBudget(slo, incidents, yearMonthStr)
		overviewRes[i].RemErrBudget = rem
		overviewRes[i].CurrentSlo = c
	}

	respond.OK(w, overviewRes)
	return nil
}

func updateFalsePositive(w http.ResponseWriter, r *http.Request) *errors.AppError {

	//var input schema.Incident
	ctx := r.Context()
	SLO, _ := ctx.Value("SLO").(*schema.SLO)

	incident, err := store.Incident().GetLatestBySLOID(SLO.ID)

	if err != nil {
		return err
	}
	if incident == nil {
		return errors.BadRequest("incident not found for SLO")
	}

	setFalsePositiveReq := &schema.SetFalsePositiveReq{}
	if err := utils.Decode(r, setFalsePositiveReq); err != nil {
		return errors.BadRequest(err.Error()).AddDebug(err)
	}

	if setFalsePositiveReq.MarkFalsePositive == nil {
		return errors.BadRequest("mark_false_positive is required")
	}

	updated := incident
	updated.MarkFalsePositive = *setFalsePositiveReq.MarkFalsePositive

	if incident, err = store.Incident().Update(incident, updated); err != nil {
		return err
	}

	respond.OK(w, incident)
	return nil
}
