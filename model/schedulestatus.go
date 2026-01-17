package model

type ScheduelStatus string

const (
	UnknownScheduleStatus ScheduelStatus = ""
	PENDING               ScheduelStatus = "PENDING"
	PAID                  ScheduelStatus = "PAID"
	PARTIAL               ScheduelStatus = "PARTIAL"
	OVERDUE               ScheduelStatus = "OVERDUE"
	CANCELLED             ScheduelStatus = "CANCELLED"
)
