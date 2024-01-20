CREATE TABLE tasks (
    task_id SERIAL PRIMARY KEY,
    task_title VARCHAR(255) NOT NULL,
    task_description TEXT,
    task_priority VARCHAR(50),
    task_due_datetime TIMESTAMP NOT NULL,
    contact_name VARCHAR(100) NOT NULL,
    is_deleted boolean DEFAULT FALSE
);


