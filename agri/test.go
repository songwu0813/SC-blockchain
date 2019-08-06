package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func main() {
	cmd := exec.Command("python3", "tmp.py", "add", ":ethereum_state.jpeg")
	out, _ := cmd.Output()
	a := strings.Split(string(out), "\n")
	fmt.Println(a[0])
	fmt.Println(a[1])
}
