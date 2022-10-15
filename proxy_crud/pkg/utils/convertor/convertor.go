package convertor

import (
	"reflect"
)

// StructsToArrays converts array of structs to array of interfaces
// It is used to bulk insert values by pgx.CopyFromRows
func StructsToArrays(proxies []any) [][]any {
	proxiesLen := len(proxies)

	arr := make([][]any, proxiesLen)
	for i := 0; i < proxiesLen; i++ {
		v := reflect.ValueOf(proxies[i])
		arr[i] = make([]any, v.NumField())
		for j := 0; j < v.NumField(); j++ {
			arr[i][j] = v.Field(j).Interface()
		}
	}
	return arr
}
