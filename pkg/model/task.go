package model

type TaskStatus int

const (
	Working TaskStatus = iota
	Ready
	Error
)
