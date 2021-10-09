package main

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"github.com/pkg/errors"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "sardor"
	password = "sardor"
	dbname   = "tasklist"
)

func main() {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		errors.Wrap(err, "StudentProfile couldn't be returned")
	}
	defer db.Close()

	var taskName []sql.NullString

	sqlStatement := `select
						array_agg(t.task_name)

					from users as u
					join task_controller as tc on tc.user_id = u.user_id
					join tasks as t on tc.task_id = t.task_id
					where tc.user_id = 'd767d830-9759-4ac8-9c36-7fcbbd6bcab2'
					group by tc.task_id, u.user_id;
					`

	if err := db.QueryRow(sqlStatement).Scan(pq.Array(&taskName)); err != nil {
		panic(err)
	}

	fmt.Println(taskName)

}

// u := pbStudent.Student{
// 	Id: studentId,
// 	FirstName: firstName,
// }

// t := pbStudent.RepeatedTask{
// 	Task: tasks,
// }

// final := pbStudent.GetStudentsTasksResponse{
// 	Student: &pbStudent.Student{
// 		Id: u.Id,
// 		FirstName: u.FirstName,
// 	},
// 	Task: []*pbStudent.Task{

// 	},
// }
