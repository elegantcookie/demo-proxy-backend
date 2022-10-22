package filter

import (
	"github.com/google/uuid"
	"testing"
)

func keyInFilterOptions(opts *Options, key string) bool {
	for _, v := range opts.FilterOptions.Fields() {
		if v.Key == key {
			return true
		}
	}
	return false
}

func TestOptions(t *testing.T) {
	t.Parallel()
	foption := &FOptions{
		apply:  true,
		limit:  10,
		page:   10,
		fields: make([]Field, 0),
	}

	options := Options{
		SortOptions:   SOptions{},
		FilterOptions: foption,
	}

	t.Run("ValidateBoolAndAdd", func(t *testing.T) {
		testKey := uuid.New().String()
		options.ValidateBoolAndAdd(testKey, "nonBoolValue", OperatorEqual)
		if len(options.FilterOptions.Fields()) != 0 {
			t.Fail()
		}

		testKey = uuid.New().String()
		options.ValidateBoolAndAdd(testKey, "true", OperatorEqual)
		if !keyInFilterOptions(&options, testKey) {
			t.Fail()
		}

		testKey = uuid.New().String()
		options.ValidateBoolAndAdd(testKey, "false", OperatorEqual)
		if !keyInFilterOptions(&options, testKey) {
			t.Fail()
		}

		testKey = uuid.New().String()
		options.ValidateBoolAndAdd(testKey, "1", OperatorEqual)
		if !keyInFilterOptions(&options, testKey) {
			t.Fail()
		}

		testKey = uuid.New().String()
		options.ValidateBoolAndAdd(testKey, "0", OperatorEqual)
		if !keyInFilterOptions(&options, testKey) {
			t.Fail()
		}
	})

	t.Run("ValidateIntAndAdd", func(t *testing.T) {
		testKey := uuid.New().String()
		options.ValidateIntAndAdd(testKey, "nonIntValue", OperatorEqual)
		if len(options.FilterOptions.Fields()) != 0 {
			t.Fail()
		}

		testKey = uuid.New().String()
		options.ValidateIntAndAdd(testKey, "1", OperatorEqual)
		if !keyInFilterOptions(&options, testKey) {
			t.Fail()
		}

		testKey = uuid.New().String()
		options.ValidateIntAndAdd(testKey, "like:-123", OperatorEqual)
		if !keyInFilterOptions(&options, testKey) {
			t.Fail()
		}
	})

	t.Run("ValidateStringAndAdd", func(t *testing.T) {

	})
}
