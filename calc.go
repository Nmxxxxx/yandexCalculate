package main

import (
	"errors"
	"strconv"
)

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
		} else if expression[i] == '(' {
			opStack = append(opStack, expression[i])
		} else if expression[i] == ')' {
			for len(opStack) > 0 && opStack[len(opStack)-1] != '(' {
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
			opStack = opStack[:len(opStack)-1]
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
