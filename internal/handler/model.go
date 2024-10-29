package handler

import "time"

type Request struct {
	Value string        `json:"value"`
	TTL   time.Duration `json:"ttl"`
}

type Response struct {
	Value string `json:"value"`
}
