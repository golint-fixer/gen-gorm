package backends

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertType(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "string", convertType("time"), "should be string")
	assert.Equal(t, "string", convertType("datetime"), "should be string")
	assert.Equal(t, "string", convertType("char"), "should be string")
	assert.Equal(t, "string", convertType("varchar"), "should be string")
	assert.Equal(t, "string", convertType("blob"), "should be string")
	assert.Equal(t, "int", convertType("integer"), "should be integer")
	assert.Equal(t, "int", convertType("int"), "should be integer")
	assert.Equal(t, "int", convertType("timestamp"), "should be integer")
	assert.Equal(t, "float", convertType("float"), "should be float")
	assert.Equal(t, "float", convertType("double"), "should be float")
	assert.Equal(t, "bool", convertType("boolean"), "should be bool")
	assert.Equal(t, "string", convertType("enum"), "should be string")
	assert.Equal(t, "string", convertType("other"), "should be string")
}

func TestCapFirst(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "String", capFirst("string"), "should be \"String\"")
	assert.Equal(t, "Int", capFirst("int"), "should be \"Int\"")
}

func TestTormatColName(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "Id", formatColName("id"), "should be Id")
	assert.Equal(t, "UserName", formatColName("user_name"), "should be UserName")
	assert.Equal(t, "Username", formatColName("username"), "should be Username")
	assert.Equal(t, "Col5", formatColName("col_5"), "should be Col5")
	assert.Equal(t, "Code", formatColName("code"), "should be Code")
}
