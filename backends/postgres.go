package backends

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/kmulvey/gen-gorm/graph"
	"github.com/kmulvey/gen-gorm/util"
	// dear gofmt, this is needed
	_ "github.com/lib/pq"
)

// Postgres is what it sounds like
type Postgres struct {
	Backend
	conn  *sql.DB
	model graph.Graph
}

// createConn creates a connection. This cannot be unit tested
func (m *Postgres) createConn(config ConnConfig) *sql.DB {
	conn, err := sql.Open("postgres", fmt.Sprintf("host=%v user=%v dbname=%v password=%v port=%v sslmode=disable", *config.Hostname, *config.Username, *config.Schema, *config.Password, *config.Port))
	if err != nil {
		log.Fatal(err)
	}

	err = conn.Ping()
	util.HandleErr(err)

	return conn
}

// createModel retrieves schema information from a mysql database
func (m *Postgres) createModel(conn *sql.DB, config ConnConfig) (database graph.Graph) {
	database.Name = *config.Schema
	database.Vertices = make(map[string]*graph.Vertex)
	// get table information
	tables, err := conn.Query("SELECT tablename FROM pg_catalog.pg_tables where tableowner=$1", *config.Username)
	util.HandleErr(err)
	for tables.Next() {
		var tableName string
		table := &graph.Vertex{}
		err = tables.Scan(&tableName)
		util.HandleErr(err)
		table.Name = formatColName(tableName)
		// get column information
		var cols = make(map[string]graph.Col)
		columns, err := conn.Query("select kc.column_name, c.data_type, c.character_maximum_length, tc.constraint_type from information_schema.table_constraints tc inner join information_schema.key_column_usage kc on kc.table_name = tc.table_name and kc.table_schema = tc.table_schema inner join information_schema.columns c on kc.column_name=c.column_name where tc.table_name=$1 and tc.table_catalog=$2", tableName, *config.Schema)
		util.HandleErr(err)

		for columns.Next() {
			var colName string
			var colType string
			var colMaxLen sql.NullInt64
			var colKey string
			err = columns.Scan(&colName, &colType, &colMaxLen, &colKey)
			util.HandleErr(err)
			if colKey == "FOREIGN KEY" {
				colKey = "MULTIPLE"
			}
			cols[formatColName(colName)] = graph.Col{Name: formatColName(colName), Type: convertType(colType), Key: colKey, MaxLen: colMaxLen}
		}
		table.Cols = cols
		database.Vertices[formatColName(tableName)] = table
	}
	// get foreign key information
	//tc.constraint_type,
	keys, err := conn.Query(`select 
		tc.table_name,
		kcu.column_name,
		ccu.table_name AS references_table,
		ccu.column_name AS references_field
		FROM information_schema.table_constraints tc
		LEFT JOIN information_schema.key_column_usage kcu
		ON tc.constraint_catalog = kcu.constraint_catalog
		AND tc.constraint_schema = kcu.constraint_schema
		AND tc.constraint_name = kcu.constraint_name
		LEFT JOIN information_schema.referential_constraints rc
		ON tc.constraint_catalog = rc.constraint_catalog
		AND tc.constraint_schema = rc.constraint_schema
		AND tc.constraint_name = rc.constraint_name
		LEFT JOIN information_schema.constraint_column_usage ccu
		ON rc.unique_constraint_catalog = ccu.constraint_catalog
		AND rc.unique_constraint_schema = ccu.constraint_schema
		AND rc.unique_constraint_name = ccu.constraint_name
		WHERE lower(tc.constraint_type) in ('foreign key')
		and tc.table_catalog=$1`, *config.Schema)

	util.HandleErr(err)
	for keys.Next() {
		var tableName string
		var colName string
		//var constraintType string
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
