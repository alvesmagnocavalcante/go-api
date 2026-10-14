package main

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

var pokemonName = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

type pokemon struct {
	ID      int             `json:"id"`
	Name    string          `json:"name"`
	Height  int             `json:"height"`
	Weight  int             `json:"weight"`
	Sprites json.RawMessage `json:"sprites"`
	Types   json.RawMessage `json:"types"`
	Stats   json.RawMessage `json:"stats"`
}

func jsonError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	writeJSON(w, map[string]string{"error": message})
}

func newPokemonHandler(client *http.Client, baseURL string) http.Handler {
	var mu sync.RWMutex
	cache := make(map[string]pokemon)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := strings.ToLower(strings.TrimSpace(r.PathValue("name")))
		if len(name) > 100 || !pokemonName.MatchString(name) {
			jsonError(w, http.StatusBadRequest, "Informe um nome válido, como pikachu ou mr-mime.")
			return
		}
		mu.RLock()
		cached, ok := cache[name]
		mu.RUnlock()
		if ok {
			writeJSON(w, cached)
			return
		}
		request, err := http.NewRequestWithContext(r.Context(), http.MethodGet, baseURL+name+"/", nil)
		if err != nil {
			jsonError(w, http.StatusInternalServerError, "Não foi possível preparar a consulta.")
			return
		}
		response, err := client.Do(request)
		if err != nil {
			jsonError(w, http.StatusBadGateway, "Não foi possível consultar a PokéAPI. Tente novamente.")
			return
		}
		defer response.Body.Close()
		if response.StatusCode == http.StatusNotFound {
			jsonError(w, http.StatusNotFound, "Pokémon não encontrado. Confira o nome e tente novamente.")
			return
		}
		if response.StatusCode != http.StatusOK {
			jsonError(w, http.StatusBadGateway, "A PokéAPI está indisponível. Tente novamente.")
			return
		}
		var data pokemon
		if err := json.NewDecoder(io.LimitReader(response.Body, 2<<20)).Decode(&data); err != nil || data.ID <= 0 || data.Name == "" {
			jsonError(w, http.StatusBadGateway, "A PokéAPI retornou uma resposta inválida.")
			return
		}
		mu.Lock()
		if len(cache) >= 128 {
			clear(cache)
		}
		cache[name] = data
		mu.Unlock()
		writeJSON(w, data)
	})
}
