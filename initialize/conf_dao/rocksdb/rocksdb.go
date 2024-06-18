package rocksdb

import (
	"github.com/linxGnu/grocksdb"
)

type Config struct {
	Path            string
	Capacity        uint64
	CreateIfMissing bool
	ErrorIfExists   bool
	ParanoidChecks  bool
	Paths           []string
	TargetSizes     []uint64
}

func (c *Config) InitBeforeInject() {

}

func (c *Config) Init() {
}

func (c *Config) Build() (*grocksdb.DB, error) {
	bbto := grocksdb.NewDefaultBlockBasedTableOptions()
	bbto.SetBlockCache(grocksdb.NewLRUCache(c.Capacity))

	opts := grocksdb.NewDefaultOptions()
	opts.SetBlockBasedTableFactory(bbto)
	opts.SetCreateIfMissing(c.CreateIfMissing)
	opts.SetErrorIfExists(c.ErrorIfExists)
	opts.SetParanoidChecks(c.ParanoidChecks)
	opts.SetDBPaths(grocksdb.NewDBPathsFromData(c.Paths, c.TargetSizes))

	return grocksdb.OpenDb(opts, c.Path)
}

type DB struct {
	*grocksdb.DB
	Conf Config
}

func (m *DB) Config() any {
	return &m.Conf
}

func (m *DB) Init() error {
	var err error
	m.DB, err = m.Conf.Build()
	return err
}

func (m *DB) Close() error {
	m.DB.Close()
	return nil
}
