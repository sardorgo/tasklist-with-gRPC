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

func (connection *server) CreateStudent(ctx context.Context, req *pbStudent.CreateStudentRequest) (*pbStudent.Student, error) {
	db := connection.conn
	id := uuid.NewV4()

	req.Student.Id = id.String()
	firstName := req.GetStudent().GetFirstName()

	sqlInsert := `insert into users (user_id, first_name) values ($1, $2);`

	if _, err := db.Exec(sqlInsert, id, firstName); err != nil {
		return nil, errors.Wrapf(err, "Student couldn't be inserted")
	}

	return req.Student, nil
}

func (connection *server) UpdateStudent(ctx context.Context, req *pbStudent.UpdateStudentRequest) (*pbStudent.Student, error) {
	db := connection.conn
	sqlStatement := `update users set first_name = $2 where user_id = $1`

	if _, err := db.Exec(sqlStatement, req.Student.Id, req.Student.FirstName); err != nil {
		return nil, err
	}

	return req.Student, nil
}

func (connection *server) DeleteStudent(ctx context.Context, req *pbStudent.DeleteStudentRequest) (*pbStudent.Empty, error) {
	db := connection.conn
	sqlStatement := `delete from users where user_id = $1`

	id := req.GetId()

	if _, err := db.Exec(sqlStatement, id); err != nil {
		return nil, err
	}

	return &pbStudent.Empty{}, nil
}

func (connection *server) GetStudentById(ctx context.Context, req *pbStudent.GetStudentByIdRequest) (*pbStudent.Student, error) {
	db := connection.conn
	id := req.GetId()

	sqlStatement := `select first_name from users where user_id = $1;`

	var firstName string

	err := db.QueryRow(sqlStatement, id).Scan(&firstName)

	if err != nil {
		errors.Wrapf(err, "Student couldn't be selected")
	}

	res := &pbStudent.Student{
		FirstName: firstName,
	}

	return res, nil
}

func (connection *server) GetAllStudents(ctx context.Context, req *pbStudent.GetAllStudentsRequest) (*pbStudent.GetAllStudentsResponse, error) {
	db := connection.conn
	sqlStatement := `select user_id, first_name from users`
	result, err := db.Query(sqlStatement)

	if err != nil {
		panic(err)
	}

	defer result.Close()

	res := []*pbStudent.Student{}
	for result.Next() {
		var firstName, id string
		if err = result.Scan(&id, &firstName); err != nil {
			errors.Wrap(err, "Student couln't be listed")
		}
		u := pbStudent.Student{
			Id:        id,
			FirstName: firstName,
		}
		res = append(res, &u)
	}
	ans := pbStudent.GetAllStudentsResponse{Student: res}
	return &ans, nil
}

func (connection *server) CreateTaskController(ctx context.Context, req *pbStudent.CreateTaskControllerRequest) (*pbStudent.Empty, error) {
	db := connection.conn
	id := uuid.NewV4()

	studentId := req.GetStudentId()
	taskId := req.GetTaskId()

	sqlInsert := `insert into task_controller(task_controller_id, user_id, task_id) values ($1, $2, $3);`

	if _, err := db.Exec(sqlInsert, id, studentId, taskId); err != nil {
		return nil, err
	}

	return &pbStudent.Empty{}, nil
}

func (connection *server) GetStudentTask(ctx context.Context, req *pbStudent.GetStudentsTasksRequest) (*pbStudent.GetStudentsTasksResponse, error) {
	db := connection.conn
	id := req.GetStudentId()

	// var ids []string

	sqlStatement := `select 
						task_id 
					from task_controller
					where user_id = $1;`

	rows, err := db.Query(sqlStatement, id)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var firstName string

	sqlSelect := `select first_name from users where user_id = $1;`

	err = db.QueryRow(sqlSelect, id).Scan(&firstName)

	if err != nil {
		log.Fatal(err)
	}

	var tasks []*pbStudent.Task

	for rows.Next() {
		var id string
		var taskName string
		if err = rows.Scan(&id); err != nil {
			errors.Wrap(err, "Id couln't be listed")
		}

		sqlStatement1 := `select
							task_name
						from
							tasks
						where task_id = $1`

		err1 := db.QueryRow(sqlStatement1, id).Scan(&taskName)

		if err1 != nil {
			log.Fatal(err1)
		}

		result := pbStudent.Task{
			Name: taskName,
		}

		tasks = append(tasks, &result)
	}

	resStudent := &pbStudent.Student{
		Id:        id,
		FirstName: firstName,
	}

	finalResult := pbStudent.GetStudentsTasksResponse{Student: resStudent, Task: tasks}

	return &finalResult, nil
}

func (connection *server) ListAllUsersTasks(ctx context.Context, req *pbStudent.ListAllUsersTasksRequest) (*pbStudent.ListAllUsersTasksResponse, error) {
	db := connection.conn
	var result []*pbStudent.GetStudentsTasksResponse

	sqlStatementForUserId := `
					select 
						distinct tc.user_id
					from 
						task_controller as tc
					where tc.task_id is not null;
					`

	rows, err := db.Query(sqlStatementForUserId)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		var tasks []*pbStudent.Task
		var studentId, firstName string

		if err = rows.Scan(&studentId); err != nil {
			errors.Wrap(err, "User couln't be listed")
		}

		sqlStatementForUsersFirstName := `select first_name from users where user_id = $1;`

		err = db.QueryRow(sqlStatementForUsersFirstName, studentId).Scan(&firstName)

		if err != nil {
			log.Fatal(err)
		}

		sqlStatementForTaskId := `select 
									tc.task_id 
								from task_controller as tc
								where tc.user_id = $1 and tc.user_id is not null;
								`

		rows1, err := db.Query(sqlStatementForTaskId, studentId)

		if err != nil {
			log.Fatal(err)
		}

		defer rows1.Close()

		for rows1.Next() {
			var taskId, taskName string

			if err = rows1.Scan(&taskId); err != nil {
				errors.Wrap(err, "TaskId Couldn't be selected")
			}

			sqlStatementForTaskName := `select task_name from tasks where task_id = $1;`

			err := db.QueryRow(sqlStatementForTaskName, taskId).Scan(&taskName)

			if err != nil {
				log.Fatal(err)
			}

			task := pbStudent.Task{
				Name: taskName,
			}

			tasks = append(tasks, &task)
		}

		res := pbStudent.GetStudentsTasksResponse{

			Student: &pbStudent.Student{
				Id:        studentId,
				FirstName: firstName,
			},

			Task: tasks,
		}

		result = append(result, &res)
	}

	ans := pbStudent.ListAllUsersTasksResponse{
		StudentTasks: result,
	}

	return &ans, nil
}

func main() {
	lis, err := net.Listen("tcp", ":5000")

	if err != nil {
		errors.Wrapf(err, "Student couldn't be returned")
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

	pbStudent.RegisterStudentsServer(s, &server{conn: db})
	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		errors.Wrap(err, "Student couldn't be returned")
	}

}
