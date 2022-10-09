package generator

import (
	"cypher-validator/pkg/parser/neo4j/hack/neo4j"
	"cypher-validator/pkg/util"
	"testing"
)

type b struct {
	bool_f bool
}

type a struct {
	int_f   int
	float_f float32
	str_f   string `omitempty:"true"`
	b       b
}

func TestMerge(t *testing.T) {
	defer util.Logger.Sync()

	var c Merger
	c.d = &dumper{}
	c.Creater = CreateSet{
		valORM: valORM{
			ValObj: a{
				int_f:   1,
				float_f: 2.0,
				str_f:   "abc",
				b: b{
					bool_f: true,
				},
			},
		},
	}
	c.Matcher = MatchSet{
		valORM: valORM{
			ValObj: a{
				int_f: 2,
			},
		},
	}
	c.Sets = make([]Set, 1)
	c.Sets[0] = Set{
		Alias: "i",
		Name:  "Set1",
		Q: &Qer{
			Field: "id",
			Val:   "val",
		},
	}

	parser.ValidateCypher(c.Dump(), "merge")

	c.d.reset()
	c.Unwinder = &Unwind{
		Alias: "$x",
		Pairs: make(map[string]interface{}, 0),
	}
	c.Unwinder.Pairs["a"] = "b"
	c.Unwinder.Pairs["c"] = 1
	c.Sets[0].Alias = "m"

	parser.ValidateCypher(c.Dump(), "unwind")

	c.d.reset()
	var innerM Merger
	innerM.Sets = make([]Set, 1)
	innerM.Sets[0] = Set{
		Alias: "m",
		Name:  "Set1",
	}
	c.innerMergers = make([]Merger, 1)
	c.innerMergers[0] = innerM
	c.Unwinder.WithAlias = "m"

	parser.ValidateCypher(c.Dump(), "merge")

}
