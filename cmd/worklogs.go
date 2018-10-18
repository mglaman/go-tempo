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
	"github.com/mglaman/tempo/pkg/tempo"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// worklogsCmd represents the logs command
var worklogsCmd = &cobra.Command{
	Use:   "worklogs",
	Short: "Retrieve todays worklogs",
	Long: `Retrieves your worklogs for the current day, limited to 50
@todo: support specifying days, date range, and pagination.`,
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

		worklogs := new(tempo.WorklogCollection)
		_ = json.NewDecoder(resp.Body).Decode(worklogs)

		width, _, _ := terminal.GetSize(0)
		table := tablewriter.NewWriter(os.Stdout)
		table.SetBorders(tablewriter.Border{Left: true, Top: true, Right: true, Bottom: false})
		table.SetCenterSeparator("|")
		table.SetColWidth(width - 25)
		table.SetAutoWrapText(true)
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
			fmt.Print(fmt.Sprintf("Watch out! You have only logged %v of 8 hours today", totalHours))
			fmt.Println()
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
