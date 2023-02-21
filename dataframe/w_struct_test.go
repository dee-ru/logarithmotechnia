package dataframe

import (
	"fmt"
	"logarithmotechnia/vector"
	"reflect"
	"testing"
	"time"
)

func TestFromStructs(t *testing.T) {
	type Finance struct {
		Money   int
		Account string
	}

	type A struct {
		Title     string
		Status    int `lth:"status"`
		OldStatus int `lto:"skip"`
		Kpi       float64
		Cpx       complex128
		IsActive  bool `lth:"is_active"`
		Date      time.Time
		OldDate   time.Time
		Misc      Finance
	}

	now := time.Now()

	stArr := []A{
		{
			Title:    "Baron",
			Status:   1,
			Kpi:      1.2,
			Cpx:      1 + 1i,
			IsActive: true,
			Date:     now,
			Misc:     Finance{1000, "br"},
		},
		{
			Title:    "Earl",
			Status:   3,
			Kpi:      2.2,
			Cpx:      1 + 3i,
			IsActive: false,
			Date:     now.Add(7 * 24 * 60 * time.Minute),
			Misc:     Finance{15000, "ct"},
		},
		{
			Title:    "King",
			Status:   5,
			Kpi:      4.45,
			Cpx:      4 + 2i,
			IsActive: true,
			Date:     now.Add(360 * 24 * 60 * time.Minute),
			Misc:     Finance{275000, "kn"},
		},
	}

	df, err := FromStructs(stArr, StructOptionHeaderMap(map[string]string{"Date": "date"}),
		StructOptionSkipFields("OldDate"))
	if err != nil {
		t.Error(err)
	}

	columnNames := []string{"Title", "status", "Kpi", "Cpx", "is_active", "date", "Misc"}
	columns := []vector.Vector{
		vector.String([]string{"Baron", "Earl", "King"}),
		vector.Integer([]int{1, 3, 5}),
		vector.Float([]float64{1.2, 2.2, 4.45}),
		vector.Complex([]complex128{1 + 1i, 1 + 3i, 4 + 2i}),
		vector.Boolean([]bool{true, false, true}),
		vector.Time([]time.Time{now, now.Add(7 * 24 * 60 * time.Minute), now.Add(360 * 24 * 60 * time.Minute)}),
		vector.Any([]any{Finance{1000, "br"}, Finance{15000, "ct"},
			Finance{275000, "kn"}}),
	}

	if !reflect.DeepEqual(df.columnNames, columnNames) {
		t.Error(fmt.Sprintf("Column names %v are not equal to expected (%v)", df.columnNames, columnNames))
	}

	if !vector.CompareVectorArrs(df.columns, columns) {
		t.Error(fmt.Sprintf("Columns %v are not equal to expected (%v)", df.columns, columns))
	}
}

func TestStructOptions(t *testing.T) {
	testData := []struct {
		name      string
		result    Option
		reference Option
	}{
		{
			name:      "StructOptionHeaderMap",
			result:    StructOptionHeaderMap(map[string]string{"name": "title"}),
			reference: ConfOption{optionStructHeaderMap, map[string]string{"name": "title"}},
		},
		{
			name:   "StructOptionDataFrameOptions",
			result: StructOptionDataFrameOptions(OptionColumnNames([]string{"id", "price"})),
			reference: ConfOption{optionStructDataframeOptions,
				[]Option{OptionColumnNames([]string{"id", "price"})}},
		},
		{
			name:      "StructOptionSkipFields",
			result:    StructOptionSkipFields("name", "dep", "salary"),
			reference: ConfOption{optionStructSkipFields, []string{"name", "dep", "salary"}},
		},
	}

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			if !reflect.DeepEqual(data.result, data.reference) {
				t.Error(fmt.Sprintf("Resulting conf option (%v) does not match reference (%v)",
					data.result, data.reference))
			}
		})
	}
}
