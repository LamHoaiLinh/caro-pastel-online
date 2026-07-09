/** API configuration. Set VITE_API_BASE_URL when building for production. */
const API_BASE_URL = (import.meta.env.VITE_API_BASE_URL || 'http://localhost:5207').replace(/\/$/, '');
export const ApiConfig = {
	baseUrl: API_BASE_URL,
	endpoints: {
		newGame: '/api/game/new',
		game: (id: string) => `/api/game/${id}`,
		move: (id: string) => `/api/game/${id}/move`,
		aiMove: (id: string) => `/api/game/${id}/ai-move`,
		undo: (id: string) => `/api/game/${id}/undo`,
		onlineCreate: '/api/online/create',
		onlineJoin: (code: string) => `/api/online/${encodeURIComponent(code)}/join`,
		onlineRoom: (code: string) => `/api/online/${encodeURIComponent(code)}`,
		onlineMove: (code: string) => `/api/online/${encodeURIComponent(code)}/move`
	} as const
} as const;
