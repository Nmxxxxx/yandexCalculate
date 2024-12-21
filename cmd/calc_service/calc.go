package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type Request struct {
	Expression string `json:"expression"`
}

type Response struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

func isOperator(c byte) bool {
	return c == '+' || c == '-' || c == '*' || c == '/'
}

func priority(op byte) int {
	switch op {
	case '+', '-':
		return 1
	case '*', '/':
		return 2
	}
	return 0
}

func applyOp(a float64, b float64, op byte) (float64, error) {
	switch op {
	case '+':
		return a + b, nil
	case '-':
		return a - b, nil
	case '*':
		return a * b, nil
	case '/':
		if b == 0 {
			return 0, errors.New("division by zero")
		}
		return a / b, nil
	}
	return 0, errors.New("unknown operator")
}

func Calc(expression string) (float64, error) {
	var numStack []float64
	var opStack []byte

	for i := 0; i < len(expression); {
		if expression[i] == ' ' {
			i++
			continue
		} else if isOperator(expression[i]) {
			for len(opStack) > 0 && priority(opStack[len(opStack)-1]) >= priority(expression[i]) {
				op := opStack[len(opStack)-1]
				opStack = opStack[:len(opStack)-1]

				if len(numStack) < 2 {
					return 0, errors.New("invalid expression")
				}
				b := numStack[len(numStack)-1]
				numStack = numStack[:len(numStack)-1]
				a := numStack[len(numStack)-1]
				numStack = numStack[:len(numStack)-1]

				result, err := applyOp(a, b, op)
				if err != nil {
					return 0, err
				}
				numStack = append(numStack, result)
			}
			opStack = append(opStack, expression[i])
		} else {
			start := i
			for i < len(expression) && ((expression[i] >= '0' && expression[i] <= '9') || expression[i] == '.') {
				i++
			}
			val, err := strconv.ParseFloat(expression[start:i], 64)
			if err != nil {
				return 0, err
			}
			numStack = append(numStack, val)
			continue
		}
		i++
	}

	for len(opStack) > 0 {
		op := opStack[len(opStack)-1]
		opStack = opStack[:len(opStack)-1]

		if len(numStack) < 2 {
			return 0, errors.New("invalid expression")
		}
		b := numStack[len(numStack)-1]
		numStack = numStack[:len(numStack)-1]
		a := numStack[len(numStack)-1]
		numStack = numStack[:len(numStack)-1]

		result, err := applyOp(a, b, op)
		if err != nil {
			return 0, err
		}
		numStack = append(numStack, result)
	}

	if len(numStack) != 1 || len(opStack) != 0 {
		return 0, errors.New("invalid expression")
	}

	return numStack[0], nil
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
		return
	}

	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Expression == "" {
		http.Error(w, "Expression is not valid", http.StatusUnprocessableEntity)
		return
	}

	result, err := Calc(req.Expression)
	if err != nil {
		if err.Error() == "division by zero" {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	resp := Response{Result: strconv.FormatFloat(result, 'f', -1, 64)}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/api/v1/calculate", calculateHandler)
	fmt.Println("Сервер запущен на порту 8080")
	http.ListenAndServe(":8080", nil)
}
