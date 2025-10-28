package main

import (
	"fmt"
	"os"
	"slo-tracker/config"
	"slo-tracker/schema"
	"slo-tracker/store"

	"github.com/gocarina/gocsv"
	"github.com/spf13/cobra"
)

type CSVRow struct {
	Name      string `csv:"name"`
	Holidays  string `csv:"holidays"`
	Weekday   int    `csv:"weekday"`
	OpenHour  string `csv:"open_hour"`
	CloseHour string `csv:"close_hour"`
}

func loadCSV(csvFile string) error {
	config.Initialize()
	conn := store.NewStore()

	file, err := os.Open(csvFile)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	rows := []*CSVRow{}
	if err := gocsv.UnmarshalFile(file, &rows); err != nil {
		return fmt.Errorf("failed to parse csv: %w", err)
	}

	for _, row := range rows {

		holidaysEnabled := row.Holidays == "Si"

		slo, err := conn.SLO().GetByName(row.Name)

		if err != nil { // store does not exists
			slo = &schema.SLO{
				SLOName:         row.Name,
				TargetSLO:       99.5,
				HolidaysEnabled: &holidaysEnabled,
			}
			conn.SLO().Create(slo)
		}

		sws := &schema.StoreWorkingSchedule{
			SLOID:     slo.ID,
			Weekday:   row.Weekday,
			OpenHour:  row.OpenHour,
			CloseHour: row.CloseHour,
		}

		conn.SLO().CreateWorkingSchedule(sws)

	}

	return nil
}

func main() {

	var rootCmd = &cobra.Command{
		Use:   "load-csv.go",
		Short: "Load data",
	}

	var migrateCmd = &cobra.Command{
		Use:   "load [csv_file]",
		Short: "Load data from a CSV file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			csvFile := args[0]
			return loadCSV(csvFile)
		},
	}

	rootCmd.AddCommand(migrateCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
