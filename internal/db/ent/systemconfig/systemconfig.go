// Code generated by ent, DO NOT EDIT.

package systemconfig

const (
	// Label holds the string label denoting the systemconfig type in the database.
	Label = "system_config"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldKey holds the string denoting the key field in the database.
	FieldKey = "key"
	// FieldValue holds the string denoting the value field in the database.
	FieldValue = "value"
	// Table holds the table name of the systemconfig in the database.
	Table = "system_configs"
)

// Columns holds all SQL columns for systemconfig fields.
var Columns = []string{
	FieldID,
	FieldKey,
	FieldValue,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultValue holds the default value on creation for the "value" field.
	DefaultValue string
)
