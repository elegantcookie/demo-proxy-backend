package filter

import (
	"fmt"
	"strconv"
	"strings"
)

func (o *Options) ValidateIntAndAdd(key, value, defaultOperator string) {
	if value != "" {
		operator := defaultOperator
		val := value
		if strings.Index(val, ":") != -1 {
			split := strings.Split(value, ":")
			operator = split[0]
			val = split[1]
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
