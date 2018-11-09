package shell

import (
	"os/exec"
	"github.com/labstack/gommon/log"
	"fmt"
)

func ExecShell(shell string) string {
	out,err := exec.Command("/bin/bash","-c",shell).Output()
	if err != nil {
		log.Fatal(err)
	}
	date := string(out)
	return date
}

func main()  {
	shell := "date +%Y-%m-%d"
	out := ExecShell(shell)
	fmt.Println(out)
}