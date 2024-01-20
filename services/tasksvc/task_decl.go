package tasksvc

import "time"

type CreateTaskReq struct {
	TaskTitle       string    `json:"task_title" validate:"required"`
	TaskDescription string    `json:"task_description" validate:"required"`
	TaskPriority    string    `json:"task_priority" validate:"required"`
	TaskDueDatetime time.Time `json:"task_due_datetime" validate:"required"`
	ContactName     string    `json:"contact_name" validate:"required"`
}

type UpdateTaskRequest struct {
	TaskId          int        `json:"task_id" validate:"required"`
	TaskTitle       *string    `json:"task_title"`
	TaskDescription *string    `json:"task_description"`
	TaskPriority    *string    `json:"task_priority"`
	TaskDueDatetime *time.Time `json:"task_due_datetime"`
	ContactName     *string    `json:"contact_name"`
}

type GetTaskReq struct {
	TaskId int `json:"taskid" validate:"required"`
}

type DeleteTaskReq struct {
	TaskId int `json:"taskid" validate:"required"`
}

// Response Struct
type DeleteTaskResp struct {
	TaskId    int  `json:"taskid"`
	IsDeleted bool `json:"isdeleted"`
}

type TaskResp struct {
	TaskTitle       string    `json:"task_title"`
	TaskId          int       `json:"task_id"`
	TaskDescription string    `json:"task_description"`
	TaskPriority    string    `json:"task_priority"`
	TaskDueDatetime time.Time `json:"task_due_datetime"`
	ContactName     string    `json:"contact_name"`
	IsDeleted       bool      `json:"is_deleted"`
}
