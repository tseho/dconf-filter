package main

import (
    "fmt"
    "flag"
    "os"
    "bufio"
    "regexp"
)

func readFileLines(path string) ([]string, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var lines []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    return lines, scanner.Err()
}

func printLines(lines []string) {
    for _, line := range lines {
        fmt.Println(line)
    }
}

func isWhitelisted(path string, rules []string) (bool) {
    var whitelisted bool = false

    var regex string
    var isExcludeRule bool

    for _, rule := range rules {
        isExcludeRule = bool(string(rule[0]) == "!")

        if isExcludeRule {
            regex = rule[1:]
        } else {
            regex = rule
        }

        match, _ := regexp.MatchString(regex, path)
        if match {
            whitelisted = !isExcludeRule
        }
    }

    return whitelisted
}

func main() {
    rulesFilePtr := flag.String("rules", "", "Path to the file with a list of rules")
    flag.Parse()

    if *rulesFilePtr == "" {
        os.Exit(1)
    }

    rules, err := readFileLines(*rulesFilePtr)
    if err != nil {
        panic(err)
    }

    // fmt.Println("Rules:")
    // fmt.Println(rules)
    // fmt.Println("Stdout:")

    stdin, err := os.Stdin.Stat()
    if err != nil {
        panic(err)
    }

    // Check if there is any content in stdin
    if stdin.Mode() & os.ModeNamedPipe != 0 {
        scanner := bufio.NewScanner(os.Stdin)
        var path string
        var lines []string

        pathRegex, _ := regexp.Compile("^\\[(.+)\\]$")
        propertyRegex, _ := regexp.Compile("^([^=]+)")

        for scanner.Scan() {
            line := scanner.Text()

            // If line is empty, it's a new block. Reset variables,
            // print the previous block and skip to the next line
            if line == "" {
                if len(lines) > 0 {
                    lines = append([]string{"[" + path + "]"}, lines...)
                    lines = append(lines, "")
                    printLines(lines)
                }

                path = ""
                lines = nil
                continue
            }

            // If the line is a new path, eg: [org/gnome/shell]
            if pathRegex.MatchString(line) {
                // Extract real path from line, eg: [org/gnome/shell] => org/gnome/shell
                path = pathRegex.FindStringSubmatch(line)[1]
                continue
            }

            // Extract the property name from the line, eg: foo=bar => foo
            propertyName := propertyRegex.FindStringSubmatch(line)[1]
            // Concatenate path and property
            propertyPath := path + "/" + propertyName

            if isWhitelisted(propertyPath, rules) {
                lines = append(lines, line)
            }
        }
    }
}
