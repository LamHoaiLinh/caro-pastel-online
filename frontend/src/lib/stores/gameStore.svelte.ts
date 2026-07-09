import type { Cell, Player } from '$lib/types/game';
import { switchPlayer } from '$lib/types/game';
import { GameConfig } from '$lib/config/gameConfig';

export interface MoveRecord {
	moveNumber: number;
	player: Player;
	x: number;
	y: number;
}

export class GameStore {
	board = $state<Cell[]>([]);
	currentPlayer = $state<Player>('red');
	moveNumber = $state(0);
	isGameOver = $state(false);
	winner = $state<Player | undefined>(undefined);
	moveHistory = $state<MoveRecord[]>([]);

	constructor() {
		this.reset();
	}

	reset() {
		this.board = Array.from({ length: GameConfig.totalCells }, (_, i) => ({
			x: i % GameConfig.boardSize,
			y: Math.floor(i / GameConfig.boardSize),
			player: 'none' as Player
		}));
		this.currentPlayer = 'red';
		this.moveNumber = 0;
		this.isGameOver = false;
		this.winner = undefined;
		this.moveHistory = [];
	}

	makeMove(x: number, y: number): boolean {
		if (this.isGameOver) return false;
		const cell = this.board[y * GameConfig.boardSize + x];
		if (!cell || cell.player !== 'none') return false;

		this.moveHistory.push({
			moveNumber: this.moveNumber + 1,
			player: this.currentPlayer,
			x,
			y
		});
		cell.player = this.currentPlayer;
		this.moveNumber++;
		this.currentPlayer = switchPlayer(this.currentPlayer);
		return true;
	}
}
