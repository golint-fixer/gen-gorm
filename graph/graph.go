package graph

import (
	"database/sql"
	"reflect"
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

	// index?
	if c.Key != "" {
		result += c.Key + " "
	}

	if c.AutoInc {
		result += `"AUTO_INCREMENT";`
	}

	result = strings.TrimRight(result, ";")
	return result + "`"
}

func (c Col) GetEdge(edges []Edge) (result Edge) {
	for _, e := range edges {
		if reflect.DeepEqual(e.OriginCol, &c) {
			result = e
			_ = "breakpoint"
		}
	}
	return result
}

func (c Col) HasEdge(edges []Edge) bool {
	for _, e := range edges {
		if e.OriginCol.Name == c.Name { // not good
			return true
		}
	}
	return false
}
