package main

import (
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

type status int

const (
	New status = iota
	Modified
	NotModified
	Unknown
)

func (s status) String() string {
	switch s {
	case New:
		return "+"
	case Modified:
		return ">"
	case NotModified:
		return "="
	default:
		return "?"
	}
}

// TimeSheet defines a list of entries between two dates
type timeSheet struct {
	From    time.Time
	To      time.Time
	Entries []timeSheetEntry
}

// TimeSheetEntry defines a task for a given date and the time spent
type timeSheetEntry struct {
	TaskName         string
	TimeDoctorTaskID string
	YouTrackIssueID  string
	SpentTime        time.Duration
	Date             time.Time
	Status           status
}

// Identifier returns a string used to identify the entry
func (t timeSheetEntry) Identifier() (output string) {
	output = fmt.Sprintf("(youdocktor:%s:%s)", t.TimeDoctorTaskID, t.Date.Format("2006-01-02"))
	return
}

func (t *timeSheetEntry) SetNew() {
	t.Status = New
}
func (t *timeSheetEntry) SetModified() {
	t.Status = Modified
}

func (t *timeSheetEntry) SetNotModified() {
	t.Status = NotModified
}

func (t *timeSheetEntry) SetUnknown() {
	t.Status = Unknown
}

// NewTimeSheetEntry Creates a new TimeSheetEntry from the input data
func NewTimeSheetEntry(date time.Time, item WorkLogItem) timeSheetEntry {
	regex := regexp.MustCompile(`^\[(?P<ID>\w{1,6}\-\d{1,6})\]`)
	r, _ := getRegexGroups(item.TaskName, regex)
	issueID := r["ID"]
	duration, _ := time.ParseDuration(item.Length + "s")

	entry := timeSheetEntry{
		TaskName:         item.TaskName,
		TimeDoctorTaskID: item.ID,
		SpentTime:        duration,
		YouTrackIssueID:  issueID,
		Date:             date,
	}
	entry.SetUnknown()
	return entry
}

func (s *timeSheet) printSummary() {

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Date", "Status", "Identifier", "Name", "Spent"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})

	for _, item := range s.Entries {
		table.Append([]string{item.Date.Format("2006-01-02"), item.Status.String(), item.YouTrackIssueID, item.TaskName, item.SpentTime.String()})
	}

	table.Render()
}

// Create private data struct to hold config options.
type config struct {
	TimeDoctor struct {
		ClientID     string `mapstructure:"clientID"`
		ClientSecret string `mapstructure:"clientSecret"`
		RedirectURL  string `mapstructure:"redirectURL"`
		AccessToken  string `mapstructure:"accessToken"`
		RefreshToken string `mapstructure:"refreshToken"`
		Expiry       string `mapstructure:"expiry"`
		CompanyID    int    `mapstructure:"companyID"`
	} `mapstructure:"timeDoctor"`
	YouTrack struct {
		Token string `mapstructure:"token"`
	} `mapstructure:"youTrack"`
}

func (c *config) ExpiryDate() (time.Time, error) {
	return time.Parse(time.RFC3339, c.TimeDoctor.Expiry)
}

func (c *config) PersistsTimeDoctorToken(token oauth2.Token) {
	c.TimeDoctor.AccessToken = token.AccessToken
	c.TimeDoctor.RefreshToken = token.RefreshToken
	c.TimeDoctor.Expiry = token.Expiry.Format(time.RFC3339)
	viper.Set("timeDoctor.accessToken", c.TimeDoctor.AccessToken)
	viper.Set("timeDoctor.refreshToken", c.TimeDoctor.RefreshToken)
	viper.Set("timeDoctor.expiry", c.TimeDoctor.Expiry)
	viper.WriteConfig()
}
