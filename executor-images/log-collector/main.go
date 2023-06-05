package main

import (
	"fmt"
	"os/exec"
)

func main() {

	cmd := exec.Command("kubectl", "get", "pods", "--output=json")

	// Run the command and capture its output
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error executing kubectl command:", err)
		return
	}

	// Print the output
	fmt.Println(string(output))

}
