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
	dumper
	Sets    []Set
	Creater CreateSet
	Matcher MatchSet
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

func init() {
	util.InitLogger()
}

func (c *Merger) Dump() string {
	c.sb.WriteString("MERGE ")
	if len(c.Sets) > 0 {
		for i := 0; i < len(c.Sets); i++ {
			c.sb.WriteString("(")
			c.Sets[i].Dump(&c.sb)
			c.sb.WriteString(")")
			if i < len(c.Sets)-1 {
				c.sb.WriteString(",")
			}
		}
	}

	c.Creater.Dump(&c.sb)
	c.Matcher.Dump(&c.sb)

	if __IS_DEBUG__ {
		util.Logger.Infof("Merge:" + c.sb.String())
	}
	return c.sb.String()
}

func (c *Set) Dump(sb *strings.Builder) {
	sb.WriteString(fmt.Sprintf("%s:%s ", c.Alias, c.Name))
	if c.Q != nil {
		sb.WriteString(fmt.Sprintf("{%s:%s}", c.Q.Field, c.Q.Val))
	}
}

func (c *valORM) Dump(sb *strings.Builder) {
	if c.ValObj != nil {
		t := reflect.TypeOf(c.ValObj)

		if __IS_DEBUG__ {
			for i := 0; i < t.NumField(); i++ {
				util.Logger.Debugf("field' name is %s, type is %s, kind is %s\n", t.Field(i).Name, t.Field(i).Type, t.Field(i).Type.Kind())
			}
		}

		if t.NumField() > 0 {
			for i := 0; i < t.NumField(); i++ {
				val := reflect.ValueOf(c.ValObj).FieldByName(t.Field(i).Name)
				if len(t.Field(i).Tag) > 0 && t.Field(i).Tag.Get("omitempty") == "true" {
					if __IS_DEBUG__ {
						util.Logger.Infof("omit zero field:" + t.Field(i).Name)
					}
				} else {
					sb.WriteString(fmt.Sprintf("i.%s = %v, ", t.Field(i).Name, val))
				}
			}
		}
	}
}

func (c *MatchSet) Dump(sb *strings.Builder) {
	sb.WriteString("ON MATCH SET ")
	c.valORM.Dump(sb)
}

func (c *CreateSet) Dump(sb *strings.Builder) {
	sb.WriteString("ON CREATE SET ")
	c.valORM.Dump(sb)
}
