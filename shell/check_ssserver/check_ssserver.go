package main

import (
	"go_instance/shell"
	"fmt"
)

func main()  {
	s := "netstat -an| grep 22 | wc -l"
	data := shell.ExecShell(s)
	fmt.Println(data)
}
