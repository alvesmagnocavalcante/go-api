package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestPokemonHandler(t *testing.T) {
	for _, tt := range []struct {
		name   string
		input  string
		status int
		body   string
		want   int
	}{
		{"success", "Pikachu", 200, `{"id":25,"name":"pikachu","height":4,"weight":60,"sprites":{},"types":[],"stats":[]}`, 200},
		{"invalid name", "bad.name", 200, `{}`, 400},
		{"not found", "unknown", 404, `Not Found`, 404},
		{"upstream failure", "pikachu", 503, `unavailable`, 502},
		{"invalid JSON", "pikachu", 200, `broken`, 502},
		{"invalid payload", "pikachu", 200, `{}`, 502},
	} {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				if r.URL.Path != "/"+strings.ToLower(tt.input)+"/" {
					t.Errorf("unexpected upstream path: %s", r.URL.Path)
				}
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte(tt.body))
			}))
			defer upstream.Close()
			mux := http.NewServeMux()
			mux.Handle("GET /api/pokemon/{name}", newPokemonHandler(upstream.Client(), upstream.URL+"/"))
			for range 2 {
				response := httptest.NewRecorder()
				mux.ServeHTTP(response, httptest.NewRequest("GET", "/api/pokemon/"+tt.input, nil))
				if response.Code != tt.want {
					t.Fatalf("status = %d, want %d: %s", response.Code, tt.want, response.Body.String())
				}
				if !strings.HasPrefix(response.Header().Get("Content-Type"), "application/json") {
					t.Fatal("expected JSON response")
				}
				if tt.want == 200 && !strings.Contains(response.Body.String(), `"name":"pikachu"`) {
					t.Fatal("missing Pokémon data")
				}
			}
			if tt.want == 400 && calls != 0 {
				t.Fatal("invalid input reached upstream")
			}
			if tt.want == 200 && calls != 1 {
				t.Fatalf("cache missed: upstream calls = %d", calls)
			}
		})
	}
}

func TestPokemonTimeout(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	}))
	defer upstream.Close()
	client := upstream.Client()
	client.Timeout = 20 * time.Millisecond
	request := httptest.NewRequest("GET", "/api/pokemon/pikachu", nil)
	request.SetPathValue("name", "pikachu")
	response := httptest.NewRecorder()
	newPokemonHandler(client, upstream.URL+"/").ServeHTTP(response, request)
	if response.Code != http.StatusBadGateway {
		t.Fatalf("status = %d, want 502", response.Code)
	}
}

func TestFrontendAssets(t *testing.T) {
	for _, path := range []string{"/", "/app.js", "/style.css"} {
		response := httptest.NewRecorder()
		newHandler().ServeHTTP(response, httptest.NewRequest("GET", path, nil))
		if response.Code != http.StatusOK || response.Body.Len() == 0 {
			t.Fatalf("asset %s unavailable: %d", path, response.Code)
		}
	}
}
