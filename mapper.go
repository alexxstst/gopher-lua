package lua

import (
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"regexp"
	"strings"
)

// Option is a configuration that is used to create a new mapper.
type Option struct {
	// Function to convert a lua table key to Go's one. This defaults to "ToUpperCamelCase".
	NameFunc func(string) string

	// Returns error if unused keys exist.
	ErrorUnused bool

	// A struct tag name for lua table keys . This defaults to "gluamapper"
	TagName string
}

// Mapper maps a lua table to a Go struct pointer.
type Mapper struct {
	Option Option
}

// NewMapper returns a new mapper.
func NewMapper(opt Option) *Mapper {
	if opt.NameFunc == nil {
		opt.NameFunc = ToUpperCamelCase
	}
	if opt.TagName == "" {
		opt.TagName = "gluamapper"
	}
	return &Mapper{opt}
}

// Map maps the lua table to the given struct pointer.
func (mapper *Mapper) Map(tbl *LTable, st interface{}) error {
	opt := mapper.Option
	mp, ok := ToGoValue(tbl, opt).(map[interface{}]interface{})
	if !ok {
		return errors.New("arguments #1 must be a table, but got an array")
	}
	config := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           st,
		TagName:          opt.TagName,
		ErrorUnused:      opt.ErrorUnused,
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	return decoder.Decode(mp)
}

// Map maps the lua table to the given struct pointer with default options.
func Map(tbl *LTable, st interface{}) error {
	return NewMapper(Option{}).Map(tbl, st)
}


// Id is an Option.NameFunc that returns given string as-is.

var camelre = regexp.MustCompile(`_([a-z])`)
// ToUpperCamelCase is an Option.NameFunc that converts strings from snake case to upper camel case.
func ToUpperCamelCase(s string) string {
	return strings.ToUpper(string(s[0])) + camelre.ReplaceAllStringFunc(s[1:], func(s string) string { return strings.ToUpper(s[1:]) })
}

// ToGoValue converts the given LValue to a Go object.
func ToGoValue(lv LValue, opt Option) interface{} {
	switch v := lv.(type) {
	case *LNilType:
		return nil
	case LBool:
		return bool(v)
	case LString:
		return string(v)
	case LNumber:
		return float64(v)
	case *LTable:
		maxn := v.MaxN()
		if maxn == 0 { // table
			ret := make(map[interface{}]interface{})
			v.ForEach(func(key, value LValue) {
				keystr := fmt.Sprint(ToGoValue(key, opt))
				ret[opt.NameFunc(keystr)] = ToGoValue(value, opt)
			})
			return ret
		} else { // array
			ret := make([]interface{}, 0, maxn)
			for i := 1; i <= maxn; i++ {
				ret = append(ret, ToGoValue(v.RawGetInt(i), opt))
			}
			return ret
		}
	default:
		return v
	}
}
