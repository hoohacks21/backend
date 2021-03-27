package main

import (
	"time"
)

type Profile struct {
	ID string 	`json:"id"`
	Name string `json:"name"`
	Coins int 	`json:"coins"`
	Organization bool `json:"organization"`
}

type Task struct {
	ID string 						`json:"id"`
	CreatedBy string 				`json:"created_by"`
	DateToComplete time.Time 		`json:"date_to_complete"`
	TaskType string 				`json:"task_type"`
	TimeToComplete time.Duration 	`json:"time_to_complete"`
	Distance int 					`json:"distance"`
	Reward int 						`json:"reward"`
	Description string				`json:"description"`
	Status int					    `json:"status"`
}

type TasksAccepted struct {
	UID string 						`json:"uid"`
	TaskID string 					`json:"task_id"`
}

type GetTaskResponse []*Task