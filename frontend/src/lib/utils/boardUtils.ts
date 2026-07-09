import { GameConfig } from '$lib/config/gameConfig';
import { UIConfig } from '$lib/config/uiConfig';

export interface Point {
	x: number;
	y: number;
}

export function calculateGhostStonePosition(x: number, y: number, offset: number = UIConfig.ghostStoneTouchOffset): Point {
	return { x, y: y - offset };
}

export function isValidCell(x: number, y: number): boolean {
	return x >= 0 && x < GameConfig.boardSize && y >= 0 && y < GameConfig.boardSize;
}

export function computeCellSize(viewportWidth: number): number {
	// Chừa chỗ cho hai cột nhãn số ở hai bên bàn cờ.
	// Công thức cũ chỉ chia cho 16 ô nên tổng bàn + nhãn rộng hơn màn hình điện thoại.
	const effectiveColumns = GameConfig.boardSize + 1.08;
	const size = Math.floor((viewportWidth * UIConfig.boardWidthFraction) / effectiveColumns);
	return Math.max(UIConfig.minCellSize, Math.min(UIConfig.maxCellSize, size));
}
