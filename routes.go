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
	err = repo.conn.QueryRow(context.Background(), selectProfileByID, &prof.ID).Scan(&prof.Name, &prof.Coins, &prof.Organization)
	if err != nil {
		c.JSON(500, err)
		return
	}
	c.JSON(200, &prof)
}

func updateProfile(c *gin.Context) {
	var prof *Profile
	err := c.Bind(&prof)
	if err != nil {
		c.JSON(501, err)
		return
	}

	var oldProfile *Profile
	_, err = repo.conn.Exec(context.Background(), selectProfileByID, &oldProfile.ID, &oldProfile.Name, &oldProfile.Coins, &oldProfile.Organization)
	if err != nil {
		c.JSON(500, err)
		return
	}

	if prof.ID == "" {
		prof.ID = oldProfile.ID
	}
	if prof.Name == "" {
		prof.Name = oldProfile.Name
	}
	if prof.Coins == 0 {
		prof.Coins = oldProfile.Coins
	}
	if prof.Organization == false {
		prof.Organization = oldProfile.Organization
	}

	_, err = repo.conn.Exec(context.Background(), updateProfilebyID, &prof.ID, &prof.Name, &prof.Coins, &prof.Organization)
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
		&task.Description,
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

	reqTask.CreatedBy = c.GetString("uid")

	_, err = repo.conn.Exec(context.Background(), postTaskQuery, 
		&reqTask.CreatedBy, 
		&reqTask.DateToComplete,
		&reqTask.TaskType,
		&reqTask.TimeToComplete,
		&reqTask.Distance,
		&reqTask.Reward,
		&reqTask.Description,
	)

	//Subtract reward from created_by user

	if err != nil {
		c.JSON(500, err)
		return
	}
	
	c.JSON(200, true)
}

func acceptTask(c *gin.Context) {
	var reqTask *TasksAccepted
	err := c.Bind(&reqTask)
	if err != nil {
		c.JSON(501, err)
		return
	}

	reqTask.UID = c.GetString("uid")

	_, err = repo.conn.Exec(context.Background(), postAcceptTask, 
		&reqTask.UID, 
		&reqTask.TaskID,
	)

	if err != nil {
		c.JSON(500, err)
		return
	}
	
	c.JSON(200, true)
	
}

func completeTask(c *gin.Context){
	var existingTask *Task
	_, err = repo.conn.QueryRow(context.Background(), selectTaskByID, c.GetString("task_id")).Scan(&existingTask.ID, &existingTask.CreatedBy, &existingTask.DateToComplete, &existingTask.TaskType, &existingTask.TimeToComplete, &existingTask.Distance, &existingTask.Reward, &existingTask.Description, &existingTask.Status)
	if err != nil {
		c.JSON(500, err)
		return
	}

	//FIXME: does this only return one entry? or can several users accept task
	var accpetedTaskEntry *TasksAccepted
	_, err = repo.conn.QueryRow(context.Background(), selectTaskByID, c.GetString("task_id")).Scan(&accpetedTaskEntry.UID, &accpetedTaskEntry.TaskID)
	if err != nil {
		c.JSON(500, err)
		return
	}

	//FIXME: Assumes only one user can accept task
	var acceptedProfile *Profile
	_, err = repo.conn.QueryRow(context.Background(), selectProfileByID, &accpetedTaskEntry.UID).Scan(&acceptedProfile.ID, &acceptedProfile.Name, &acceptedProfile.Coins, &acceptedProfile.Organization)
	if err != nil {
		c.JSON(500, err)
		return
	}

	var updatedProfile *Profile
	updatedProfile.coins = acceptedProfile.Coins + existingTask.reward
	_, err = repo.conn.Exec(context.Background(), updateProfilebyID, &updatedProfile.ID, &updatedProfile.Name, &updatedProfile.Coins, &updatedProfile.Organization)
	if err != nil {
		c.JSON(500, err)
		return
	}
	


	var taskID *string
	_, err := repo.conn.Exec(context.Background(), updateTaskByID, c.GetString("id"), &taskID)
	if err != nil {
		c.JSON(500, err)
		return
	}
	
	
	// Get user who completed the task by their uid, update task in task table, give suer reward.
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
	rows, err := repo.conn.Query(context.Background(), getTasksQuery)
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
			&task.Description,
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
	type Donation struct {
		UID string `json:'uid'`
		Coins string `json:'coins'`
	}
	var reqDonation *Donation

	err := c.Bind(&reqDonation)
	if err != nil {
		c.JSON(501, err)
		return
	}
	// FINISH'
}


const (
	selectProfileByID = "SELECT uid, name, coins. organization FROM profiles WHERE uid $1;"
	updateProfilebyID = "UPDATE profiles SET (name, coins, organization) WHERE uid $1"
	selectTaskByID = "SELECT uid, created_by, date_to_complete, task_type, time_to_complete, distance, reward, description FROM tasks WHERE id $1;"
	updateTaskByID = "UPDATE tasks SET (status, reward) WHERE uid $1;"
	postTaskQuery = "INSERT_INTO tasks (uid, created_by, date_to_complete, task_type, time_to_complete, distance, reward, description) VALUES ($1,$2,$3,$4,$5,$6,$7);"
	deleteTaskByID = "DELETE FROM tasks WHERE uid = $1;"
	getTasksQuery = "SELECT * FROM tasks WHERE ID NOT IN (SELECT TaskID FROM tasks_accepted);"
	// postTaskSubtract = ""
	// postDonateAdd = ""
	// addInitialOrgCoins = ""
	selectAcceptedTask = "SELECT uid, task_id FROM tasks_accepted WHERE task_id $1"
	postAcceptTask = "INSERT_INTO tasks_accepted (uid, task_id) VALUES ($1,$2)"
)