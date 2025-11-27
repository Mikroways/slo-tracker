package config

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

const (
	// EnvDev const represents dev environment
	EnvDev = "dev"
	// EnvStaging const represents staging environment
	EnvStaging = "staging"
	// EnvProduction const represents production environment
	EnvProduction = "production"
)

type Holiday struct {
	Date string `json:"fecha"`
	Type string `json:"tipo"`
	Name string `json:"nombre"`
}

// Initialize ...
func Initialize() {

	viper.AutomaticEnv()

	viper.SetDefault("ENV", EnvDev)
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("DB_DRIVER", "postgres")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "root")
	viper.SetDefault("DB_PASS", "SecretPassword")
	viper.SetDefault("DB_NAME", "slotracker")
	viper.SetDefault("HOLIDAYS_ENDPOINT", "https://api.argentinadatos.com/v1/feriados/")
	viper.SetDefault("HOLIDAYS_DATES", []string{})
	viper.SetDefault("TZ", "America/Argentina/Buenos_Aires")

	currentYear := time.Now().Year()
	viper.SetDefault("HOLIDAY_YEAR", currentYear)

	_, err := refreshHolidays(currentYear)

	if err != nil {
		log.Println("Failing refreshing holidays: ", err)
	}

}

func FetchHolidays(year int) []string {

	if year != viper.Get("HOLIDAY_YEAR") {
		_, err := refreshHolidays(year)

		if err != nil {
			log.Println("Failing fetching holidays: ", err)
		}

		viper.SetDefault("HOLIDAY_YEAR", year)
	}

	return viper.Get("HOLIDAYS_DATES").([]string)
}

func refreshHolidays(year int) ([]Holiday, error) {

	resp, err := http.Get(viper.Get("HOLIDAYS_ENDPOINT").(string) + strconv.Itoa(year))

	if err != nil {
		return nil, fmt.Errorf("error fetching API: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad response: %s", resp.Status)
	}

	var holidays []Holiday

	if err := json.NewDecoder(resp.Body).Decode(&holidays); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	dates := make([]string, len(holidays))

	for i, val := range holidays {
		dates[i] = val.Date
	}

	viper.SetDefault("HOLIDAYS_DATES", dates)

	return holidays, nil
}
