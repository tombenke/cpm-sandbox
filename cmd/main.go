package main

import (
	"fmt"
	cpm "github.com/tombenke/cpm_sandbox/internal"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run bin/cpm.go <path-to-file>")
		os.Exit(1)
	}

	// The complete path to the YAML format tasks description
	filePath := os.Args[1]

	// The complete path without extension
	filePathWoExt := cpm.GetFilePathWithoutExtension(filePath)

	// Load the tasks from the YAML file
	tasks := cpm.LoadTasks(filePath)

	// Calculate the critical path
	criticalPath := tasks.CalculateCriticalPath([]string{"A", "E"}, []string{"I"})

	// Print the results in different formats
	tasks.PrintTasks(criticalPath, filePathWoExt)
	tasks.PrintGanttDiagram(criticalPath, filePathWoExt)
}
