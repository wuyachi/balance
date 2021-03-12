package config

// Account ...
type Account struct {
	Eth  []string
	BSC  []string
	Heco []string
	Ont  []string
	Zil  []string
}

// Threshold ...
type Threshold struct {
	Eth  uint
	BSC  uint
	Heco uint
	Ont  uint
	Zil  uint
}

// Node ...
type Node struct {
	Eth  []string
	BSC  []string
	Heco []string
	Ont  []string
	Zil  []string
}

// Config ...
type Config struct {
	Account   Account
	Threshold Threshold
	Node      Node
}
