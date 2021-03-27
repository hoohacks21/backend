package main

import (
	"time"
)

type Profile struct {
	string Name `json:'name'`
	int Coins	`json:'coins'`
}

type Task struct {
	string CreatedBy `json:'created_by'`
	time.Time DateToComplete `json:'date_to_complete'`
	string TaskType `json:'task_type'`
	time.Duration TimeToComplete `json:'time_to_complete'`
	int Distance `json:'distance'`
	int Reward `json:'reward'`
}

type GetTaskResponse []*Task

