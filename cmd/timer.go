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
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mglaman/tempo/pkg/tempo"
	"github.com/mglaman/tempo/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/cheggaaa/pb.v1"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var running = false

// timerCmd represents the timer command
var timerCmd = &cobra.Command{
	Use:   "timer",
	Short: "Create a worklog timer",
	Long:  `Creates a timer that can be converted into a worklog`,
	Run: func(cmd *cobra.Command, args []string) {
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			for _ = range c {
				if running {
					running = false
				} else {
					os.Exit(1)
				}
			}
		}()

		start := time.Now()
		running = true

		bar := pb.StartNew(10)
		bar.SetRefreshRate(time.Second)
		for running == true {
			bar.Increment()
			time.Sleep(time.Second)
		}
		bar.Finish()
		elapsed := time.Since(start)

		// Round up to 15 minutes if less than 15 minutes..
		if elapsed.Minutes() < 15 {
			elapsed = time.Duration(int64(time.Minute) * 15)
		}
		// Round off all timers to 15 minute intervals.
		elapsed = elapsed.Round(time.Minute * 15)

		fmt.Println(fmt.Sprintf("Logging %s", elapsed))
		fmt.Println()
		issueKey := util.Prompt("Enter the issue key")
		description := util.Prompt("Worklog description")
		time := time.Now().Local()

		workLog := tempo.WorklogPayload{
			IssueKey:         issueKey,
			TimeSpentSeconds: elapsed.Seconds(),
			BillableSeconds:  elapsed.Seconds(),
			StartDate:        time.Format("2006-01-02"),
			StartTime:        time.Format("15:04:05"),
			Description:      description,
			AuthorUsername:   viper.GetString("username"),
		}

		url := "https://api.tempo.io/2/worklogs"
		j, _ := json.Marshal(workLog)
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(j))
		req.Header.Add("Authorization", "Bearer "+viper.GetString("token"))
		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			panic(err.Error())
		}

		worklog := new(tempo.Worklog)
		_ = json.NewDecoder(resp.Body).Decode(worklog)

		fmt.Println("Time log submitted!")
	},
}

func init() {
	rootCmd.AddCommand(timerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// timerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// timerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
