<script lang="ts">
	import { onMount } from 'svelte';
	import { base } from '$app/paths';
	import Board from '$lib/components/Board.svelte';
	import PlayerTimerStrip from '$lib/components/PlayerTimerStrip.svelte';
	import MoveNotation from '$lib/components/MoveNotation.svelte';
	import GameSettings from '$lib/components/GameSettings.svelte';
	import GameResultBanner from '$lib/components/GameResultBanner.svelte';
	import { GameStore } from '$lib/stores/gameStore.svelte';
	import { soundManager } from '$lib/utils/sound';
	import { ApiConfig } from '$lib/config/apiConfig';
	import { GameConfig } from '$lib/config/gameConfig';
	import { switchPlayer } from '$lib/types/game';
	import type { Player, Cell, GameMode, TimeControl, DifficultyLevel, OnlineRole, OnlineRoomState } from '$lib/types/game';
	import { difficultyName } from '$lib/types/game';
	import { recommendedMoveTimeLimit, timeControlDescription, timeControlShort } from '$lib/utils/timeControl';
	import RulesGuide from '$lib/components/RulesGuide.svelte';

	type PlayMode = 'ai' | 'local' | 'online';

	let store = new GameStore();
	let playMode = $state<PlayMode>('ai');
	let gameId = $state('');
	let loading = $state(true);
	let error = $state('');
	let errorMessage = $state('');
	let winningLine = $state<Array<{ x: number; y: number }>>([]);
	let lastMove = $state<{ x: number; y: number } | null>(null);
	let redTime = $state(420);
	let blueTime = $state(420);
	let gameMode = $state<GameMode>('pvai');
	let timeControl = $state<TimeControl>('7+5');
	let moveTimeLimit = $state(30);
	let turnTimeRemaining = $state(30);
	let clockRunning = $state(true);
	let timeoutReason = $state<'move' | 'total' | ''>('');
	let aiSide = $state<'red' | 'blue'>('blue');
	let difficulty = $state<DifficultyLevel>(3);
	let redDifficulty = $state<number | null>(null);
	let blueDifficulty = $state<number | null>(null);
	let isAiThinking = $state(false);
	let moveInProgress = $state(false);

	let onlineLobby = $state(false);
	let playerName = $state('');
	let joinCode = $state('');
	let onlineCode = $state('');
	let onlineRole = $state<OnlineRole>('spectator');
	let playerToken = $state('');
	let redName = $state('Người chơi 1');
	let blueName = $state('Người chơi 2');
	let opponentJoined = $state(false);
	let onlineConnected = $state(false);
	let pollTimer: ReturnType<typeof setInterval> | null = null;


	const boardInteractive = $derived(() => {
		if (store.isGameOver || moveInProgress || isAiThinking) return false;
		if (playMode === 'online') {
			return opponentJoined && onlineRole !== 'spectator' && onlineRole === store.currentPlayer;
		}
		if (gameMode === 'pvai') return store.currentPlayer !== aiSide;
		if (gameMode === 'aivai') return false;
		return true;
	});

	const shareLink = $derived(() => {
		if (!onlineCode || typeof window === 'undefined') return '';
		return `${window.location.origin}${base}/game?mode=online&room=${onlineCode}`;
	});

	const selectedTimeSummary = $derived(timeControlShort(timeControl, moveTimeLimit));
	const selectedTimeDetail = $derived(timeControlDescription(timeControl, moveTimeLimit));

	function aiLabel(side: 'red' | 'blue'): string {
		if (playMode === 'online') return side === 'red' ? redName : blueName;
		if (gameMode === 'pvp') return side === 'red' ? 'Người chơi 1' : 'Người chơi 2';
		const diff = side === 'red' ? redDifficulty : blueDifficulty;
		if (diff == null) return side === 'red' ? 'Người chơi' : 'Người chơi';
		return `AI (${difficultyName(diff as DifficultyLevel)})`;
	}

	function showError(message: string) {
		errorMessage = message;
		setTimeout(() => {
			if (errorMessage === message) errorMessage = '';
		}, 5000);
	}

	async function responseMessage(response: Response): Promise<string> {
		try {
			const data = await response.json();
			return data.message || data.error || `Lỗi ${response.status}`;
		} catch {
			return (await response.text()) || `Lỗi ${response.status}`;
		}
	}

	function syncGameState(state: Record<string, any>, recordNewMove = false) {
		const oldMoveNumber = store.moveNumber;
		const oldBoard = store.board;
		let discovered: { x: number; y: number; player: Player } | null = null;

		if (recordNewMove && state.moveNumber > oldMoveNumber && Array.isArray(state.board)) {
			for (let i = 0; i < oldBoard.length; i++) {
				if (oldBoard[i]?.player === 'none' && state.board[i]?.player !== 'none') {
					discovered = { x: state.board[i].x, y: state.board[i].y, player: state.board[i].player };
					break;
				}
			}
		}

		store.board = state.board;
		store.currentPlayer = state.currentPlayer;
		store.moveNumber = state.moveNumber;
		store.isGameOver = state.isGameOver;
		store.winner = state.winner && state.winner !== 'none' ? state.winner : undefined;
		redTime = state.redTimeRemaining ?? redTime;
		blueTime = state.blueTimeRemaining ?? blueTime;
		turnTimeRemaining = state.turnTimeRemaining ?? turnTimeRemaining;
		clockRunning = state.clockRunning ?? true;
		timeoutReason = state.timeoutReason ?? '';
		if (state.timeControl) timeControl = state.timeControl as TimeControl;
		if (state.moveTimeLimit > 0) moveTimeLimit = state.moveTimeLimit;
		winningLine = state.winningLine ?? [];
		if (state.redDifficulty != null) redDifficulty = state.redDifficulty;
		if (state.blueDifficulty != null) blueDifficulty = state.blueDifficulty;

		if (discovered) {
			lastMove = { x: discovered.x, y: discovered.y };
			store.moveHistory.push({
				moveNumber: state.moveNumber,
				player: discovered.player,
				x: discovered.x,
				y: discovered.y
			});
			soundManager.playStoneSound(discovered.player === 'red' ? 'red' : 'blue');
		}
	}

	function findNewMove(oldBoard: Cell[], newBoard: Cell[]): { x: number; y: number } {
		for (let i = 0; i < oldBoard.length; i++) {
			if (oldBoard[i].player === 'none' && newBoard[i].player !== 'none') return { x: newBoard[i].x, y: newBoard[i].y };
		}
		return { x: 0, y: 0 };
	}

	function resetVisualState() {
		store.reset();
		winningLine = [];
		lastMove = null;
		redDifficulty = null;
		blueDifficulty = null;
		isAiThinking = false;
		moveInProgress = false;
		timeoutReason = '';
		error = '';
	}

	onMount(() => {
		playerName = localStorage.getItem('caro-player-name') || '';
		const params = new URLSearchParams(window.location.search);
		const requested = params.get('mode');
		const room = (params.get('room') || '').trim().toUpperCase();

		if (requested === 'online') {
			playMode = 'online';
			gameMode = 'pvp';
			loading = false;
			onlineLobby = !room;
			if (room) joinOnlineRoom(room);
		} else {
			playMode = requested === 'local' ? 'local' : 'ai';
			gameMode = playMode === 'local' ? 'pvp' : 'pvai';
			createNewGame();
		}

		return () => {
			if (pollTimer) clearInterval(pollTimer);
		};
	});

	async function createNewGame() {
		if (playMode === 'online') {
			await createOnlineRoom();
			return;
		}
		loading = true;
		resetVisualState();
		try {
			const response = await fetch(`${ApiConfig.baseUrl}${ApiConfig.endpoints.newGame}`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({
					timeControl,
					moveTimeLimit,
					gameMode,
					...(gameMode === 'pvai'
						? { [aiSide === 'red' ? 'redDifficulty' : 'blueDifficulty']: difficulty }
						: gameMode === 'aivai'
							? { redDifficulty: difficulty, blueDifficulty: difficulty }
							: {})
				})
			});
			if (!response.ok) throw new Error(await responseMessage(response));
			const data = await response.json();
			gameId = data.gameId;
			syncGameState(data.state);
			if ((gameMode === 'aivai' || (gameMode === 'pvai' && aiSide === 'red')) && !store.isGameOver) {
				setTimeout(makeAiMove, 250);
			}
		} catch (err) {
			error = err instanceof Error ? err.message : 'Không thể tạo ván chơi';
		} finally {
			loading = false;
		}
	}

	async function handleMove(x: number, y: number) {
		if (!boardInteractive()) {
			if (playMode === 'online' && !opponentJoined) showError('Đang chờ người chơi thứ hai vào phòng.');
			else if (playMode === 'online' && onlineRole === 'spectator') showError('Bạn đang ở chế độ xem.');
			else if (playMode === 'online') showError('Chưa đến lượt của bạn.');
			return;
		}
		if (playMode === 'online') {
			await makeOnlineMove(x, y);
			return;
		}
		if (!gameId) return;

		const cell = store.board[y * GameConfig.boardSize + x];
		if (!cell || cell.player !== 'none') return;
		moveInProgress = true;
		const previousPlayer = store.currentPlayer;
		cell.player = previousPlayer;
		soundManager.playStoneSound(previousPlayer === 'red' ? 'red' : 'blue');

		try {
			const response = await fetch(`${ApiConfig.baseUrl}${ApiConfig.endpoints.move(gameId)}`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ x, y })
			});
			if (!response.ok) {
				cell.player = 'none';
				showError(await responseMessage(response));
				return;
			}
			const data = await response.json();
			syncGameState(data.state);
			store.moveHistory.push({ moveNumber: data.state.moveNumber, player: previousPlayer, x, y });
			lastMove = { x, y };
			if (data.state.isGameOver && data.state.winner) soundManager.playWinSound(data.state.winner);
		} catch {
			cell.player = 'none';
			showError('Mất kết nối với máy chủ.');
		} finally {
			moveInProgress = false;
		}

		const aiPlayer = gameMode === 'pvai' && aiSide === 'red' ? 'red' : 'blue';
		if ((gameMode === 'pvai' || gameMode === 'aivai') && !store.isGameOver && store.currentPlayer === aiPlayer) {
			setTimeout(makeAiMove, 180);
		}
	}

	async function makeAiMove() {
		if (!gameId || store.isGameOver || isAiThinking) return;
		isAiThinking = true;
		const previousBoard = store.board.map((cell) => ({ ...cell }));
		const aiPlayer = store.currentPlayer;
		try {
			const response = await fetch(`${ApiConfig.baseUrl}${ApiConfig.endpoints.aiMove(gameId)}`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: '{}'
			});
			if (!response.ok) {
				showError(await responseMessage(response));
				return;
			}
			const data = await response.json();
			const aiMove = findNewMove(previousBoard, data.state.board);
			syncGameState(data.state);
			store.moveHistory.push({ moveNumber: data.state.moveNumber, player: aiPlayer, x: aiMove.x, y: aiMove.y });
			lastMove = aiMove;
			soundManager.playStoneSound(aiPlayer === 'red' ? 'red' : 'blue');
			if (data.state.isGameOver && data.state.winner) soundManager.playWinSound(data.state.winner);
			if (gameMode === 'aivai' && !store.isGameOver) setTimeout(makeAiMove, 260);
		} catch {
			showError('AI không phản hồi. Hãy thử tạo ván mới.');
		} finally {
			isAiThinking = false;
		}
	}

	async function handleUndo() {
		if (!gameId || store.isGameOver || playMode === 'online') return;
		try {
			const response = await fetch(`${ApiConfig.baseUrl}${ApiConfig.endpoints.undo(gameId)}`, { method: 'POST' });
			if (!response.ok) {
				showError(await responseMessage(response));
				return;
			}
			const data = await response.json();
			syncGameState(data.state);
			store.moveHistory = store.moveHistory.slice(0, Math.max(0, store.moveHistory.length - 1));
			winningLine = [];
			lastMove = null;
		} catch {
			showError('Không thể hoàn tác nước đi.');
		}
	}

	function handleTimeOut(player: string) {
		if (store.isGameOver || playMode === 'online') return;
		store.isGameOver = true;
		store.winner = switchPlayer(player as Player);
	}

	function updateOnlineTimeControl(value: TimeControl) {
		timeControl = value;
		moveTimeLimit = recommendedMoveTimeLimit(value);
	}

	function savePlayerName() {
		playerName = playerName.trim().slice(0, 24);
		if (playerName) localStorage.setItem('caro-player-name', playerName);
	}

	async function createOnlineRoom() {
		savePlayerName();
		if (!playerName) {
			showError('Hãy nhập tên của bạn trước khi tạo phòng.');
			return;
		}
		loading = true;
		resetVisualState();
		try {
			const response = await fetch(`${ApiConfig.baseUrl}${ApiConfig.endpoints.onlineCreate}`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ timeControl, moveTimeLimit, playerName })
			});
			if (!response.ok) throw new Error(await responseMessage(response));
			const room: OnlineRoomState = await response.json();
			applyOnlineRoom(room, true);
			onlineLobby = false;
			updateRoomUrl(room.code);
			connectRoomSocket(room.code);
		} catch (err) {
			showError(err instanceof Error ? err.message : 'Không thể tạo phòng online.');
			onlineLobby = true;
		} finally {
			loading = false;
		}
	}

	async function joinOnlineRoom(codeValue = joinCode) {
		const code = codeValue.trim().toUpperCase();
		if (!code) {
			showError('Hãy nhập mã phòng.');
			return;
		}
		savePlayerName();
		loading = true;
		resetVisualState();
		try {
			const savedToken = localStorage.getItem(`caro-room-token-${code}`) || '';
			const response = await fetch(`${ApiConfig.baseUrl}${ApiConfig.endpoints.onlineJoin(code)}`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ playerName, playerToken: savedToken })
			});
			if (!response.ok) throw new Error(await responseMessage(response));
			const room: OnlineRoomState = await response.json();
			applyOnlineRoom(room, true);
			onlineLobby = false;
			joinCode = code;
			updateRoomUrl(code);
			connectRoomSocket(code);
		} catch (err) {
			showError(err instanceof Error ? err.message : 'Không thể vào phòng.');
			onlineLobby = true;
		} finally {
			loading = false;
		}
	}

	function applyOnlineRoom(room: OnlineRoomState, preserveCredentials = false) {
		onlineCode = room.code;
		redName = room.redName || 'Người chơi 1';
		blueName = room.blueName || 'Người chơi 2';
		opponentJoined = room.opponentJoined;
		if (preserveCredentials && room.role) onlineRole = room.role;
		if (preserveCredentials && room.playerToken) {
			playerToken = room.playerToken;
			localStorage.setItem(`caro-room-token-${room.code}`, room.playerToken);
		}
		syncGameState(room.state as unknown as Record<string, any>, true);
	}

	function updateRoomUrl(code: string) {
		const url = `${base}/game?mode=online&room=${encodeURIComponent(code)}`;
		window.history.replaceState({}, '', url);
	}

	async function refreshOnlineRoom(code: string) {
		try {
			const response = await fetch(`${ApiConfig.baseUrl}${ApiConfig.endpoints.onlineRoom(code)}`, { cache: 'no-store' });
			if (!response.ok) throw new Error(await responseMessage(response));
			const room: OnlineRoomState = await response.json();
			applyOnlineRoom(room, false);
			onlineConnected = true;
		} catch {
			onlineConnected = false;
		}
	}

	function connectRoomSocket(code: string) {
		if (pollTimer) clearInterval(pollTimer);
		refreshOnlineRoom(code);
		pollTimer = setInterval(() => {
			if (onlineCode === code) refreshOnlineRoom(code);
		}, 900);
	}

	async function makeOnlineMove(x: number, y: number) {
		if (!onlineCode || !playerToken) return;
		moveInProgress = true;
		try {
			const response = await fetch(`${ApiConfig.baseUrl}${ApiConfig.endpoints.onlineMove(onlineCode)}`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json', 'X-Player-Token': playerToken },
				body: JSON.stringify({ x, y, playerToken })
			});
			if (!response.ok) {
				showError(await responseMessage(response));
				return;
			}
			const room: OnlineRoomState = await response.json();
			applyOnlineRoom(room, false);
			if (room.state.isGameOver && room.state.winner) soundManager.playWinSound(room.state.winner as 'red' | 'blue');
		} catch {
			showError('Mất kết nối khi gửi nước đi.');
		} finally {
			moveInProgress = false;
		}
	}

	async function copyShareLink() {
		try {
			await navigator.clipboard.writeText(shareLink());
			showError('Đã sao chép link phòng.');
		} catch {
			showError('Không thể sao chép tự động. Hãy chọn và sao chép link.');
		}
	}

	function leaveOnlineRoom() {
		if (pollTimer) clearInterval(pollTimer);
		pollTimer = null;
		onlineConnected = false;
		onlineCode = '';
		onlineRole = 'spectator';
		playerToken = '';
		opponentJoined = false;
		onlineLobby = true;
		resetVisualState();
		window.history.replaceState({}, '', `${base}/game?mode=online`);
	}
</script>

{#if errorMessage}
	<div class="fixed top-20 left-1/2 -translate-x-1/2 z-50 w-[min(92vw,640px)] rounded-2xl border border-rose-200 bg-rose-50/95 px-4 py-3 shadow-xl flex items-center justify-between gap-3">
		<p class="text-sm font-semibold text-rose-800">{errorMessage}</p>
		<button onclick={() => errorMessage = ''} class="text-rose-700 font-black text-xl">×</button>
	</div>
{/if}

{#if loading}
	<div class="min-h-[70vh] flex items-center justify-center px-4">
		<div class="glass-panel rounded-3xl px-8 py-7 text-center">
			<div class="mx-auto h-9 w-9 rounded-full border-4 border-emerald-200 border-t-emerald-700 animate-spin"></div>
			<p class="mt-4 font-semibold text-emerald-900">Đang kết nối ván chơi...</p>
		</div>
	</div>
{:else if error}
	<div class="min-h-[70vh] flex items-center justify-center px-4">
		<div class="glass-panel max-w-lg rounded-3xl p-7 text-center">
			<h1 class="text-xl font-extrabold text-rose-700">Không thể mở ván chơi</h1>
			<p class="mt-3 text-emerald-950/75">{error}</p>
			<p class="mt-2 text-xs text-emerald-900/55">Backend: {ApiConfig.baseUrl}</p>
		</div>
	</div>
{:else if playMode === 'online' && onlineLobby}
	<section class="px-3 py-6 sm:py-10">
		<div class="max-w-4xl mx-auto glass-panel rounded-[28px] p-5 sm:p-8">
			<div class="text-center max-w-2xl mx-auto">
				<h1 class="text-3xl sm:text-4xl font-black text-emerald-950">Phòng Caro online</h1>
				<p class="mt-3 text-emerald-950/70">Tạo phòng mới hoặc nhập mã phòng được gửi từ người chơi khác.</p>
				<div class="mt-4 flex justify-center"><RulesGuide /></div>
			</div>
			<div class="mt-7 grid md:grid-cols-2 gap-4">
				<div class="rounded-2xl border border-emerald-200 bg-white/75 p-5">
					<h2 class="text-lg font-extrabold text-emerald-900">Tạo phòng</h2>
					<label class="block mt-4 text-sm font-semibold text-emerald-900">
						Tên của bạn
						<input bind:value={playerName} maxlength="24" placeholder="Ví dụ: Anh Linh" class="mint-input mt-1.5 w-full rounded-xl px-3 py-2.5" />
					</label>
					<label class="block mt-3 text-sm font-semibold text-emerald-900">
						Tổng giờ và cộng giờ
						<select value={timeControl} onchange={(event) => updateOnlineTimeControl((event.currentTarget as HTMLSelectElement).value as TimeControl)} class="mint-input mt-1.5 w-full rounded-xl px-3 py-2.5">
							<option value="1+0">1 min + 0 giây/nước</option>
							<option value="3+0">3 min + 0 giây/nước</option>
							<option value="3+2">3 min + 2 giây/nước</option>
							<option value="7+5">7 min + 5 giây/nước</option>
							<option value="10+0">10 min + 0 giây/nước</option>
							<option value="15+10">15 min + 10 giây/nước</option>
						</select>
					</label>
					<label class="block mt-3 text-sm font-semibold text-emerald-900">
						Giới hạn mỗi lượt
						<select bind:value={moveTimeLimit} class="mint-input mt-1.5 w-full rounded-xl px-3 py-2.5">
							<option value={10}>10 giây/lượt</option>
							<option value={15}>15 giây/lượt</option>
							<option value={20}>20 giây/lượt</option>
							<option value={30}>30 giây/lượt</option>
							<option value={45}>45 giây/lượt</option>
							<option value={60}>60 giây/lượt</option>
							<option value={90}>90 giây/lượt</option>
						</select>
						<p class="mt-2 rounded-xl border border-emerald-200 bg-emerald-50/80 px-3 py-2 text-xs leading-5 text-emerald-800/80">{selectedTimeDetail}</p>
					</label>
					<button onclick={createOnlineRoom} class="mint-button mt-5 w-full rounded-xl px-4 py-3 font-extrabold">Tạo phòng và lấy link</button>
				</div>
				<div class="rounded-2xl border border-amber-200 bg-amber-50/80 p-5">
					<h2 class="text-lg font-extrabold text-amber-900">Vào phòng có sẵn</h2>
					<label class="block mt-4 text-sm font-semibold text-amber-950">
						Tên của bạn
						<input bind:value={playerName} maxlength="24" placeholder="Tên hiển thị" class="mint-input mt-1.5 w-full rounded-xl px-3 py-2.5" />
					</label>
					<label class="block mt-3 text-sm font-semibold text-amber-950">
						Mã phòng 6 ký tự
						<input bind:value={joinCode} maxlength="6" oninput={() => joinCode = joinCode.toUpperCase()} placeholder="ABC123" class="mint-input mt-1.5 w-full rounded-xl px-3 py-2.5 uppercase tracking-[0.25em] font-black" />
					</label>
					<button onclick={() => joinOnlineRoom()} class="mt-5 w-full rounded-xl px-4 py-3 font-extrabold bg-amber-400 text-amber-950 hover:bg-amber-300 transition">Vào phòng</button>
				</div>
			</div>
		</div>
	</section>
{:else}
	<div class="flex flex-col items-center px-1.5 sm:px-4 py-3 sm:py-5 gap-2.5">
		{#if playMode === 'online'}
			<div class="w-full max-w-[900px] glass-panel rounded-2xl p-3 sm:p-4 overflow-hidden sticky top-[58px] z-30">
				<div class="grid grid-cols-1 sm:grid-cols-[1fr_auto] items-center gap-3">
					<div>
						<p class="text-xs font-bold uppercase tracking-wider text-emerald-700">Mã phòng</p>
						<div class="flex items-center gap-2 mt-1">
							<strong class="text-2xl tracking-[0.18em] text-emerald-950">{onlineCode}</strong>
							<span class="inline-flex h-2.5 w-2.5 rounded-full {onlineConnected ? 'bg-emerald-500' : 'bg-amber-500'}"></span>
						</div>
					</div>
					<div class="grid grid-cols-2 gap-2">
						<button onclick={copyShareLink} class="mint-button rounded-xl px-3 py-2 text-sm font-bold">Sao chép link</button>
						<button onclick={leaveOnlineRoom} class="soft-button rounded-xl px-3 py-2 text-sm font-bold">Rời phòng</button>
					</div>
				</div>
				<div class="mt-3 rounded-xl border border-emerald-200 bg-white/80 p-2.5">
					<p class="text-[11px] font-bold uppercase tracking-wide text-emerald-700 mb-1">Link mời bạn chơi</p>
					<input readonly value={shareLink()} class="mint-input w-full min-w-0 rounded-lg px-2.5 py-2 text-xs" aria-label="Link chia sẻ phòng" />
					<p class="mt-1.5 text-[11px] leading-4 text-emerald-800/70">Bấm “Sao chép link”, gửi cho bạn. Người nhận chỉ cần mở link và nhập tên.</p>
				</div>
				<div class="mt-2 flex flex-wrap items-center gap-2 text-xs">
					<span class="rounded-lg bg-emerald-100 px-2.5 py-1.5 font-bold text-emerald-900">Thời gian: {selectedTimeSummary}</span>
					<RulesGuide compact={true} />
				</div>
				<p class="mt-2 text-sm font-semibold {opponentJoined ? 'text-emerald-800' : 'text-amber-800'}">
					{#if onlineRole === 'spectator'}Bạn đang xem ván đấu.{:else}Bạn cầm quân {onlineRole === 'red' ? 'Đỏ (O)' : 'Xanh (X)'}.{/if}
					{opponentJoined ? ' Hai người đã sẵn sàng, đồng hồ bắt đầu chạy.' : ' Đang chờ người thứ hai mở link. Đồng hồ đang tạm dừng và chưa trừ thời gian.'}
				</p>
			</div>
		{:else}
			<GameSettings bind:gameMode bind:timeControl bind:moveTimeLimit bind:aiSide bind:difficulty moveNumber={store.moveNumber} onNewGame={createNewGame} />
		{/if}

		{#if playMode === 'online' && !opponentJoined}
			<div class="w-full max-w-[900px] rounded-xl border border-amber-300 bg-amber-50/95 px-3 py-2.5 text-sm font-bold text-amber-900 text-center">
				Đang chờ người chơi thứ hai · cả tổng giờ và giới hạn mỗi lượt đều chưa chạy
			</div>
		{/if}
		<div class="w-full max-w-[900px] rounded-xl border border-emerald-200 bg-emerald-50/90 px-3 py-2 text-xs sm:text-sm text-emerald-900 flex flex-wrap items-center justify-between gap-2">
			<span><b>Thời gian:</b> {selectedTimeSummary}</span>
			<span class="text-emerald-800/70">{selectedTimeDetail}</span>
		</div>

		<div class="w-full max-w-[900px] glass-panel rounded-2xl px-3 py-2 flex items-center justify-between gap-2 text-sm">
			<div class="font-semibold text-emerald-950">
				Lượt <span class={store.currentPlayer === 'red' ? 'text-[#cf5f57]' : 'text-emerald-700'}>{store.currentPlayer === 'red' ? 'Đỏ O' : 'Xanh X'}</span>
				<span class="text-emerald-950/50"> · Nước {store.moveNumber}</span>
			</div>
			<div class="flex items-center gap-2">
				{#if isAiThinking}<span class="text-xs font-bold text-amber-700 animate-pulse">AI đang nghĩ...</span>{/if}
				{#if playMode !== 'online' && store.moveNumber > 0}
					<button onclick={handleUndo} disabled={store.isGameOver || isAiThinking} class="soft-button rounded-lg px-2.5 py-1.5 text-xs font-bold disabled:opacity-40">Hoàn tác</button>
				{/if}
			</div>
		</div>

		<PlayerTimerStrip player="blue" timeRemaining={blueTime} {turnTimeRemaining} {moveTimeLimit} isActive={clockRunning && (playMode !== 'online' || opponentJoined) && store.currentPlayer === 'blue' && !store.isGameOver} onTimeOut={() => handleTimeOut('blue')} label={aiLabel('blue')} />
		<Board board={store.board} onMove={handleMove} {winningLine} {lastMove} interactive={boardInteractive()} />
		<PlayerTimerStrip player="red" timeRemaining={redTime} {turnTimeRemaining} {moveTimeLimit} isActive={clockRunning && (playMode !== 'online' || opponentJoined) && store.currentPlayer === 'red' && !store.isGameOver} onTimeOut={() => handleTimeOut('red')} label={aiLabel('red')} />
		<MoveNotation moves={store.moveHistory} currentMoveNumber={store.moveNumber} />
	</div>

	{#if store.isGameOver}
		<GameResultBanner winner={store.winner} {timeoutReason} onNewGame={createNewGame} />
	{/if}
{/if}
