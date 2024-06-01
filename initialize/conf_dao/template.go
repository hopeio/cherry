package conf_dao

type CloseFunc func() error

type DaoConfig[D any] interface {
	Build() (*D, CloseFunc)
}

type DaoT[C DaoConfig[D], D any] struct {
	Conf   C
	Client *D
	close  CloseFunc
}

func (d *DaoT[C, D]) Config() any {
	return d.Conf
}

func (d *DaoT[C, D]) Set() {
	d.Client, d.close = d.Conf.Build()
}

func (d *DaoT[C, D]) Close() error {
	if d.close != nil {
		return d.close()
	}
	return nil
}
