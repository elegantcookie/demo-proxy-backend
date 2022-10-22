package filter

import (
	"fmt"
	"strconv"
	"strings"
)

// TODO: refactor

//func parseValueForOptions(value string) (operator, val string, ok bool) {
//	if strings.Index(val, ":") != -1 {
//		split := strings.Split(value, ":")
//		operator = split[0]
//		val = split[1]
//		return operator, val, true
//	}
//	return operator, val, false
//}

func isIn(operatorName string, operators []string) bool {
	for i := 0; i < len(operators); i++ {
		if operators[i] == operatorName {
			return true
		}
	}
	return false
}

func (o *Options) ValidateIntAndAdd(key, value, defaultOperator string) {

	// may be faster to pass map and use it as forbidden operators
	forbiddenOperators := []string{
		OperatorLike,
	}

	if value != "" {
		operator := defaultOperator
		val := value
		if strings.Index(val, ":") != -1 {
			split := strings.Split(value, ":")
			operator = split[0]
			val = split[1]
		}
		if in := isIn(operator, forbiddenOperators); in {
			fmt.Printf("invalid: %v", val)
			return
		}
		var intVal int
		var err error
		if intVal, err = strconv.Atoi(val); err != nil {
			fmt.Printf("invalid: %v", val)
			return
		}
		o.FilterOptions.AddField(key, intVal, operator, dataTypeInt)
	}
}

func (o *Options) ValidateStringAndAdd(key, value, defaultOperator string) {
	if value != "" {
		operator := defaultOperator
		val := value
		if strings.Index(val, ":") != -1 {
			split := strings.Split(value, ":")
			operator = split[0]
			val = split[1]
		}
		o.FilterOptions.AddField(key, value, operator, dataTypeString)

	}
}

func (o *Options) ValidateBoolAndAdd(key, value, defaultOperator string) {
	forbiddenOperators := []string{
		OperatorLike,
	}
	if value != "" {
		operator := defaultOperator
		val := value
		if strings.Index(val, ":") != -1 {
			split := strings.Split(value, ":")
			operator = split[0]
			val = split[1]
		}
		if in := isIn(operator, forbiddenOperators); in {
			fmt.Printf("invalid: %v", val)
			return
		}
		var boolVal bool
		var err error
		if boolVal, err = strconv.ParseBool(val); err != nil {
			fmt.Printf("invalid: %v", val)
			return
		}
		o.FilterOptions.AddField(key, boolVal, operator, dataTypeBool)
	}
}
