package conf_dao

// template
type Config[T any, V any] struct {
	Config *T
}

func (c *Config[T, V]) Init() {

}

func (c *Config[T, V]) Build() *V {
	return new(V)
}

type Dao[T, V any] struct {
	Conf   Config[T, V]
	Entity *V
}

func (d *Dao[T, V]) Config() any {
	return d.Conf
}

func (d *Dao[T, V]) SetEntity() {
	d.Entity = d.Conf.Build()
}

func (d *Dao[T, V]) Close() error {
	return nil
}
