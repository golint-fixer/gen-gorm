package graph

import (
	"database/sql"
	"strconv"
	"strings"
)

// Graph is a databse
type Graph struct {
	Vertices map[string]*Vertex
	Name     string
}

// Vertex is a table
type Vertex struct {
	Name    string
	Cols    map[string]Col
	HasMany string
	Edges   []Edge
}

// Col is a column ... duh?!?
type Col struct {
	Name       string
	Type       string
	MaxLen     sql.NullInt64
	AutoInc    bool
	Key        string
	Constraint string
}

// Edge is a foreign key
type Edge struct {
	DestinationTable *Vertex
	DestinationCol   *Col
	OriginCol        *Col
}

// GetMeta generates the gorm meta string
func (c Col) GetMeta() string {
	// bail
	if !c.MaxLen.Valid && !c.AutoInc && c.Constraint == "" || c.Constraint == "FOREIGN KEY" {
		return ""
	}
	var result = "`gorm:"

	if c.MaxLen.Valid {
		result += "size:" + strconv.FormatInt(c.MaxLen.Int64, 10) + ";"
	}
	if c.Constraint != "" && c.Constraint != "FOREIGN KEY" {
		result += c.Constraint + ";"
	}
	if c.AutoInc {
		result += `"AUTO_INCREMENT";`
	}

	result = strings.TrimRight(result, ";")
	return result + "`"
}
