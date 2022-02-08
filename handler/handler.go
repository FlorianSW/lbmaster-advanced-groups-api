package handler

import (
	"code.cloudfoundry.org/lager"
	"encoding/json"
	"lbmaster-advanced-groups-api/internal/domain"
	"net/http"
	"strconv"
	"strings"
)

type webHandler struct {
	repo   domain.PrefixGroupRepository
	apiKey string
	logger lager.Logger
}

func NewHandler(repo domain.PrefixGroupRepository, apiKey string, logger lager.Logger) *webHandler {
	return &webHandler{
		repo:   repo,
		apiKey: apiKey,
		logger: logger,
	}
}

func (h webHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	l := h.logger.Session("serve", lager.Data{"url": r.URL.String()})
	defer func() {
		err := r.Body.Close()
		if err != nil {
			l.Error("close-body", err)
		}
	}()

	if r.Header.Get("Authorization") == "" || r.Header.Get("Authorization") != "Bearer "+h.apiKey {
		w.WriteHeader(401)
		return
	}

	path := r.URL.Path
	if strings.HasPrefix(path, "/api/prefixGroups") {
		h.handlePrefixGroups(w, r, strings.TrimPrefix(path, "/api/prefixGroups"))
		return
	}

	w.WriteHeader(404)
}

func (h webHandler) handlePrefixGroups(w http.ResponseWriter, r *http.Request, part string) {
	if part == "" {
		h.handleListPrefixGroups(w, r)
		return
	}
	parts := strings.Split(strings.TrimPrefix(part, "/"), "/")
	if len(parts) == 1 {
		h.handleListPrefixGroupMembers(w, r, parts)
	} else if len(parts) == 2 {
		h.handlePrefixGroupMembers(w, r, parts)
	} else {
		w.WriteHeader(404)
	}
}

func (h webHandler) handleListPrefixGroupMembers(w http.ResponseWriter, r *http.Request, parts []string) {
	idx, err := strconv.Atoi(parts[0])
	if err != nil {
		w.WriteHeader(404)
		return
	}
	if r.Method != "GET" {
		w.WriteHeader(405)
	} else {
		m, err := h.repo.Members(domain.PrefixGroup{Index: idx})
		if err != nil {
			w.WriteHeader(500)
			return
		}
		c, err := json.Marshal(m)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(200)
		_, _ = w.Write(c)
	}
}

func (h webHandler) handlePrefixGroupMembers(w http.ResponseWriter, r *http.Request, parts []string) {
	idx, err := strconv.Atoi(parts[0])
	if err != nil {
		w.WriteHeader(404)
		return
	}
	pg := domain.PrefixGroup{Index: idx}
	if r.Method == "PUT" {
		err = h.repo.AddMember(pg, domain.SteamUID(parts[1]))
		if err != nil {
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(204)
	} else if r.Method == "DELETE" {
		err = h.repo.RemoveMember(pg, domain.SteamUID(parts[1]))
		if err != nil {
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(204)
	} else if r.Method == "GET" {
		m, err := h.repo.Members(pg)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		for _, m := range m {
			if m.String() == parts[1] {
				w.WriteHeader(204)
				return
			}
		}
		w.WriteHeader(404)
	} else {
		w.WriteHeader(405)
	}
}

func (h webHandler) handleListPrefixGroups(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(405)
		return
	}
	l, err := h.repo.List()
	if err != nil {
		w.WriteHeader(500)
		return
	}
	var result []map[string]interface{}
	for _, group := range l {
		result = append(result, map[string]interface{}{
			"index":  group.Index,
			"prefix": group.Prefix,
		})
	}
	c, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	_, _ = w.Write(c)
}
