package metastg

import (
	"github.com/gocraft/dbr"
)

type MetaStg struct {
	DBConn *dbr.Connection
}

func (p *MetaStg) Init(dbDriver, dsn string) error {
	var err error

	p.DBConn, err = dbr.Open(dbDriver, dsn, nil)
	if err != nil {
		return err
	}

	switch dbDriver {
	case "mysql":
		err = p.InstallMysqlSchema()
	case "sqlite3":
		err = p.InstallSqlite3Schema()
	}

	return nil
}

func (p *MetaStg) Close() error {
	var err error

	err = p.DBConn.Close()
	if err != nil {
		return err
	}

	return nil
}