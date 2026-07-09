<script lang="ts">
	import type { MoveRecord } from '$lib/stores/gameStore.svelte';
	interface Props { moves: MoveRecord[]; currentMoveNumber?: number; }
	let { moves, currentMoveNumber }: Props = $props();
	let scrollContainer: HTMLDivElement | undefined = $state();
	$effect(() => { if (scrollContainer) scrollContainer.scrollLeft = scrollContainer.scrollWidth; });
	function toNotation(x: number, y: number): string { return `${String.fromCharCode(97 + x)}${y + 1}`; }
</script>

<div class="w-full max-w-[900px] mx-auto glass-panel rounded-2xl px-2" data-testid="move-notation">
	{#if moves.length > 0}
		<div bind:this={scrollContainer} class="flex items-center gap-1.5 overflow-x-auto py-2 scrollbar-none">
			{#each moves as move (move.moveNumber)}
				<span class="shrink-0 rounded-lg px-2 py-1 text-xs font-mono font-bold {move.moveNumber === currentMoveNumber ? 'bg-amber-100 text-amber-900' : move.player === 'red' ? 'bg-rose-50 text-rose-700' : 'bg-emerald-50 text-emerald-800'}">
					{move.moveNumber}.{toNotation(move.x, move.y)}
				</span>
			{/each}
		</div>
	{:else}
		<p class="py-2 text-center text-xs font-semibold text-emerald-950/45">Chưa có nước đi</p>
	{/if}
</div>
<style>
	.scrollbar-none { scrollbar-width: none; }
	.scrollbar-none::-webkit-scrollbar { display: none; }
</style>
