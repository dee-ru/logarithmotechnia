package dataframe

import (
	"logarithmotechnia/vector"
)

func (df *Dataframe) Mutate(columns []Column, options ...vector.Option) *Dataframe {
	conf := vector.MergeOptions(options)

	afterColumnIndex := df.colNum - 1

	if conf.HasOption(vector.KeyOptionAfterColumn) {
		pos := df.Names().Find(conf.Value(vector.KeyOptionAfterColumn))
		if pos > 0 {
			afterColumnIndex = pos
		}
	}

	if conf.HasOption(vector.KeyOptionBeforeColumn) {
		pos := df.Names().Find(conf.Value(vector.KeyOptionBeforeColumn))
		if pos > 0 {
			afterColumnIndex = pos - 1
		}
	}

	columnMap := map[string]vector.Vector{}
	for i, column := range df.columns {
		columnMap[df.config.columnNames[i]] = column
	}

	uniqueNewNames := []string{}
	for _, column := range columns {
		if _, ok := columnMap[column.name]; ok {
			uniqueNewNames = append(uniqueNewNames, column.name)
		}
		columnMap[column.name] = column.vector
	}

	newNames := []string{}
	newNames = append(newNames, df.config.columnNames[:afterColumnIndex]...)
	newNames = append(newNames, uniqueNewNames...)
	newNames = append(newNames, df.config.columnNames[afterColumnIndex:]...)

	newColumns := []vector.Vector{}
	for _, name := range newNames {
		newColumns = append(newColumns, columnMap[name])
	}

	return New(newColumns, OptionColumnNames(newNames))
}

func (df *Dataframe) Transmute(map[string]vector.Vector) *Dataframe {
	df.Mutate([]Column{})

	return nil
}
