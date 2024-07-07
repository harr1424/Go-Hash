package main

import (
	"sync"
)

var state = AppState{}

type AppState struct {
	EnImageHash  string
	EnPImageHash string
	EsImageHash  string
	EsPImageHash string
	FrImageHash  string
	PoImageHash  string
	ItImageHash  string
	DeImageHash  string
	mu           sync.Mutex
}
