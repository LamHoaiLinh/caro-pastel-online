<script lang="ts">
	import type { Player } from '$lib/types/game';
	interface Props { winner: Player | undefined; onNewGame: () => void; }
	let { winner, onNewGame }: Props = $props();
</script>

<div class="fixed inset-0 z-50 bg-emerald-950/30 backdrop-blur-sm flex items-start justify-center px-4 pt-24">
	<div class="glass-panel rounded-3xl p-6 sm:p-8 w-full max-w-md text-center shadow-2xl">
		<p class="text-xs font-bold uppercase tracking-[0.22em] text-emerald-700">Kết thúc ván</p>
		<h2 class="mt-3 text-3xl font-black {winner === 'red' ? 'text-[#c8564e]' : 'text-emerald-800'}">
			{winner === 'red' ? 'Quân Đỏ chiến thắng' : winner === 'blue' ? 'Quân Xanh chiến thắng' : 'Ván đấu kết thúc'}
		</h2>
		<button onclick={onNewGame} class="mint-button mt-6 rounded-xl px-6 py-3 font-extrabold w-full">Tạo ván mới</button>
	</div>
</div>
