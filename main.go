package jiraclient

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/andygrunwald/go-jira"
	"github.com/spf13/viper"
)

func GetJiraContent(issueId string) (*jira.Issue, error) {
	tp := jira.BasicAuthTransport{
		Username: viper.GetString("jira.username"),
		Password: viper.GetString("jira.token"),
	}
	jiraClient, err := jira.NewClient(tp.Client(), viper.GetString("jira.domain"))
	if err != nil {
		return nil, err
	}

	issue, _, err := jiraClient.Issue.Get(issueId, nil)
	if issue == nil || err != nil {
		fmt.Println("Issue not found")
		return nil, err
	}

	/**
	* 參考 issue.go 程式，將 issue 的資訊印出來
	 */
	// issue 單號 跟 摘要
	fmt.Printf("%s: %+v\n", issue.Key, issue.Fields.Summary)
	// issue 類型
	fmt.Printf("Type: %s\n", issue.Fields.Type.Name)
	// issue 優先權
	fmt.Printf("Priority: %s\n", issue.Fields.Priority.Name)
	// issue 描述
	fmt.Printf("Description: %s\n", issue.Fields.Description)

	// 回傳 issue
	return issue, nil
}

func GetJiraIssueFromContentLink(content string) ([]string, error) {
	jiraDomain := viper.GetString("jira.domain")
	regexPattern := `https:\/\/` + jiraDomain + `\/browse\/([A-Z]+-\d+)`
	re := regexp.MustCompile(regexPattern)

	matches := re.FindAllStringSubmatch(content, -1)
	if matches == nil {
		return nil, errors.New("no JIRA issue IDs found in content")
	}

	var issueIDs []string
	for _, match := range matches {
		if len(match) > 1 {
			issueIDs = append(issueIDs, match[1])
		}
	}

	return issueIDs, nil
}
