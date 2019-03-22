package main

import (
	"flag"
	"log"
	"time"

	logger "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const dateFormatLayout = "2006-01-02"

var (
	conf *config

	fromString string
	toString   string
	dryRun     bool
	verbose    bool
	summary    bool
)

func init() {

	logger.SetLevel(logger.WarnLevel)

	// Load Configuration
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		logger.Fatalf("Fatal error config file: %s", err)
	}

	err = viper.Unmarshal(&conf)
	if err != nil {
		logger.Fatalf("Unable to decode into config struct, %v", err)
	}

}

func main() {
	// Flags
	flag.StringVar(&fromString, "from", "2019-03-01", "From date (YYYY-MM-DD)")
	flag.StringVar(&toString, "to", "2019-03-20", "To date (YYYY-MM-DD)")
	flag.BoolVar(&dryRun, "dry-run", true, "Dry Run (no changes done)")
	flag.BoolVar(&verbose, "verbose", false, "Verbose logs")
	flag.BoolVar(&summary, "summary", true, "Print Summary")
	flag.Parse()

	if verbose {
		logger.SetLevel(logger.TraceLevel)
	}

	startTime, err := time.Parse(dateFormatLayout, fromString)
	if err != nil {
		logger.WithFields(logger.Fields{
			"from": fromString,
		}).Fatal("From date is not valid")
	}
	endTime, err := time.Parse(dateFormatLayout, toString)
	if err != nil {
		logger.WithFields(logger.Fields{
			"to": toString,
		}).Fatal("To date is not valid")
	}

	logger.WithFields(logger.Fields{
		"from":   fromString,
		"to":     toString,
		"dryRun": dryRun,
	}).Info("Running YouDocktor")

	timeDoctorAPI := NewTimeDoctorAPI(conf)
	youTrackAPI := NewYouTrackAPI(conf)

	sheet := populateSheet(startTime, endTime, timeDoctorAPI)
	updateFromSheet(sheet, youTrackAPI)

	if summary {
		sheet.printSummary()
	}
}

func populateSheet(start time.Time, end time.Time, timeDoctorAPI *TimeDoctorAPI) (sheet *timeSheet) {

	sheet = &timeSheet{
		From:    start,
		To:      end,
		Entries: make([]timeSheetEntry, 0),
	}
	for rd := rangeDate(start, end); ; {
		date := rd()
		if date.IsZero() {
			break
		}

		wl, err := timeDoctorAPI.GetWorkLog(date, date)
		if err != nil {
			log.Fatal(err)
			return
		}

		sheet.Entries = populateEntries(sheet.Entries, date, wl.WorkLogs.Items)
	}
	return
}

func populateEntries(input []timeSheetEntry, date time.Time, items []WorkLogItem) (output []timeSheetEntry) {
	output = input
	for _, item := range items {
		entry := NewTimeSheetEntry(date, item)
		output = append(output, entry)
	}
	return
}

func updateFromSheet(sheet *timeSheet, youTrackAPI *YouTrackAPI) {
	for index, item := range sheet.Entries {

		if item.YouTrackIssueID == "" {
			logger.WithFields(logger.Fields{
				"taskName": item.TaskName,
			}).Info("Not a valid YouTrack entry, it will be ignored")
			continue
		}

		// to modify a item we need to access it via index
		err := youTrackAPI.CreateOrUpdate(&sheet.Entries[index], dryRun)
		if err != nil {
			logger.WithFields(logger.Fields{
				"issueID": item.YouTrackIssueID,
			}).Warnf("Error with WorkItem: %s", err)
		}
	}
}
