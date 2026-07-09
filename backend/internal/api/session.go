package api

import (
	"caro-ai-pvp/internal/domain"
	"caro-ai-pvp/internal/engine"
	"log/slog"
	"sync"
	"time"
)

type GameSession struct {
	mu              sync.Mutex
	game            domain.GameState
	redTimeMs       int64
	blueTimeMs      int64
	lastMoveAt      time.Time
	moveTimeLimitMs int64
	clockRunning    bool
	timeoutReason   string
	redDifficulty   *int
	blueDifficulty  *int
	logger          *slog.Logger
	activeGameCount func() int
	redAI           *engine.MinimaxAI
	blueAI          *engine.MinimaxAI
}

func NewGameSession(
	timeControl string,
	initialTimeMs int64,
	incrementSeconds int,
	mode domain.GameMode,
	redDiff, blueDiff *int,
	logger *slog.Logger,
	activeGameCount func() int,
) *GameSession {
	return &GameSession{
		game:            domain.NewGameState(mode, timeControl, initialTimeMs, incrementSeconds),
		redTimeMs:       initialTimeMs,
		blueTimeMs:      initialTimeMs,
		lastMoveAt:      time.Now(),
		clockRunning:    true,
		redDifficulty:   redDiff,
		blueDifficulty:  blueDiff,
		logger:          logger,
		activeGameCount: activeGameCount,
	}
}

// SetMoveTimeLimit configures the maximum thinking time for one turn.
// A value <= 0 disables the per-turn limit.
func (s *GameSession) SetMoveTimeLimit(seconds int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if seconds <= 0 {
		s.moveTimeLimitMs = 0
		return
	}
	s.moveTimeLimitMs = int64(seconds) * 1000
}

// PauseClock freezes both total clocks and the current-turn clock.
func (s *GameSession) PauseClock() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.game.IsGameOver || !s.clockRunning {
		return
	}
	now := time.Now()
	s.redTimeMs, s.blueTimeMs = s.currentTimesLocked(now)
	s.clockRunning = false
	s.lastMoveAt = now
}

// StartClock starts or resumes the clock from the current instant.
func (s *GameSession) StartClock() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.game.IsGameOver || s.clockRunning {
		return
	}
	s.clockRunning = true
	s.lastMoveAt = time.Now()
}

func (s *GameSession) GetResponse() GameResponse {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	s.expireIfNeededLocked(now)
	return s.buildResponseAt(now)
}

func (s *GameSession) IsGameOver() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.expireIfNeededLocked(time.Now())
	return s.game.IsGameOver
}

func (s *GameSession) ResetTurnClock() {
	s.mu.Lock()
	s.lastMoveAt = time.Now()
	s.mu.Unlock()
}

func (s *GameSession) LastActivityAt() time.Time {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.lastMoveAt
}

func (s *GameSession) ExtractForAI() (domain.Board, domain.Player, bool, int64, int, int, *int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	s.expireIfNeededLocked(now)
	redTime, blueTime := s.currentTimesLocked(now)
	timeRemaining := redTime
	diff := s.redDifficulty
	if s.game.CurrentPlayer == domain.PlayerBlue {
		timeRemaining = blueTime
		diff = s.blueDifficulty
	}
	if turnRemaining := s.currentTurnTimeLocked(now); turnRemaining > 0 && turnRemaining < timeRemaining {
		timeRemaining = turnRemaining
	}

	return s.game.Board, s.game.CurrentPlayer, s.game.IsGameOver,
		timeRemaining, s.game.IncrementSeconds, s.game.MoveNumber, diff
}

func (s *GameSession) GetOrCreateAI(player domain.Player) *engine.MinimaxAI {
	threads := engine.GetEngineThreadsForLoad(s.activeGameCount())
	if player == domain.PlayerRed {
		if s.redAI == nil {
			s.redAI = engine.NewMinimaxAI(s.logger, threads)
		}
		return s.redAI
	}
	if s.blueAI == nil {
		s.blueAI = engine.NewMinimaxAI(s.logger, threads)
	}
	return s.blueAI
}

func (s *GameSession) ApplyMove(x, y int) (GameResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	s.expireIfNeededLocked(now)
	if s.game.IsGameOver {
		return GameResponse{}, domain.ErrGameOver
	}

	newGame, err := s.game.WithMove(x, y)
	if err != nil {
		return GameResponse{}, err
	}

	result := domain.CheckWinFromMove(newGame.Board, x, y)
	if result.HasWinner {
		newGame = newGame.WithGameOver(result.Winner, result.WinningLine)
	}

	elapsed := int64(0)
	if s.clockRunning {
		elapsed = now.Sub(s.lastMoveAt).Milliseconds()
	}
	inc := int64(newGame.IncrementSeconds) * 1000
	if s.game.CurrentPlayer == domain.PlayerRed {
		s.redTimeMs = max(0, s.redTimeMs-elapsed+inc)
	} else {
		s.blueTimeMs = max(0, s.blueTimeMs-elapsed+inc)
	}
	s.lastMoveAt = now
	s.timeoutReason = ""

	s.game = newGame

	if newGame.IsGameOver {
		s.clockRunning = false
		s.DisposeAI()
	}

	return s.buildResponseAt(now), nil
}

func (s *GameSession) UndoLastMove() (GameResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	newGame, err := s.game.UndoMove()
	if err != nil {
		return GameResponse{}, err
	}
	s.game = newGame
	s.lastMoveAt = time.Now()
	s.timeoutReason = ""
	return s.buildResponseAt(s.lastMoveAt), nil
}

func (s *GameSession) DisposeAI() {
	if s.redAI != nil {
		s.redAI.Dispose()
		s.redAI = nil
	}
	if s.blueAI != nil {
		s.blueAI.Dispose()
		s.blueAI = nil
	}
}

func (s *GameSession) currentTimesLocked(now time.Time) (int64, int64) {
	redTime := s.redTimeMs
	blueTime := s.blueTimeMs
	if s.game.IsGameOver || !s.clockRunning {
		return redTime, blueTime
	}

	elapsed := now.Sub(s.lastMoveAt).Milliseconds()
	if s.game.CurrentPlayer == domain.PlayerRed {
		redTime = max(0, redTime-elapsed)
	} else if s.game.CurrentPlayer == domain.PlayerBlue {
		blueTime = max(0, blueTime-elapsed)
	}
	return redTime, blueTime
}

func (s *GameSession) currentTurnTimeLocked(now time.Time) int64 {
	if s.moveTimeLimitMs <= 0 {
		return 0
	}
	if s.game.IsGameOver {
		return 0
	}
	if !s.clockRunning {
		return s.moveTimeLimitMs
	}
	return max(0, s.moveTimeLimitMs-now.Sub(s.lastMoveAt).Milliseconds())
}

func (s *GameSession) expireIfNeededLocked(now time.Time) {
	if s.game.IsGameOver || !s.clockRunning {
		return
	}
	redTime, blueTime := s.currentTimesLocked(now)
	turnTime := s.currentTurnTimeLocked(now)
	turnExpired := s.moveTimeLimitMs > 0 && turnTime <= 0

	if s.game.CurrentPlayer == domain.PlayerRed && (redTime <= 0 || turnExpired) {
		s.redTimeMs = redTime
		s.blueTimeMs = blueTime
		s.timeoutReason = "total"
		if turnExpired && redTime > 0 {
			s.timeoutReason = "move"
		}
		s.game = s.game.WithGameOver(domain.PlayerBlue, nil)
		s.clockRunning = false
		s.lastMoveAt = now
		s.DisposeAI()
	} else if s.game.CurrentPlayer == domain.PlayerBlue && (blueTime <= 0 || turnExpired) {
		s.redTimeMs = redTime
		s.blueTimeMs = blueTime
		s.timeoutReason = "total"
		if turnExpired && blueTime > 0 {
			s.timeoutReason = "move"
		}
		s.game = s.game.WithGameOver(domain.PlayerRed, nil)
		s.clockRunning = false
		s.lastMoveAt = now
		s.DisposeAI()
	}
}

func (s *GameSession) buildResponseAt(now time.Time) GameResponse {
	redTime, blueTime := s.currentTimesLocked(now)
	turnTime := s.currentTurnTimeLocked(now)
	cells := make([]CellResponse, 0, domain.BoardSize*domain.BoardSize)
	for y := range domain.BoardSize {
		for x := range domain.BoardSize {
			player := s.game.Board.GetPlayerAt(x, y)
			cells = append(cells, CellResponse{X: x, Y: y, Player: player.String()})
		}
	}

	winningLine := make([]PositionResponse, len(s.game.WinningLine))
	for i, p := range s.game.WinningLine {
		winningLine[i] = PositionResponse{X: p.X, Y: p.Y}
	}

	return GameResponse{
		Board:             cells,
		CurrentPlayer:     s.game.CurrentPlayer.String(),
		MoveNumber:        s.game.MoveNumber,
		IsGameOver:        s.game.IsGameOver,
		Winner:            s.game.Winner.String(),
		WinningLine:       winningLine,
		RedTimeRemaining:  float64(redTime) / 1000.0,
		BlueTimeRemaining: float64(blueTime) / 1000.0,
		TurnTimeRemaining: float64(turnTime) / 1000.0,
		TimeControl:       s.game.TimeControl,
		InitialTime:       int(s.game.InitialTimeMs / 1000),
		Increment:         s.game.IncrementSeconds,
		MoveTimeLimit:     int(s.moveTimeLimitMs / 1000),
		ClockRunning:      s.clockRunning,
		TimeoutReason:     s.timeoutReason,
		GameMode:          s.game.GameMode.String(),
		RedDifficulty:     s.redDifficulty,
		BlueDifficulty:    s.blueDifficulty,
	}
}
