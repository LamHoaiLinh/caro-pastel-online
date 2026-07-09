<script lang="ts">
	import { base } from '$app/paths';
</script>

<section class="px-3 py-8 sm:py-14">
	<div class="max-w-6xl mx-auto glass-panel rounded-[28px] overflow-hidden">
		<div class="grid lg:grid-cols-[1.08fr_0.92fr] items-stretch">
			<div class="p-6 sm:p-10 lg:p-14">
				<span class="inline-flex rounded-full bg-emerald-100/90 border border-emerald-200 px-3 py-1 text-xs font-bold uppercase tracking-wider text-emerald-800">
					Caro 16 × 16 · Exact 5
				</span>
				<h1 class="mt-5 text-4xl sm:text-5xl lg:text-6xl font-black tracking-tight text-emerald-950 leading-[1.05]">
					Cờ Caro pastel cho điện thoại và máy tính
				</h1>
				<p class="mt-5 max-w-2xl text-base sm:text-lg leading-7 text-emerald-950/75">
					Chơi với AI, hai người trên cùng thiết bị hoặc tạo phòng online rồi gửi link cho bạn bè.
				</p>
				<div class="mt-8 grid sm:grid-cols-3 gap-3">
					<a href={`${base}/game?mode=ai`} class="mint-button rounded-2xl p-4 text-center font-extrabold">
						<span class="block text-lg">Chơi với AI</span>
						<span class="block mt-1 text-xs font-medium text-white/80">5 cấp độ</span>
					</a>
					<a href={`${base}/game?mode=local`} class="soft-button rounded-2xl p-4 text-center font-extrabold">
						<span class="block text-lg">2 người cùng máy</span>
						<span class="block mt-1 text-xs font-medium text-emerald-900/65">Không cần tài khoản</span>
					</a>
					<a href={`${base}/game?mode=online`} class="rounded-2xl p-4 text-center font-extrabold bg-amber-100/90 border border-amber-200 text-amber-900 hover:bg-amber-200/90 transition">
						<span class="block text-lg">Chơi online</span>
						<span class="block mt-1 text-xs font-medium text-amber-900/65">Tạo mã phòng</span>
					</a>
				</div>
			</div>
			<div class="bg-emerald-900/92 p-6 sm:p-10 text-emerald-50">
				<h2 class="text-xl font-extrabold">Luật đang áp dụng</h2>
				<div class="mt-5 space-y-4 text-sm sm:text-base leading-6 text-emerald-50/85">
					<p><strong class="text-white">Bàn cờ:</strong> 16 × 16 ô.</p>
					<p><strong class="text-white">Thắng:</strong> đúng 5 quân liên tiếp; 6 quân trở lên không tính.</p>
					<p><strong class="text-white">Hai đầu bị chặn:</strong> không tính thắng.</p>
					<p><strong class="text-white">Luật mở:</strong> nước thứ hai của bên Đỏ phải cách nước đầu tối thiểu 3 giao điểm.</p>
				</div>
				<div class="mt-8 rounded-2xl bg-white/10 border border-white/10 p-4 text-sm leading-6">
					Giao diện tự co theo màn hình, hỗ trợ cảm ứng, âm thanh, rung nhẹ và lưu hình nền đã chọn trên thiết bị.
				</div>
			</div>
		</div>
	</div>
</section>
