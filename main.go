package main

import (
	"flag"
	"html/template"
	"os"
	"os/exec"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kmulvey/gen-gorm/backends"
	"github.com/kmulvey/gen-gorm/graph"
	"github.com/kmulvey/gen-gorm/util"

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
	var dbConfig = backends.ConnConfig{
		Hostname: flag.String("hostname", "", "hostname"),
		Username: flag.String("username", "", "username"),
		Password: flag.String("password", "", "password"),
		Schema:   flag.String("schema", "", "schema"),
		Port:     flag.String("port", "3306", "port"),
	}
	output := flag.String("output", "", "output")
	engine := flag.String("engine", "", "engine")
	flag.Parse()

	// connect to db
	//db, err := sql.Open("postgres", fmt.Sprintf("host=%v user=%v dbname=%v password=%v port=%v sslmode=disable", hostname, username, schema, password, port))

	// get table structure from DB
	data := backends.GetTableInfo(dbConfig, *engine)

	// get 'er done
	processTemplates(data, *output)
}

// processTemplates fills in the templates with data, puts them in the output
// directory and fmt them
func processTemplates(data graph.Graph, output string) {

	// parse templates
	modelsTemplate, err := template.ParseFiles("templates/models.tmpl")
	util.HandleErr(err)

	// create directory
	if _, err := os.Stat(output); os.IsNotExist(err) {
		os.Mkdir(output, 0755)
	}

	// create the files
	modelsGo, err := os.Create(output + "/models.go")
	util.HandleErr(err)
	defer modelsGo.Close()

	// exec templates
	err = modelsTemplate.Execute(modelsGo, data)
	util.HandleErr(err)

	// format the file
	cmd := exec.Command("gofmt", "-w", output)
	err = cmd.Run()
	util.HandleErr(err)
}
