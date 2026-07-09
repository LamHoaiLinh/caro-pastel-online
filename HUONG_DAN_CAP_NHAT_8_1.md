# HƯỚNG DẪN CẬP NHẬT LỖI 8.1

Lỗi `open rule violation` xuất hiện vì backend giữ một luật khai cuộc ẩn: nước thứ hai của quân Đỏ phải cách quân Đỏ đầu tiên ít nhất 3 ô theo một trục. Bản 8.1 đã bỏ luật này để chơi Caro thông thường.

## Cách cập nhật nhanh trên GitHub
1. Giải nén file ZIP bản 8.1.
2. Trong repository GitHub, chọn **Add file → Upload files**.
3. Kéo toàn bộ nội dung bên trong thư mục bản 8.1 vào vùng upload và cho phép ghi đè các file cũ.
4. Commit với nội dung: `Fix open rule violation v8.1`.
5. Render sẽ tự deploy lại backend từ nhánh `main`.
6. GitHub Actions sẽ tự build lại frontend.
7. Chờ cả Render và workflow GitHub Pages hoàn tất, sau đó nhấn `Ctrl + F5` để bỏ cache.

## Các file cốt lõi đã sửa
- `backend/internal/domain/game.go`
- `backend/internal/engine/search.go`
- `backend/internal/engine/parallel.go`
- `frontend/src/routes/game/+page.svelte`
- `frontend/src/lib/components/Board.svelte`
- `frontend/src/lib/components/Cell.svelte`
- `frontend/package-lock.json`

## Kiểm tra sau khi cập nhật
- Đánh quân O lần đầu.
- Chờ AI đánh quân X.
- Đánh quân O lần thứ hai ngay cạnh hoặc gần quân O đầu tiên.
- Không còn thông báo `open rule violation`.
