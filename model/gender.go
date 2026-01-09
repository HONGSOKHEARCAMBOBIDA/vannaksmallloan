package model

type Gender int

const (
	GenderUnknown Gender = iota
	GenderMale
	GenderFemale
	//iota starts at 0
	//1 → Male
	//2 → Femal
)
