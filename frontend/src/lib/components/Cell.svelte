<script lang="ts">
	import type { Player } from '$lib/types/game';

	interface Props {
		x: number;
		y: number;
		player: Player;
		isLastMove?: boolean;
		isOpenRuleInvalid?: boolean;
		cellSize: number;
		interactive?: boolean;
		onclick?: () => void;
		onkeydown?: (e: KeyboardEvent) => void;
	}

	let {
		x,
		y,
		player,
		isLastMove = false,
		isOpenRuleInvalid = false,
		cellSize,
		interactive = true,
		onclick,
		onkeydown
	}: Props = $props();
</script>

<button
	onclick={onclick}
	onkeydown={onkeydown}
	disabled={!interactive || player !== 'none'}
	class="caro-cell {isLastMove ? 'last-move' : ''} {isOpenRuleInvalid ? 'open-rule-invalid' : ''}"
	style="width: {cellSize}px; height: {cellSize}px; min-width: {cellSize}px; min-height: {cellSize}px; font-size: {cellSize * 0.56}px;"
	aria-label="Ô {x + 1},{y + 1}"
	data-x={x}
	data-y={y}
>
	{#if player === 'red'}
		<span class="stone stone-red">O</span>
	{:else if player === 'blue'}
		<span class="stone stone-blue">X</span>
	{/if}
</button>

<style>
	.caro-cell {
		display: flex;
		align-items: center;
		justify-content: center;
		position: relative;
		border: 1px solid rgba(92, 156, 118, 0.42);
		background: rgba(245, 253, 247, 0.9);
		transition: background 120ms ease, transform 120ms ease;
	}
	.caro-cell:not(:disabled):hover {
		background: rgba(211, 241, 222, 0.98);
	}
	.caro-cell:not(:disabled):active {
		transform: scale(0.94);
	}
	.caro-cell.last-move {
		background: rgba(255, 236, 174, 0.92);
		box-shadow: inset 0 0 0 2px rgba(204, 151, 36, 0.54);
	}
	.caro-cell.open-rule-invalid {
		background: rgba(255, 224, 224, 0.65);
		opacity: 0.58;
	}
	.stone {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		width: 82%;
		height: 82%;
		font-weight: 900;
		line-height: 1;
		text-shadow: 0 1px 0 rgba(255,255,255,0.65);
	}
	.stone-red { color: #dd6b62; }
	.stone-blue { color: #247861; }
</style>
