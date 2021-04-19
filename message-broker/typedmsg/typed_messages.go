package typedmsg

type Unsubscribe struct {
	Topics  Topics  `json:"topics"`
	Address Address `json:"address"`
}

type Subscribe struct {
	Topics  Topics  `json:"topics"`
	Address Address `json:"address"`
}

type Address struct {
	Addresses Addresses `json:"addresses"`
}
type Addresses []string
type Topics []string
