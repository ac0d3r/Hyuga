package db

import (
	"reflect"
	"strings"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
)

type Config struct {
	DSN string `mapstructure:"dsn"`
}

type DB struct {
	*leveldb.DB
}

func NewDB(cfg *Config) (*DB, error) {
	db, err := leveldb.OpenFile(cfg.DSN, nil)
	if err != nil {
		return nil, err
	}

	return &DB{DB: db}, nil
}

func (db *DB) Close() {
	db.DB.Close()
}

func (db *DB) Create(m Model) error {
	rv := reflect.ValueOf(m)
	if reflect.TypeOf(m).Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	createdf := rv.FieldByName("Created")
	if createdf.Int() == 0 {
		createdf.SetInt(time.Now().Unix())
	}

	updatedf := rv.FieldByName("Updated")
	if updatedf.Int() == 0 {
		updatedf.SetInt(createdf.Int())
	}

	return db.save(m)
}

func (db *DB) Update(m Model) error {
	rv := reflect.ValueOf(m)
	if reflect.TypeOf(m).Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	updatedf := rv.FieldByName("Updated")
	if updatedf.Int() == 0 {
		updatedf.SetInt(time.Now().Unix())
	}

	return db.save(m)
}

func (db *DB) save(m Model) error {
	d, err := m.encode()
	if err != nil {
		return err
	}
	return db.Put(m.id(), d, nil)
}

func (db *DB) get(m Model) error {
	v, err := db.Get(m.id(), nil)
	if err != nil {
		return err
	}
	if err = m.decode(v); err != nil {
		return err
	}
	return nil
}

type Model interface {
	encode() ([]byte, error)
	decode([]byte) error
	id() []byte
	pre() []byte
}

type BaseModel struct {
	Created int64 `json:"created"`
	Updated int64 `json:"updated"`
}

func modelName(a any) string {
	rt := reflect.TypeOf(a)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	return strings.ToLower(rt.Name())
}
