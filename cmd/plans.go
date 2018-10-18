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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"time"
)

type PlanPeriod struct {
	From string
	To   string
	// Define as a float so we can convert it to hours and allow a decimal.
	TimePlannedSeconds float32
}

type Plans struct {
	Self        string
	Id          int
	StartDate   string
	EndDate     string
	CreatedAt   string
	Description string
	UpdatedAt   string
	Assignee    struct {
		Self string
		Type string
	}
	PlanItem struct {
		Self string
		Type string
	}
	Recurrence struct {
		Rule              string
		RecurrenceEndDate string
	}
	Dates struct {
		Metadata struct {
			Count int
			All   string
		}
		Values []PlanPeriod
	}
}
type PlansResponse struct {
	Self     string
	Metadata struct {
		Count int
	}
	Results []Plans
}

// plansCmd represents the plans command
var plansCmd = &cobra.Command{
	Use:   "plans",
	Short: "Retrieve your work plans for the day",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		time := time.Now().Local().Format("2006-01-02");
		token := viper.GetString("token")
		username := viper.GetString("username")
		url := "https://api.tempo.io/2/plans/user/" + username + "?from=" + time + "&to=" + time

		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("Authorization", "Bearer "+token)

		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			panic(err.Error())
		}

		totalHours := float32(0)

		plans := new(PlansResponse)
		_ = json.NewDecoder(resp.Body).Decode(plans)
		//fmt.Println(plans.Results)
		table := tablewriter.NewWriter(os.Stdout)
		table.SetBorders(tablewriter.Border{Left: true, Top: true, Right: true, Bottom: false})
		table.SetCenterSeparator("|")
		table.SetColMinWidth(0, 50)
		table.SetColWidth( 50)
		table.SetHeader([]string{"Description", "Hours"})
		for _, plan := range plans.Results {
			for _, planItem := range plan.Dates.Values {
				table.Append([]string{
					plan.Description,
					fmt.Sprintf("%v", planItem.TimePlannedSeconds/60/60),
				})
				totalHours += planItem.TimePlannedSeconds / 60 / 60;
			}
		}
		table.SetFooter([]string{
			"", fmt.Sprintf("%v", totalHours),
		})
		table.Render()
	},
}

func init() {
	rootCmd.AddCommand(plansCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// plansCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// plansCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
