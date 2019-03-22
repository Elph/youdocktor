package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

type WorkLogResponse struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Total     int    `json:"total"`
	WorkLogs  struct {
		Count  int           `json:"count"`
		Offset int           `json:"offse"`
		Limit  int           `json:"limit"`
		Items  []WorkLogItem `json:"items"`
	} `json:"worklogs"`
}

type WorkLogItem struct {
	ID       string `json:"id"`
	Length   string `json:"length"`
	TaskID   string `json:"task_id"`
	TaskName string `json:"task_name"`
}

type TimeDoctorAPI struct {
	config *config
}

func NewTimeDoctorAPI(c *config) *TimeDoctorAPI {
	return &TimeDoctorAPI{
		config: c,
	}
}

// GetWorkLog gets the list of worklogitems between two dares
func (api *TimeDoctorAPI) GetWorkLog(start time.Time, end time.Time) (response WorkLogResponse, err error) {

	oauthConf := oauth2.Config{
		ClientID:     api.config.TimeDoctor.ClientID,
		ClientSecret: api.config.TimeDoctor.ClientSecret,
		Scopes:       []string{"refresh_token"},
		RedirectURL:  api.config.TimeDoctor.RedirectURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://webapi.timedoctor.com/oauth/v2/auth",
			TokenURL: "https://webapi.timedoctor.com/oauth/v2/token",
		},
	}

	t := new(oauth2.Token)
	t.TokenType = "Bearer"
	t.AccessToken = api.config.TimeDoctor.AccessToken
	t.RefreshToken = api.config.TimeDoctor.RefreshToken
	t.Expiry, _ = api.config.ExpiryDate()
	tokenSource := oauthConf.TokenSource(oauth2.NoContext, t)

	// internally will refresh if expired
	validToken, err := tokenSource.Token()
	client := oauthConf.Client(context.Background(), t)

	// if token has changed, persist it on the struct so it can be saved later
	if validToken.AccessToken != api.config.TimeDoctor.AccessToken {
		api.config.PersistsTimeDoctorToken(*validToken)
	}

	startDate := start.Format("2006-01-02")
	endDate := end.Format("2006-01-02")
	url := fmt.Sprintf("https://webapi.timedoctor.com/v1.1/companies/%d/worklogs?_format=json&start_date=%s&end_date=%s&consolidated=1", api.config.TimeDoctor.CompanyID, startDate, endDate)

	response = WorkLogResponse{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		b, _ := ioutil.ReadAll(resp.Body)
		return response, errors.New(string(b))
	}

	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Println(err)
	}
	return response, nil
}
