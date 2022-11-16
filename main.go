package main

import (
	"encoding/json"
	"fmt"

	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/TylerBrock/colorjson"
	"github.com/gocolly/colly"
)

type GetContestsResponse struct {
	PresentContests []PresentContests `json:"present_contests"`
	FutureContests  []FutureContests  `json:"future_contests"`
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
}

type CustomTime struct {
	time.Time
}

type AtcoderName struct {
	Name string
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

	Contests["AtCoder"] = make(map[string]map[string]map[string]string)
	Contests["AtCoder"]["FutureContests"] = make(map[string]map[string]string)

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

	for id, p := range jsonbody.FutureContests {
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
	}

	AtcoderFunc(Contests)
	// fmt.Println(Contests)

	jsonStr, err := json.Marshal(Contests)
	if err != nil {
		panic(err)
	}

	var obj map[string]interface{}
	json.Unmarshal([]byte(jsonStr), &obj)

	writeJSON(obj)

	// Make a custom formatter with indent set
	f := colorjson.NewFormatter()
	f.Indent = 4

}

func writeJSON(data map[string]interface{}) {
	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println("Unable to create json file")
		return
	}

	_ = ioutil.WriteFile("contests.json", file, 0644)
	fmt.Println("Wrote Contest data to contests.json")
}

func AtcoderFunc(Contests map[string]map[string]map[string]map[string]string) {
	collector := colly.NewCollector(
		colly.AllowedDomains("atcoder.jp", "www.atcoder.jp"),
	)

	format := "2006-01-02 15:04:05-0700"
	loc, _ := time.LoadLocation("Asia/Calcutta")

	for i := 1; i < 10; i++ {

		rawI := strconv.Itoa(i)
		Contests["AtCoder"]["FutureContests"][rawI] = make(map[string]string)

		ContestSelTime := fmt.Sprintf("#contest-table-upcoming  div  div  table  tbody  tr:nth-child(%d)  td:nth-child(1)  a", i)
		ContestSelName := fmt.Sprintf("#contest-table-upcoming  div  div  table  tbody  tr:nth-child(%d)  td:nth-child(2)", i)

		// for contest name
		collector.OnHTML(ContestSelName, func(element *colly.HTMLElement) {
			ContestName := element.ChildText("a")
			Contests["AtCoder"]["FutureContests"][rawI]["Name"] = ContestName

		})

		// for contestTime
		collector.OnHTML(ContestSelTime, func(element *colly.HTMLElement) {
			ContestStartTime := element.ChildText("time")
			parsed_time, _ := time.Parse(format, ContestStartTime)
			IST_time := parsed_time.In(loc)
			Contests["AtCoder"]["FutureContests"][rawI]["Start"] = fmt.Sprint(IST_time)
		})

	}

	collector.OnRequest(func(request *colly.Request) {
		fmt.Println("Visiting : ", request.URL.String())
	})

	collector.Visit("https://atcoder.jp/contests")

}
