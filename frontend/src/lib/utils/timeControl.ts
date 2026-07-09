import type { TimeControl } from '$lib/types/game';

export function parseTimeControl(value: TimeControl | string): { initialMin: number; incrementSec: number } {
	const [initial, increment] = String(value).split('+').map((part) => Number.parseInt(part, 10));
	return {
		initialMin: Number.isFinite(initial) ? initial : 0,
		incrementSec: Number.isFinite(increment) ? increment : 0
	};
}

export function timeControlDescription(value: TimeControl | string): string {
	const { initialMin, incrementSec } = parseTimeControl(value);
	return `Mỗi bên có ${initialMin} min tổng thời gian; sau mỗi nước đi được cộng thêm ${incrementSec} giây.`;
}

export function timeControlShort(value: TimeControl | string): string {
	const { initialMin, incrementSec } = parseTimeControl(value);
	return `${initialMin} min + ${incrementSec} giây/nước`;
}
