package internal

import (
	"fmt"
	"sort"
	"strings"
)

// PrintGanttDiagram Prints the tasks array in the format of a PlantUML gantt diagram source code
func (tasks Tasks) PrintGanttDiagram(criticalPath []string, fileName string) {
	tasks.calculateSubOrdRanks()
	tasksOrdered := tasks.orderBySubOrdRanks()

	var sb strings.Builder
	sb.WriteString("@startgantt\n")

	for _, task := range tasksOrdered {

		sb.WriteString(fmt.Sprintf("[%s] requires %d days\n",
			task.Name,
			task.Duration,
		))

		if contains(criticalPath, task.Name) {
			sb.WriteString(fmt.Sprintf("[%s] is colored in LightGray/Red\n",
				task.Name,
			))
		}

		if len(task.Deps) > 0 {
			for _, dep := range task.Deps {
				sb.WriteString(fmt.Sprintf("[%s] starts at [%s]'s end\n", task.Name, dep))
			}
		}
	}

	sb.WriteString("@endgantt\n")
	WriteFile(fileName+".uml", sb.String())
}

func (tasks Tasks) getDepsSubOrdRank(deps []string) (error, int) {
	err := error(nil)
	highestSubOrdRank := -1
	for _, dep := range deps {
		subOrdRank := tasks[dep].subOrdRank
		if subOrdRank < 0 {
			err = fmt.Errorf("Missing subOrdRank")
		}
		if subOrdRank > highestSubOrdRank {
			highestSubOrdRank = subOrdRank
		}
	}
	return err, highestSubOrdRank
}

func (tasks Tasks) calculateTaskSubOrdRank(task *Task) *Task {
	if len(task.Deps) == 0 {
		task.subOrdRank = 1
	} else {
		if err, subOrdRank := tasks.getDepsSubOrdRank(task.Deps); err == nil {
			task.subOrdRank = subOrdRank + 1
		}
	}

	return task
}

func (tasks Tasks) calculateSubOrdRanks() {
	numTasks := len(tasks)
	for t := 0; t < numTasks; t++ {
		for ti, _ := range tasks {
			tasks[ti] = tasks.calculateTaskSubOrdRank(tasks[ti])
		}
	}
}

type TasksBySubOrdRank []Task

func (a TasksBySubOrdRank) Len() int           { return len(a) }
func (a TasksBySubOrdRank) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a TasksBySubOrdRank) Less(i, j int) bool { return a[i].subOrdRank < a[j].subOrdRank }

func (tasks Tasks) orderBySubOrdRanks() []Task {
	tasksArray := []Task{}
	for _, task := range tasks {
		tasksArray = append(tasksArray, (*task))
	}
	sort.Sort(TasksBySubOrdRank(tasksArray))
	return tasksArray
}
