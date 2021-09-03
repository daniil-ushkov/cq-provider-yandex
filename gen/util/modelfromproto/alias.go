package modelfromproto

// Alias is applied to Column or Table and changes it.
type Alias interface {
	ApplyToColumn(column *Column)
	ApplyToTable(table *Table)
}

type UnimplementedAlias struct{}

func (u UnimplementedAlias) ApplyToColumn(*Column) {}

func (u UnimplementedAlias) ApplyToTable(*Table) {}

type changeName struct {
	UnimplementedAlias
	Name string
}

// ChangeName returns Alias which changes name of Column or Table
func ChangeName(name string) Alias {
	return changeName{Name: name}
}

func (c changeName) ApplyToColumn(column *Column) {
	column.Name = c.Name
}

func (c changeName) ApplyToTable(table *Table) {
	table.Alias = c.Name
}

type changeColumn struct {
	UnimplementedAlias
	column *Column
}

// ReplaceColumn replaces Column with Column from param.
func ReplaceColumn(column *Column) Alias {
	return changeColumn{column: column}
}

func (r changeColumn) ApplyToColumn(column *Column) {
	*column = *r.column
}
