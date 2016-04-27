package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"os/exec"
	"unicode"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kmulvey/gen-gorm/graph"

	_ "github.com/lib/pq"
)

type column struct {
	ColumnName   string
	ColumnType   string
	DBColumnName string
}
type table struct {
	TableName string
	Cols      []column
}

func main() {
	hostname := flag.String("hostname", "", "hostname")
	username := flag.String("username", "", "username")
	password := flag.String("password", "", "password")
	schema := flag.String("schema", "", "schema")
	port := flag.String("port", "3306", "port")
	output := flag.String("output", "", "output")
	flag.Parse()

	// connect to db
	//db, err := sql.Open("postgres", fmt.Sprintf("host=%v user=%v dbname=%v password=%v port=%v sslmode=disable", hostname, username, schema, password, port))
	conn, err := sql.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8", *username, *password, *hostname, *port, *schema))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	err = conn.Ping()
	handleErr(err)

	// get table structure from DB
	data := getTableInfo(conn, *schema)

	// get 'er done
	processTemplates(data, *output)
}

// processTemplates fills in the templates with data, puts them in the output
// directory and fmt them
func processTemplates(data graph.Graph, output string) {

	// parse templates
	modelsTemplate, err := template.ParseFiles("templates/models.tmpl")
	handleErr(err)

	// create directory
	if _, err := os.Stat(output); os.IsNotExist(err) {
		os.Mkdir(output, 0755)
	}

	// create the files
	modelsGo, err := os.Create(output + "/models.go")
	handleErr(err)
	defer modelsGo.Close()

	// exec templates
	err = modelsTemplate.Execute(modelsGo, data)
	handleErr(err)

	// format the file
	cmd := exec.Command("gofmt", "-w", output)
	err = cmd.Run()
	handleErr(err)
}

// getTableInfo retrieves schema information from the database
func getTableInfo(conn *sql.DB, schema string) (database graph.Graph) {
	database.Name = schema
	database.Vertices = make(map[string]graph.Vertex)
	// get table information
	tables, err := conn.Query(fmt.Sprintf("SELECT table_name FROM information_schema.tables WHERE table_schema = '%v' ORDER BY table_name DESC;", schema))
	handleErr(err)
	for tables.Next() {
		var tableName string
		var table graph.Vertex
		err = tables.Scan(&tableName)
		handleErr(err)
		table.Name = formatColName(tableName)

		// get column information
		var cols = make(map[string]graph.Col)
		columns, err := conn.Query(fmt.Sprintf("SELECT COLUMN_NAME, DATA_TYPE FROM INFORMATION_SCHEMA.COLUMNS WHERE table_name = '%v' AND table_schema = '%v';", tableName, schema))
		handleErr(err)

		for columns.Next() {
			var colName string
			var colType string
			err = columns.Scan(&colName, &colType)
			handleErr(err)
			cols[formatColName(colName)] = graph.Col{Name: formatColName(colName), Type: convertType(colType)}
		}
		table.Cols = cols
		database.Vertices[formatColName(tableName)] = table
	}
	// get foreign key information
	keys, err := conn.Query("SELECT TABLE_NAME,COLUMN_NAME,REFERENCED_TABLE_NAME,REFERENCED_COLUMN_NAME FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE WHERE TABLE_SCHEMA=? and REFERENCED_TABLE_NAME is not null", schema)
	handleErr(err)
	for keys.Next() {
		var tableName string
		var colName string
		var refTableName string
		var refColName string

		err = keys.Scan(&tableName, &colName, &refTableName, &refColName)
		handleErr(err)

		var originTable = database.Vertices[formatColName(tableName)]
		var destTable = database.Vertices[formatColName(refTableName)]
		var e = graph.Edge{DestinationTable: destTable, DestinationCol: destTable.Cols[formatColName(refColName)], OriginCol: originTable.Cols[formatColName(colName)]}
		originTable.Edges = append(originTable.Edges, e)
	}
	return database
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

// handleErr handles errors in a consistent way
func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
