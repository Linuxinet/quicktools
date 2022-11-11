package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/tidwall/gjson"
)

func main() {

	c_name := make([]string, 0, 500)
	c_start_date := make([]string, 0, 500)
	c_code := make([]string, 0, 500)

	//////future
	f_c_name := make([]string, 0, 500)
	f_c_start_date := make([]string, 0, 500)
	f_c_end_date := make([]string, 0, 500)
	f_c_duration := make([]string, 0, 500)
	f_c_code := make([]string, 0, 500)

	URL := "https://www.codechef.com/api/list/contests/all?sort_by=START&sorting_order=asc&offset=0&mode=premium"

	resp, err := http.Get(URL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	apibody, _ := io.ReadAll(resp.Body)

	con_name := gjson.Get(string(apibody), "present_contests.#.contest_name")
	con_code := gjson.Get(string(apibody), "present_contests.#.contest_code")
	con_start_date := gjson.Get(string(apibody), "present_contests.#.contest_start_date")

	////////////////////////////////////future
	f_con_name := gjson.Get(string(apibody), "future_contests.#.contest_name")
	f_con_code := gjson.Get(string(apibody), "future_contests.#.contest_code")
	f_con_start_date := gjson.Get(string(apibody), "future_contests.#.contest_start_date")
	f_con_end_date := gjson.Get(string(apibody), "future_contests.#.contest_end_date")
	f_con_duration := gjson.Get(string(apibody), "future_contests.#.contest_duration")

	for i := range con_name.Array() {
		c_name = append(c_name, con_name.Array()[i].Str)
		c_code = append(c_code, con_code.Array()[i].Str)
		c_start_date = append(c_start_date, con_start_date.Array()[i].Str)
	}

	for i := range f_con_name.Array() {
		f_c_name = append(f_c_name, f_con_name.Array()[i].Str)
		f_c_code = append(f_c_code, f_con_code.Array()[i].Str)
		f_c_start_date = append(f_c_start_date, f_con_start_date.Array()[i].Str)
		f_c_end_date = append(f_c_end_date, f_con_end_date.Array()[i].Str)
		f_c_duration = append(f_c_duration, f_con_duration.Array()[i].Str)
	}

	fmt.Println("##########  PRESENT CONTESTS  ##########")
	fmt.Println("                    ")
	fmt.Println("                    ")
	for a := range c_name {
		fmt.Printf("##%d\n\nContest Name: %s\n\nContest Code: %s\n\nStart Time: %s\n\n\n", a+1, c_name[a], c_code[a], c_start_date[a])
	}

	//future
	fmt.Println("##########  FUTURE CONTESTS  ##########")
	fmt.Println("                    ")
	fmt.Println("                    ")

	for a := range f_c_name {
		fmt.Printf("##%d\n\nContest Name: %s\n\nContest Code: %s\n\nStart Time: %s\n\nEnd Time: %s\n\nDuration: %s\n\n\n", a+1, f_c_name[a], f_c_code[a], f_c_start_date[a], f_c_end_date[a], f_c_duration[a])
	}

}
