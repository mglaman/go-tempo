// Copyright © 2018 Matt Glaman <nmd.matt@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
)

type WorkAttributeValue struct {
	Key   string
	Value string
}
type User struct {
	Self        string
	Username    string
	DisplayName string
}
type Worklog struct {
	Self           string
	TempoWorklogId int
	JiraWorklogId  int
	Issue          struct {
		Self string
		Key  string
	}
	// Set to float so we can perform math
	TimeSpentSeconds float32
	StartDate        string
	StartTime        string
	Description      string
	CreatedAt        string
	UpdatedAt        string
	Author           User
	Attributes       struct {
		Self  string
		Items []WorkAttributeValue
	}
}
type LogsResponse struct {
	Self     string
	Metadata struct {
		Count    int
		Offset   int
		Limit    int
		Next     string
		Previous string
	}
	Results []Worklog
}

// worklogsCmd represents the logs command
var worklogsCmd = &cobra.Command{
	Use:   "worklogs",
	Short: "Retrieve todays worklogs",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		time := time.Now().Local().Format("2006-01-02")
		token := viper.GetString("token")
		username := viper.GetString("username")
		url := "https://api.tempo.io/2/worklogs/user/" + username + "?from=" + time + "&to=" + time

		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("Authorization", "Bearer "+token)

		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			panic(err.Error())
		}
		totalHours := float32(0)

		worklogs := new(LogsResponse)
		_ = json.NewDecoder(resp.Body).Decode(worklogs)

		table := tablewriter.NewWriter(os.Stdout)
		table.SetBorders(tablewriter.Border{Left: true, Top: true, Right: true, Bottom: false})
		table.SetCenterSeparator("|")
		table.SetColMinWidth(1, 75)
		table.SetColWidth(75)
		table.SetHeader([]string{"Issue", "Description", "Hours"})
		for _, worklog := range worklogs.Results {
			table.Append([]string{
				worklog.Issue.Key,
				worklog.Description,
				fmt.Sprintf("%v", worklog.TimeSpentSeconds/60/60),
			})
			totalHours += worklog.TimeSpentSeconds / 60 / 60
		}
		table.SetFooter([]string{
			"", "", fmt.Sprintf("%v", totalHours),
		})
		table.Render()

		if totalHours < 8 {
			fmt.Println()
			fmt.Printf("Watch out! You have only logged %v of 8 hours today", totalHours)
		}
	},
}

func init() {
	rootCmd.AddCommand(worklogsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// logsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// logsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
