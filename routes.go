package main

func getProfile(c *gin.Context) {
	var prof *Profile
	err := c.Bind(&prof)
	if err != nil {
		c.JSON(501, err)
		return
	}
	err = repo.conn.QueryRow(context.Background(), selectProfileByID, &prof.ID).Scan(&prof.Lat, &prof.Long, &prof.Interests)
	if err != nil {
		c.JSON(500, err)
		return
	}
	c.JSON(200, &prof)
}

func getTask(c *gin.Context) {
}

func postTask(c *gin.Context) {
}

func deleteTask(c *gin.Context) {
}

func getTasks(c *gin.Context) {
}

func postDonate(c *gin.Context) {
}


