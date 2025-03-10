package internal

import (
	"gopkg.in/yaml.v2"
)

// Load tasks from a YAML format file
func LoadTasks(path string) Tasks {
	tasks := make(Tasks)
	content := LoadFile(path)
	err := yaml.Unmarshal([]byte(content), &tasks)
	if err != nil {
		panic(err)
	}

	for k, _ := range tasks {
		tasks[k].Name = k
	}

	return tasks
}
