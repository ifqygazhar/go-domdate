package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type Service1 struct {
	URL   string
	Agent string
}

func NewService1() *Service1 {
	return &Service1{
		URL:   "https://www.cubdomain.com/domains-registered-by-date/",
		Agent: "GoogleBot v3",
	}
}

func (s *Service1) CountPages(date string) int {
	c := colly.NewCollector()

	var lastPage int

	c.OnHTML("a.page-link", func(e *colly.HTMLElement) {
		pageText := e.Text
		page, err := strconv.Atoi(pageText)
		if err == nil && page > lastPage {
			lastPage = page
		}
	})

	c.Visit(s.URL + date + "/1")

	return lastPage
}

func (s *Service1) Dump(date, page string, exts []string) []string {
	c := colly.NewCollector()

	var data []string

	c.OnHTML("div.col-md-4", func(e *colly.HTMLElement) {
		siteText := e.Text
		if len(exts) == 0 {
			data = append(data, strings.ReplaceAll(siteText, "\n", ""))
		} else {
			for _, ext := range exts {
				if match, _ := regexp.MatchString("("+ext+")$", strings.ReplaceAll(siteText, "\n", "")); match {
					data = append(data, strings.ReplaceAll(siteText, "\n", ""))
					break
				}
			}
		}
	})

	c.Visit(s.URL + date + "/" + page)

	return data
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("There are 2 servers we provided, which every server providing their own data.")
	fmt.Println("Server 1, serve domains by date registration.")
	fmt.Println("Server 2, serve domains by their extensions. So, choose it first.")
	fmt.Print("⚙ Then your choice is: ")
	serviceChoice, _ := reader.ReadString('\n')
	serviceChoice = strings.TrimSpace(serviceChoice)

	if serviceChoice == "1" {
		fmt.Println("⚙ Start using server 1")
		fmt.Println("⚙ info")
		fmt.Println(`"""
    This Server is serve domains by dates, you must input valid date formats which are correct by this bot.
    Valid Formats:
      @ 2023-08-01

"""`)
		fmt.Print("⚙ From date (YYYY-MM-DD / TAHUN-BULAN-HARI): ")
		fromDate, _ := reader.ReadString('\n')
		fromDate = strings.TrimSpace(fromDate)

		fmt.Print("⚙ To date (YYYY-MM-DD / TAHUN-BULAN-HARI): ")
		toDate, _ := reader.ReadString('\n')
		toDate = strings.TrimSpace(toDate)

		service1 := NewService1()

		currentDate, _ := time.Parse("2006-01-02", fromDate)
		endDate, _ := time.Parse("2006-01-02", toDate)
		var allData []string

		for currentDate.Before(endDate) || currentDate.Equal(endDate) {
			validDate := currentDate.Format("2006-01-02")
			totalPages := service1.CountPages(validDate)

			fmt.Printf("⚙ Date: %s | Total pages: %d\n", validDate, totalPages)

			for page := 1; page <= totalPages; page++ {
				data := service1.Dump(validDate, strconv.Itoa(page), nil) // You can specify extensions here
				allData = append(allData, data...)
			}

			currentDate = currentDate.AddDate(0, 0, 1)
		}

		fmt.Println("⚙ Dumping all data to grablist.txt..")
		if err := saveDataToFile("grablist.txt", allData); err != nil {
			fmt.Println("⚙ Error saving data:", err)
		} else {
			fmt.Println("⚙ Data saved to grablist.txt")
		}

	} else if serviceChoice == "2" {
		// Implement Server 2 logic here
	} else {
		fmt.Println("Invalid service choice.")
	}

	fmt.Println("Done... Click Enter to exit!")
	fmt.Scanln()
}

func saveDataToFile(filename string, data []string) error {
	content := strings.Join(data, "\n")
	return ioutil.WriteFile(filename, []byte(content), 0644)
}
