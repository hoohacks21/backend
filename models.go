package main

import (
	"time"
)

type Profile struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Coins        int    `json:"coins"`
	Organization bool   `json:"organization"`
}

type Task struct {
	ID             int       `json:"id"`
	CreatedBy      string    `json:"created_by"`
	DateToComplete time.Time `json:"date_to_complete"`
	TaskType       string    `json:"task_type"`
	TimeToComplete int       `json:"time_to_complete"`
	Lat            float64   `json:"lat"`
	Long           float64   `json:"long"`
	Reward         int       `json:"reward"`
	Description    string    `json:"description"`
}

// type TaskResponse struct {
// 	Task
// 	Distance int 					`json:"distance"`
// }

// type TaskRequest struct {
// 	Task
// 	Lat float64 						`json:"lat"`
// 	Long float64 				`json:"long"`
// }

type TasksAccepted struct {
	UID    string `json:"uid"`
	TaskID int    `json:"task_id"`
}

type GetTaskResponse []*Task

type TaskCompleteRequest struct {
	UID    string `json:"uid"`
	TaskID int    `json:"task_id"`
}
