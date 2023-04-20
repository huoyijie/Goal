package main

import (
	"fmt"

	"github.com/huoyijie/Goal"
	"github.com/huoyijie/Goal/examples/cdn"
	"github.com/huoyijie/Goal/util"
)

func main() {
	config := goal.Config{
		AllowOrigins:   []string{"http://127.0.0.1:4000"},
		TrustedProxies: nil,
		Cookie: goal.Cookie{
			Domain: "127.0.0.1",
			Secure: false,
		},
	}
	db := util.OpenSqliteDB()
	router := goal.New(config, db, &cdn.Resource{}).Router()
	router.Static("uploads", "uploads")
	router.Run(fmt.Sprintf("%s:%d", "127.0.0.1", 8100))
}
