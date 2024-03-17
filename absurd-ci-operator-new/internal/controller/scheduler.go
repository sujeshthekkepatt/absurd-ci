package controller

import (
	"container/list"
)

type Scheduler struct {
	StepsList map[string][]string
}

func InitScheduler(steps []string) *Scheduler {

	scheduler := Scheduler{
		StepsList: make(map[string][]string),
	}

	for _, val := range steps {

		scheduler.StepsList[val] = nil

	}

	return &scheduler
}

func (s *Scheduler) AddDependencies(parentStep, childStep string) {

	s.StepsList[parentStep] = append(s.StepsList[parentStep], childStep)

}

func (s *Scheduler) TopoSort(stepName string, stack *list.List, visited map[string]bool) {

	visited[stepName] = true
	elems := s.StepsList[stepName]

	for _, val := range elems {

		if !visited[val] {

			s.TopoSort(val, stack, visited)
		}

	}
	stack.PushBack(stepName)

}

func (s *Scheduler) Schedule() []string {

	var scheduleOrder []string

	visited := make(map[string]bool)

	for key, _ := range s.StepsList {

		visited[key] = false
	}

	stack := list.New()

	for key, val := range visited {

		if !val {

			s.TopoSort(key, stack, visited)
		}

	}

	for stack.Len() != 0 {

		val := stack.Back()

		stackItem := val.Value.(string)
		scheduleOrder = append(scheduleOrder, stackItem)

		stack.Remove(val)
	}

	return scheduleOrder
}
