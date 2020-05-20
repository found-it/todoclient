
package main

import (
    "os"
    // "io"
    "fmt"
    "log"
    "flag"
    "bytes"
    // "strings"
    "os/user"
    "strconv"
    "hash/fnv"
    "net/http"
    "io/ioutil"
    "path/filepath"
    "encoding/json"
    "gopkg.in/yaml.v2"
    "github.com/fatih/color"
)



type config struct {
    Url     string  `yaml:"url"`
    Timeout int     `yaml:"timeout"`
}

type Tasks struct {
    Todo []Task `json:"todo"`
    // Done []Task `json:"done"`
}

type Task struct {
    Id       string     `json:"id"`
    Name     string     `json:"name"`
    Complete bool       `json:"complete"`
    // Tags    []string    `json:"tags"`
}


func hash(s string) string {
    h := fnv.New32a()
    h.Write([]byte(s))
    return strconv.FormatUint(uint64(h.Sum32()), 10)
}


//
//  Helper function to print out the todos
//
func printer(task []Task) {

    var verbose = true
    var c config
    c.fetch()

    yellow  := color.New(color.FgYellow, color.Bold).SprintFunc()
    green   := color.New(color.FgGreen, color.Bold).SprintFunc()
    magenta := color.New(color.FgMagenta).SprintFunc()
    blue    := color.New(color.FgBlue).SprintFunc()

    fmt.Println()
    fmt.Printf("%s [%s]\n", green("Todo Server"), magenta(system(c)))
    fmt.Println("---------------------------------")
    fmt.Println()

    for i := 0; i < len(task); i++ {
        var status string = "TODO"
        if task[i].Complete {
            status = "DONE"
        }
        if verbose {
            fmt.Printf("%-10s |  %s: %s\n", task[i].Id, yellow(status), blue(task[i].Name))// , task[i].Tags)
        } else {
            fmt.Printf("%s: %s\n", yellow(status), blue(task[i].Name))// , task[i].Tags)
        }
    }
    fmt.Println()
}


//
//  Helper function to check errors
//
func check(e error) {
    if e != nil {
        log.Fatal(e)
    }
}



//
//  Initialize the configuration
//
//  TODO: --force initialize if there is already a file there
//  TODO: --url to initialize with
//  TODO: --timeout to initialize with
//
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
#url: http://localhost:9000
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

    const path = "/api/tasks/"

    resp, err := http.Get(c.Url + path)
    check(err)

    buf, err := ioutil.ReadAll(resp.Body)
    defer resp.Body.Close()

    check(err)

    var tasks []Task
    err = json.Unmarshal(buf, &tasks)
    check(err)

    printer(tasks)
}


//
//  Print out the system info
//
func system(c config) string {

    const path = "/api/system"

    resp, err := http.Get(c.Url + path)
    check(err)

    buf, err := ioutil.ReadAll(resp.Body)
    defer resp.Body.Close()
    check(err)

    type SystemInfo struct {
        Hostname string `json:"hostname"`
    }

    var si SystemInfo
    json.Unmarshal(buf, &si)

    return si.Hostname
}


//
//  Add function
//
//  Arguments
//    c     Configuration for server
//    name  Name of task to add
//
func add(c config, name string) {

    const path = "/api/create"

    fmt.Printf("Sending '%s' to '%s'", name, c.Url + path)

    task := &Task {
        Id: hash(name),
        Complete: false,
        Name: name,
    }

    jsontask, err := json.Marshal(task)
    check(err)

    resp, err := http.Post(c.Url + path, "application/json", bytes.NewReader(jsontask))
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
    initCmd  := flag.NewFlagSet("init",  flag.ExitOnError)
    // sysCmd   := flag.NewFlagSet("system",  flag.ExitOnError)
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

    case "system":
        fmt.Println("Hostname: ", system(c))

    default:
        fmt.Println("Expected different command than [", os.Args[1], "]")
        log.Fatal("Need a command [ list, add ]")
    }

}
