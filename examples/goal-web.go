package main

import (
	"fmt"

	goal "github.com/huoyijie/Goal"
	"github.com/huoyijie/Goal/examples/cdn"
	"github.com/huoyijie/Goal/util"
)

func main() {
	router := goal.NewGoal(util.OpenSqliteDB(), &cdn.Resource{}).Router()
	router.Static("uploads", "uploads")
	router.Run(fmt.Sprintf("%s:%d", "127.0.0.1", 8100))
}
