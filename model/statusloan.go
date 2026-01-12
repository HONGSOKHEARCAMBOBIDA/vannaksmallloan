package model

type Status int

const (
	Unknown Status = iota
	Pending
	Checked
	Approved
	Active
	WritenOff
	ReJected
)
