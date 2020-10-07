package main

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/openshift/gmarkley-VI/jiraSosRepot/functions"
	"log"
)

func main() {
	//Setup
	totalPoints := 0.000000
	jiraURL := "https://issues.redhat.com"
	username, password := functions.ReadCredentials()
	var jiraJQL [1]string
	jiraJQL[0] = "project = WINC AND (type = Bug and status = Done AND \"Story Points\" is not null )OR project = OCPBUGSM AND (component in (\"Windows Containers\") AND \"Story Points\" is not EMPTY )"

	//Create the client
	client, _ := functions.CreatTheClient(username, password, jiraURL)

	//Loop over the jiraJQL array and Request the issue objects
	for z := 0; z < len(jiraJQL); z++ {

		var issues []jira.Issue

		// append the jira issues to []jira.Issue
		appendFunc := func(i jira.Issue) (err error) {
			issues = append(issues, i)
			return err
		}

		// In this example, we'll search for all the issues with the provided JQL filter and Print the Story Points
		err := client.Issue.SearchPages(fmt.Sprintf(`%s`, jiraJQL[z][0]), nil, appendFunc)
		if err != nil {
			log.Fatal(err)
		}

		for _, i := range issues {
			storyPoint := i.Fields.Unknowns["customfield_12310243"]
			totalPoints += storyPoint.(float64)
		}

		fmt.Printf("Number of Closed Bugs - %d \n", len(issues))
		fmt.Printf("Total of Points - %d \n", int(totalPoints))
		avaragePoints := totalPoints / float64(len(issues))
		fmt.Printf("Avarage Points per Bug - %f \n", avaragePoints)
	}
}
