import { beforeEach, describe, expect, it } from 'vitest';
import { GameStore } from './gameStore.svelte';

describe('GameStore', () => {
	let store: GameStore;

	beforeEach(() => {
		store = new GameStore();
	});

	it('khởi tạo bàn 16x16 và lượt Đỏ', () => {
		expect(store.board).toHaveLength(256);
		expect(store.currentPlayer).toBe('red');
		expect(store.moveNumber).toBe(0);
	});

	it('ghi nhận nước đi hợp lệ và đổi lượt', () => {
		expect(store.makeMove(7, 7)).toBe(true);
		expect(store.board[7 * 16 + 7].player).toBe('red');
		expect(store.currentPlayer).toBe('blue');
		expect(store.moveHistory[0]).toEqual({ moveNumber: 1, player: 'red', x: 7, y: 7 });
	});

	it('từ chối đánh vào ô đã có quân', () => {
		expect(store.makeMove(7, 7)).toBe(true);
		expect(store.makeMove(7, 7)).toBe(false);
		expect(store.moveHistory).toHaveLength(1);
	});

	it('reset xóa trạng thái ván cũ', () => {
		store.makeMove(7, 7);
		store.reset();
		expect(store.moveHistory).toEqual([]);
		expect(store.board.every((cell) => cell.player === 'none')).toBe(true);
		expect(store.currentPlayer).toBe('red');
	});
});
