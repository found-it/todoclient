
package main

import (
    "os"
    "fmt"
    "log"
    "flag"
    "net/http"
    "io/ioutil"
    "encoding/json"
)

type Tasks struct {
    Todo []Task `json:"todo"`
    // Done []Task `json:"done"`
}

type Task struct {
    Id       uint32     `json:"id"`
    Name     string     `json:"name"`
    Complete bool       `json:"complete"`
    // Tags    []string    `json:"tags"`
}

func printer(task []Task) {
    for i := 0; i < len(task); i++ {
        var status string = "TODO"
        if task[i].Complete {
            status = "DONE"
        }
        fmt.Println("  " + status + ": " + task[i].Name)// , task[i].Tags)
    }
}


func list(url string) {

    resp, err := http.Get(url)
    if err != nil {
        log.Fatal(err)
    }

    buf, err := ioutil.ReadAll(resp.Body)
    defer resp.Body.Close()

    if err != nil {
        log.Fatal(err)
    }

    var tasks []Task
    json.Unmarshal(buf, &tasks)

    printer(tasks)

}


func main() {

    const url string = `http://104.43.236.196:5000/api/Tasks/`

    listCmd  := flag.NewFlagSet("list",  flag.ExitOnError)
    addCmd   := flag.NewFlagSet("add",   flag.ExitOnError)
    argCount := len(os.Args[1:])

    if argCount < 1 {
        log.Fatal("Need a command [ list, add ]")
    }

    switch os.Args[1] {

    case "list":
        listCmd.Parse(os.Args[2:])
        list(url)

    case "add":
        if argCount < 2 {
            log.Fatal("Need a task to add")
        }
        addCmd.Parse(os.Args[2:])
        fmt.Println(addCmd.Args())

    default:
        fmt.Println("Expected different command than [", os.Args[1], "]")
        log.Fatal("Need a command [ list, add ]")
    }

}
