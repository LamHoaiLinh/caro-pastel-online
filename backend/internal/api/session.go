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
		redDifficulty:   redDiff,
		blueDifficulty:  blueDiff,
		logger:          logger,
		activeGameCount: activeGameCount,
	}
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

	elapsed := now.Sub(s.lastMoveAt).Milliseconds()
	inc := int64(newGame.IncrementSeconds) * 1000
	if s.game.CurrentPlayer == domain.PlayerRed {
		s.redTimeMs = max(0, s.redTimeMs-elapsed+inc)
	} else {
		s.blueTimeMs = max(0, s.blueTimeMs-elapsed+inc)
	}
	s.lastMoveAt = now

	s.game = newGame

	if newGame.IsGameOver {
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
	if s.game.IsGameOver {
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

func (s *GameSession) expireIfNeededLocked(now time.Time) {
	if s.game.IsGameOver {
		return
	}
	redTime, blueTime := s.currentTimesLocked(now)
	if s.game.CurrentPlayer == domain.PlayerRed && redTime <= 0 {
		s.redTimeMs = 0
		s.game = s.game.WithGameOver(domain.PlayerBlue, nil)
		s.lastMoveAt = now
		s.DisposeAI()
	} else if s.game.CurrentPlayer == domain.PlayerBlue && blueTime <= 0 {
		s.blueTimeMs = 0
		s.game = s.game.WithGameOver(domain.PlayerRed, nil)
		s.lastMoveAt = now
		s.DisposeAI()
	}
}

func (s *GameSession) buildResponseAt(now time.Time) GameResponse {
	redTime, blueTime := s.currentTimesLocked(now)
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
		TimeControl:       s.game.TimeControl,
		InitialTime:       int(s.game.InitialTimeMs / 1000),
		Increment:         s.game.IncrementSeconds,
		GameMode:          s.game.GameMode.String(),
		RedDifficulty:     s.redDifficulty,
		BlueDifficulty:    s.blueDifficulty,
	}
}
