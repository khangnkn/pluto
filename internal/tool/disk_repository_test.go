package tool

import (
	"database/sql"
	"testing"

	gomocket "github.com/Selvatico/go-mocket"
	"github.com/jinzhu/gorm"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/suite"
)

func TestDiskRepository(t *testing.T) {
	suite.Run(t, new(diskRepoTest))
}

type diskRepoTest struct {
	suite.Suite
	gd goldie.Tester

	gormDB *gorm.DB
	sqlDb  *sql.DB
}

func (r *diskRepoTest) SetupTest() {
	t := r.T()
	gomocket.Catcher.Register()
	gomocket.Catcher.Logging = true
	conn, err := sql.Open(gomocket.DriverName, "connection_string")
	if err != nil {
		t.Fatal(err)
	}
	r.sqlDb = conn
	db, err := gorm.Open("mysql", conn)
	if err != nil {
		t.Fatal("Failed creating gormDB")
	}

	r.gormDB = db
}

func (r *diskRepoTest) TearDownTest() {
	_ = r.gormDB.Close()
	_ = r.sqlDb.Close()
	gomocket.Catcher.Reset()
}

func TestDiskRepository_GetAll(t *testing.T) {

}
