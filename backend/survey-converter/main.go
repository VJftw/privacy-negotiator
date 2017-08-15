package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// SELECT surveys.user_id, users.long_lived_token, surveys.photo_id, surveys.raw_json FROM surveys INNER JOIN users on users.id = surveys.user_id

func main() {
	file, err := os.Open("surveys.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	generalSheet := [][]string{[]string{"FacebookID", "Gender", "AgeRange", "Q1", "Q2", "Q4", "Q5", "Q7", "Q8"}}
	generalQ3 := [][]string{[]string{"FacebookID", "untag", "crop", "remove", "blur", "nothing", "other"}}
	generalQ6 := [][]string{[]string{"FacebookID", "politics", "religion", "work", "sports", "family", "location", "education", "favourite-teams", "inspirational-people", "languages", "music", "movies", "likes", "groups", "events"}}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		r := csv.NewReader(strings.NewReader(line))

		record, err := r.Read()
		if err != nil {
			panic(err)
		}
		// user_id, long_lived_access_token, photo_id, raw_json
		userID := record[0]
		accessToken := record[1]
		photoID := record[2]
		rawJSON := record[3]

		if photoID == "general" {
			fmt.Print(userID + "...")
			var responses generalJSON
			json.Unmarshal([]byte(rawJSON), &responses)
			generalSheet = append(generalSheet, processGeneral(userID, accessToken, responses))
			generalQ3 = append(generalQ3, processGeneralQ3(userID, responses))
			generalQ6 = append(generalQ6, processGeneralQ6(userID, responses, generalQ6[0]))
			fmt.Println("Done.")
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Processed %d surveys.\n", len(generalSheet)-1)

	writeCSV("general.csv", generalSheet)
	writeCSV("generalQ3.csv", generalQ3)
	writeCSV("generalQ6.csv", generalQ6)

}

func writeCSV(filename string, records [][]string) {
	file, _ := os.Create(filename)
	defer file.Close()
	w := csv.NewWriter(file)
	w.WriteAll(records)
	w.Flush()
	fmt.Printf("Wrote %s.\n", filename)
}

func processGeneralQ6(userID string, responses generalJSON, headingOrder []string) []string {

	var results [15]string
	for _, part := range responses.Q6Parts {
		for i, headerID := range headingOrder {
			if part.ID == headerID {
				results[i] = strconv.Itoa(part.Value)
				break
			}
		}
	}

	return append([]string{userID}, results[:]...)
}

func processGeneralQ3(userID string, responses generalJSON) []string {
	untag := "no"
	crop := "no"
	remove := "no"
	blur := "no"
	nothing := "no"
	other := "no"

	for _, choice := range responses.Q3Checkboxes {
		if choice.Selected {
			switch choice.ID {
			case "untag":
				untag = "yes"
				break
			case "crop":
				crop = "yes"
				break
			case "remove":
				remove = "yes"
				break
			case "blur":
				blur = "yes"
				break
			case "nothing":
				nothing = "yes"
				break
			case "other":
				other = responses.Q3Other
				break
			}
		}
	}

	return []string{
		userID,
		untag,
		crop,
		remove,
		blur,
		nothing,
		other,
	}
}

func processGeneral(userID string, accessToken string, responses generalJSON) []string {
	var q1 string
	var q2 string
	var q4 string
	var q5 string
	url := fmt.Sprintf(
		"https://graph.facebook.com/v2.9/%s?fields=gender,age_range&access_token=%s",
		userID,
		accessToken,
	)
	res, err := http.Get(url)

	if err != nil {
		panic(err)
	}

	resMe := &responseMe{}
	err = json.NewDecoder(res.Body).Decode(resMe)
	defer res.Body.Close()
	if err != nil {
		log.Printf("Error: %s", err)
	}

	// q1
	for _, choice := range responses.Q1Choices {
		if responses.Q1Answer == choice.Description {
			q1 = choice.ID
			break
		}
	}

	// q2
	for _, choice := range responses.Q2Choices {
		if responses.Q2Answer == choice.Description {
			q2 = choice.ID
			break
		}
	}

	// q4
	if responses.Q4Answer == "No" {
		q4 = "no - " + responses.Q4Why
	} else {
		q4 = "yes"
	}

	if responses.Q5Answer == "No" {
		q5 = "no - " + responses.Q5Why
	} else {
		q5 = "yes"
	}

	return []string{
		userID,
		resMe.Gender,
		strconv.Itoa(resMe.AgeRange.Min),
		q1,
		q2,
		q4,
		q5,
		responses.Q7Answer,
		responses.Q8Answer,
	}
}

type responseMe struct {
	Gender   string   `json:"gender"`
	AgeRange ageRange `json:"age_range"`
}

type ageRange struct {
	Min int `json:"min"`
}

type generalJSON struct {
	Q1Answer     string     `json:"q1Answer"`
	Q1Choices    []choice   `json:"q1Choices"`
	Q2Answer     string     `json:"q2Answer"`
	Q2Choices    []choice   `json:"q2Choices"`
	Q3Answer     []string   `json:"q3Answer"`
	Q3Checkboxes []checkbox `json:"q3Checkboxes"`
	Q3Other      string     `json:"q3Other"`
	Q4Answer     string     `json:"q4Answer"`
	Q4Why        string     `json:"q4Why"`
	Q5Answer     string     `json:"q5Answer"`
	Q5Why        string     `json:"q5Why"`
	Q6Parts      []weight   `json:"q6Parts"`
	Q7Answer     string     `json:"q7Answer"`
	Q8Answer     string     `json:"q8Answer"`
}

type choice struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

type checkbox struct {
	choice
	Selected bool `json:"selected"`
}

type weight struct {
	choice
	Value int `json:"value"`
}
