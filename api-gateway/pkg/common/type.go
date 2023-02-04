package common

type Response struct {
	Status uint        `json:"Status"`
	Data   interface{} `json:"Data"`
	Error  string      `json:"Error"`
}

type TokenData struct {
	User  interface{} `json:"User"`
	Token string      `json:"Token"`
}
