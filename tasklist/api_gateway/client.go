package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	pbStudent "task/student_proto"
	pbTask "task/task_proto"

	"google.golang.org/grpc"
)

var studentConn pbStudent.StudentsClient
var taskConn pbTask.TasksClient

func PostStudent(c *gin.Context) {
	var newStudent pbStudent.Student

	if err := c.ShouldBindJSON(&newStudent); err != nil {
		return
	}

	req := &pbStudent.CreateStudentRequest{
		Student: &pbStudent.Student{
			FirstName: newStudent.FirstName,
			Id:        newStudent.Id,
		},
	}

	res, err := studentConn.CreateStudent(context.Background(), req)

	if err != nil {
		log.Fatalf("%v", err)
	}

	c.IndentedJSON(http.StatusCreated, res)
}

func GETStudentById(c *gin.Context) {
	id := c.Param("id")

	req := &pbStudent.GetStudentByIdRequest{
		Id: id,
	}

	res, err := studentConn.GetStudentById(context.Background(), req)

	if err != nil {
		log.Fatalf("%v", err)
	}

	c.IndentedJSON(http.StatusOK, res)
}

func PutStudent(c *gin.Context) {
	var updateWanted pbStudent.Student

	if err := c.ShouldBindJSON(&updateWanted); err != nil {
		return
	}
	req := &pbStudent.UpdateStudentRequest{
		Student: &pbStudent.Student{
			FirstName: updateWanted.FirstName,
			Id:        updateWanted.Id,
		},
	}
	res, err := studentConn.UpdateStudent(context.Background(), req)
	if err != nil {
		panic(err)
	}
	c.IndentedJSON(http.StatusOK, res)
}

func DELETEStudent(c *gin.Context) {
	id := c.Param("id")
	req := &pbStudent.DeleteStudentRequest{
		Id: id,
	}
	res, err := studentConn.DeleteStudent(context.Background(), req)

	if err != nil {
		panic(err)
	}
	c.IndentedJSON(http.StatusOK, res)
}

func ListAll(c *gin.Context) {
	req := &pbStudent.GetAllStudentsRequest{}
	res, err := studentConn.GetAllStudents(context.Background(), req)

	if err != nil {
		log.Fatalf("%v", err)
	}

	c.IndentedJSON(http.StatusOK, res)
}

func PostTask(c *gin.Context) {
	var newTask pbTask.Task

	if err := c.ShouldBindJSON(&newTask); err != nil {
		return
	}

	req := &pbTask.CreateTaskRequest{
		Task: &pbTask.Task{
			Name: newTask.Name,
			Id:   newTask.Id,
		},
	}

	res, err := taskConn.CreateTask(context.Background(), req)

	if err != nil {
		log.Fatalf("%v", err)
	}

	c.IndentedJSON(http.StatusCreated, res)
}

func GETTaskById(c *gin.Context) {
	id := c.Param("id")

	req := &pbTask.GetTaskByIdRequest{
		Id: id,
	}

	res, err := taskConn.GetTaskById(context.Background(), req)

	if err != nil {
		log.Fatalf("%v", err)
	}

	c.IndentedJSON(http.StatusOK, res)
}

func DELETETask(c *gin.Context) {
	id := c.Param("id")

	req := &pbTask.DeleteTaskRequest{
		Id: id,
	}

	res, err := taskConn.DeleteTask(context.Background(), req)

	if err != nil {
		log.Fatalf("%v", err)
	}

	c.IndentedJSON(http.StatusOK, res)
}

func PostTaskController(c *gin.Context) {
	var newTaskForStudent pbStudent.CreateTaskControllerRequest

	if err := c.ShouldBindJSON(&newTaskForStudent); err != nil {
		log.Fatalf("%v", err)
	}

	req := &pbStudent.CreateTaskControllerRequest{
		StudentId: newTaskForStudent.StudentId,
		TaskId:    newTaskForStudent.TaskId,
	}

	res, err := studentConn.CreateTaskController(context.Background(), req)

	if err != nil {
		log.Fatalf("%v", err)
	}

	c.IndentedJSON(http.StatusCreated, res)
}

func GetStudentTasks(c *gin.Context) {
	id := c.Param("id")

	req := &pbStudent.GetStudentsTasksRequest{
		StudentId: id,
	}

	res, err := studentConn.GetStudentTask(context.Background(), req)

	if err != nil {
		panic(err)
	}
	c.IndentedJSON(http.StatusOK, res)
}

func GetAllStudentTasks(c *gin.Context) {
	req := &pbStudent.ListAllUsersTasksRequest{}
	res, err := studentConn.ListAllUsersTasks(context.Background(), req)

	if err != nil {
		log.Fatalf("%v", err)
	}

	log.Printf("%v", res)

	c.IndentedJSON(http.StatusOK, res)
}

func main() {
	// Student Client

	conn, err := grpc.Dial("localhost:5000", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("%v", err)
	}
	defer conn.Close()
	studentConn = pbStudent.NewStudentsClient(conn)

	// Task Client

	conn1, err := grpc.Dial("localhost:5500", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("%v", err)
	}
	defer conn.Close()

	taskConn = pbTask.NewTasksClient(conn1)

	// --------------------------------------------------------------

	router := gin.Default()

	router.POST("/students", PostStudent)
	router.POST("/students/task", PostTaskController)
	router.GET("/students/:id", GETStudentById)
	router.GET("/students", ListAll)
	router.PUT("/students", PutStudent)
	router.DELETE("/students/:id", DELETEStudent)
	router.GET("/students/tasks/:id", GetStudentTasks)
	router.GET("/students/tasks/all", GetAllStudentTasks)

	router.POST("/tasks", PostTask)
	router.GET("/tasks/:id", GETTaskById)
	router.DELETE("/tasks/:id", DELETETask)

	router.Run("localhost:9000")
}
