package main

import (
	"context"
	"github.com/gin-gonic/gin"
)

func getProfile(c *gin.Context) {
	var prof *Profile
	err := c.Bind(&prof)
	if err != nil {
		c.JSON(501, err)
		return
	}
	err = repo.conn.QueryRow(context.Background(), selectProfileByID, &prof.ID).Scan(&prof.Name, &prof.Coins)
	if err != nil {
		c.JSON(500, err)
		return
	}
	c.JSON(200, &prof)
}

func getTask(c *gin.Context) {
	var task *Task
	err := c.Bind(&task)
	if err != nil {
		c.JSON(501, err)
		return
	}
	err = repo.conn.QueryRow(context.Background(), selectTaskByID, &task.ID).Scan(
		&task.CreatedBy, 
		&task.DateToComplete,
		&task.TaskType,
		&task.TimeToComplete,
		&task.Distance,
		&task.Reward,
	)
	if err != nil {
		c.JSON(500, err)
		return
	}
	c.JSON(200, &task)

}

func postTask(c *gin.Context) {
	var reqTask *Task

	err := c.Bind(&reqTask)
	if err != nil {
		c.JSON(501, err)
		return
	}
	// FINISH

}

func deleteTask(c *gin.Context) {
	var targetID *string
	_, err := repo.conn.Exec(context.Background(), deleteTaskByID, c.GetString("id"), &targetID)
	if err != nil {
		c.JSON(500, err)
		return
	}

	c.JSON(200, true)
}

func getTasks(c *gin.Context) {
	rows, err := repo.conn.Query(context.Background(), getTasks)
	if err != nil {
		c.JSON(500, err)
		return
	}

	tasks := make([]*Task, 0)

	for rows.Next() {
		task := &Task{}
		err = rows.Scan(
			&task.CreatedBy, 
			&task.DateToComplete,
			&task.TaskType,
			&task.TimeToComplete,
			&task.Distance,
			&task.Reward,
		)
		if err != nil {
			c.JSON(501, err)
			return
		}
		tasks = append(tasks, task)
	}

	c.JSON(200, &tasks)
}

func postDonate(c *gin.Context) {
	type donation struct {
		uid string `json:'uid'`
		coins string `json:'coins'`
	}
	var reqDonation *donation

	err := c.Bind(&reqDonation)
	if err != nil {
		c.JSON(501, err)
		return
	}
	// FINISH'
}


const (
	selectProfileByID = "SELECT uid, name, coins FROM profiles WHERE uid $1"
	selectTaskByID = "SELECT uid, created_by, date_to_complete, task_type, time_to_complete, distance, reward FROM tasks WHERE uid $1"
	postTasks = ""
	deleteTaskByID = "DELETE FROM tasks WHERE uid = $1"
	getTasks = "SELECT * FROM tasks"
	postDonate = ""
)