package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
)

func main() {

	URL := "https://codeforces.com/contests"

	resp, err := http.Get(URL)
	if err != nil {
		panic(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	re := regexp.MustCompile(`<td class="dark left">(.*?)</td>`)
	// fmt.Println("Program : ", PRGM)
	var days string
	submatchall := re.FindAllStringSubmatch(string(body), 2)
	for _, element := range submatchall {
		// fmt.Println(element[1])
		days = element[1]
		fmt.Println(days)
	}
}
