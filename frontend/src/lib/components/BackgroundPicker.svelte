<script lang="ts">
	import { base } from '$app/paths';

	const backgrounds = [
		{ id: 'workspace', name: 'Góc làm việc', file: 'workspace.webp' },
		{ id: 'mint-family', name: 'Gia đình xanh', file: 'mint-family.webp' },
		{ id: 'pastel-garden', name: 'Vườn cổ tích', file: 'pastel-garden.webp' },
		{ id: 'cat-tulips', name: 'Mèo và tulip', file: 'cat-tulips.webp' }
	];

	interface Props {
		selected: string;
		onSelect: (id: string) => void;
	}

	let { selected, onSelect }: Props = $props();
	let open = $state(false);
</script>

<div class="relative">
	<button
		onclick={() => open = !open}
		class="soft-button rounded-xl px-3 py-2 text-xs sm:text-sm font-semibold"
		aria-label="Đổi hình nền"
	>
		Hình nền
	</button>
	{#if open}
		<div class="absolute right-0 mt-2 w-64 rounded-2xl p-3 glass-panel z-50">
			<p class="text-xs font-bold uppercase tracking-wide text-emerald-800 mb-2">Chọn hình nền</p>
			<div class="grid grid-cols-2 gap-2">
				{#each backgrounds as bg}
					<button
						onclick={() => { onSelect(bg.id); open = false; }}
						class="overflow-hidden rounded-xl border-2 text-left transition {selected === bg.id ? 'border-emerald-500' : 'border-white/70'}"
					>
						<img src={`${base}/backgrounds/${bg.file}`} alt={bg.name} class="h-16 w-full object-cover" />
						<span class="block bg-white/90 px-2 py-1 text-[11px] font-semibold text-emerald-900">{bg.name}</span>
					</button>
				{/each}
			</div>
		</div>
	{/if}
</div>
