/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
    "log"
    "net/http"
    "io/ioutil"
    "encoding/json"

	"github.com/spf13/cobra"
    "github.com/fatih/color"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks",
	Long: `List all the tasks in the database`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("list called")
        list()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


type Task struct {
    Id       string     `json:"id"`
    Name     string     `json:"name"`
    Complete bool       `json:"complete"`
    // Tags    []string    `json:"tags"`
}


//
//  Helper function to check errors
//
func check(e error) {
    if e != nil {
        log.Fatal(e)
    }
}


func list() {

    const Url = "http://localhost:9000"

    const path = "/api/tasks/"

    resp, err := http.Get(Url + path)
    check(err)

    buf, err := ioutil.ReadAll(resp.Body)
    defer resp.Body.Close()

    check(err)

    var tasks []Task
    err = json.Unmarshal(buf, &tasks)
    check(err)

    printer(tasks, true)

}


//
//  Helper function to print out the todos
//
func printer(task []Task, all bool) {

    var verbose = true

    yellow  := color.New(color.FgYellow, color.Bold).SprintFunc()
    green   := color.New(color.FgGreen, color.Bold).SprintFunc()
    magenta := color.New(color.FgMagenta).SprintFunc()
    blue    := color.New(color.FgBlue).SprintFunc()

    fmt.Println()
    fmt.Printf("%s [%s]\n", green("Todo Server"), magenta("localhost"))
    fmt.Println("---------------------------------")
    fmt.Println()

    max := 0
    for _, t := range task {
        if len(t.Name) > max {
            max = len(t.Name)
        }
    }

    max = max + 10

    for _, t := range task {
        if !t.Complete {
            if verbose {
                fmt.Printf("%s: %-*s  %s\n", yellow("TODO"), max, blue(t.Name), t.Id);
            } else {
                fmt.Printf("%s: %s\n", yellow("TODO"), blue(t.Name))// , task[i].Tags)
            }
        }
    }
    if all {
        fmt.Println()
        for _, t := range task {
            if t.Complete {
                if verbose {
                    fmt.Printf("%s: %-*s  %s\n", yellow("DONE"), max, blue(t.Name), t.Id);
                } else {
                    fmt.Printf("%s: %s\n", yellow("DONE"), blue(t.Name))// , task[i].Tags)
                }
            }
        }
    }
    fmt.Println()
}

