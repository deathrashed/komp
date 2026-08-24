package recents

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"time"
)

const cap = 20

type Entry struct {
	Path     string    `json:"path"`
	LastUsed time.Time `json:"last_used"`
	Uses     int       `json:"uses"`
}

type Store struct{ file string }

func Load(file string) (*Store, error) {
	s := &Store{file: file}
	b, err := os.ReadFile(file)
	if err != nil {
		return s, nil // fresh
	}
	var es []Entry
	if json.Unmarshal(b, &es) != nil {
		return s, nil // corrupt → warn upstream, start fresh
	}
	s.save(es)
	return s, nil
}

func (s *Store) Touch(path string) error {
	es := s.all()
	found := false
	for i := range es {
		if es[i].Path == path {
			es[i].Uses++
			es[i].LastUsed = time.Now()
			found = true
		}
	}
	if !found {
		es = append(es, Entry{Path: path, LastUsed: time.Now(), Uses: 1})
	}
	sort.Slice(es, func(i, j int) bool { return es[i].LastUsed.After(es[j].LastUsed) })
	if len(es) > cap {
		es = es[:cap]
	}
	return s.save(es)
}

func (s *Store) Recent(n int) []Entry {
	es := s.all()
	if len(es) > n {
		es = es[:n]
	}
	return es
}

func (s *Store) all() []Entry {
	b, err := os.ReadFile(s.file)
	if err != nil {
		return nil
	}
	var es []Entry
	if json.Unmarshal(b, &es) != nil {
		return nil
	}
	return es
}

func (s *Store) save(es []Entry) error {
	_ = os.MkdirAll(filepath.Dir(s.file), 0o755)
	b, _ := json.MarshalIndent(es, "", "  ")
	return os.WriteFile(s.file, b, 0o600)
}
