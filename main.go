package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/TylerBrock/colorjson"
)

type GetContestsResponse struct {
	Status           string             `json:"status"`
	Message          string             `json:"message"`
	PresentContests  []PresentContests  `json:"present_contests"`
	FutureContests   []FutureContests   `json:"future_contests"`
	PracticeContests []PracticeContests `json:"practice_contests"`
	PastContests     []PastContests     `json:"past_contests"`
}
type PresentContests struct {
	ContestCode         string     `json:"contest_code"`
	ContestName         string     `json:"contest_name"`
	ContestStartDate    string     `json:"contest_start_date"`
	ContestStartDateIso CustomTime `json:"contest_start_date_iso"`
}
type FutureContests struct {
	ContestCode         string     `json:"contest_code,omitempty"`
	ContestName         string     `json:"contest_name,omitempty"`
	ContestStartDate    string     `json:"contest_start_date,omitempty"`
	ContestEndDate      string     `json:"contest_end_date,omitempty"`
	ContestStartDateIso CustomTime `json:"contest_start_date_iso,omitempty"`
	ContestEndDateIso   CustomTime `json:"contest_end_date_iso,omitempty"`
	ContestDuration     string     `json:"contest_duration,omitempty"`
	DistinctUsers       int        `json:"distinct_users,omitempty"`
}
type PracticeContests struct {
	ContestCode         string     `json:"contest_code"`
	ContestName         string     `json:"contest_name"`
	ContestStartDate    string     `json:"contest_start_date"`
	ContestEndDate      string     `json:"contest_end_date"`
	ContestStartDateIso CustomTime `json:"contest_start_date_iso"`
	ContestEndDateIso   CustomTime `json:"contest_end_date_iso"`
	ContestDuration     string     `json:"contest_duration"`
	DistinctUsers       int        `json:"distinct_users"`
}
type PastContests struct {
	ContestCode         string     `json:"contest_code"`
	ContestName         string     `json:"contest_name"`
	ContestStartDate    string     `json:"contest_start_date"`
	ContestEndDate      string     `json:"contest_end_date"`
	ContestStartDateIso CustomTime `json:"contest_start_date_iso"`
	ContestEndDateIso   CustomTime `json:"contest_end_date_iso"`
	ContestDuration     string     `json:"contest_duration"`
	DistinctUsers       int        `json:"distinct_users"`
}

type CustomTime struct {
	time.Time
}

const TimeFormat = time.RFC3339

func (ct *CustomTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		ct.Time = time.Time{}
		return
	}
	ct.Time, err = time.Parse(TimeFormat, s)
	if err != nil {
		panic(err)
	}
	return
}

func main() {

	Contests := make(map[string]map[string]map[string]map[string]string)
	Contests["CodeChef"] = make(map[string]map[string]map[string]string)
	Contests["CodeChef"]["PresentContests"] = make(map[string]map[string]string)
	Contests["CodeChef"]["FutureContests"] = make(map[string]map[string]string)

	URL := "https://www.codechef.com/api/list/contests/all?sort_by=START&sorting_order=asc&offset=0&mode=premium"

	resp, err := http.Get(URL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	apibody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var jsonbody GetContestsResponse
	if err := json.Unmarshal(apibody, &jsonbody); err != nil {
		log.Println(err)
	}

	// fmt.Println(len(jsonbody.FutureContests))
	for id, p := range jsonbody.FutureContests {

		// Contests["CodeChef"]["FutureContests"][string(rune(id+1))] = make(map[string]string)

		// Contests["CodeChef"]["FutureContests"][string(rune(id+1))]["Name"] = p.ContestName
		// Contests["CodeChef"]["FutureContests"][string(rune(id+1))]["Code"] = p.ContestCode
		// Contests["CodeChef"]["FutureContests"][string(rune(id+1))]["Start"] = p.ContestStartDate
		// Contests["CodeChef"]["FutureContests"][string(rune(id+1))]["End"] = p.ContestEndDate
		// Contests["CodeChef"]["FutureContests"][string(rune(id+1))]["Duration"] = p.ContestDuration

		Contests["CodeChef"]["FutureContests"][string(strconv.Itoa(id+1))] = map[string]string{

			"Name":     p.ContestName,
			"Code":     p.ContestCode,
			"Start":    p.ContestStartDate,
			"End":      p.ContestEndDate,
			"Duration": p.ContestDuration,
		}
	}

	for id, p := range jsonbody.PresentContests {

		Contests["CodeChef"]["PresentContests"][string(strconv.Itoa(id+1))] = map[string]string{
			"Name":  p.ContestName,
			"Code":  p.ContestCode,
			"Start": p.ContestStartDate,
		}

		// Contests["CodeChef"]["PresentContests"][string(rune(id+1))]["Name"] = p.ContestName
		// Contests["CodeChef"]["PresentContests"][string(rune(id+1))]["Code"] = p.ContestCode
		// Contests["CodeChef"]["PresentContests"][string(rune(id+1))]["Start"] = p.ContestStartDate
	}

	// fmt.Println(Contests)

	jsonStr, err := json.Marshal(Contests)
	if err != nil {
		panic(err)
	}

	var obj map[string]interface{}
	json.Unmarshal([]byte(jsonStr), &obj)

	// Make a custom formatter with indent set
	f := colorjson.NewFormatter()
	f.Indent = 4

	// Marshall the Colorized JSON
	s, _ := f.Marshal(obj)
	fmt.Println(string(s))

}
