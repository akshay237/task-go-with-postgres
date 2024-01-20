package taskhdlr

import (
	"log"
	"task-go-with-postgres/services/tasksvc"

	"github.com/go-chi/chi/v5"
)

type TaskHdlr struct {
	tasksvc tasksvc.TaskSvc
	logger  *log.Logger
}

func (o *TaskHdlr) RegisterRoutes(routes chi.Router) {

	// CRUD operations are
	routes.Post("/add", o.addTask)
	routes.Post("/get", o.getTask)
	routes.Post("/update", o.updateTask)
	routes.Post("/delete", o.deleteTask)
	routes.Post("/getall", o.getAllTasks)

}

func NewTaskHdlr(taskSvcI tasksvc.TaskSvc, loggerI *log.Logger) *TaskHdlr {
	return &TaskHdlr{
		tasksvc: taskSvcI,
		logger:  loggerI,
	}
}
