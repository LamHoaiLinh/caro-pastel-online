package persistence

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type GameRecord struct {
	ID             string     `json:"id"`
	GameMode       string     `json:"gameMode"`
	TimeControl    string     `json:"timeControl"`
	RedType        string     `json:"redType"`
	BlueType       string     `json:"blueType"`
	RedDifficulty  *int       `json:"redDifficulty,omitempty"`
	BlueDifficulty *int       `json:"blueDifficulty,omitempty"`
	Winner         string     `json:"winner"`
	MoveCount      int        `json:"moveCount"`
	CreatedAt      time.Time  `json:"createdAt"`
	CompletedAt    *time.Time `json:"completedAt,omitempty"`
}

type MoveRecord struct {
	GameID          string   `json:"gameId"`
	MoveNumber      int      `json:"moveNumber"`
	Player          string   `json:"player"`
	PosX            int      `json:"posX"`
	PosY            int      `json:"posY"`
	IsBot           bool     `json:"isBot"`
	Difficulty      *int     `json:"difficulty,omitempty"`
	ThinkTimeMs     *int64   `json:"thinkTimeMs,omitempty"`
	RemainingTimeMs *int64   `json:"remainingTimeMs,omitempty"`
	SearchDepth     *int     `json:"searchDepth,omitempty"`
	NodesSearched   *int64   `json:"nodesSearched,omitempty"`
	NPS             *float64 `json:"nps,omitempty"`
	TTHitRate       *float64 `json:"ttHitRate,omitempty"`
	SearchScore     *int     `json:"searchScore,omitempty"`
	ThreadsUsed     *int     `json:"threadsUsed,omitempty"`
	AllocatedTimeMs *int64   `json:"allocatedTimeMs,omitempty"`
	MoveType        *string  `json:"moveType,omitempty"`
	MasterPct       *float64 `json:"masterPct,omitempty"`
	SlaveDepth      *int     `json:"slaveDepth,omitempty"`
	SlaveNodes      *int64   `json:"slaveNodes,omitempty"`
	PonderDepth     *int     `json:"ponderDepth,omitempty"`
	PonderNodes     *int64   `json:"ponderNodes,omitempty"`
}

type snapshot struct {
	Games map[string]GameRecord   `json:"games"`
	Moves map[string][]MoveRecord `json:"moves"`
}

type MatchStore struct {
	mu     sync.RWMutex
	path   string
	games  map[string]GameRecord
	moves  map[string][]MoveRecord
	closed bool
}

func NewMatchStore(path string) (*MatchStore, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	store := &MatchStore{
		path:  path,
		games: make(map[string]GameRecord),
		moves: make(map[string][]MoveRecord),
	}
	data, err := os.ReadFile(path)
	if err == nil && len(data) > 0 {
		var saved snapshot
		if json.Unmarshal(data, &saved) == nil {
			if saved.Games != nil {
				store.games = saved.Games
			}
			if saved.Moves != nil {
				store.moves = saved.Moves
			}
		}
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	return store, nil
}

func (s *MatchStore) CreateGame(g GameRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if g.CreatedAt.IsZero() {
		g.CreatedAt = time.Now().UTC()
	}
	if g.Winner == "" {
		g.Winner = "none"
	}
	s.games[g.ID] = g
	return s.persistLocked()
}

func (s *MatchStore) RecordMove(m MoveRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.moves[m.GameID] = append(s.moves[m.GameID], m)
	return s.persistLocked()
}

func (s *MatchStore) CompleteGame(gameID, winner string, moveCount int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	g, ok := s.games[gameID]
	if !ok {
		return errors.New("game not found")
	}
	now := time.Now().UTC()
	g.Winner = winner
	g.MoveCount = moveCount
	g.CompletedAt = &now
	s.games[gameID] = g
	return s.persistLocked()
}

func (s *MatchStore) GetGame(gameID string) (*GameRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	g, ok := s.games[gameID]
	if !ok {
		return nil, errors.New("game not found")
	}
	copy := g
	return &copy, nil
}

func (s *MatchStore) GetMoves(gameID string) ([]MoveRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := s.moves[gameID]
	result := make([]MoveRecord, len(items))
	copy(result, items)
	return result, nil
}

func (s *MatchStore) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return
	}
	_ = s.persistLocked()
	s.closed = true
}

func (s *MatchStore) persistLocked() error {
	if s.path == "" || s.closed {
		return nil
	}
	data, err := json.Marshal(snapshot{Games: s.games, Moves: s.moves})
	if err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}
