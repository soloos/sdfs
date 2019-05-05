package metastg

import (
	"database/sql"
	"soloos/common/sdbapi"
	"soloos/sdfs/types"
)

func (p *NetINodeDriver) NetINodeCommitSizeInDB(uNetINode types.NetINodeUintptr, size uint64) error {
	var (
		sess sdbapi.Session
		err  error
	)

	err = p.dbConn.InitSession(&sess)
	if err != nil {
		return err
	}

	_, err = sess.Update("b_netinode").
		Set("netinode_size", size).
		Where("netinode_id=?", uNetINode.Ptr().IDStr()).
		Exec()
	if err != nil {
		return err
	}

	uNetINode.Ptr().Size = size

	return nil
}

func (p *NetINodeDriver) FetchNetINodeFromDB(pNetINode *types.NetINode) error {
	var (
		sess    sdbapi.Session
		sqlRows *sql.Rows
		err     error
	)

	err = p.dbConn.InitSession(&sess)
	if err != nil {
		goto QUERY_DONE
	}

	sqlRows, err = sess.Select("netinode_size", "netblock_cap", "memblock_cap").
		From("b_netinode").
		Where("netinode_id=?", pNetINode.IDStr()).Rows()
	if err != nil {
		goto QUERY_DONE
	}

	if sqlRows.Next() == false {
		err = types.ErrObjectNotExists
		goto QUERY_DONE
	}

	err = sqlRows.Scan(&pNetINode.Size, &pNetINode.NetBlockCap, &pNetINode.MemBlockCap)
	if err != nil {
		goto QUERY_DONE
	}

QUERY_DONE:
	if sqlRows != nil {
		sqlRows.Close()
	}
	return err
}

func (p *NetINodeDriver) StoreNetINodeInDB(pNetINode *types.NetINode) error {
	var (
		sess          sdbapi.Session
		netINodeIDStr = pNetINode.IDStr()
		err           error
	)

	err = p.dbConn.InitSession(&sess)
	if err != nil {
		return err
	}

	err = sess.ReplaceInto("b_netinode").
		PrimaryColumns("netinode_id").PrimaryValues(netINodeIDStr).
		Columns("netinode_size", "netblock_cap", "memblock_cap").
		Values(pNetINode.Size, pNetINode.NetBlockCap, pNetINode.MemBlockCap).
		Exec()
	if err != nil {
		return err
	}

	return nil
}
