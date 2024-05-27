package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Blxssy/YandexLyceum-2sprint/agent/pkg"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	computingPower, err := strconv.Atoi(os.Getenv("COMPUTING_POWER"))
	if err != nil {
		computingPower = 1
	}

	for i := 0; i < computingPower; i++ {
		go worker()
	}

	select {}
}

func worker() {
	for {
		task, err := getTask()
		if err != nil {
			time.Sleep(2 * time.Second)
			continue
		}
		result := executeTask(task)
		submitResult(result)
	}
}

func getTask() (*pkg.Task, error) {
	resp, err := http.Get("http://localhost:8080/internal/task")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("no tasks available")
	}

	var response struct {
		Task pkg.Task `json:"task"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	return &response.Task, nil
}

func executeTask(task *pkg.Task) *pkg.TaskResult {
	time.Sleep(time.Duration(task.OperationTime) * time.Millisecond)
	var result float64
	switch task.Operation {
	case "+":
		result = task.Arg1 + task.Arg2
	case "-":
		result = task.Arg1 - task.Arg2
	case "*":
		result = task.Arg1 * task.Arg2
	case "/":
		result = task.Arg1 / task.Arg2
	}
	return &pkg.TaskResult{ID: task.ID, Result: result}
}

func submitResult(result *pkg.TaskResult) {
	data, _ := json.Marshal(result)
	resp, err := http.Post("http://localhost:8080/internal/task/result", "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("Failed to submit result:", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Failed to submit result, status code:", resp.StatusCode)
	}
}
