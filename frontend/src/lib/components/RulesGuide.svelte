<script lang="ts">
	interface Props { compact?: boolean; }
	let { compact = false }: Props = $props();
	let open = $state(false);

	const win = ['O','O','O','O','O'];
	const overline = ['O','O','O','O','O','O'];
	const blocked = ['X','O','O','O','O','O','X'];
</script>

<button onclick={() => open = true} class="soft-button rounded-xl px-3 py-2 text-sm font-bold whitespace-nowrap">
	{compact ? 'Luật chơi' : 'Hướng dẫn luật chơi'}
</button>

{#if open}
	<div class="fixed inset-0 z-[100] bg-emerald-950/55 p-3 sm:p-6 overflow-y-auto" role="presentation" onclick={(e) => e.currentTarget === e.target && (open = false)}>
		<div class="mx-auto max-w-3xl rounded-3xl bg-[#f6fff8] border border-emerald-200 shadow-2xl overflow-hidden" role="dialog" aria-modal="true" aria-label="Hướng dẫn luật chơi Caro">
			<header class="sticky top-0 z-10 flex items-center justify-between gap-3 bg-emerald-900 px-4 sm:px-6 py-4 text-white">
				<div>
					<h2 class="text-xl sm:text-2xl font-black">Luật chơi Caro 16 × 16</h2>
					<p class="text-xs sm:text-sm text-emerald-100 mt-1">Giải thích bằng hình minh họa dễ hiểu</p>
				</div>
				<button onclick={() => open = false} class="rounded-full bg-white/15 h-10 w-10 text-2xl font-bold" aria-label="Đóng">×</button>
			</header>
			<div class="p-4 sm:p-6 space-y-4 text-emerald-950">
				<article class="rule-card">
					<div class="rule-number">1</div>
					<div>
						<h3>Hai bên lần lượt đánh</h3>
						<p>Quân <b class="text-[#cf5f57]">Đỏ O</b> đi trước, quân <b class="text-emerald-700">Xanh X</b> đi sau. Mỗi lượt chỉ đặt một quân vào ô trống.</p>
						<div class="mini-board mt-3"><span class="red">O</span><span></span><span class="blue">X</span><span></span><span></span><span></span><span></span><span></span><span></span></div>
					</div>
				</article>
				<article class="rule-card">
					<div class="rule-number">2</div>
					<div>
						<h3>Thắng khi có đúng 5 quân liên tiếp</h3>
						<p>Năm quân cùng màu nằm liền nhau theo hàng ngang, dọc hoặc chéo sẽ thắng.</p>
						<div class="line-demo valid">{#each win as stone}<span>{stone}</span>{/each}</div>
					</div>
				</article>
				<article class="rule-card">
					<div class="rule-number">3</div>
					<div>
						<h3>Sáu quân trở lên không tính thắng</h3>
						<p>Phiên bản này dùng luật “Exact 5”: phải đúng 5 quân, không phải 6 hoặc nhiều hơn.</p>
						<div class="line-demo invalid">{#each overline as stone}<span>{stone}</span>{/each}</div>
					</div>
				</article>
				<article class="rule-card">
					<div class="rule-number">4</div>
					<div>
						<h3>Bị chặn cả hai đầu thì không thắng</h3>
						<p>Dù có đúng 5 quân, nếu hai đầu đều bị quân đối phương chặn thì không được tính thắng.</p>
						<div class="line-demo invalid blocked">{#each blocked as stone, i}<span class={i === 0 || i === blocked.length - 1 ? 'enemy' : ''}>{stone}</span>{/each}</div>
					</div>
				</article>
				<article class="rule-card">
					<div class="rule-number">5</div>
					<div>
						<h3>Cách đọc thời gian</h3>
						<p>Ví dụ <b>7 min + 5 giây/nước</b>: mỗi bên bắt đầu với 7 phút tổng thời gian; sau khi đặt xong mỗi quân, đồng hồ bên đó được cộng thêm 5 giây.</p>
				</div>
				</article>
				<div class="rounded-2xl bg-amber-50 border border-amber-200 p-4 text-sm leading-6 text-amber-950">
					<b>Lưu ý:</b> Khi hết giờ, bên còn lại thắng. Ở chế độ online, hãy giữ tab mở và mạng ổn định để đồng hồ đồng bộ đúng.
				</div>
			</div>
		</div>
	</div>
{/if}

<style>
	.rule-card{display:grid;grid-template-columns:2.25rem 1fr;gap:.8rem;padding:1rem;border:1px solid #ccebd6;border-radius:1rem;background:rgba(255,255,255,.86)}
	.rule-number{display:flex;align-items:center;justify-content:center;width:2.25rem;height:2.25rem;border-radius:999px;background:#2f8a64;color:white;font-weight:900}
	h3{font-weight:900;font-size:1.05rem;margin:0 0 .35rem}
	p{font-size:.92rem;line-height:1.55;color:rgba(6,78,59,.82)}
	.mini-board{display:grid;grid-template-columns:repeat(3,2.15rem);width:max-content;border:1px solid #9acbb0}
	.mini-board span{display:flex;align-items:center;justify-content:center;height:2.15rem;border:1px solid #b8dcc5;font-weight:900;font-size:1.25rem;background:#effaf3}
	.red{color:#cf5f57}.blue{color:#187b58}
	.line-demo{display:flex;flex-wrap:wrap;gap:.2rem;margin-top:.75rem}
	.line-demo span{display:flex;align-items:center;justify-content:center;width:2rem;height:2rem;border-radius:.45rem;font-weight:900;background:#ffe8e5;color:#c84f49;border:1px solid #efb5b0}
	.line-demo.valid{padding:.45rem;border:2px solid #31a66f;border-radius:.8rem;width:max-content;max-width:100%}
	.line-demo.invalid{padding:.45rem;border:2px dashed #dc6b63;border-radius:.8rem;width:max-content;max-width:100%;position:relative}
	.line-demo.invalid::after{content:'Không tính';margin-left:.45rem;align-self:center;color:#b42318;font-size:.78rem;font-weight:900}
	.line-demo .enemy{background:#e1f4e8;color:#16714f;border-color:#a8d9bb}
	@media(max-width:420px){.rule-card{grid-template-columns:1.9rem 1fr;padding:.8rem}.rule-number{width:1.9rem;height:1.9rem}.line-demo span{width:1.75rem;height:1.75rem}}
</style>
