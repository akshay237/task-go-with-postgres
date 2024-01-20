package tasksvc

import (
	"log"
	"task-go-with-postgres/database"
)

type TaskSvc interface {
	AddTask(req *CreateTaskReq) (*TaskResp, error)
	GetTask(req *GetTaskReq) (*TaskResp, error)
	UpdateTask(req *UpdateTaskRequest) (*TaskResp, error)
	DeleteTask(req *DeleteTaskReq) (*DeleteTaskResp, error)
	GetAllTasks() ([]*TaskResp, error)
}

type TaskSvcImpl struct {
	postgres *database.PostgresDBService
	logger   *log.Logger
}

func NewTaskSvc(postgresI *database.PostgresDBService, logger *log.Logger) TaskSvc {
	return &TaskSvcImpl{
		postgres: postgresI,
		logger:   logger,
	}
}
