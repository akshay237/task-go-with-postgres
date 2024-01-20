package tasksvc

type CreateTaskReq struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	Priority    string `json:"priority" validate:"required"`
	TimeStamp   int64  `json:"timestamp" validate:"required"`
}

type UpdateTaskRequest struct {
	TaskId      int    `json:"taskid" validate:"required"`
	Description string `json:"description" validate:"required"`
	Priority    string `json:"priority" validate:"required"`
	TimeStamp   int64  `json:"timestamp" validate:"required"`
}

type GetTaskReq struct {
	TaskId int `json:"taskid" validate:"required"`
}

type DeleteTaskReq struct {
	TaskId int `json:"taskid" validate:"required"`
}

// Response Struct
type TaskResp struct {
	TaskId      int    `json:"taskid"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	TimeStamp   int64  `json:"timestamp"`
}

type DeleteTaskResp struct {
	TaskId    int  `json:"taskid"`
	IsDeleted bool `json:"isdeleted"`
}
