package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	logger "github.com/sirupsen/logrus"
)

type workItem struct {
	ID       string `json:"id"`
	Date     int64  `json:"date"`
	Text     string `json:"text"`
	Duration struct {
		Minutes int `json:"minutes"`
	} `json:"duration"`
}

// YouTrackAPI defines a struct to access the Youtrack methods
type YouTrackAPI struct {
	config *config
}

// NewYouTrackAPI creates a new YouTrackAPI with the associated configuration
func NewYouTrackAPI(c *config) *YouTrackAPI {
	return &YouTrackAPI{
		config: c,
	}
}

// CreateOrUpdate Creates or updatesa a workitem entry on Youtrack
func (api *YouTrackAPI) CreateOrUpdate(entry *timeSheetEntry, dryRun bool) (err error) {

	textIdentifier := entry.Identifier()
	contextLogger := logger.WithFields(logger.Fields{
		"issueID":  entry.YouTrackIssueID,
		"date":     entry.Date.Format("2006-01-02"),
		"spent":    entry.SpentTime,
		"taskName": entry.TaskName,
		"dryRun":   dryRun,
	})

	// get the list of workitems of the issue check on the list if exists one with the text starting as XXXXX
	items, err := api.getWorkItems(entry.YouTrackIssueID)
	if err != nil {
		contextLogger.Fatal("Could not retrieve WorkItems")
	}

	var currentWorkItem workItem
	for _, item := range items {
		if strings.HasPrefix(item.Text, textIdentifier) {
			currentWorkItem = item
			break
		}
	}

	minutes := int(entry.SpentTime.Minutes())
	if currentWorkItem.Duration.Minutes == minutes {
		entry.SetNotModified()
		contextLogger.Info("No time variation, nothing to update")
		return
	}

	toBeReplaced := "[" + entry.YouTrackIssueID + "] "
	taskName := strings.Replace(entry.TaskName, toBeReplaced, "", 1)

	if currentWorkItem.ID == "" {
		contextLogger.Info("New WorkItem will be created")
		entry.SetNew()
	} else {
		contextLogger.Info("Already existin WorkItem will be updated")
		entry.SetModified()
	}

	if !dryRun {
		err = api.createOrUpdateWorkItem(entry.YouTrackIssueID, currentWorkItem.ID, minutes, textIdentifier+taskName, entry.Date)
	}
	return
}

// getWorkLog gets the list of worklogitems between two dares
func (api *YouTrackAPI) getWorkItems(issueID string) (response []workItem, err error) {

	url := "https://youtrack.lodgify.net/api/issues/" + issueID + "/timeTracking/workItems?fields=id,date,text,duration%28minutes%29"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return
	}

	client := &http.Client{}
	req.Header.Add("Authorization", "Bearer "+api.config.YouTrack.Token)

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Println(err)
	}
	return
}

// createOrUpdateWorkItem creates or updates a new track time entry for the youtrac card
func (api *YouTrackAPI) createOrUpdateWorkItem(issueID string, workItemID string, minutes int, text string, date time.Time) (err error) {
	url := "https://youtrack.lodgify.net/api/issues/" + issueID + "/timeTracking/workItems"

	// append the ID for updating the entry
	if workItemID != "" {
		url = url + "/" + workItemID
	}

	timestamp := makeTimestamp(date)
	jsonStr := fmt.Sprintf("{\"date\":%d,\"duration\":{\"minutes\":%d}, \"text\":\"%s\"}", timestamp, minutes, text)
	json := []byte(jsonStr)

	req, err := http.NewRequest("POST", url, bytes.NewReader(json))
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return
	}

	client := &http.Client{}
	req.Header.Add("Authorization", "Bearer "+api.config.YouTrack.Token)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)

	if res.StatusCode != 200 {
		b, _ := ioutil.ReadAll(res.Body)
		log.Fatal(string(b))
	}

	return
}
