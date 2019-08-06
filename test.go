package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	// "strings"
)

func main() {
	fmt.Println("vim-go")
	cmd := exec.Command("ssh", "-i", "awsEC-ubuntu.pem", "ubuntu@18.219.71.129", "source ~/.bashrc", ";", "python3", "CalCadAddr.py", "-a", "121.53264", "-b", "25.04228", "-c", "0.00000025", "-d", "3.8")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	var k map[string]interface{}
	err = json.Unmarshal(out, &k)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(k["result"].(map[string]interface{})["cadaddr"])
}
