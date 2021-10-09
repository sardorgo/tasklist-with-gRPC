package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/satori/uuid"

	pbStudent "task/student_proto"
	pbTask "task/task_proto"

	"google.golang.org/grpc"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "sardor"
	password = "sardor"
	dbname   = "tasklist"
)

type server struct {
	conn *sql.DB
	pbStudent.UnimplementedStudentsServer
	pbTask.UnimplementedTasksServer
}

func (connection *server) CreateTask(ctx context.Context, req *pbTask.CreateTaskRequest) (*pbTask.Task, error) {
	db := connection.conn
	id := uuid.NewV4()

	req.Task.Id = id.String()
	taskName := req.GetTask().GetName()

	sqlInsert := `insert into tasks (task_id, task_name) values ($1, $2);`

	if _, err := db.Exec(sqlInsert, id, taskName); err != nil {
		return nil, errors.Wrapf(err, "Task couldn't be inserted")
	}

	return req.Task, nil
}

func (connection *server) DeleteTask(ctx context.Context, req *pbTask.DeleteTaskRequest) (*pbTask.Empty, error) {
	db := connection.conn
	sqlStatement := `delete from tasks where task_id = $1`

	id := req.GetId()

	if _, err := db.Exec(sqlStatement, id); err != nil {
		return nil, err
	}

	return &pbTask.Empty{}, nil
}

func (connection *server) GetTaskById(ctx context.Context, req *pbTask.GetTaskByIdRequest) (*pbTask.Task, error) {
	db := connection.conn
	id := req.GetId()

	sqlStatement := `select * from tasks where task_id = $1;`

	var taskName, taskId string

	err := db.QueryRow(sqlStatement, id).Scan(&taskId, &taskName)

	if err != nil {
		errors.Wrapf(err, "Task couldn't be selected")
	}

	res := &pbTask.Task{
		Id:   taskId,
		Name: taskName,
	}

	return res, nil
}

// func (connection *server) GetIdAndRespondTask(ctx context.Context, req *pbTask.GetIdAndRespondTaskRequest) (*pbTask.TaskResponse, error) {
// 	db := connection.conn
// 	id := req.GetId()

// 	var taskName string

// 	sqlStatement := `select task_name from tasks where task_id = $1;`
// 	err := db.QueryRow(sqlStatement, id).Scan(&taskName)

// 	if err != nil {
// 		errors.Wrapf(err, "Task couldn't be selected")
// 	}

// 	res := &pbTask.TaskResponse{
// 		TaskName: taskName,
// 	}

// 	return res, nil
// }

func main() {
	lis, err := net.Listen("tcp", ":5500")

	if err != nil {
		log.Fatalf("%v", err)
	}

	s := grpc.NewServer()

	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		errors.Wrap(err, "StudentProfile couldn't be returned")
	}
	defer db.Close()

	err = db.Ping()

	if err != nil {
		errors.Wrap(err, "Student couldn't be listed")
	}

	pbTask.RegisterTasksServer(s, &server{conn: db})

	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		errors.Wrap(err, "Task couldn't be returned")
	}
}
