package backends

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/kmulvey/gen-gorm/graph"
	"github.com/kmulvey/gen-gorm/util"
)

type Mysql struct {
	Backend
	conn  *sql.DB
	model graph.Graph
}

// createConn
func (m *Mysql) createConn(config ConnConfig) *sql.DB {
	conn, err := sql.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8", *config.Username, *config.Password, *config.Hostname, *config.Port, *config.Schema))
	if err != nil {
		log.Fatal(err)
	}
	//defer conn.Close()

	err = conn.Ping()
	util.HandleErr(err)

	return conn
}

// createModel retrieves schema information from a mysql database
func (m *Mysql) createModel(conn *sql.DB, config ConnConfig) (database graph.Graph) {
	database.Name = *config.Schema
	database.Vertices = make(map[string]graph.Vertex)
	// get table information
	tables, err := conn.Query(fmt.Sprintf("SELECT table_name FROM information_schema.tables WHERE table_schema = '%v' ORDER BY table_name DESC;", *config.Schema))
	util.HandleErr(err)
	for tables.Next() {
		var tableName string
		var table graph.Vertex
		err = tables.Scan(&tableName)
		util.HandleErr(err)
		table.Name = formatColName(tableName)
		// get column information
		var cols = make(map[string]graph.Col)
		columns, err := conn.Query(fmt.Sprintf("SELECT COLUMN_NAME, DATA_TYPE FROM INFORMATION_SCHEMA.COLUMNS WHERE table_name = '%v' AND table_schema = '%v';", tableName, *config.Schema))
		util.HandleErr(err)

		for columns.Next() {
			var colName string
			var colType string
			err = columns.Scan(&colName, &colType)
			util.HandleErr(err)
			cols[formatColName(colName)] = graph.Col{Name: formatColName(colName), Type: convertType(colType)}
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
		var e = graph.Edge{DestinationTable: destTable, DestinationCol: destTable.Cols[formatColName(refColName)], OriginCol: originTable.Cols[formatColName(colName)]}
		originTable.Edges = append(originTable.Edges, e)

	}
	return database
}
