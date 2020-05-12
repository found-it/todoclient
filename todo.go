
package main

import (
    "os"
    "fmt"
    "log"
    "flag"
    "bytes"
    "strings"
    "os/user"
    "net/http"
    "io/ioutil"
    "path/filepath"
    "encoding/json"
)

import "gopkg.in/yaml.v2"


type config struct {
    Url     string  `yaml:"url"`
    Timeout int     `yaml:"timeout"`
}

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

func check(e error) {
    if e != nil {
        log.Fatal(e)
    }
}

func initialize() {

    usr, err := user.Current()
    check(err)

    configfile := filepath.Join(usr.HomeDir, ".todo.cfg")

    if _, err := os.Stat(configfile); err == nil {
        // Config exists, let the user know
        fmt.Println("Configuration already exists at:", configfile)
        os.Exit(1)

    } else if os.IsNotExist(err) {

        cfg := `
#url: http://localhost/api/Tasks
#timeout: 10000
`
        err = ioutil.WriteFile(configfile, []byte(cfg), 0644)
        check(err)
        fmt.Println("A configuration file has been created at:", configfile)

    } else {
        log.Fatal(err)
    }
}



//
//  Fetch configuration from `~/.todo.config`
//
//  The config file must be in the format
//      url: http://<url of server>/api/Tasks
//      timeout: <number>
//
func (c *config) fetch() *config {

    usr, err := user.Current()
    check(err)

    configfile := filepath.Join(usr.HomeDir, ".todo.cfg")

    if _, err := os.Stat(configfile); err == nil {
        // Config exists so use it

        file, err := ioutil.ReadFile(configfile)
        check(err)

        err = yaml.Unmarshal(file, c)
        check(err)

        if c.Url == "" {
            log.Fatal("Must specify a url in ~/.todo.config")
        }

        if !strings.Contains(c.Url, "/api/Tasks") {
            log.Fatal("Url must contain the API path '/api/Tasks'")
        }

    } else if os.IsNotExist(err) {

        fmt.Println("No configuration found at:", configfile)
        fmt.Println("Please run 'todo init'")
        os.Exit(1)

    } else {
        log.Fatal(err)
    }

    return c
}


//
//  List function
//
//  Arguments
//    c     Configuration for server
//
func list(c config) {

    resp, err := http.Get(c.Url)
    check(err)

    buf, err := ioutil.ReadAll(resp.Body)
    defer resp.Body.Close()

    check(err)

    var tasks []Task
    json.Unmarshal(buf, &tasks)

    printer(tasks)

}


//
//  Add function
//
//  Arguments
//    c     Configuration for server
//    name  Name of task to add
//
func add(c config, name string) {

    fmt.Printf("Sending '%s' to '%s'", name, c.Url)

    task := &Task {
        Complete: false,
        Name: name,
    }

    jsontask, err := json.Marshal(task)
    check(err)

    resp, err := http.Post(c.Url, "application/json", bytes.NewReader(jsontask))
    check(err)

    if resp.StatusCode != 201 {
        log.Fatalf("[%s] Error posting task '%s'", resp.Status, name)
    }
}


//
//  Main function
//
//  Implements
//    - list
//    - add "todo task"
//
func main() {

    listCmd  := flag.NewFlagSet("list",  flag.ExitOnError)
    addCmd   := flag.NewFlagSet("add",   flag.ExitOnError)
    initCmd   := flag.NewFlagSet("init",  flag.ExitOnError)
    argCount := len(os.Args[1:])

    if argCount < 1 {
        log.Fatal("Need a command [ init, list, add ]")
    }

    if os.Args[1] == "init" {
        initCmd.Parse(os.Args[2:])
        initialize()
        os.Exit(0)
    }


    var c config
    c.fetch()

    switch os.Args[1] {

    case "list":
        listCmd.Parse(os.Args[2:])
        list(c)

    case "add":
        if argCount < 2 {
            log.Fatal("Need a task to add")
        }
        addCmd.Parse(os.Args[2:])
        add(c, addCmd.Arg(0))

    default:
        fmt.Println("Expected different command than [", os.Args[1], "]")
        log.Fatal("Need a command [ list, add ]")
    }

}
