import { describe, it, expect } from 'vitest';
import { ApiConfig } from '$lib/config/apiConfig';

describe('ApiConfig', () => {
	it('có baseUrl HTTP', () => {
		expect(ApiConfig.baseUrl).toMatch(/^http/);
	});

	it('tạo đúng endpoint ván AI', () => {
		expect(ApiConfig.endpoints.newGame).toBe('/api/game/new');
		expect(ApiConfig.endpoints.move('abc123')).toBe('/api/game/abc123/move');
		expect(ApiConfig.endpoints.aiMove('abc123')).toBe('/api/game/abc123/ai-move');
		expect(ApiConfig.endpoints.undo('abc123')).toBe('/api/game/abc123/undo');
	});

	it('tạo đúng endpoint phòng online và mã hóa mã phòng', () => {
		expect(ApiConfig.endpoints.onlineCreate).toBe('/api/online/create');
		expect(ApiConfig.endpoints.onlineJoin('AB C')).toBe('/api/online/AB%20C/join');
		expect(ApiConfig.endpoints.onlineRoom('ABC123')).toBe('/api/online/ABC123');
		expect(ApiConfig.endpoints.onlineMove('ABC123')).toBe('/api/online/ABC123/move');
	});
});
