package metastg

import (
	"soloos/common/log"
)

func (p *MetaStg) installSchema() error {
	var (
		sqls []string
		err  error
	)

	sqls = p.prepareNetINodesSqls()
	for _, sql := range sqls {
		_, err = p.dbConn.Exec(sql)
		if err != nil {
			log.Error(err, sql)
		}
	}

	return nil
}
