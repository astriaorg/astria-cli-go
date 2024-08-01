package cmd

import (
	"testing"
)

func TestGetFieldValueByTag(t *testing.T) {
	type testStruct struct {
		Field1 string `flag:"field-1"`
		Field2 int    `tag2:"field-2"`
		Field3 bool   `b:"field-3"`
	}

	s := testStruct{
		Field1: "value1",
		Field2: 1,
		Field3: true,
	}

	t.Run("Get field by tag", func(t *testing.T) {
		val1, ok := GetFieldValueByTag(s, "flag", "field-1")
		if !ok {
			t.Errorf("Expected to find field by tag")
		}
		if val1.String() != "value1" {
			t.Errorf("Expected field value to be %s, got %s", s.Field1, val1.String())
		}

		val2, ok := GetFieldValueByTag(s, "tag2", "field-2")
		if !ok {
			t.Errorf("Expected to find field by tag")
		}
		if val2.Int() != 1 {
			t.Errorf("Expected field value to be %d, got %s", s.Field2, val2.String())
		}

		val3, ok := GetFieldValueByTag(s, "b", "field-3")
		if !ok {
			t.Errorf("Expected to find field by tag")
		}
		if val3.Bool() != true {
			t.Errorf("Expected field value to be %t, got %s", s.Field3, val3.String())
		}
	})

}
