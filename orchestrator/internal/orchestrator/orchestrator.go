package orchestrator

import (
	"encoding/json"
	"github.com/Blxssy/YandexLyceum-2sprint/orchestrator/pkg"
	"net/http"
	"sync"
)

type Orchestrator struct {
	mu          sync.Mutex
	expressions map[string]*pkg.Expression
	tasks       []*pkg.Task
}

func NewOrchestrator() *Orchestrator {
	return &Orchestrator{
		expressions: make(map[string]*pkg.Expression),
		tasks:       []*pkg.Task{},
	}
}

func (o *Orchestrator) AddExpression(w http.ResponseWriter, r *http.Request) {
	var expr pkg.Expression
	if err := json.NewDecoder(r.Body).Decode(&expr); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	o.mu.Lock()
	defer o.mu.Unlock()
	o.expressions[expr.ID] = &expr

	//TODO:
	// Разбить выражение на задачи и добавить их в очередь
	// Нужно написать логику разбора выражения на отдельные задачи
	// Пример:
	o.tasks = append(o.tasks, &pkg.Task{ID: "1", Arg1: 2, Arg2: 2, Operation: "+", OperationTime: 1000})

	w.WriteHeader(http.StatusCreated)
}

func (o *Orchestrator) GetExpressions(w http.ResponseWriter, r *http.Request) {
	o.mu.Lock()
	defer o.mu.Unlock()
	//TODO:
	// Выводит не все записи
	expressions := make([]*pkg.Expression, 0, len(o.expressions))
	for _, expr := range o.expressions {
		expressions = append(expressions, expr)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"expressions": expressions})
}

func (o *Orchestrator) GetExpressionByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/api/v1/expressions/"):]
	o.mu.Lock()
	defer o.mu.Unlock()
	expr, exists := o.expressions[id]
	if !exists {
		http.Error(w, "Expression not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"expression": expr})
}

func (o *Orchestrator) GetTask(w http.ResponseWriter, r *http.Request) {
	o.mu.Lock()
	defer o.mu.Unlock()
	if len(o.tasks) == 0 {
		http.Error(w, "No tasks available", http.StatusNotFound)
		return
	}
	task := o.tasks[0]
	o.tasks = o.tasks[1:]
	json.NewEncoder(w).Encode(map[string]interface{}{"task": task})
}

func (o *Orchestrator) SubmitTaskResult(w http.ResponseWriter, r *http.Request) {
	var result pkg.TaskResult
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	o.mu.Lock()
	defer o.mu.Unlock()

	//TODO:
	// Сделать обновление состояния выражения по результату задачи
	// Нужно сопоставить задачу и обновить результат выражения
	for _, expr := range o.expressions {
		expr.Status = pkg.StatusCompleted
		expr.Result = result.Result
	}
	w.WriteHeader(http.StatusOK)
}
