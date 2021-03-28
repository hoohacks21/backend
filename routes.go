package main

import (
	"log"
	"strconv"
	"context"

	"github.com/gin-gonic/gin"
	
	"github.com/jackc/pgx/v4"
)

func getProfile(c *gin.Context) {
	prof := &Profile{}
	
	prof.ID = c.GetString("uid")
	log.Println(selectProfileByID)
	err := repo.conn.QueryRow(context.Background(), selectProfileByID, prof.ID).Scan(&prof.Name, &prof.Coins, &prof.Organization)
	if err != nil  {
		if err == pgx.ErrNoRows{
			log.Printf("[GET PROFILE - NOT FOUND] %v | %v", prof.ID, err)
			c.JSON(500, "Profile does not exist")
			return
		}
		log.Printf("[GET PROFILE] %v", err)
		c.JSON(500, err)
		return
	}
	c.JSON(200, &prof)
}

func updateProfile(c *gin.Context) {
	var prof *Profile
	err := c.Bind(&prof)
	if err != nil &&  err != pgx.ErrNoRows   {
		log.Printf("[UPDATE PROFILE] %v", err)
		c.JSON(501, err)
		return
	}

	prof.ID = c.GetString("uid")

	oldProfile := &Profile{}
	err = repo.conn.QueryRow(context.Background(), selectProfileByID, &prof.ID).Scan(&oldProfile.ID, &oldProfile.Name, &oldProfile.Coins, &oldProfile.Organization)
	if err != nil && err != pgx.ErrNoRows {
		log.Printf("[UPDATE PROFILE] %v", err)
		c.JSON(500, err)
		return
	}

	if prof.ID == "" {
		prof.ID = oldProfile.ID
	}
	if prof.Name == "" {
		prof.Name = oldProfile.Name
	}
	
	if prof.Organization == false {
		prof.Organization = oldProfile.Organization
	}

	_, err = repo.conn.Exec(context.Background(), updateProfilebyID, &prof.ID, &prof.Name)
	if err != nil {
		c.JSON(500, err)
		return
	}

	c.JSON(200, &prof)
}


func getTask(c *gin.Context) {
	taskID, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		log.Printf("[GET TASK 1] %v | %v", taskID, err)
		c.JSON(500, err)
		return
	}

	task := &Task{}
	err = repo.conn.QueryRow(context.Background(), selectTaskByID, taskID).Scan(
		&task.ID, 
		&task.CreatedBy, 
		&task.DateToComplete,
		&task.TaskType,
		&task.TimeToComplete,
		&task.Lat,
		&task.Long,
		&task.Reward,
		&task.Description,
	)
	if err != nil {
		log.Printf("[GET TASK 2] %v | %v", taskID, err)
		c.JSON(500, err)
		return
	}
	c.JSON(200, &task)

}

func postTask(c *gin.Context) {
	var reqTask *Task
	err := c.Bind(&reqTask)
	if err != nil {
		log.Printf("[POST TASK] %v", err)
		c.JSON(501, err)
		return
	}

	reqTask.CreatedBy = c.GetString("uid")

	acceptedProfile := &Profile{}
	err = repo.conn.QueryRow(context.Background(), selectProfileByID, &reqTask.CreatedBy).Scan(&acceptedProfile.ID, &acceptedProfile.Name, &acceptedProfile.Coins, &acceptedProfile.Organization)
	if err != nil {
		log.Printf("[POST TASK] %v", err)
		c.JSON(500, err)
		return
	}

	if acceptedProfile.Coins > reqTask.Reward {
		log.Printf("[POST TASK] not enough coins")
		c.JSON(500, err)
		return
	}

	_, err = repo.conn.Exec(context.Background(), postTaskQuery, 
		&reqTask.CreatedBy, 
		&reqTask.DateToComplete,
		&reqTask.TaskType,
		&reqTask.TimeToComplete,
		&reqTask.Long,
		&reqTask.Lat,
		&reqTask.Reward,
		&reqTask.Description,
		1,		
	)
	if err != nil {
		log.Printf("[POST TASK] %v", err)
		c.JSON(500, err)
		return
	}

	updatedProfile := &Profile{}
	updatedProfile.ID = c.GetString("uid")
	updatedProfile.Coins = acceptedProfile.Coins - reqTask.Reward
	_, err = repo.conn.Exec(context.Background(), updateProfilebyID, &updatedProfile.ID, &updatedProfile.Name, &updatedProfile.Coins, &updatedProfile.Organization)
	if err != nil {
		log.Printf("[POST TASK] %v", err)
		c.JSON(500, err)
		return
	}

	
	c.JSON(200, true)
}

func acceptTask(c *gin.Context) {
	var reqTask *TasksAccepted
	err := c.Bind(&reqTask)
	if err != nil {
		log.Printf("[ACCEPT TASK] %v", err)
		c.JSON(501, err)
		return
	}

	reqTask.UID = c.GetString("uid")

	_, err = repo.conn.Exec(context.Background(), postAcceptTask, 
		&reqTask.UID, 
		&reqTask.TaskID,
		1,
	)

	if err != nil {
		log.Printf("[ACCEPT TASK] %v", err)
		c.JSON(500, err)
		return
	}
	
	c.JSON(200, true)
	
}

func completeTask(c *gin.Context){
	var taskCompleteRequest *TaskCompleteRequest
	err := c.Bind(&taskCompleteRequest)
	if err != nil {
		log.Printf("[COMPLETE TASK] %v", err)
		c.JSON(501, err)
		return
	}

	taskCompleteRequest.UID = c.GetString("uid")

	existingTask := &Task{}

	err = repo.conn.QueryRow(context.Background(), selectTaskByID, taskCompleteRequest.TaskID).Scan(&existingTask.ID, &existingTask.CreatedBy, &existingTask.DateToComplete, &existingTask.TaskType, &existingTask.TimeToComplete, &existingTask.Lat, &existingTask.Long, &existingTask.Reward, &existingTask.Description)
	if err != nil { 
		log.Printf("[COMPLETE TASK] %v", err)
		c.JSON(500, err)
		return
	}

	_, err = repo.conn.Exec(context.Background(), addRewardByID, taskCompleteRequest.UID, existingTask.Reward)
	if err != nil {
		log.Printf("[COMPLETE TASK] %v", err)
		c.JSON(500, err)
		return
	}

	_, err = repo.conn.Exec(context.Background(), completeTaskByID, taskCompleteRequest.TaskID, 2)
	if err != nil {
		log.Printf("[COMPLETE TASK] %v", err)
		c.JSON(500, err)
		return
	}
	
	
	// Get user who completed the task by their uid, update task in task table, give suer reward.
}
func deleteTask(c *gin.Context) {
	taskID, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		log.Printf("[DELETE TASK 1] %v | %v", taskID, err)
		c.JSON(500, err)
		return
	}

	if err != nil {
		log.Printf("[DELETE TASK 2] %v", err)
		c.JSON(501, err)
		return
	}
	
	_, err = repo.conn.Exec(context.Background(), deleteTaskByID, c.GetString("uid"), taskID)
	if err != nil {
		log.Printf("[DELETE TASK 3] %v", err)
		c.JSON(500, err)
		return
	}

	c.JSON(200, true)
}

func getTasks(c *gin.Context) {
	rows, err := repo.conn.Query(context.Background(), getTasksQuery)
	if err != nil {
		log.Printf("[GET TASKS] %v", err)
		c.JSON(500, err)
		return
	}

	tasks := make([]*Task, 0)

	for rows.Next() {
		task := &Task{}
		err = rows.Scan(
			&task.ID,
			&task.CreatedBy, 
			&task.DateToComplete,
			&task.TaskType,
			&task.TimeToComplete,
			&task.Lat,
			&task.Long,
			&task.Reward,
			&task.Description,
		)
		if err != nil {
			log.Printf("[GET TASKS 2] %v", err)
			c.JSON(501, err)
			return
		}
		tasks = append(tasks, task)
	}

	c.JSON(200, &tasks)
}

func verifiedOrganization(c *gin.Context){
	uid := c.GetString("uid")

	_, err := repo.conn.Exec(context.Background(), promoteToOrg, uid, 100, true)
	if err != nil {
		c.JSON(500, err)
		return
	}

	c.JSON(200, true)
}


const (
	promoteToOrg = "UPDATE profiles SET coins = coins + $2, organization = $3 WHERE uid = $1"
	completeTaskByID = "UPDATE tasks_accepted SET status = $2 WHERE task_id = $1"
	addRewardByID = "UPDATE profiles SET coins = coins + $2 WHERE uid = $1"
	selectProfileByID = "SELECT name, coins, organization FROM profiles WHERE uid = $1"
	updateProfilebyID = "INSERT INTO profiles (uid, name, coins, organization) "+
	"VALUES ($1, $2, 0, false) " +
	"ON CONFLICT (uid)" +
	"DO UPDATE SET name = $2;"
	selectTaskByID = "SELECT * FROM tasks WHERE id = $1;"
	updateTaskByID = "UPDATE tasks SET reward = $2 WHERE uid = $1;"
	postTaskQuery = "INSERT INTO tasks (created_by, date_to_complete, task_type, time_to_complete, lat, long, reward, description) VALUES ($1,$2,$3,$4,$5,$6,$7,$8);"
	deleteTaskByID = "DELETE FROM tasks WHERE created_by = $1 and id = $2;"
	getTasksQuery = "SELECT * FROM tasks WHERE ID NOT IN (SELECT task_id FROM tasks_accepted);"
	selectAcceptedTask = "SELECT uid, task_id FROM tasks_accepted WHERE task_id = $1"
	postAcceptTask = "INSERT INTO tasks_accepted VALUES ($1,$2,$3)"
)