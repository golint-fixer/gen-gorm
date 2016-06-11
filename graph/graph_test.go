package graph

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMeta(t *testing.T) {
	t.Parallel()
	var colOne = Col{Name: "Id", Type: "int"}
	var colTwo = Col{Name: "Id", Type: "int", Key: "primary_key", AutoInc: true}
	var colThree = Col{Name: "Name", Type: "string", Key: "", MaxLen: sql.NullInt64{Int64: 45, Valid: true}}
	assert.Equal(t, "", colOne.GetMeta(), "should be equal")
	assert.Equal(t, "`gorm:primary_key \"AUTO_INCREMENT\"`", colTwo.GetMeta(), "should be equal")
	assert.Equal(t, "`gorm:size:45`", colThree.GetMeta(), "should be equal")
}

func TestHasEdge(t *testing.T) {
	t.Parallel()
	var colOne = Col{Name: "Id", Type: "int", Key: "primary_key", MaxLen: sql.NullInt64{Int64: 0, Valid: false}}
	var colTwo = Col{Name: "Name", Type: "string", Key: "", MaxLen: sql.NullInt64{Int64: 45, Valid: true}}
	assert.Equal(t, "", colOne.GetMeta(), "should be equal")
	assert.Equal(t, "`gorm:size:45`", colTwo.GetMeta(), "should be equal")
}
