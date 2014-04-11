// Chef client command-line tool.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/marpaia/chef-golang"
	"github.com/shurcooL/go-goon"
)

var _ = goon.Dump

func chefConnect() *chef.Chef {
	c, err := chef.Connect()
	if err != nil {
		panic(err)
	}
	c.SSLNoVerify = true
	return c
}

func main() {
	flag.Parse()
	args := flag.Args()

	switch {
	case len(args) == 1:
		c := chefConnect()

		results, err := c.Search("node", "role:"+args[0])
		if err != nil {
			panic(err)
		}

		//fmt.Println(results.Total)
		for _, row := range results.Rows {
			row := row.(map[string]interface{})

			fmt.Println(row["name"])
		}
	/*case false:
	c := chefConnect()

	nodes, err := c.GetNodes()
	if err != nil {
		panic(err)
	}

	goon.DumpExpr(nodes)*/
	default:
		flag.PrintDefaults()
		os.Exit(2)
	}
}
