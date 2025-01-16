package main

import (
	"fmt"

	"github.com/RulezKT/nodes"
)

func main() {

	nodes := nodes.Nodes{}
	nodes.Load("files")
	// fmt.Println("sec = ", nodes.SecArr)
	// fmt.Println("lng = ", nodes.LngArr)
	nodes.Calc(-682470731)
	fmt.Println(nodes.North, nodes.South)

	nodes.Calc(682470731)
	fmt.Println(nodes.North, nodes.South)

	// date_in_seconds := int64(-682470731)

	// nodesSec, nodesLng := nodes.Load("files")

	// // fmt.Println(nodesSec)
	// // fmt.Println(nodesLng)

	// nNode, sNode := nodes.Nodes(date_in_seconds, nodesSec, nodesLng)

	// fmt.Println(nNode, sNode)

}
