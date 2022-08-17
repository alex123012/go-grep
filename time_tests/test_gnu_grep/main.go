package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/alex123012/go-grep/time_tests/utils"
)

func main() {
	pattern, file := utils.ParseArgs()

	username := ""

	start := time.Now()
	o := checkStaticAddressIsFree(pattern, username, file)
	fmt.Printf("GNU-grep: %d ms, result = %t\n", time.Since(start).Microseconds(), o)
}

func checkStaticAddressIsFree(staticAddress string, username string, ccdDir string) bool {
	o := runBash(fmt.Sprintf("grep -rl '%s' %s | grep -vx %s/%s | wc -l", staticAddress, ccdDir, ccdDir, username))

	return strings.TrimSpace(o) == "0"
}

func runBash(script string) string {
	// log.Debug(script)
	cmd := exec.Command("bash", "-c", script)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return (fmt.Sprint(err) + " : " + string(stdout))
	}
	return string(stdout)
}
