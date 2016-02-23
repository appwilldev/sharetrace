package models

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"

	"github.com/appwilldev/sharetrace/conf"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
)

var (
	dbEngineDefault    *xorm.Engine
	dbEngineDefaultRaw *sql.DB
)

type ModelSession struct {
	*xorm.Session
}

func init() {
	var err error
	dsn := fmt.Sprintf("user=%s dbname=%s host=%s port=%d sslmode=disable",
		conf.DatabaseConfig.User,
		conf.DatabaseConfig.DBName,
		conf.DatabaseConfig.Host,
		conf.DatabaseConfig.Port)

	dbEngineDefault, err = xorm.NewEngine("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to init db engine: " + err.Error())
	}
	dbEngineDefault.SetMaxOpenConns(100)
	dbEngineDefault.SetMaxIdleConns(50)
	dbEngineDefault.ShowErr = true
	//dbEngineDefault.ShowSQL = conf.DebugMode
	dbEngineDefault.ShowSQL = true

	dbEngineDefaultRaw, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to init db engine: " + err.Error())
	}
	dbEngineDefaultRaw.SetMaxIdleConns(10)
	dbEngineDefaultRaw.SetMaxOpenConns(20)
}

func NewModelSession() *ModelSession {
	ms := new(ModelSession)
	ms.Session = dbEngineDefault.NewSession()

	return ms
}

func newAutoCloseModelsSession() *ModelSession {
	ms := new(ModelSession)
	ms.Session = dbEngineDefault.NewSession()
	ms.IsAutoClose = true

	return ms
}

type DBModel interface {
	TableName() string
	UniqueCond() (string, []interface{})
}

func InsertDBModel(s *ModelSession, m DBModel) (err error) {
	if s == nil {
		s = newAutoCloseModelsSession()
	}
	_, err = s.InsertOne(m)

	return
}

func UpdateDBModel(s *ModelSession, m DBModel) (err error) {
	whereStr, whereArgs := m.UniqueCond()
	if s == nil {
		s = newAutoCloseModelsSession()
	}

	_, err = s.Where(whereStr, whereArgs...).Update(m)

	return
}

func DeleteDBModel(s *ModelSession, m DBModel) (err error) {
	whereStr, whereArgs := m.UniqueCond()

	if s == nil {
		s = newAutoCloseModelsSession()
	}

	_, err = s.Where(whereStr, whereArgs...).Delete(m)

	return
}

func rawSqlQuery(sqlStr string, columnTypes []reflect.Type, queryArgs ...interface{}) ([][]interface{}, error) {
	rows, err := dbEngineDefaultRaw.Query(sqlStr, queryArgs...)
	if err != nil {
		return nil, err
	}

	res := make([][]interface{}, 0)
	scanDestValue := make([]interface{}, len(columnTypes))
	for i, columnType := range columnTypes {
		scanDestValue[i] = reflect.New(columnType).Interface()
	}

	for rows.Next() {
		err = rows.Scan(scanDestValue...)
		if err != nil {
			return nil, err
		}

		scanRow := make([]interface{}, len(columnTypes))
		for i := range scanRow {
			scanRow[i] = reflect.Indirect(reflect.ValueOf(scanDestValue[i])).Interface()
		}

		res = append(res, scanRow)
	}

	return res, nil
}

func RawSqlQuery(sqlStr string, columnTypes []reflect.Type, queryArgs ...interface{}) ([][]interface{}, error) {
	rows, err := dbEngineDefaultRaw.Query(sqlStr, queryArgs...)
	if err != nil {
		return nil, err
	}

	scanDestValue := make([]interface{}, len(columnTypes))
	for i, columnType := range columnTypes {
		scanDestValue[i] = reflect.New(columnType).Interface()
	}

	res := make([][]interface{}, 0)
	for rows.Next() {

		err = rows.Scan(scanDestValue...)
		if err != nil {
			//TODO exclude null error
			//log.Println("----", "Null error?", err.Error())
			//return nil, err
		}

		scanRow := make([]interface{}, len(columnTypes))
		for i := range scanRow {
			scanRow[i] = reflect.Indirect(reflect.ValueOf(scanDestValue[i])).Interface()
		}
		res = append(res, scanRow)
	}

	return res, nil
}

func generateSequenceValue(sequenceName string) (int64, error) {
	var sql = fmt.Sprintf("SELECT nextval('%s')", sequenceName)
	columnTypes := []reflect.Type{reflect.TypeOf(int64(1))}

	res, err := rawSqlQuery(sql, columnTypes)
	if err != nil {
		log.Println("gen %s sequence error: %s", sequenceName, err.Error())
		return 0, fmt.Errorf("gen %s sequence error: %s", sequenceName, err.Error())
	}
	if len(res) == 0 {
		log.Println("gen %s sequence error: failed to increase id", sequenceName)
		return 0, fmt.Errorf("gen %s sequence error: failed to increase id", sequenceName)
	}

	return res[0][0].(int64), nil
}
