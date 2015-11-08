package config

import "strconv"
import "strings"

type Config interface {
	Bool(query string) (val bool, found bool)
	BoolOrDefault(query string, defVal bool) (val bool)
	Float(query string) (val float64, found bool)
	FloatOrDefault(query string, defVal float64) (val float64)
	Int(query string) (val int64, found bool)
	IntOrDefault(query string, defVal int64) (val int64)
	Name() string
	Query(query string) (cfg Config, found bool)
	QueryAll(query string) (cfgs []Config)
	String(query string) (val string, found bool)
	StringOrDefault(query string, defVal string) (val string)
	Value() string
}

type config struct {
	name     string
	value    string
	children []*config
}

var _ Config = new(config)

func (this *config) Bool(query string) (val bool, found bool) {
	if str, ok := this.String(query); ok {
		if val, err := strconv.ParseBool(str); err == nil {
			return val, true
		}
	}
	return false, false
}

func (this *config) BoolOrDefault(query string, defVal bool) (val bool) {
	if val, ok := this.Bool(query); ok {
		return val
	}
	return defVal
}

func (this *config) Float(query string) (val float64, found bool) {
	if str, ok := this.String(query); ok {
		if val, err := strconv.ParseFloat(str, 64); err == nil {
			return val, true
		}
	}
	return 0.0, false
}

func (this *config) FloatOrDefault(query string, defVal float64) (val float64) {
	if val, ok := this.Float(query); ok {
		return val
	}
	return defVal
}

func (this *config) Int(query string) (val int64, found bool) {
	if str, ok := this.String(query); ok {
		if val, err := strconv.ParseInt(str, 10, 64); err == nil {
			return val, true
		}
	}
	return 0, false
}

func (this *config) IntOrDefault(query string, defVal int64) (val int64) {
	if val, ok := this.Int(query); ok {
		return val
	}
	return defVal
}

func (this *config) Name() string {
	return this.name
}

func (this *config) Query(path string) (cfg Config, found bool) {
	query, conds := this.parse(path)
	return this.queryLoop(query, conds, 0)
}

func (this *config) queryLoop(query []string, conds []string, level int) (cfg *config, found bool) {
	for _, child := range this.children {
		if child.name == query[level] {
			if conds[level] == "*" || conds[level] == child.value {
				if level == (len(query) - 1) {
					return child, true
				} else {
					return child.queryLoop(query, conds, level+1)
				}
			}
		}
	}
	return nil, false
}

func (this *config) QueryAll(path string) (cfgs []Config) {
	query, conds := this.parse(path)
	return this.queryAllLoop(query, conds, 0, make([]Config, 0, 16))
}

func (this *config) queryAllLoop(query []string, conds []string, level int, cfgs []Config) {
	for _, child := range this.children {
		if child.name == query[level] {
			if conds[level] == "*" || conds[level] == child.value {
				if level == (len(query) - 1) {
					cfgs = append(cfgs, child)
				} else {
					cfgs = child.queryAllLoop(query, conds, level+1, cfgs)
				}
			}
		}
	}
	return cfgs
}

func (this *config) parse(path string) (query []string, conds []string) {
	query = strings.Split(path, "/")
	conds = make([]string, len(query))
	for key, part := range query {
		if index := strings.IndexByte(part, ':'); index >= 0 {
			query[key] = part[:index]
			conds[key] = part[index+1:]
		} else {
			conds[key] = "*"
		}
	}
	return query, conds
}

func (this *config) String(query string) (val string, found bool) {
	if cfg, ok := this.Query(query); ok {
		return cfg.Value(), true
	}
	return "", false
}

func (this *config) StringOrDefault(query string, defVal string) (val string) {
	if val, ok := this.String(query); ok {
		return val
	}
	return defVal
}

func (this *config) Value() string {
	return this.value
}
