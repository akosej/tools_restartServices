package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	path, _         = os.Getwd()
	time_restart, _ = strconv.Atoi(Config("time"))
	to_go           = Config("to_go")
)

// Config -- EXTRACT VARIABLES FROM THE CONFIGURATION FILE
func Config(data string) string {
	var value string
	lines, err := File2lines(path + "/gardservices.conf")
	if err != nil {
		fmt.Println("Could not find configuration file")
	} else {
		for _, line := range lines {
			if strings.Contains(line, "gs."+data) {
				corte := strings.Split(line, "=")
				value = corte[1]
			}
		}
	}
	return value
}

// run -- Execute command
func run(arg ...string) {
	head := arg[0]
	parts := arg[1:len(arg)]
	run := exec.Command(head, parts...)
	run.Run()
}

func main() {
	if runtime.GOOS == "windows" {
		fmt.Println("Can't Execute this on a linux machine")
	} else {
		services := strings.Split(to_go, ",")
		for {
			for _, serv := range services {

				result := isActive(serv)
				if result {
					fmt.Println("The service " + serv + " is running")
				} else {
					run("systemctl", "start", serv+".service")
					fmt.Println("Executed the service: " + serv)
				}

			}
			time.Sleep(time.Duration(time_restart) * time.Second)
		}
	}
}

//--------------------------------------------------------------------------------
//--------------------------------------------------------------------------------

func File2lines(filePath string) ([]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return LinesFromReader(f)
}

func LinesFromReader(r io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func isActive(service string) bool {
	out, _ := exec.Command("systemctl", "is-active", service+".service").Output()
	//if err != nil {
	//	fmt.Printf("%s", err)
	//}

	output := string(out[:])
	//fmt.Println(out)
	if output == "active\n" {
		return true
	} else {
		return false
	}

}
