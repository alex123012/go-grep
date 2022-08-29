package main

import (
	"crypto/md5"
	"fmt"
	"os/exec"
	"time"

	"github.com/alex123012/go-grep/time_tests/utils"
)

func main() {
	pattern, file := utils.ParseArgs()

	username := ""

	start := time.Now()
	o := checkStaticAddressIsFree(pattern, username, file)
	fmt.Printf("GNU-grep: %d ms, result = %t\n", time.Since(start).Microseconds(), o == "exit status 1 : ")
	// fmt.Println(o)
}

func checkStaticAddressIsFree(staticAddress string, username string, ccdDir string) string {
	o := runBash(fmt.Sprintf("grep -rl '%s' %s", staticAddress, ccdDir))

	return o
}

func runBash(script string) string {
	cmd := exec.Command("bash", "-c", script)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return (fmt.Sprint(err) + " : " + string(stdout))
	}
	md5.Sum([]byte("lol"))
	return string(stdout)
}
