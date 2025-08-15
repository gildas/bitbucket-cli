package common

import "github.com/gildas/go-core"

type Column[T any] struct {
	Name          string
	DefaultSorter bool
	Compare       func(a, b T) bool
}

type Columns[T any] []Column[T]

func (columns Columns[T]) Columns() []string {
	return core.Map(columns, func(column Column[T]) string { return column.Name })
}

func (columns Columns[T]) Sorters() []string {
	return core.Map(columns, func(column Column[T]) string {
		if column.DefaultSorter {
			return "+" + column.Name
		}
		return column.Name
	})
}

func (columns Columns[T]) SortBy(sorter string) func(a, b T) bool {
	for _, column := range columns {
		if column.Name == sorter {
			return column.Compare
		}
	}
	return func(a, b T) bool { return false } // We should never get here!
}
