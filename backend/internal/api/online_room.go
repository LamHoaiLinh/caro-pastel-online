package api

import (
	"caro-ai-pvp/internal/domain"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	maxOnlineRooms      = 100
	onlineRoomIdleLimit = 45 * time.Minute
)

var (
	errRoomNotFound    = errors.New("phòng chơi không tồn tại hoặc đã hết hạn")
	errRoomFull        = errors.New("phòng đã đủ hai người chơi")
	errNotYourTurn     = errors.New("chưa đến lượt của bạn")
	errWaitingOpponent = errors.New("đang chờ người chơi thứ hai vào phòng")
	errReadOnly        = errors.New("bạn đang xem phòng và không thể đánh cờ")
)

type CreateOnlineRequest struct {
	TimeControl   string `json:"timeControl"`
	MoveTimeLimit int    `json:"moveTimeLimit"`
	PlayerName    string `json:"playerName"`
}

type JoinOnlineRequest struct {
	PlayerName  string `json:"playerName"`
	PlayerToken string `json:"playerToken"`
}

type OnlineMoveRequest struct {
	X           int    `json:"x"`
	Y           int    `json:"y"`
	PlayerToken string `json:"playerToken"`
}

type OnlineRoomResponse struct {
	Code           string       `json:"code"`
	Role           string       `json:"role,omitempty"`
	PlayerToken    string       `json:"playerToken,omitempty"`
	RedName        string       `json:"redName"`
	BlueName       string       `json:"blueName"`
	OpponentJoined bool         `json:"opponentJoined"`
	State          GameResponse `json:"state"`
}

type OnlineRoom struct {
	mu           sync.Mutex
	code         string
	session      *GameSession
	redToken     string
	blueToken    string
	redName      string
	blueName     string
	createdAt    time.Time
	lastActivity time.Time
}

type OnlineRoomStore struct {
	mu    sync.RWMutex
	rooms map[string]*OnlineRoom
}

func NewOnlineRoomStore() *OnlineRoomStore {
	return &OnlineRoomStore{rooms: make(map[string]*OnlineRoom)}
}

func (s *OnlineRoomStore) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.rooms)
}

func (s *OnlineRoomStore) Get(code string) (*OnlineRoom, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	room, ok := s.rooms[normalizeRoomCode(code)]
	return room, ok
}

func (s *OnlineRoomStore) Create(timeControl string, moveTimeLimit int, playerName string, activeGameCount func() int) (*OnlineRoom, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.rooms) >= maxOnlineRooms {
		return nil, domain.ErrTooManyGames
	}

	code := ""
	for attempts := 0; attempts < 20; attempts++ {
		candidate := newRoomCode()
		if _, exists := s.rooms[candidate]; !exists {
			code = candidate
			break
		}
	}
	if code == "" {
		return nil, errors.New("không thể tạo mã phòng")
	}

	canonical, initialTimeMs, incrementSeconds := parseTimeControl(timeControl)
	session := NewGameSession(canonical, initialTimeMs, incrementSeconds, domain.GameModePvP, nil, nil, nil, activeGameCount)
	session.SetMoveTimeLimit(normalizeMoveTimeLimit(moveTimeLimit, canonical))
	// An online room is only a waiting room until player 2 joins. Neither the
	// total clock nor the per-move clock may run during that waiting period.
	session.PauseClock()
	now := time.Now()
	room := &OnlineRoom{
		code:         code,
		session:      session,
		redToken:     newPlayerToken(),
		redName:      cleanPlayerName(playerName, "Người chơi 1"),
		createdAt:    now,
		lastActivity: now,
	}
	s.rooms[code] = room
	return room, nil
}

func (s *OnlineRoomStore) CleanupExpired() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	removed := 0
	for code, room := range s.rooms {
		room.mu.Lock()
		idle := now.Sub(room.lastActivity)
		gameOver := room.session.IsGameOver()
		room.mu.Unlock()
		if idle > onlineRoomIdleLimit || (gameOver && idle > 10*time.Minute) {
			room.session.DisposeAI()
			delete(s.rooms, code)
			removed++
		}
	}
	return removed
}

func (s *OnlineRoomStore) CleanupAll() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	count := len(s.rooms)
	for code, room := range s.rooms {
		room.session.DisposeAI()
		delete(s.rooms, code)
	}
	return count
}

func (r *OnlineRoom) join(name, existingToken string) (role, token string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.lastActivity = time.Now()

	if existingToken != "" {
		if existingToken == r.redToken {
			if strings.TrimSpace(name) != "" {
				r.redName = cleanPlayerName(name, r.redName)
			}
			return "red", r.redToken
		}
		if existingToken == r.blueToken && r.blueToken != "" {
			if strings.TrimSpace(name) != "" {
				r.blueName = cleanPlayerName(name, r.blueName)
			}
			return "blue", r.blueToken
		}
	}

	if r.blueToken == "" {
		r.blueToken = newPlayerToken()
		r.blueName = cleanPlayerName(name, "Người chơi 2")
		// Start both clocks exactly when the second player successfully joins.
		r.session.StartClock()
		return "blue", r.blueToken
	}
	return "spectator", ""
}

func (r *OnlineRoom) response(role, token string) OnlineRoomResponse {
	r.mu.Lock()
	code := r.code
	redName := r.redName
	blueName := r.blueName
	opponentJoined := r.blueToken != ""
	r.lastActivity = time.Now()
	r.mu.Unlock()

	return OnlineRoomResponse{
		Code:           code,
		Role:           role,
		PlayerToken:    token,
		RedName:        redName,
		BlueName:       blueName,
		OpponentJoined: opponentJoined,
		State:          r.session.GetResponse(),
	}
}

func (r *OnlineRoom) applyMove(token string, x, y int) (OnlineRoomResponse, error) {
	// Keep the room lock for the full authorization + move transaction. This
	// prevents two near-simultaneous requests from one player being accepted
	// on both consecutive turns.
	r.mu.Lock()
	defer r.mu.Unlock()

	role := "spectator"
	if token != "" && token == r.redToken {
		role = "red"
	} else if token != "" && token == r.blueToken {
		role = "blue"
	}
	if role == "spectator" {
		return OnlineRoomResponse{}, errReadOnly
	}
	if r.blueToken == "" {
		return OnlineRoomResponse{}, errWaitingOpponent
	}

	current := r.session.GetResponse()
	if current.CurrentPlayer != role {
		return OnlineRoomResponse{}, errNotYourTurn
	}

	state, err := r.session.ApplyMove(x, y)
	if err != nil {
		return OnlineRoomResponse{}, err
	}
	r.lastActivity = time.Now()

	return OnlineRoomResponse{
		Code:           r.code,
		Role:           role,
		PlayerToken:    token,
		RedName:        r.redName,
		BlueName:       r.blueName,
		OpponentJoined: r.blueToken != "",
		State:          state,
	}, nil
}

func (h *Handler) CreateOnlineRoom(w http.ResponseWriter, req *http.Request) {
	var body CreateOnlineRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "bad_request", Message: "Dữ liệu tạo phòng không hợp lệ"})
		return
	}
	room, err := h.rooms.Create(body.TimeControl, body.MoveTimeLimit, body.PlayerName, func() int {
		return h.store.ActiveGameCount() + h.rooms.Count()
	})
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, room.response("red", room.redToken))
}

func (h *Handler) JoinOnlineRoom(w http.ResponseWriter, req *http.Request) {
	room, ok := h.rooms.Get(req.PathValue("code"))
	if !ok {
		writeJSON(w, http.StatusNotFound, ErrorResponse{Error: "room_not_found", Message: errRoomNotFound.Error()})
		return
	}
	var body JoinOnlineRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "bad_request", Message: "Dữ liệu vào phòng không hợp lệ"})
		return
	}
	role, token := room.join(body.PlayerName, body.PlayerToken)
	writeJSON(w, http.StatusOK, room.response(role, token))
}

func (h *Handler) GetOnlineRoom(w http.ResponseWriter, req *http.Request) {
	room, ok := h.rooms.Get(req.PathValue("code"))
	if !ok {
		writeJSON(w, http.StatusNotFound, ErrorResponse{Error: "room_not_found", Message: errRoomNotFound.Error()})
		return
	}
	writeJSON(w, http.StatusOK, room.response("spectator", ""))
}

func (h *Handler) MakeOnlineMove(w http.ResponseWriter, req *http.Request) {
	room, ok := h.rooms.Get(req.PathValue("code"))
	if !ok {
		writeJSON(w, http.StatusNotFound, ErrorResponse{Error: "room_not_found", Message: errRoomNotFound.Error()})
		return
	}
	var body OnlineMoveRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "bad_request", Message: "Nước đi không hợp lệ"})
		return
	}
	token := strings.TrimSpace(body.PlayerToken)
	if token == "" {
		token = strings.TrimSpace(req.Header.Get("X-Player-Token"))
	}
	resp, err := room.applyMove(token, body.X, body.Y)
	if err != nil {
		switch err {
		case errReadOnly:
			writeJSON(w, http.StatusForbidden, ErrorResponse{Error: "read_only", Message: err.Error()})
		case errNotYourTurn:
			writeJSON(w, http.StatusConflict, ErrorResponse{Error: "not_your_turn", Message: err.Error()})
		case errWaitingOpponent:
			writeJSON(w, http.StatusConflict, ErrorResponse{Error: "waiting_opponent", Message: err.Error()})
		default:
			writeError(w, err)
		}
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) CleanupOnlineRooms() int {
	return h.rooms.CleanupExpired()
}

func (h *Handler) CleanupAllOnlineRooms() int {
	return h.rooms.CleanupAll()
}

func normalizeRoomCode(code string) string {
	return strings.ToUpper(strings.TrimSpace(code))
}

func cleanPlayerName(name, fallback string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return fallback
	}
	runes := []rune(name)
	if len(runes) > 24 {
		runes = runes[:24]
	}
	return string(runes)
}

func newRoomCode() string {
	const alphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	buf := make([]byte, 6)
	random := make([]byte, 6)
	if _, err := rand.Read(random); err != nil {
		return "ROOM01"
	}
	for i := range buf {
		buf[i] = alphabet[int(random[i])%len(alphabet)]
	}
	return string(buf)
}

func newPlayerToken() string {
	buf := make([]byte, 24)
	if _, err := rand.Read(buf); err != nil {
		return newGameID() + newGameID()
	}
	return base64.RawURLEncoding.EncodeToString(buf)
}

func parseTimeControl(value string) (string, int64, int) {
	switch value {
	case "1+0", "bullet":
		return "1+0", 60_000, 0
	case "3+0":
		return "3+0", 180_000, 0
	case "3+2", "blitz":
		return "3+2", 180_000, 2
	case "10+0":
		return "10+0", 600_000, 0
	case "15+10", "classical":
		return "15+10", 900_000, 10
	default:
		return "7+5", 420_000, 5
	}
}

func normalizeMoveTimeLimit(requested int, timeControl string) int {
	if requested <= 0 {
		switch timeControl {
		case "1+0":
			return 10
		case "3+0", "3+2":
			return 20
		case "10+0":
			return 45
		case "15+10":
			return 60
		default:
			return 30
		}
	}
	if requested < 5 {
		return 5
	}
	if requested > 300 {
		return 300
	}
	return requested
}
