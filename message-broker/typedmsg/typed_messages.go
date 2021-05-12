package typedmsg

type Unsubscribe struct {
	Topics  Topics  `json:"topics"`
	Address Address `json:"address"`
}

type Subscribe struct {
	Topics  Topics  `json:"topics"`
	Command Command `json:"command"`
	Durable Durable `json:"durable"`
	Address Address `json:"address"`
}

type Message struct {
	Topics             []Topic       `json:"topics,omitempty"`
	Command            Command       `json:"command"`
	Address            ClientAddress `json:"address,omitempty"`
	UniqueIDForDurable string        `json:"unique_id_for_durable,omitempty"`
}

type MessageInfo struct {
	Address Address `json:"address"`
	Topics  []Topic `json:"topics"`
}

type Command string
type ClientAddress string

//////////////////////////////////////////////////
type Address struct {
	Addresses Addresses `json:"addresses"`
}

type Durable bool
type Addresses []string
type Topics []Topic

type DurableTopicsValue []string

type Topic struct {
	Value     string  `json:"value"`
	IsDurable Durable `json:"is_durable,omitempty"`
}

type StopMessage struct {
	ClientAddress     string
	UniqueClientId    string
	OnlyDurableTopics DurableTopicsValue
	MyMagicChan       chan string
	Name              string
}

type ClientId struct {
	Value string `json:"value"`
}

type UniqueIdAndAddress struct {
	UniqueId      string
	ClientAddress ClientAddress
}
