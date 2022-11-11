package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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
	
	var id int

	Contests := make(map[string]map[string]map[string]map[string]string)
	Contests["CodeChef"] = make(map[string]map[string]map[string]string)
	Contests["CodeChef"]["PresentContests"] = make(map[string]map[string]string)
	Contests["CodeChef"]["FutureContests"] = make(map[string]map[string]string)
	Contests["CodeChef"]["FutureContests"][string(rune(id))] = make(map[string]string)
	Contests["CodeChef"]["PresentContests"][string(rune(id))] = make(map[string]string)

	// c_name := make([]string, 0, 500)
	// c_start_date := make([]string, 0, 500)
	// c_code := make([]string, 0, 500)

	// //////future
	// f_c_name := make([]string, 0, 500)
	// f_c_start_date := make([]string, 0, 500)
	// f_c_end_date := make([]string, 0, 500)
	// f_c_duration := make([]string, 0, 500)
	// f_c_code := make([]string, 0, 500)

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

	fmt.Println(len(jsonbody.FutureContests))
	for id, p := range jsonbody.FutureContests {

		Contests["CodeChef"]["FutureContests"][string(rune(id+1))]["Name"] = p.ContestName
		Contests["CodeChef"]["FutureContests"][string(rune(id+1))]["Code"] = p.ContestCode
		Contests["CodeChef"]["FutureContests"][string(rune(id+1))]["Start"] = p.ContestStartDate
		Contests["CodeChef"]["FutureContests"][string(rune(id+1))]["End"] = p.ContestEndDate
		Contests["CodeChef"]["FutureContests"][string(rune(id+1))]["Duration"] = p.ContestDuration
	}

	for id, p := range jsonbody.PresentContests {
		// c_name = append(c_name, p.ContestName)
		// c_code = append(c_code, p.ContestCode)
		// c_start_date = append(c_start_date, p.ContestStartDate)

		Contests["CodeChef"]["PresentContests"][string(rune(id+1))]["Name"] = p.ContestName
		Contests["CodeChef"]["PresentContests"][string(rune(id+1))]["Code"] = p.ContestCode
		Contests["CodeChef"]["PresentContests"][string(rune(id+1))]["Start"] = p.ContestStartDate
	}

	// fmt.Println(contests)

	// fmt.Println(len(f_c_name))
	// fmt.Println(len(c_name))
	// for k := range f_c_name {
	// 	Contests["CodeChef"]["FutureContests"]["Name"] = f_c_name[k]
	// 	Contests["CodeChef"]["FutureContests"]["Code"] = f_c_code[k]
	// 	Contests["CodeChef"]["FutureContests"]["Start"] = f_c_start_date[k]
	// 	Contests["CodeChef"]["FutureContests"]["End"] = f_c_end_date[k]
	// 	Contests["CodeChef"]["FutureContests"]["Duration"] = f_c_duration[k]
	// }

	// for k := range c_name {
	// 	Contests["CodeChef"]["PresentContests"]["Name"] = c_name[k]
	// 	Contests["CodeChef"]["PresentContests"]["Code"] = c_code[k]
	// 	Contests["CodeChef"]["PresentContests"]["Start"] = c_start_date[k]
	// }

	fmt.Println(Contests)

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
