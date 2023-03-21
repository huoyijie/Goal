package main

import (
	"fmt"

	"github.com/huoyijie/goal"
	"github.com/huoyijie/goal/util"
)

func main() {
	router := goal.NewGoal(util.OpenSqliteDB()).Router()
	router.Run(fmt.Sprintf("%s:%d", "127.0.0.1", 8100))
}
