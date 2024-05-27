package main

import (
	"github.com/Blxssy/YandexLyceum-2sprint/orchestrator/internal/orchestrator"
	"net/http"
)

func main() {
	//router := gin.Engine{}

	o := orchestrator.NewOrchestrator()

	http.HandleFunc("/api/v1/calculate", o.AddExpression)
	http.HandleFunc("/api/v1/expressions", o.GetExpressions)
	http.HandleFunc("/api/v1/expressions/", o.GetExpressionByID)
	http.HandleFunc("/internal/task", o.GetTask)
	http.HandleFunc("/internal/task/result", o.SubmitTaskResult)

	http.ListenAndServe(":8080", nil)
}
