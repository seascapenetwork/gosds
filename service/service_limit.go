package service

type Limit uint8

const (
	REQUEST   Limit = 1
	REPLY     Limit = 2
	SUBSCRIBE Limit = 3
	BROADCAST Limit = 4
)
