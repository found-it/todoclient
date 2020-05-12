
package main

import (
    "fmt"
    "log"
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

    var url string = `http://104.43.236.196:5000/api/Tasks/`
    list(url)

}
