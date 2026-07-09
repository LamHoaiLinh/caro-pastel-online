<script lang="ts">
	import type { GameMode, TimeControl, DifficultyLevel } from '$lib/types/game';
	import { difficultyName } from '$lib/types/game';

	interface Props {
		gameMode: GameMode;
		timeControl: TimeControl;
		aiSide: 'red' | 'blue';
		difficulty: DifficultyLevel;
		moveNumber: number;
		onNewGame: () => void;
	}

	let {
		gameMode = $bindable(),
		timeControl = $bindable(),
		aiSide = $bindable(),
		difficulty = $bindable(),
		moveNumber,
		onNewGame
	}: Props = $props();

	let isOpen = $state(true);

	$effect(() => {
		if (moveNumber > 0) isOpen = false;
	});

	const modeText = $derived(gameMode === 'pvai' ? 'Người vs AI' : gameMode === 'aivai' ? 'AI vs AI' : 'Hai người cùng máy');
</script>

<div class="w-full max-w-[900px] mx-auto">
	<button
		onclick={() => isOpen = !isOpen}
		class="glass-panel rounded-2xl px-4 py-3 w-full flex items-center justify-between gap-3 text-sm font-semibold text-emerald-900"
	>
		<span>{modeText} · {timeControl}</span>
		<span class="transition-transform {isOpen ? 'rotate-180' : ''}">⌄</span>
	</button>

	{#if isOpen}
		<div class="glass-panel rounded-2xl mt-2 p-4 space-y-4">
			<div>
				<p class="text-xs font-bold uppercase tracking-wider text-emerald-800 mb-2">Chế độ</p>
				<div class="grid grid-cols-3 gap-2">
					<button onclick={() => gameMode = 'pvp'} disabled={moveNumber > 0} class="rounded-xl px-3 py-2 text-sm font-bold {gameMode === 'pvp' ? 'mint-button' : 'soft-button'}">2 người</button>
					<button onclick={() => gameMode = 'pvai'} disabled={moveNumber > 0} class="rounded-xl px-3 py-2 text-sm font-bold {gameMode === 'pvai' ? 'mint-button' : 'soft-button'}">Với AI</button>
					<button onclick={() => gameMode = 'aivai'} disabled={moveNumber > 0} class="rounded-xl px-3 py-2 text-sm font-bold {gameMode === 'aivai' ? 'mint-button' : 'soft-button'}">AI đấu AI</button>
				</div>
			</div>

			<div class="grid sm:grid-cols-3 gap-3">
				<label class="text-sm font-semibold text-emerald-900">
					<span class="block mb-1.5">Thời gian</span>
					<select bind:value={timeControl} disabled={moveNumber > 0} class="mint-input rounded-xl px-3 py-2.5 w-full">
						<option value="1+0">1+0</option>
						<option value="3+0">3+0</option>
						<option value="3+2">3+2</option>
						<option value="7+5">7+5</option>
						<option value="10+0">10+0</option>
						<option value="15+10">15+10</option>
					</select>
				</label>
				{#if gameMode !== 'pvp'}
					<label class="text-sm font-semibold text-emerald-900">
						<span class="block mb-1.5">Độ khó AI</span>
						<select bind:value={difficulty} disabled={moveNumber > 0} class="mint-input rounded-xl px-3 py-2.5 w-full">
							{#each [1,2,3,4,5] as level}
								<option value={level}>{level} · {difficultyName(level as DifficultyLevel)}</option>
							{/each}
						</select>
					</label>
				{/if}
				{#if gameMode === 'pvai'}
					<label class="text-sm font-semibold text-emerald-900">
						<span class="block mb-1.5">AI cầm quân</span>
						<select bind:value={aiSide} disabled={moveNumber > 0} class="mint-input rounded-xl px-3 py-2.5 w-full">
							<option value="blue">Xanh · đi sau</option>
							<option value="red">Đỏ · đi trước</option>
						</select>
					</label>
				{/if}
			</div>

			<button onclick={onNewGame} class="mint-button rounded-xl px-5 py-2.5 font-bold w-full sm:w-auto">Bắt đầu ván mới</button>
		</div>
	{/if}
</div>
