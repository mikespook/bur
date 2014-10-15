package bur

import "sync"

var users = make(map[string]User)
var usersLock sync.RWMutex

var total = State{}
var totalLock sync.RWMutex
