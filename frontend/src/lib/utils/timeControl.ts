import type { TimeControl } from '$lib/types/game';

export function parseTimeControl(value: TimeControl | string): { initialMin: number; incrementSec: number } {
	const [initial, increment] = String(value).split('+').map((part) => Number.parseInt(part, 10));
	return {
		initialMin: Number.isFinite(initial) ? initial : 7,
		incrementSec: Number.isFinite(increment) ? increment : 0
	};
}

export function recommendedMoveTimeLimit(value: TimeControl | string): number {
	switch (String(value)) {
		case '1+0': return 10;
		case '3+0':
		case '3+2': return 20;
		case '10+0': return 45;
		case '15+10': return 60;
		default: return 30;
	}
}

export function timeControlDescription(value: TimeControl | string, moveTimeLimit = recommendedMoveTimeLimit(value)): string {
	const { initialMin, incrementSec } = parseTimeControl(value);
	return `Mỗi bên có ${initialMin} min tổng thời gian. Sau mỗi nước hợp lệ được cộng ${incrementSec} giây. Mỗi lượt chỉ được nghĩ tối đa ${moveTimeLimit} giây; hết giới hạn này sẽ thua dù tổng giờ vẫn còn.`;
}

export function timeControlShort(value: TimeControl | string, moveTimeLimit = recommendedMoveTimeLimit(value)): string {
	const { initialMin, incrementSec } = parseTimeControl(value);
	return `${initialMin} min/bên · cộng ${incrementSec} giây/nước · tối đa ${moveTimeLimit} giây/lượt`;
}
