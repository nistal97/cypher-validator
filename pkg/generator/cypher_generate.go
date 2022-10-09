package generator

import (
	"cypher-validator/pkg/util"
	"fmt"
	"reflect"
	"strings"
)

const __IS_DEBUG__ bool = true

type dumper struct {
	sb strings.Builder
}

type Qer struct {
	Field string
	Val   interface{}
}

type Set struct {
	Alias string
	Name  string
	Q     *Qer
}

type Merger struct {
	d            *dumper
	innerMergers []Merger
	Sets         []Set
	Creater      CreateSet
	Matcher      MatchSet
	Unwinder     *Unwind
}

type valORM struct {
	ValObj interface{}
}

type CreateSet struct {
	valORM
}

type MatchSet struct {
	valORM
}

type Unwind struct {
	WithAlias string
	Alias     string
	Pairs     map[string]interface{}
}

func init() {
	util.InitLogger()
}

func (c *dumper) reset() {
	c.sb.Reset()
}

func (c *Merger) Dump() string {
	if len(c.innerMergers) > 0 {
		for _, m := range c.innerMergers {
			m.d = c.d
			m.Dump()
		}
	}

	if c.Unwinder != nil {
		c.Unwinder.Dump(c.d)
	}

	c.d.sb.WriteString("MERGE ")
	if len(c.Sets) > 0 {
		for i := 0; i < len(c.Sets); i++ {
			c.d.sb.WriteString("(")
			c.Sets[i].Dump(c.d)
			c.d.sb.WriteString(")")
			if i < len(c.Sets)-1 {
				c.d.sb.WriteString(",")
			}
		}
	}
	if c.Unwinder == nil {
		c.Creater.Dump(c.d)
		c.Matcher.Dump(c.d)
	} else {
		c.Unwinder.DumpFields(c.d)
	}

	if __IS_DEBUG__ {
		util.Logger.Infof("Merge:" + c.d.sb.String())
	}
	return c.d.sb.String()
}

func (c *Unwind) Dump(d *dumper) {
	if c.WithAlias != "" {
		d.sb.WriteString(fmt.Sprintf(" WITH %s ", c.WithAlias))
	}
	d.sb.WriteString(fmt.Sprintf("UNWIND %s as m ", c.Alias))
}

func (c *Unwind) DumpFields(d *dumper) {
	i := 0
	d.sb.WriteString(" SET ")
	for k, _ := range c.Pairs {
		d.sb.WriteString(fmt.Sprintf("m.%s = %s.%s", k, c.Alias, k))
		if i < len(c.Pairs)-1 {
			d.sb.WriteString(fmt.Sprintf(", "))
		}
		i++
	}
}

func (c *Set) Dump(d *dumper) {
	d.sb.WriteString(fmt.Sprintf("%s:%s ", c.Alias, c.Name))
	if c.Q != nil {
		d.sb.WriteString(fmt.Sprintf("{%s:"+dump_interface_fmt(c.Q.Val)+"}", c.Q.Field, c.Q.Val))
	}
}

func (c *valORM) Dump(d *dumper) {
	if c.ValObj != nil {
		t := reflect.TypeOf(c.ValObj)

		if __IS_DEBUG__ {
			for i := 0; i < t.NumField(); i++ {
				util.Logger.Debugf("field' name is %s, type is %s, kind is %s\n", t.Field(i).Name, t.Field(i).Type, t.Field(i).Type.Kind())
			}
		}

		if t.NumField() > 0 {
			conditions := make([]string, 0)
			for i := 0; i < t.NumField(); i++ {
				val := reflect.ValueOf(c.ValObj).FieldByName(t.Field(i).Name)
				if val.IsZero() && len(t.Field(i).Tag) > 0 && t.Field(i).Tag.Get("omitempty") == "true" {
					if __IS_DEBUG__ {
						util.Logger.Infof("omit zero field:" + t.Field(i).Name)
					}
				} else if val.Type().Kind() == reflect.Struct {
					if __IS_DEBUG__ {
						util.Logger.Infof("omit struct field:" + t.Field(i).Name)
					}
				} else {
					fmt_str := "i.%s = " + dump_val_fmt(val) + " "
					conditions = append(conditions, fmt.Sprintf(fmt_str, t.Field(i).Name, val))
				}
			}
			for i := 0; i < len(conditions); i++ {
				d.sb.WriteString(conditions[i])
				if i < len(conditions)-1 {
					d.sb.WriteString(", ")
				}
			}
		}
	}
}

func (c *MatchSet) Dump(d *dumper) {
	if c.ValObj != nil {
		d.sb.WriteString("ON MATCH SET ")
		c.valORM.Dump(d)
	}
}

func (c *CreateSet) Dump(d *dumper) {
	if c.ValObj != nil {
		d.sb.WriteString("ON CREATE SET ")
		c.valORM.Dump(d)
	}
}

func dump_val_fmt(val reflect.Value) string {
	if val.Type().Kind() == reflect.String {
		return "'%v'"
	}
	return "%v"
}

func dump_interface_fmt(iter interface{}) string {
	switch iter.(type) {
	case string:
		return "'%v'"
	}
	return "%v"
}
