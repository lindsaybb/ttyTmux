package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	helpFlag   = flag.Bool("h", false, "Show this help")
	baudFlag   = flag.Int("b", 115200, "Baudate (for all devices, sorry!)")
	todaysDate string
	bashPath   = "/bin/bash"
	lsDevs     = []string{bashPath, "-c", "ls /dev/ttyUSB*"}
	tmuxPath   = "/usr/bin/tmux"
	tmuxLs     = []string{tmuxPath, "ls"}
	tmuxNew    = []string{tmuxPath, "new-session", "-d", "-s"}
	tmuxKill   = []string{tmuxPath, "kill-session", "-t"}
	tmuxSes    []string
	ttyDevs    []string
	result     string
	response   []string
)

func init() {
	flag.Parse()
	if *helpFlag {
		flag.PrintDefaults()
		os.Exit(0)
	}
	todaysDate = generateDatestamp()
}

func listTmuxSessions() (ses []string) {
	result = execCommand(tmuxLs)
	// error here is typically no running sessions
	if strings.HasPrefix(result, "ERR") {
		return
	}
	response = strings.Split(result, "\n")
	for i := 0; i < len(response); i++ {
		s := strings.Split(response[i], ":")
		ses = append(ses, s[0])
	}
	return
}

func killTmuxSessions(session string) {
	tmuxKill = append(tmuxKill, session)
	result = execCommand(tmuxKill)
	tmuxKill = tmuxKill[:len(tmuxKill)-1]
	// ignore error but fix list first
	if strings.HasPrefix(result, "ERR") {
		fmt.Println(result)
	}
}

func listTtyDevs() (devs []string) {
	result = execCommand(lsDevs)
	if strings.HasPrefix(result, "ERR") {
		fmt.Println(result)
		os.Exit(1)
	}
	if len(result) < 1 {
		return
	}
	devs = strings.Fields(result)
	return
}

func main() {
	// kill all previously running sessions that use this program's syntax
	tmuxSes = listTmuxSessions()
	for i := 0; i < len(tmuxSes); i++ {
		killTmuxSessions(tmuxSes[i])
	}
	// look for all active serial adapters
	ttyDevs = listTtyDevs()
	// if a process is using the device, it won't allow connect
	for i := 0; i < len(ttyDevs); i++ {
		// uses device ID to allow for non-linear naming (skips)
		tmuxNew = append(tmuxNew, fmt.Sprintf("tty%s", ttyDevs[i][len(ttyDevs[i])-1:]))
		result = execCommand(tmuxNew)
		tmuxNew = tmuxNew[:len(tmuxNew)-1]
		if strings.HasPrefix(result, "ERR") {
			fmt.Println(result)
		}
	}
	tmuxSes = listTmuxSessions()
	for i := 0; i < len(ttyDevs); i++ {
		result = execCommand([]string{
			"/usr/bin/tmux",
			"send-keys",
			"-t",
			tmuxSes[i],
			fmt.Sprintf("minicom -D %s -b %d -C %s_%s.log", ttyDevs[i], *baudFlag, todaysDate, tmuxSes[i]),
			"Enter"})
		if strings.HasPrefix(result, "ERR") {
			fmt.Printf("[%s] %s", ttyDevs[i], result)
		}
	}
	// how can I verify/ what do I need to know all is good?
}

func execCommand(call []string) string {
	cmd := exec.Command(call[0], call[1:]...)
	output, err := cmd.CombinedOutput()
	//fmt.Printf("%s :: %s\n", output, err)
	if err != nil {
		return fmt.Sprintf("ERR: %s\n", err)
	} else {
		return string(output)
	}
}

func generateDatestamp() string {
	t := time.Now()
	tf := t.Format(time.RFC3339)
	tf = strings.Replace(tf, "-", "", -1)
	return tf[:8]
}
