# HƯỚNG DẪN BUILD VÀ ĐƯA CARO PASTEL LÊN MẠNG

## 1. Kiến trúc đã dùng
- `frontend/`: SvelteKit, build thành website tĩnh và đưa lên GitHub Pages.
- `backend/`: Go API, xử lý luật cờ, AI và phòng online; đưa lên Render.
- Phòng online dùng mã 6 ký tự và đồng bộ trạng thái bằng HTTP polling khoảng 0,9 giây. Cách này ổn định trên gói miễn phí và không phụ thuộc WebSocket.
- Dữ liệu phòng nằm trong RAM. Khi Render deploy lại hoặc khởi động lại, các phòng đang mở sẽ mất.

## 2. Chạy thử trên Windows
Yêu cầu:
- Node.js 22 trở lên.
- Go 1.23 trở lên.

Cách nhanh:
1. Giải nén dự án.
2. Nhấp đúp `run_local_windows.bat`.
3. Chờ hai cửa sổ Backend và Frontend mở.
4. Mở `http://localhost:5173`.

Chạy thủ công:
```bat
cd backend
go run ./cmd/server
```
Mở cửa sổ CMD thứ hai:
```bat
cd frontend
npm ci
set VITE_API_BASE_URL=http://localhost:5207
npm run dev -- --host 0.0.0.0
```

## 3. Build bản chạy nội bộ trên Windows
Nhấp đúp `build_all_windows.bat`.
Kết quả:
- `dist/caro-server.exe`: backend.
- `dist/frontend/`: website tĩnh.

Lưu ý: mở trực tiếp `dist/frontend/index.html` bằng `file://` không phải cách chạy đúng. Cần một web server tĩnh hoặc dùng `npm run preview`.

## 4. Upload toàn bộ mã nguồn lên GitHub
1. Tạo repository mới trên GitHub, ví dụ `caro-pastel-online`.
2. Không upload các thư mục `frontend/node_modules`, `frontend/.svelte-kit`, `frontend/build`, `dist`.
3. Upload toàn bộ phần còn lại, gồm cả `.github/workflows/deploy-pages.yml` và `render.yaml`.
4. Nhánh chính nên là `main`.

## 5. Deploy backend lên Render
Cách khuyến nghị dùng Blueprint:
1. Đăng nhập Render và kết nối tài khoản GitHub.
2. Chọn **New > Blueprint**.
3. Chọn repository vừa upload.
4. Render đọc file `render.yaml` ở thư mục gốc.
5. Chọn **Deploy Blueprint**.
6. Chờ trạng thái chuyển sang **Live**.
7. Sao chép URL dạng `https://caro-pastel-api-xxxx.onrender.com`.
8. Mở `URL/health`. Nếu hiện `{"status":"ok"}` là backend hoạt động.

Nếu tạo Web Service thủ công:
- Root Directory: `backend`
- Build Command: `go build -trimpath -ldflags="-s -w" -o bin/caro-server ./cmd/server`
- Start Command: `./bin/caro-server`
- Health Check Path: `/health`
- Environment variable: `MATCH_DB_PATH=/tmp/caro-pastel-matches.json`

## 6. Khai báo URL Render trong GitHub
1. Vào repository GitHub.
2. Chọn **Settings > Secrets and variables > Actions > Variables**.
3. Tạo repository variable:
   - Name: `VITE_API_BASE_URL`
   - Value: URL Render, không có dấu `/` ở cuối.
4. Ví dụ: `https://caro-pastel-api-xxxx.onrender.com`.

## 7. Bật GitHub Pages
1. Vào **Settings > Pages**.
2. Trong **Build and deployment**, chọn Source là **GitHub Actions**.
3. Vào tab **Actions**.
4. Chạy workflow **Deploy frontend to GitHub Pages**, hoặc sửa một file rồi push lên nhánh `main`.
5. Khi job màu xanh, mở URL Pages được GitHub cung cấp.

Workflow tự nhận biết repository dạng project site và tự thêm `BASE_PATH`, nên link, ảnh nền và trang `/game` vẫn hoạt động khi website nằm dưới `https://ten-user.github.io/ten-repo/`.

## 8. Cách kiểm tra online bằng hai thiết bị
1. Thiết bị 1 mở website và chọn **Chơi online**.
2. Nhập tên, chọn thời gian, bấm **Tạo phòng và lấy link**.
3. Bấm **Sao chép link** rồi gửi sang thiết bị 2.
4. Thiết bị 2 mở link, nhập tên và vào phòng.
5. Hai thiết bị phải thấy cùng mã phòng và trạng thái “Hai người đã sẵn sàng”.
6. Thử đánh luân phiên. Máy chủ sẽ từ chối người đánh sai lượt.

## 9. Các lỗi thường gặp
### GitHub Actions báo thiếu VITE_API_BASE_URL
Chưa tạo repository variable ở bước 6 hoặc đặt sai tên.

### Website mở được nhưng báo mất kết nối
- Render chưa Live.
- URL trong `VITE_API_BASE_URL` sai hoặc có dấu `/` cuối.
- Sau khi đổi variable cần chạy lại workflow Pages.

### Render ngủ sau thời gian không dùng
Gói miễn phí có thể cần thời gian khởi động lại ở lần truy cập đầu. Chờ một lúc rồi thử lại.

### Phòng mất sau khi Render deploy
Đây là thiết kế hiện tại: phòng online lưu trong RAM. Muốn giữ phòng qua restart cần chuyển RoomStore sang Redis hoặc PostgreSQL.

### AI quá nặng trên gói miễn phí
Bản này đã giảm Transposition Table mặc định xuống 32 MB mỗi AI. Nên để độ khó mặc định mức 3; mức 5 có thể phản hồi chậm hơn khi Render đang tải cao.
