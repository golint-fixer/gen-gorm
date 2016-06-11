package graph

import (
	"database/sql"
	"testing"

	"github.com/kmulvey/gen-gorm/graph"
	"github.com/stretchr/testify/assert"
)

func TestGetMeta(t *testing.T) {
	t.Parallel()
	var colOne = graph.Col{Name: "Id", Type: "int", Key: "primary_key", MaxLen: sql.NullInt64{Int64: 0, Valid: false}}
	var colTwo = graph.Col{Name: "Name", Type: "string", Key: "", MaxLen: sql.NullInt64{Int64: 45, Valid: true}}
	assert.Equal(t, "", colOne.GetMeta(), "should be equal")
	assert.Equal(t, "`gorm:size:45`", colTwo.GetMeta(), "should be equal")
}
