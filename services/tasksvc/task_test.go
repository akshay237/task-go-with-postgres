package tasksvc

import (
	"encoding/json"
	"log"
	"os"
	"task-go-with-postgres/database"
	"testing"
	"time"
)

func getPostgresConfig(t *testing.T) *database.ConfigPostgres {
	config := &database.ConfigPostgres{
		Host:     "127.0.0.1",
		User:     "akshay",
		Password: "password",
		Database: "akshay",
	}
	return config
}

func createTaskTest(t *testing.T, taskSvcI TaskSvc, title, desc, priority, name string, timestamep time.Time) {
	req := &CreateTaskReq{
		TaskTitle:       title,
		TaskDescription: desc,
		TaskPriority:    priority,
		ContactName:     name,
		TaskDueDatetime: timestamep,
	}
	taskRsp, err := taskSvcI.AddTask(req)
	if err != nil {
		t.Log("Error is", err)
	}
	taskRspBytes, _ := json.Marshal(taskRsp)
	t.Log("Add Task Resp", string(taskRspBytes))
}

func getTaskTest(t *testing.T, taskSvcI TaskSvc, taskid int) {
	req := &GetTaskReq{
		TaskId: taskid,
	}
	taskRsp, err := taskSvcI.GetTask(req)
	if err != nil {
		t.Log("Error is", err)
	}
	taskRspBytes, _ := json.Marshal(taskRsp)
	t.Log("Get Task Resp", string(taskRspBytes))
}

func updateTaskTest(t *testing.T, taskSvcI TaskSvc, priority string, timestamep time.Time) {
	req := &CreateTaskReq{
		TaskPriority:    priority,
		TaskDueDatetime: timestamep,
	}
	taskRsp, err := taskSvcI.AddTask(req)
	if err != nil {
		t.Log("Error is", err)
	}
	taskRspBytes, _ := json.Marshal(taskRsp)
	t.Log("Add Task Resp", string(taskRspBytes))
}

func getAllTaskTest(t *testing.T, taskSvcI TaskSvc) {

	taskRsp, err := taskSvcI.GetAllTasks()
	if err != nil {
		t.Log("Error is", err)
	}
	taskRspBytes, _ := json.Marshal(taskRsp)
	t.Log("Get All Tasks Resp", string(taskRspBytes))
}

func deleteTaskTest(t *testing.T, taskSvcI TaskSvc, taskid int) {
	req := &DeleteTaskReq{
		TaskId: taskid,
	}
	taskRsp, err := taskSvcI.DeleteTask(req)
	if err != nil {
		t.Log("Error is", err)
	}
	taskRspBytes, _ := json.Marshal(taskRsp)
	t.Log("Deleted Task Resp", string(taskRspBytes))
}

func TestTaskSvc(t *testing.T) {

	postgresConfig := getPostgresConfig(t)
	postgresI := database.NewPostgresDBService(postgresConfig)
	logger := log.New(os.Stdout, "TASKSVC", log.LstdFlags|log.Lmicroseconds)

	taskSvcI := NewTaskSvc(postgresI, logger)

	// call the test functions
	// add or create task
	var timestamp time.Time = time.Now().AddDate(0, 0, 4)
	createTaskTest(t, taskSvcI, "Set Alarm", "Discuss Project Details", "MEDIUM", "Akshay", timestamp)
	getTaskTest(t, taskSvcI, 2)
	timestamp = time.Now().AddDate(0, 0, 2)
	updateTaskTest(t, taskSvcI, "HIGH", timestamp)
	getAllTaskTest(t, taskSvcI)
	//deleteTaskTest(t, taskSvcI, 1)

}
