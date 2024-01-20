# How to run
1. Clone the repo
2. Navigate to the directory.
3. go run main.go --configfile config.json

# Find API curls and Response

# 1. Get Task Info
curl --location 'http://localhost:5000/api/v1/task/get' \
--header 'Content-Type: application/json' \
--data '{
    "task_id": 1
}'

Response:
{
    "err": null,
    "data": {
        "task_title": "Call John",
        "task_id": 1,
        "task_description": "Discuss project details",
        "task_priority": "High",
        "task_due_datetime": "2024-01-23T14:30:00Z",
        "contact_name": "John",
        "is_deleted": false
    },
    "msg": "fetched task info successfully"
}

# 2. Get All Tasks
curl --location --request POST 'http://localhost:5000/api/v1/task/getall'

Response:
{
    "err": null,
    "data": [
        {
            "task_title": "Call John",
            "task_id": 1,
            "task_description": "Discuss project details",
            "task_priority": "High",
            "task_due_datetime": "2024-01-23T14:30:00Z",
            "contact_name": "John",
            "is_deleted": false
        },
        {
            "task_title": "Set Alarm",
            "task_id": 2,
            "task_description": "Discuss Project Details",
            "task_priority": "MEDIUM",
            "task_due_datetime": "2024-01-24T22:03:22.07461Z",
            "contact_name": "Akshay",
            "is_deleted": false
        }
    ],
    "msg": "fetch all tasks successfully"
}

# 3. Update a Task
curl --location 'http://localhost:5000/api/v1/task/update' \
--header 'Content-Type: application/json' \
--data '{
    "task_id": 2,
    "task_priority": "HIGH"
}'

Response:
{
    "err": null,
    "data": {
        "task_title": "Set Alarm",
        "task_id": 2,
        "task_description": "Discuss Project Details",
        "task_priority": "HIGH",
        "task_due_datetime": "2024-01-24T22:03:22.07461Z",
        "contact_name": "Akshay",
        "is_deleted": false
    },
    "msg": "updated task info successfully"
}

# 4. Delete Task
curl --location 'http://localhost:5000/api/v1/task/delete' \
--header 'Content-Type: application/json' \
--data '{
    "task_id":4
}'

Response:
{
    "err": null,
    "data": {
        "task_id": 4,
        "isdeleted": true
    },
    "msg": "Deleted task successfully"
}

# 5. Create Task
curl --location 'http://localhost:5000/api/v1/task/add' \
--header 'Content-Type: application/json' \
--data '{
    "task_title":"Complete Task",
    "task_description": "Implement CRUD operations for the task",
    "task_priority": "HIGH",
    "task_due_datetime": "2024-01-21T13:00:00Z",
    "contact_name":"akshay"
}'

Response:
{
    "err": null,
    "data": {
        "task_title": "Complete Task",
        "task_id": 6,
        "task_description": "Implement CRUD operations for the task",
        "task_priority": "HIGH",
        "task_due_datetime": "2024-01-21T13:00:00Z",
        "contact_name": "akshay",
        "is_deleted": false
    },
    "msg": "added task successfully"
}
