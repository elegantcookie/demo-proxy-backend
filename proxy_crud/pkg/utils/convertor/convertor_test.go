package convertor

import (
	"fmt"
	"testing"
	"time"
)

type TestStruct struct {
	Field1 string    `json:"field_1"`
	Field2 int       `json:"field_2"`
	Field3 time.Time `json:"field_3"`
}

func TestStructsToArrays(t *testing.T) {
	testArr := []TestStruct{
		{"1", 1, time.Now()},
		{"2", 2, time.Now()},
		{"3", 3, time.Now()},
	}

	s := make([]any, len(testArr))
	fmt.Println(s)
	// [<nil> <nil> <nil>]
	for i, v := range testArr {
		s[i] = v
	}
	arr := StructsToArrays(s)
	fmt.Println(arr)
	// [
	//	[1 1 2022-10-15 13:47:08.7220984 +0300 MSK m=+0.001544001],
	//	[2 2 2022-10-15 13:47:08.7220984 +0300 MSK m=+0.001544001],
	//	[3 3 2022-10-15 13:47:08.7220984 +0300 MSK m=+0.001544001]
	//]

}
