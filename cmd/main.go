package main

import (
	"fmt"

	"github.com/RulezKT/nodes"
)

func main() {

	nodes := nodes.Nodes{}
	nodes.Load("files")

	fmt.Println(nodes.Calc(-682470731)) //185.01141648192737 5.011416481927357

	fmt.Println(nodes.Calc(682470731)) // 	67.2131172448477 247.21311724484767

	// date_in_seconds := int64(-682470731)

	// nodesSec, nodesLng := nodes.Load("files")

	// // fmt.Println(nodesSec)
	// // fmt.Println(nodesLng)

	// nNode, sNode := nodes.Nodes(date_in_seconds, nodesSec, nodesLng)

	// fmt.Println(nNode, sNode)

}
