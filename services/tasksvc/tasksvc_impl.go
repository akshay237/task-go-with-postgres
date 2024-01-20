package tasksvc

func (o *TaskSvcImpl) AddTask(req *CreateTaskReq) (*TaskResp, error) {
	conn, err := o.postgres.GetDBInstance()
	if err != nil {
		o.logger.Println("Error in getting postgres connection is", err)
		return nil, err
	}

	query := `insert into tasks(task_title, task_description, task_priority, task_due_datetime, contact_name) values (?, ?, ?, ?, ?);`
	res, err := conn.Exec(query, req.TaskTitle, req.TaskDescription, req.TaskPriority, req.TaskDueDatetime, req.ContactName)
	if err != nil {
		o.logger.Println("Error in eexcuting insert query is", err)
		return nil, err
	}
	nrows, err := res.RowsAffected()
	if nrows == 0 || err != nil {
		o.logger.Println("Error in insert is", err, "and rows affected are", nrows)
		return nil, err
	}
	taskid, err := res.LastInsertId()
	if err != nil {
		o.logger.Println("Error in getting lastid", err)
		return nil, err
	}
	return o.GetTask(&GetTaskReq{TaskId: int(taskid)})
}

func (o *TaskSvcImpl) GetTask(req *GetTaskReq) (*TaskResp, error) {
	conn, err := o.postgres.GetDBInstance()
	if err != nil {
		o.logger.Println("Error in getting postgres connection is", err)
		return nil, err
	}

	resp := &TaskResp{}
	query := `select task_id, task_title, task_description, task_priority, task_due_datetime, contact_name, is_deleted from tasks where task_id = ? and is_deleted != 1;`
	err = conn.QueryRow(query, req.TaskId).Scan(&resp.TaskId, &resp.TaskTitle, &resp.TaskDescription, &resp.TaskPriority, &resp.TaskDueDatetime, &resp.ContactName, &resp.IsDeleted)
	if err != nil {
		o.logger.Println("Scan error", err)
		return nil, err
	}
	return resp, nil
}

func (o *TaskSvcImpl) UpdateTask(req *UpdateTaskRequest) (*TaskResp, error) {
	conn, err := o.postgres.GetDBInstance()
	if err != nil {
		o.logger.Println("Error in getting postgres connection is", err)
		return nil, err
	}

	taskInfo, err := o.GetTask(&GetTaskReq{TaskId: req.TaskId})
	if err != nil {
		o.logger.Println("Error in getting task info", err)
		return nil, err
	}

	if req.TaskTitle == nil {
		req.TaskTitle = &taskInfo.TaskTitle
	}

	if req.TaskDescription == nil {
		req.TaskDescription = &taskInfo.TaskDescription
	}

	if req.TaskPriority == nil {
		req.TaskPriority = &taskInfo.TaskPriority
	}

	if req.TaskDueDatetime == nil {
		req.TaskDueDatetime = &taskInfo.TaskDueDatetime
	}

	if req.ContactName == nil {
		req.ContactName = &taskInfo.ContactName
	}

	query := `update tasks set task_title=? and task_description= ? and task_priority=? and task_due_datetime=? and contact_name = ? where task_id = ?;`
	res, err := conn.Exec(query, req.TaskTitle, req.TaskDescription, req.TaskPriority, req.TaskDueDatetime, req.ContactName, req.TaskId)
	if err != nil {
		o.logger.Println("Error in updating task", err)
		return nil, err
	}

	nrows, err := res.RowsAffected()
	if nrows == 0 || err != nil {
		o.logger.Println("Error in updating task", err)
		return nil, err
	}

	return o.GetTask(&GetTaskReq{TaskId: req.TaskId})
}

func (o *TaskSvcImpl) DeleteTask(req *DeleteTaskReq) (*DeleteTaskResp, error) {
	conn, err := o.postgres.GetDBInstance()
	if err != nil {
		o.logger.Println("Error in getting postgres connection is", err)
		return nil, err
	}

	query := `update tasks set is_deleted = 1 where task_id = ?;`
	res, err := conn.Exec(query, req.TaskId)
	if err != nil {
		o.logger.Println("Error in updating task", err)
		return nil, err
	}

	nrows, err := res.RowsAffected()
	if nrows == 0 || err != nil {
		o.logger.Println("Error in updating task", err)
		return nil, err
	}

	return &DeleteTaskResp{TaskId: req.TaskId, IsDeleted: true}, nil
}

func (o *TaskSvcImpl) GetAllTasks() ([]*TaskResp, error) {
	conn, err := o.postgres.GetDBInstance()
	if err != nil {
		o.logger.Println("Error in getting postgres connection is", err)
		return nil, err
	}

	rsp := []*TaskResp{}
	query := `select task_id, task_title, task_description, task_priority, task_due_datetime, contact_name, is_deleted from tasks where  is_deleted != 1;`
	rows, err := conn.Query(query)
	if err != nil {
		o.logger.Println("Error in getting task info", err)
		return nil, err
	}

	for rows.Next() {
		resp := &TaskResp{}
		scanerr := rows.Scan(&resp.TaskId, &resp.TaskTitle, &resp.TaskDescription, &resp.TaskPriority, &resp.TaskDueDatetime, &resp.ContactName, &resp.IsDeleted)
		if scanerr != nil {
			o.logger.Println("Error in scan", scanerr)
			continue
		}
		rsp = append(rsp, resp)
	}
	return rsp, nil
}
