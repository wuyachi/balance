package config

// Account ...
type Account struct {
	Eth     []string
	BSC     []string
	Heco    []string
	Polygon []string
	OK      []string
	Arb     []string
	Ont     []string
	Zil     []string
	Neo     []string
	Metis   []string
}

// Threshold ...
type Threshold struct {
	Eth     uint
	BSC     uint
	Heco    uint
	Polygon uint
	Ont     uint
	Zil     uint
	Neo     uint
	OK      uint
	Arb     uint
	Metis   uint
}

// Node ...
type Node struct {
	Eth     []string
	BSC     []string
	Heco    []string
	OK      []string
	Arb     []string
	Ont     []string
	Zil     []string
	Neo     []string
	Polygon []string
	Metis   []string
}

// Config ...
type Config struct {
	Account   Account
	Threshold Threshold
	Node      Node
}
