# Scripts

This section explains the utilization of custom scripts that complement slo-tracker

## load-csv

This script reads from a CSV file and inserts the data into `slos` and `store_working_schedules`, this way we can bulk insert a bunch of SLOS ready to be used

Example:
```sh
go run scripts/load-csv.go load path/to/file.csv
```

#### Struct of the CSV

In order to work, the CSV must have the folling columns:

| Column Name |  Meaning |
|-------------|----------|
| name        |  Name of SLO |
| holidays    |  Whether or not consider incidents happening on a Holiday |
| weekdays    |  Day of the week that the SLO works, 0 for Sunday, 1 for Monday, 2 for Tuesday and so on until 6 for Saturday |
| open_hour   | Time that the SLO starts working, format 24hs |
| close_hour  | Time that the SLO stops working, format 24hs |

As reference, you can refer to [this csv example](csv-example.csv)

#### DB connection

This script loads env vars, the same way slo-tracker does to interact with the database, so the script expects the following variables:

* 	DB_DRIVER
*   DB_HOST
*   DB_HOST
*   DB_PASS
*   DB_NAME

If any of them is not present, the scripts uses default values