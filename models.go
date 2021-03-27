package main

import (
	"time"
)

type Profile struct {
	ID string 	`json:'ID'`
	Name string `json:'name'`
	Coins int 	`json:'coins'`
	Organization bool `json:'organization'`
}

type Task struct {
	ID string 						`json:'ID'`
	CreatedBy string 				`json:'created_by'`
	DateToComplete time.Time 		`json:'date_to_complete'`
	TaskType string 				`json:'task_type'`
	TimeToComplete time.Duration 	`json:'time_to_complete'`
	Distance int 					`json:'distance'`
	Reward int 						`json:'reward'`
	Description string				`json:'description'`
}

type TasksAccepted struct {
	UID string 						`json:'uid'`
	TaskID string 					`json:'task_id'`
}

type GetTaskResponse []*Task