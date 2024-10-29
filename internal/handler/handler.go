package handler

import (
	"encoding/json"
	"net/http"
	"net/url"

	"go.uber.org/zap"
)

func parseKey(u *url.URL) string {
	params := u.Query()
	key := params.Get("key")

	return key
}

func AddSetHandler(f func(key, value string, ttl int64) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			key := parseKey(r.URL)
			if key == "" {
				http.Error(w, "Key is empty", http.StatusBadRequest)
				zap.S().Errorf("Key is empty")
				return
			}

			decoder := json.NewDecoder(r.Body)
			var request Request
			err := decoder.Decode(&request)

			if err != nil {
				http.Error(w, "Bad request", http.StatusBadRequest)
				zap.S().Error(err)
				return
			}

			err := f(key, request.Value, request.TTL)
			if err != nil {
				http.Error(w, err, http.StatusBadRequest)
				zap.S().Error(err)
			}
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			zap.S().Errorf("Method not allowed")
		}
	}
}

func DelHandler(f func(key string) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			key := parseKey(r.URL)
			if key == "" {
				http.Error(w, "Key is empty", http.StatusBadRequest)
				zap.S().Errorf("Key is empty")
				return
			}

			err := f(key)
			if err != nil {
				http.Error(w, err, http.StatusBadRequest)
				zap.S().Error(err)
			}
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			zap.S().Errorf("Method not allowed")
		}
	}
}

func GetHandler(f func(key string) (string, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			key := parseKey(r.URL)
			if key == "" {
				http.Error(w, "Key is empty", http.StatusBadRequest)
				zap.S().Errorf("Key is empty")
				return
			}

			value, err := f(key)

			if err != nil {
				http.Error(w, "Key not found", http.StatusBadRequest)
				zap.S().Errorf("Key not found")
				return
			}

			response := Response{
				Value: value,
			}
			json.NewEncoder(w).Encode(response)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			zap.S().Errorf("Method not allowed")
		}
	}
}
