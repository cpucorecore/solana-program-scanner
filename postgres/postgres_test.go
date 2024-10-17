package postgres_test

//
//import (
//	"fmt"
//	"github.com/lib/pq"
//	_ "github.com/lib/pq"
//	"github.com/stretchr/testify/require"
//	"inj_extractor/config"
//	"testing"
//	"time"
//	"xorm.io/xorm"
//)
//
//var (
//	dataSourceName string
//	engine         *xorm.Engine
//)
//
//const (
//	driverName = "postgres"
//)
//
//func init() {
//	config.LoadFromYaml("../config.yaml")
//	dataSourceName = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
//		config.G.Postgres.User,
//		config.G.Postgres.Passwd,
//		config.G.Postgres.Host,
//		config.G.Postgres.Port,
//		config.G.Postgres.Db)
//
//	var err error
//	engine, err = xorm.NewEngine(driverName, dataSourceName)
//	if err != nil {
//		panic(err)
//	}
//
//	engine.DropTables(TestUser{})
//	err = engine.Sync(TestUser{})
//	if err != nil {
//		panic(err)
//	}
//}
//
//type TestUser struct {
//	Id        int64
//	Name      string    `xorm:"unique not null"`
//	CreatedAt time.Time `xorm:"created"`
//	UpdatedAt time.Time `xorm:"updated"`
//}
//
//func TestPostgresUniqueConstrain(t *testing.T) {
//	effected, err := engine.Insert(TestUser{Name: "abc11"})
//	require.Nil(t, err)
//	t.Log(effected)
//	effected, err = engine.Insert(TestUser{Name: "abc11"})
//	if err != nil {
//		if pgErr, ok := err.(*pq.Error); ok {
//			require.Equal(t, "23505", string(pgErr.Code))
//		} else {
//			t.Error(err)
//		}
//	}
//	t.Log(effected)
//
//	engine.DropTables(TestUser{})
//}
