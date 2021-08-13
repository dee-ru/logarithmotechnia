package dataframe

import (
	"fmt"
	"github.com/dee-ru/logarithmotechnia/vector"
	"reflect"
	"strconv"
	"testing"
)

func TestNew(t *testing.T) {
	testData := []struct {
		name      string
		columns   []vector.Vector
		config    []Config
		dfColumns []vector.Vector
		dfConfig  Config
	}{
		{
			name:      "empty",
			columns:   []vector.Vector{},
			config:    []Config{},
			dfColumns: []vector.Vector{},
			dfConfig:  Config{columnNames: []string{}},
		},
		{
			name:      "empty with column names",
			columns:   []vector.Vector{},
			config:    []Config{OptionColumnNames([]string{"one", "two", "three"})},
			dfColumns: []vector.Vector{},
			dfConfig:  Config{columnNames: []string{}},
		},
		{
			name: "normal",
			columns: []vector.Vector{
				vector.Integer([]int{1, 2, 3}, nil),
				vector.String([]string{"1", "2", "3"}, nil),
				vector.Boolean([]bool{true, true, false}, nil),
			},
			config: []Config{},
			dfColumns: []vector.Vector{
				vector.Integer([]int{1, 2, 3}, nil),
				vector.String([]string{"1", "2", "3"}, nil),
				vector.Boolean([]bool{true, true, false}, nil),
			},
			dfConfig: Config{columnNames: []string{"1", "2", "3"}},
		},
		{
			name: "normal with column names",
			columns: []vector.Vector{
				vector.Integer([]int{1, 2, 3}, nil),
				vector.String([]string{"1", "2", "3"}, nil),
				vector.Boolean([]bool{true, true, false}, nil),
			},
			config: []Config{OptionColumnNames([]string{"int", "string", "bool"})},
			dfColumns: []vector.Vector{
				vector.Integer([]int{1, 2, 3}, nil),
				vector.String([]string{"1", "2", "3"}, nil),
				vector.Boolean([]bool{true, true, false}, nil),
			},
			dfConfig: Config{columnNames: []string{"int", "string", "bool"}},
		},
		{
			name: "normal with partial column names",
			columns: []vector.Vector{
				vector.Integer([]int{1, 2, 3}, nil),
				vector.String([]string{"1", "2", "3"}, nil),
				vector.Boolean([]bool{true, true, false}, nil),
			},
			config: []Config{OptionColumnNames([]string{"int", "string"})},
			dfColumns: []vector.Vector{
				vector.Integer([]int{1, 2, 3}, nil),
				vector.String([]string{"1", "2", "3"}, nil),
				vector.Boolean([]bool{true, true, false}, nil),
			},
			dfConfig: Config{columnNames: []string{"int", "string", "3"}},
		},
		{
			name: "different columns' length",
			columns: []vector.Vector{
				vector.Integer([]int{1, 2, 3, 4, 5}, nil),
				vector.String([]string{"1", "2", "3"}, nil),
				vector.Boolean([]bool{true, true}, nil),
			},
			config: []Config{},
			dfColumns: []vector.Vector{
				vector.Integer([]int{1, 2, 3, 4, 5}, nil),
				vector.String([]string{"1", "2", "3", "", ""}, []bool{false, false, false, true, true}),
				vector.Boolean([]bool{true, true, false, false, false}, []bool{false, false, true, true, true}),
			},
			dfConfig: Config{columnNames: []string{"1", "2", "3"}},
		},
	}

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			df := New(data.columns, data.config...)

			if !reflect.DeepEqual(df.columns, data.dfColumns) {
				t.Error(fmt.Sprintf("Columns (%v) are not equal to expected (%v)", df.columns, data.dfColumns))
			}
			if !reflect.DeepEqual(df.config, data.dfConfig) {
				t.Error(fmt.Sprintf("Config (%v) are not equal to expected (%v)",
					df.config, data.dfConfig))
			}
		})
	}
}

func TestDataframe_ByIndices(t *testing.T) {
	df := New([]vector.Vector{
		vector.Integer([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, nil),
		vector.String([]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}, nil),
		vector.Boolean([]bool{true, false, true, false, true, false, true, false, true, false}, nil),
	})

	testData := []struct {
		name      string
		indices   []int
		dfColumns []vector.Vector
	}{
		{
			name:    "normal",
			indices: []int{1, 3, 5, 8, 10},
			dfColumns: []vector.Vector{
				vector.Integer([]int{1, 3, 5, 8, 10}, nil),
				vector.String([]string{"1", "3", "5", "8", "10"}, nil),
				vector.Boolean([]bool{true, true, true, false, false}, nil),
			},
		},
		{
			name:    "with invalid",
			indices: []int{-1, 0, 1, 3, 5, 8, 10, 11, 100},
			dfColumns: []vector.Vector{
				vector.Integer([]int{1, 3, 5, 8, 10}, nil),
				vector.String([]string{"1", "3", "5", "8", "10"}, nil),
				vector.Boolean([]bool{true, true, true, false, false}, nil),
			},
		},
		{
			name:    "reverse",
			indices: []int{10, 8, 5, 3, 1},
			dfColumns: []vector.Vector{
				vector.Integer([]int{10, 8, 5, 3, 1}, nil),
				vector.String([]string{"10", "8", "5", "3", "1"}, nil),
				vector.Boolean([]bool{false, false, true, true, true}, nil),
			},
		},
		{
			name:    "empty",
			indices: []int{},
			dfColumns: []vector.Vector{
				vector.Integer([]int{}, nil),
				vector.String([]string{}, nil),
				vector.Boolean([]bool{}, nil),
			},
		},
	}

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			newDf := df.ByIndices(data.indices)

			if !reflect.DeepEqual(newDf.columns, data.dfColumns) {
				t.Error(fmt.Sprintf("Columns (%v) are not equal to expected (%v)", newDf.columns, data.dfColumns))
			}
		})
	}
}

func TestDataframe_ColNum(t *testing.T) {
	df := New([]vector.Vector{
		vector.Integer([]int{1, 2, 3, 4, 5}, nil),
		vector.String([]string{"1", "2", "3"}, nil),
		vector.Boolean([]bool{true, true}, nil),
	})

	if df.ColNum() != 3 {
		t.Error("Column number is incorrect!")
	}
}

func TestDataframe_RowNum(t *testing.T) {
	df := New([]vector.Vector{
		vector.Integer([]int{1, 2, 3, 4, 5}, nil),
		vector.String([]string{"1", "2", "3"}, nil),
		vector.Boolean([]bool{true, true}, nil),
	})

	if df.RowNum() != 5 {
		t.Error("Row number is incorrect!")
	}
}

func TestDataframe_Ci(t *testing.T) {
	df := New([]vector.Vector{
		vector.Integer([]int{1, 2, 3, 4, 5}, nil),
		vector.String([]string{"1", "2", "3"}, nil),
		vector.Boolean([]bool{true, true}, nil),
	}, OptionColumnNames([]string{"int", "string", "bool"}))

	testData := []struct {
		index  int
		column vector.Vector
	}{
		{1, vector.Integer([]int{1, 2, 3, 4, 5}, nil)},
		{2, vector.String([]string{"1", "2", "3", "", ""}, []bool{false, false, false, true, true})},
		{3, vector.Boolean([]bool{true, true, false, false, false}, []bool{false, false, true, true, true})},
		{0, nil},
		{-1, nil},
		{4, nil},
	}

	for i, data := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			column := df.Ci(data.index)

			if !reflect.DeepEqual(column, data.column) {
				t.Error(fmt.Sprintf("Columns (%v) are not equal to expected (%v)", column, data.column))
			}
		})
	}
}

func TestDataframe_Cn(t *testing.T) {
	df := New([]vector.Vector{
		vector.Integer([]int{1, 2, 3, 4, 5}, nil),
		vector.String([]string{"1", "2", "3"}, nil),
		vector.Boolean([]bool{true, true}, nil),
	}, OptionColumnNames([]string{"int", "string", "bool"}))

	testData := []struct {
		name   string
		column vector.Vector
	}{
		{"int", vector.Integer([]int{1, 2, 3, 4, 5}, nil)},
		{"string", vector.String([]string{"1", "2", "3", "", ""}, []bool{false, false, false, true, true})},
		{"bool", vector.Boolean([]bool{true, true, false, false, false}, []bool{false, false, true, true, true})},
		{"", nil},
		{"some", nil},
	}

	for i, data := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			column := df.Cn(data.name)

			if !reflect.DeepEqual(column, data.column) {
				t.Error(fmt.Sprintf("Columns (%v) are not equal to expected (%v)", column, data.column))
			}
		})
	}
}

func TestDataframe_C(t *testing.T) {
	df := New([]vector.Vector{
		vector.Integer([]int{1, 2, 3}, nil),
		vector.String([]string{"1", "2", "3"}, nil),
		vector.Boolean([]bool{true, false, true}, nil),
	}, OptionColumnNames([]string{"int", "string", "bool"}))

	testData := []struct {
		selector interface{}
		column   vector.Vector
	}{
		{"int", vector.Integer([]int{1, 2, 3}, nil)},
		{2, vector.String([]string{"1", "2", "3"}, nil)},
		{0, nil},
		{4, nil},
		{"some", nil},
	}

	for i, data := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			column := df.C(data.selector)

			if !reflect.DeepEqual(column, data.column) {
				t.Error(fmt.Sprintf("Columns (%v) are not equal to expected (%v)", column, data.column))
			}
		})
	}
}

func TestDataframe_ColumnNames(t *testing.T) {
	testData := []struct {
		name        string
		dataframe   *Dataframe
		columnNames []string
	}{
		{
			name: "with names",
			dataframe: New([]vector.Vector{
				vector.Integer([]int{1, 2, 3}, nil),
				vector.String([]string{"1", "2", "3"}, nil),
				vector.Boolean([]bool{true, false, true}, nil),
			}, OptionColumnNames([]string{"int", "string", "bool"})),
			columnNames: []string{"int", "string", "bool"},
		},
	}

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			columnNames := data.dataframe.ColumnNames()

			if !reflect.DeepEqual(columnNames, data.columnNames) {
				t.Error(fmt.Sprintf("Columns names (%v) are not equal to expected (%v)",
					columnNames, data.columnNames))
			}
		})
	}
}

func TestDataframe_IsEmpty(t *testing.T) {
	testData := []struct {
		name      string
		dataframe *Dataframe
		isEmpty   bool
	}{
		{
			name:      "empty",
			dataframe: New([]vector.Vector{}),
			isEmpty:   true,
		},
		{
			name: "non-empty",
			dataframe: New([]vector.Vector{
				vector.Integer([]int{1, 2, 3}, nil),
				vector.String([]string{"1", "2", "3"}, nil),
				vector.Boolean([]bool{true, false, true}, nil),
			}, OptionColumnNames([]string{"int", "string", "bool"})),
			isEmpty: false,
		},
	}

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			isEmpty := data.dataframe.IsEmpty()

			if isEmpty != data.isEmpty {
				t.Error(fmt.Sprintf("IsEmpty (%t) is not equal to expected (%t)",
					isEmpty, data.isEmpty))
			}
		})
	}
}

func TestDataframe_SetColumnName(t *testing.T) {
	df := New([]vector.Vector{
		vector.Integer([]int{1, 2, 3}, nil),
		vector.String([]string{"1", "2", "3"}, nil),
		vector.Boolean([]bool{true, false, true}, nil),
	})

	testData := []struct {
		name        string
		index       int
		columnNames []string
	}{
		{
			name:        "int",
			index:       1,
			columnNames: []string{"int", "2", "3"},
		},
		{
			name:        "string",
			index:       2,
			columnNames: []string{"int", "string", "3"},
		},
		{
			name:        "bool",
			index:       3,
			columnNames: []string{"int", "string", "bool"},
		},
		{
			name:        "zero",
			index:       0,
			columnNames: []string{"int", "string", "bool"},
		},
		{
			name:        "four",
			index:       4,
			columnNames: []string{"int", "string", "bool"},
		},
	}

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			df.SetColumnName(data.index, data.name)

			if !reflect.DeepEqual(df.config.columnNames, data.columnNames) {
				t.Error(fmt.Sprintf("Column names (%v) is not equal to expected (%v)",
					df.config.columnNames, data.columnNames))
			}
		})
	}
}

func TestDataframe_SetColumnNames(t *testing.T) {
	testData := []struct {
		name        string
		columnNames []string
		resultNames []string
	}{
		{"all", []string{"int", "string", "bool"}, []string{"int", "string", "bool"}},
		{"partial", []string{"int", "string"}, []string{"int", "string", "3"}},
		{"excess", []string{"int", "string", "bool", "some"}, []string{"int", "string", "bool"}},
	}

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			df := New([]vector.Vector{
				vector.Integer([]int{1, 2, 3}, nil),
				vector.String([]string{"1", "2", "3"}, nil),
				vector.Boolean([]bool{true, false, true}, nil),
			})

			df.SetColumnNames(data.columnNames)

			if !reflect.DeepEqual(df.config.columnNames, data.resultNames) {
				t.Error(fmt.Sprintf("Column names (%v) is not equal to expected (%v)",
					df.config.columnNames, data.resultNames))
			}
		})
	}
}

func TestDataframe_Clone(t *testing.T) {
	testData := []struct {
		name      string
		dataframe *Dataframe
		isEmpty   bool
	}{
		{
			name:      "empty",
			dataframe: New([]vector.Vector{}),
			isEmpty:   true,
		},
		{
			name: "non-empty",
			dataframe: New([]vector.Vector{
				vector.Integer([]int{1, 2, 3}, nil),
				vector.String([]string{"1", "2", "3"}, nil),
				vector.Boolean([]bool{true, false, true}, nil),
			}, OptionColumnNames([]string{"int", "string", "bool"})),
			isEmpty: false,
		},
	}

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			newDf := data.dataframe.Clone()

			if !reflect.DeepEqual(newDf.columns, data.dataframe.columns) {
				t.Error(fmt.Sprintf("Columns (%v) are not equal to expected (%v)",
					newDf.columns, data.dataframe.columns))
			}
			if !reflect.DeepEqual(newDf.config, data.dataframe.config) {
				t.Error(fmt.Sprintf("Config (%v) is not equal to expected (%v)",
					newDf.config, data.dataframe.config))
			}
		})
	}
}

func TestDataframe_Columns(t *testing.T) {

}
