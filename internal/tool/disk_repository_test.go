package tool

import (
	"database/sql"
	"testing"

	"github.com/sebdah/goldie/v2"

	gomocket "github.com/Selvatico/go-mocket"
	"github.com/jinzhu/gorm"
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

func TestDiskRepo(t *testing.T) {
	suite.Run(t, new(diskRepoTest))
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
	r.gd = goldie.New(t)
}

func (r *diskRepoTest) TearDownTest() {
	_ = r.gormDB.Close()
	_ = r.sqlDb.Close()
	gomocket.Catcher.Reset()
}

func (r *diskRepoTest) TestDiskRepository_GetAll() {
	gomocket.Catcher.NewMock().
		WithQuery("SELECT * FROM `tools`  WHERE `tools`.`deleted_at` IS NULL").
		WithReply([]map[string]interface{}{
			{
				"id":   1,
				"name": "rectangle",
			},
		})
	repo := NewDiskRepository(r.gormDB)
	tools, err := repo.GetAll()
	r.NoError(err)
	r.NotNil(tools)
	r.gd.AssertJson(r.T(), "get_all_successfully", tools)
}
