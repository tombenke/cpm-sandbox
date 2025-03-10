package internal

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

type Task struct {
	// The unique ID of the task
	Name string

	// The duration of the task
	Duration int `yaml:"duration"`

	// The Dependants of the task
	Deps []string `yaml:"deps"`

	// The Earliest Start of the task
	ES int

	// The Earliest Finish of the task
	EF int

	// The Latest Start of the task
	LS int

	// The Latest Finish of the task
	LF int

	// The subordinal ranking of the tasks,
	// used for GANTT chart generation
	subOrdRank int
}

type Tasks map[string]*Task

// CalculateCriticalPath Calculates the critical path or the tasks
// according to the dependencies and durations assigned to them.
// Returns with a string array, that holds the IDs of the criticval path, ordered by their dependencies.
func (tasks Tasks) CalculateCriticalPath(start []string, end []string) []string {

	// Do the forward calculation of ES and EF values
	tasks.calculateCriticalPathFwd(start)

	// Do the backward calculation of LS and SF values
	tasks.calculateCriticalPathBwd(end)

	// Find the nodes of the critical path
	var criticalPath []*Task
	for _, task := range tasks {
		if task.ES == task.LS && task.EF == task.LF {
			criticalPath = append(criticalPath, task)
		}
	}

	// Sort the nodes of the critical path, and get their IDs in the right order
	sort.Sort(CriticalPath(criticalPath))
	var criticalPathNames []string
	for _, task := range criticalPath {
		criticalPathNames = append(criticalPathNames, task.Name)
	}

	return criticalPathNames
}

// contains checks if a string is contained by an array
func contains(elems []string, s string) bool {
	for _, v := range elems {
		if v == s {
			return true
		}
	}
	return false
}

func (tasks Tasks) getDependants(node *Task) []string {
	dependants := []string{}
	for _, task := range tasks {
		if contains(task.Deps, node.Name) {
			dependants = append(dependants, task.Name)
		}
	}
	return dependants
}

func (tasks Tasks) calculateNodeFwd(node *Task, EF int) {
	if EF >= node.ES {
		node.ES = EF
		node.EF = node.ES + node.Duration

		// Prepare for backward calculation
		node.LF = math.MaxInt64
	}
	dependants := tasks.getDependants(node)
	for _, dep := range dependants {
		tasks.calculateNodeFwd(tasks[dep], node.EF)
	}
}

func (tasks Tasks) calculateCriticalPathFwd(starts []string) {
	for _, start := range starts {
		node := tasks[start]
		tasks.calculateNodeFwd(node, 0)
	}
}

func (tasks Tasks) calculateNodeBwd(node *Task, LF int) {
	if LF <= node.LF {
		node.LF = LF
		node.LS = node.LF - node.Duration
	}
	for _, dep := range node.Deps {
		tasks.calculateNodeBwd(tasks[dep], node.LS)
	}
}

func (tasks Tasks) calculateCriticalPathBwd(ends []string) {
	for _, end := range ends {
		node := tasks[end]
		tasks.calculateNodeBwd(node, node.EF)
	}
}

type CriticalPath []*Task

func (a CriticalPath) Len() int           { return len(a) }
func (a CriticalPath) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a CriticalPath) Less(i, j int) bool { return a[i].ES <= a[j].ES }

// PrintTasks Prints the actual state of the tasks array in plain text format
func (tasks Tasks) PrintTasks(criticalPath []string, fileName string) {
	var sb strings.Builder
	for _, task := range tasks {
		sb.WriteString(fmt.Sprintf("[%s - D:%d, ES:%d, EF:%d, LS:%d, LF:%d, Slack:%d]\n",
			task.Name,
			task.Duration,
			task.ES,
			task.EF,
			task.LS,
			task.LF,
			task.LS-task.ES,
		))
	}
	sb.WriteString(fmt.Sprintf("Critical Path: %+v\n", criticalPath))
	WriteFile(fileName+".txt", sb.String())
}
