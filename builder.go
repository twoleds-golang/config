package config

import "strconv"

type Builder interface {
	Bool(name string, val bool) Builder
	CloseSection() Builder
	Config() Config
	Float(name string, val float64) Builder
	Int(name string, val int64) Builder
	Section(name string, val string) Builder
	String(name string, val string) Builder
}

type builder struct {
	stack []*config
}

func NewBuilder() Builder {
	b := new(builder)
	b.stack = make([]*config, 1, 16)
	b.stack[0] = b.create("", "", true)
	return b
}

func (this *builder) append(cfg *config) *config {
	cur := this.current()
	cur.children = append(cur.children, cfg)
	return cfg
}

func (this *builder) create(name string, val string, isSection bool) *config {
	cfg := new(config)
	cfg.name = name
	cfg.value = val
	if isSection {
		cfg.children = make([]*config, 0, 16)
	}
	return cfg
}

func (this *builder) current() *config {
	return this.stack[len(this.stack)-1]
}

func (this *builder) Bool(name string, val bool) Builder {
	this.append(this.create(name, strconv.FormatBool(val), false))
	return this
}

func (this *builder) CloseSection() Builder {
	this.stack = this.stack[0 : len(this.stack)-1]
	return this
}

func (this *builder) Config() Config {
	return this.stack[0]
}

func (this *builder) Float(name string, val float64) Builder {
	this.append(this.create(name, strconv.FormatFloat(val, 'g', -1, 64), false))
	return this
}

func (this *builder) Int(name string, val int64) Builder {
	this.append(this.create(name, strconv.FormatInt(val, 10), false))
	return this
}

func (this *builder) Section(name string, val string) Builder {
	this.stack = append(this.stack, this.append(this.create(name, val, true)))
	return this
}

func (this *builder) String(name string, val string) Builder {
	this.append(this.create(name, val, false))
	return this
}
