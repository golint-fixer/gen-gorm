package backends

import (
	"database/sql"
	"fmt"
	"strings"

	// dear gofmt, this is needed
	_ "github.com/go-sql-driver/mysql"
	"github.com/kmulvey/gen-gorm/graph"
	"github.com/kmulvey/gen-gorm/util"
)

// Mysql is what it sounds like
type Mysql struct {
	Backend
	conn  *sql.DB
	model graph.Graph
}

// createConn creates a connection. This cannot be unit tested
func (m *Mysql) createConn(config ConnConfig) *sql.DB {
	conn, err := sql.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8", *config.Username, *config.Password, *config.Hostname, *config.Port, *config.Schema))
	util.HandleErr(err)

	err = conn.Ping()
	util.HandleErr(err)

	return conn
}

// createModel retrieves schema information from a mysql database
func (m *Mysql) createModel(conn *sql.DB, config ConnConfig) (database graph.Graph) {
	database.Name = *config.Schema
	database.Vertices = make(map[string]*graph.Vertex)
	// get table information
	tables, err := conn.Query("SELECT table_name FROM information_schema.tables WHERE table_type='BASE TABLE' and table_schema=?", *config.Schema)
	util.HandleErr(err)
	for tables.Next() {
		var tableName string
		table := &graph.Vertex{}
		err = tables.Scan(&tableName)
		util.HandleErr(err)
		table.Name = formatColName(tableName)
		// get column information
		var cols = make(map[string]graph.Col)
		columns, err := conn.Query("SELECT COLUMN_NAME, DATA_TYPE, CHARACTER_MAXIMUM_LENGTH, EXTRA, COLUMN_KEY FROM INFORMATION_SCHEMA.COLUMNS c inner join INFORMATION_SCHEMA.TABLE_CONSTRAINTS tc on c.table_schema=tc.table_schema and c.table_name=tc.table_name WHERE c.table_name=? AND c.table_schema=?", tableName, *config.Schema)
		// might be better for constrains
		// SELECT non_unique, index_name, seq_in_index, column_name
		// FROM INFORMATION_SCHEMA.STATISTICS
		// WHERE table_name = 'entitlement'
		// AND table_schema = 'auth_service'
		// order by index_name, seq_in_index

		util.HandleErr(err)

		for columns.Next() {
			var colName string
			var colType string
			var colMaxLen sql.NullInt64
			var colExtra string
			var colKey string
			//var conType string
			err = columns.Scan(&colName, &colType, &colMaxLen, &colExtra, &colKey)
			util.HandleErr(err)
			var autoInc bool
			if strings.Contains(colExtra, "auto_increment") {
				autoInc = true
			}
			if colKey == "MUL" {
				colKey = "MULTIPLE"
			} else if colKey == "PRI" {
				colKey = "primary_key"
			}
			cols[formatColName(colName)] = graph.Col{Name: formatColName(colName), Type: convertType(colType), MaxLen: colMaxLen, AutoInc: autoInc, Key: colKey}
		}
		table.Cols = cols
		database.Vertices[formatColName(tableName)] = table
	}
	// get foreign key information
	keys, err := conn.Query("SELECT TABLE_NAME,COLUMN_NAME,REFERENCED_TABLE_NAME,REFERENCED_COLUMN_NAME FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE WHERE TABLE_SCHEMA=? and REFERENCED_TABLE_NAME is not null", *config.Schema)
	util.HandleErr(err)
	for keys.Next() {
		var tableName string
		var colName string
		var refTableName string
		var refColName string

		err = keys.Scan(&tableName, &colName, &refTableName, &refColName)
		util.HandleErr(err)

		var originTable = database.Vertices[formatColName(tableName)]
		var destTable = database.Vertices[formatColName(refTableName)]
		var originKey = originTable.Cols[formatColName(colName)].Key

		if originKey == "MULTIPLE" {
			destTable.HasMany = originTable.Name
		}

		dC := destTable.Cols[formatColName(refColName)]
		oC := originTable.Cols[formatColName(colName)]
		var e = graph.Edge{DestinationTable: destTable, DestinationCol: &dC, OriginCol: &oC}
		originTable.Edges = append(originTable.Edges, e)
	}
	return database
}
