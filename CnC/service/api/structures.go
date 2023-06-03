package api

type Zombie struct {
	Ip             string
	Id             string
	Active         bool
	Ports          []int
	CustomUsername string
}

type Action struct {
	Id   string `json:"id"`
	Ddos struct {
		Target string `json:"target"`
		Rounds int    `json:"rounds"`
	} `json:"ddos"`
	Info  []string `json:"info"`
	Email bool     `json:"email"`
}
