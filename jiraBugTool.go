package main

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/openshift/gmarkley-VI/jiraBugtool/functions"
	"log"
	"strings"
)

func exportConsole(key string, output string) {
	fmt.Printf("%s - %s\n", key, output)
}
func exportJira(client *jira.Client, id string, key string, output string) *jira.Comment {
	com := jira.Comment{
		ID:           id,
		Self:         "",
		Name:         "",
		Author:       jira.User{},
		Body:         output,
		UpdateAuthor: jira.User{},
		Updated:      "",
		Created:      "",
		Visibility:   jira.CommentVisibility{},
	}
	commentOUT, _, err := client.Issue.AddComment(key, &com)
	if err != nil {
		panic(err)
	}
	return commentOUT
}

func main() {
	//Setup
	jiraURL := "https://issues.redhat.com"
	username, password := functions.ReadCredentials()
	var jiraJQL [1][2]string
	jiraJQL[0][0] = "project = WINC AND type = Bug and status = Done AND \"Story Points\" is not null"

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
			options := &jira.GetQueryOptions{Expand: "renderedFields"}
			u, _, err := client.Issue.Get(i.Key, options)
			if err != nil {
				fmt.Printf("\n==> error: %v\n", err)
				return
			}

			if len(u.RenderedFields.Comments.Comments) >= 1 {
				c := u.RenderedFields.Comments.Comments[len(u.RenderedFields.Comments.Comments)-1]
				if strings.Contains(c.Updated, "days ago") {
					commentString := fmt.Sprintf("%s Please comment/update - Last update was %+v", name, c.Updated)
					exportConsole(i.Key, commentString)
				}
			} else {
				commentString := fmt.Sprintf("%s Please add a comment.", name)
				exportConsole(i.Key, commentString)
			}
		}
	}
}
