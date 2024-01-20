package taskhdlr

import (
	"net/http"
	"task-go-with-postgres/apiserver"
	"task-go-with-postgres/services/tasksvc"
)

func (o *TaskHdlr) addTask(w http.ResponseWriter, r *http.Request) {
	var req *tasksvc.CreateTaskReq
	errData, err := apiserver.JsonBodyParser(r.Body, &req)
	if err != nil {
		o.logger.Println("ERR : ", err)
		apiserver.APIResponseBadRequest(w, r, "MALFORMED_REQUEST", errData, err.Error())
		return
	}

	apiresp, apierr := o.tasksvc.AddTask(req)
	if apierr != nil {
		o.logger.Println("ERR : ", apierr)
		apiserver.APIResponseInternalServerError(w, r, "MALFORMED_REQUEST", nil, apierr.Error())
		return
	}
	apiserver.APIResponseOK(w, r, apiresp, "added task successfully")

}

func (o *TaskHdlr) getTask(w http.ResponseWriter, r *http.Request) {
	var req *tasksvc.GetTaskReq
	errData, err := apiserver.JsonBodyParser(r.Body, &req)
	if err != nil {
		o.logger.Println("ERR : ", err)
		apiserver.APIResponseBadRequest(w, r, "MALFORMED_REQUEST", errData, err.Error())
		return
	}

	apiresp, apierr := o.tasksvc.GetTask(req)
	if apierr != nil {
		o.logger.Println("ERR : ", apierr)
		apiserver.APIResponseInternalServerError(w, r, "MALFORMED_REQUEST", nil, apierr.Error())
		return
	}
	apiserver.APIResponseOK(w, r, apiresp, "fetched task info successfully")
}

func (o *TaskHdlr) updateTask(w http.ResponseWriter, r *http.Request) {
	var req *tasksvc.UpdateTaskRequest
	errData, err := apiserver.JsonBodyParser(r.Body, &req)
	if err != nil {
		o.logger.Println("ERR : ", err)
		apiserver.APIResponseBadRequest(w, r, "MALFORMED_REQUEST", errData, err.Error())
		return
	}

	apiresp, apierr := o.tasksvc.UpdateTask(req)
	if apierr != nil {
		o.logger.Println("ERR : ", apierr)
		apiserver.APIResponseInternalServerError(w, r, "MALFORMED_REQUEST", nil, apierr.Error())
		return
	}
	apiserver.APIResponseOK(w, r, apiresp, "updated task info successfully")
}

func (o *TaskHdlr) deleteTask(w http.ResponseWriter, r *http.Request) {
	var req *tasksvc.DeleteTaskReq
	errData, err := apiserver.JsonBodyParser(r.Body, &req)
	if err != nil {
		o.logger.Println("ERR : ", err)
		apiserver.APIResponseBadRequest(w, r, "MALFORMED_REQUEST", errData, err.Error())
		return
	}

	apiresp, apierr := o.tasksvc.DeleteTask(req)
	if apierr != nil {
		o.logger.Println("ERR : ", apierr)
		apiserver.APIResponseInternalServerError(w, r, "MALFORMED_REQUEST", nil, apierr.Error())
		return
	}
	apiserver.APIResponseOK(w, r, apiresp, "Deleted task successfully")
}

func (o *TaskHdlr) getAllTasks(w http.ResponseWriter, r *http.Request) {

	apiresp, apierr := o.tasksvc.GetAllTasks()
	if apierr != nil {
		o.logger.Println("ERR : ", apierr)
		apiserver.APIResponseInternalServerError(w, r, "MALFORMED_REQUEST", nil, apierr.Error())
		return
	}
	apiserver.APIResponseOK(w, r, apiresp, "added task successfully")
}
