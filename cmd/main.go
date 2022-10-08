package main

import (
	"cypher-validator/pkg/util"
	"flag"
	"os"
)

type FUNC_MASK uint64

const (
	NEO4J_CYPHER   FUNC_MASK = 1
	GRAPHQL_CYPHER FUNC_MASK = NEO4J_CYPHER << 1

	ARG_EXCLUDE_MODE = "EXCLUDE_MODE"
)

var FUNC_MASK_BITMAP = [...]FUNC_MASK{GRAPHQL_CYPHER, NEO4J_CYPHER}

func main() {
	exclude_func := flag.Uint64(ARG_EXCLUDE_MODE, 0, `exclude function mask:
		neo4j cypher  1
		graphql cypher 2
	`)
	flag.Parse()

	var exitMsg string
	util.InitLogger()
	defer util.Logger.Sync()
	goto exec

exit:
	util.Logger.Infof("Exit! Reason:%s", exitMsg)
	os.Exit(2)

exec:
	if *exclude_func != 0 && FUNC_MASK(*exclude_func) >= 1<<len(FUNC_MASK_BITMAP) {
		exitMsg = "invalid arg:" + ARG_EXCLUDE_MODE
		goto exit
	}
	util.Logger.Infof("Exclude Mode:%b", *exclude_func)

}
