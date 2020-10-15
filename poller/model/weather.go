package model

type Weather struct {
	Temperature struct {
		Temp float32 `json:"temp"`
	} `json:"main"`
}
