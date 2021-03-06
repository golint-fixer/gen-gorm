package backends

import (
	"database/sql"
	"log"
	"unicode"

	"github.com/kmulvey/gen-gorm/graph"
)

// Backend is a database
type Backend interface {
	createModel() graph.Graph
	createConn() *sql.DB
}

// ConnConfig db params
type ConnConfig struct {
	Hostname *string
	Username *string
	Password *string
	Schema   *string
	Port     *string
}

// BackendFactory selects the correct engine and generates the model. This should not be unit tested.
func BackendFactory(config ConnConfig, engine string) graph.Graph {
	switch engine {
	case "mysql":
		m := Mysql{}
		conn := m.createConn(config)
		defer conn.Close()
		return m.createModel(conn, config)
	case "postgres":
		p := Postgres{}
		conn := p.createConn(config)
		defer conn.Close()
		return p.createModel(conn, config)
	default:
		log.Fatal("We dont support \"" + engine + "\" :(")
		return graph.Graph{} // this is stupid
	}
}

// capFirst capitalized the first character of a string
func capFirst(input string) string {
	arr := []byte(input)
	arr[0] = byte(unicode.ToUpper(rune(arr[0])))
	return string(arr)
}

// formatColName formats the column name into camel case
func formatColName(name string) string {
	arr := []byte(name)
	var output []byte
	capNextChar := false
	for i, char := range arr {
		if i == 0 || capNextChar {
			output = append(output, byte(unicode.ToUpper(rune(char))))
			capNextChar = false
		} else if char == '_' {
			capNextChar = true
		} else {
			output = append(output, char)
			capNextChar = false
		}
	}
	return string(output)
}

// convertType converts the db col type to the corresponding go type
func convertType(dbType string) string {
	switch dbType {
	// Dates represented as strings
	case "time", "date", "datetime":
		fallthrough

	// Buffers represented as strings
	case "bit", "blob", "tinyblob", "longblob", "mediumblob", "binary", "varbinary":
		fallthrough

	// Numbers that may exceed float precision, repesent as string
	case "bigint", "decimal", "numeric", "geometry", "bigserial":
		fallthrough

	// Network addresses represented as strings
	case "cidr", "inet", "macaddr":
		fallthrough

	// Strings
	case "set", "char", "text", "uuid", "varchar", "nvarchar", "tinytext", "longtext", "character", "mediumtext":
		return "string"
	// Integers
	case "int", "year", "serial", "integer", "tinyint", "smallint", "mediumint", "timestamp":
		return "int"
	// Floats
	case "real", "float", "double", "double precision":
		return "float"

	// Booleans
	case "boolean":
		return "bool"

	// Enum special case
	case "enum":
		return "string"

	default:
		return "string"
	}
}
