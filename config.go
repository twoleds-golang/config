package config

import "strconv"
import "strings"

// Config represents a node in a hierarchical configuration data.
// Configuration data are organizes as a tree.
type Config interface {
	// Bool returns a boolean value for the specified query.
	Bool(query string) (val bool, found bool)
	// BoolOrDefault returns a boolean value for the specified query if match.
	// Otherwise returns the default value.
	BoolOrDefault(query string, defVal bool) (val bool)
	// Float returns a float number for the specified query.
	Float(query string) (val float64, found bool)
	// FloatOrDefault returns a float number for the specified query if match.
	// Otherwise returns the default value.
	FloatOrDefault(query string, defVal float64) (val float64)
	// Int returns a integer value for the specified query.
	Int(query string) (val int64, found bool)
	// IntOrDefault returns a integer value for the specified query if match.
	// Otherwise returns the default value.
	IntOrDefault(query string, defVal int64) (val int64)
	// Name returns name of this configuration node.
	Name() string
	// Query returns a configuration node for the specified query.
	Query(query string) (cfg Config, found bool)
	// Query returns all configuration nodes which match the specified query.
	QueryAll(query string) (cfgs []Config)
	// String returns a string value for the specified query.
	String(query string) (val string, found bool)
	// StringOrDefault returns a string value for the specified query if match.
	// Otherwise returns the default value.
	StringOrDefault(query string, defVal string) (val string)
	// Value returns value of this configuration node.
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

func (this *config) queryAllLoop(query []string, conds []string, level int, cfgs []Config) []Config {
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
