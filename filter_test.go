package validate

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddFilter(t *testing.T) {
	is := assert.New(t)
	is.Panics(func() {
		AddFilter("myFilter", "invalid")
	})
	is.Panics(func() {
		AddFilter("myFilter", func() {})
	})
	is.Panics(func() {
		AddFilter("bad-name", func() {})
	})
	is.Panics(func() {
		AddFilter("", func() {})
	})
	is.Panics(func() {
		AddFilter("myFilter", func(v string) (bool, int) { return false, 0 })
	})
	is.Panics(func() {
		AddFilter("myFilter", func() interface{} { return nil })
	})

	AddFilters(M{
		"myFilter0": func(val interface{}) string { return "myFilter0" },
	})
	AddFilter("myFilter1", func(val interface{}) string { return "myFilter1" })

	v := New(map[string]interface{}{
		"name": " inhere ",
		"age":  " 50 ",
		"key0": "val0",
		"key1": "val1",
		"tags": "go,php",
	})
	v.AddFilters(M{
		"myFilter2": func(val interface{}) (string, error) { return "myFilter2", nil },
	})
	v.FilterRule("key0", "myFilter0")
	v.FilterRules(MS{
		"key1": "myFilter2",
		"name": "trim|upper",
		"tags": "str2arr:,",
		//
		"age, not-exist": "trim|int",
	})

	is.Panics(func() {
		v.FilterRule("", "")
	})

	v.Sanitize() // do filtering
	is.True(v.IsOK())
	is.Equal(50, v.Filtered("age"))
	is.Equal("INHERE", v.Filtered("name"))
	is.Equal("myFilter0", v.Filtered("key0"))
	is.Equal("myFilter2", v.Filtered("key1"))
	is.Contains(fmt.Sprint(v.FilteredData()), "key0:myFilter0")

	// filter fail
	v = New(SValues{
		"name": {"inhere"},
	})
	v.FilterRules(MS{
		"name": "invalid|int",
	})
	v.Filtering()
	is.True(v.IsFail())
	is.Contains(v.Errors, "_filter")
}
