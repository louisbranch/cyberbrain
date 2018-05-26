package sqlite

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/luizbranco/srs/web"
)

func TestDBInterface(t *testing.T) {
	var _ web.Database = &Database{}
}

func testDB() (*Database, string) {
	tmpfile, err := ioutil.TempFile("", "srs_db")
	if err != nil {
		log.Fatal(err)
	}
	name := tmpfile.Name()
	db, err := New(name)
	if err != nil {
		log.Fatal(err)
	}
	return db, name
}
