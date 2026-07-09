# Caro Pastel Online

Phiên bản 8.1 giao diện xanh lá pastel cho game Caro 16×16, tối ưu để chơi trên điện thoại và máy tính.

## Tính năng chính

- Giao diện pastel xanh lá, responsive cho màn hình nhỏ và lớn.
- 4 ảnh nền đi kèm, đổi nền ngay trong giao diện và tự ghi nhớ lựa chọn.
- Chơi với AI gồm 5 cấp độ.
- Hai người chơi chung một thiết bị.
- Chơi online bằng mã phòng 6 ký tự hoặc gửi đường dẫn phòng cho người khác.
- Người thứ ba có thể vào phòng ở chế độ xem.
- Luật thắng đúng 5 quân liên tiếp và luật khai cuộc của dự án gốc được giữ lại.
- Frontend xuất thành website tĩnh bằng SvelteKit để triển khai trên GitHub Pages.
- Backend Go có thể triển khai bằng Render Blueprint từ file `render.yaml`.

## Khởi động nhanh trên Windows

Nhấp đúp:

```text
run_local_windows.bat
```

Sau khi cửa sổ backend và frontend mở, truy cập:

```text
http://localhost:5173
```

## Build kiểm tra toàn bộ

Nhấp đúp:

```text
build_all_windows.bat
```

Kết quả:

- Frontend tĩnh: `frontend/build/`
- Backend Windows: `dist/caro-server.exe`

## Triển khai online

Xem hướng dẫn từng bước tại:

```text
HUONG_DAN_BUILD_DEPLOY.md
```

Luồng triển khai đề nghị:

1. Đưa toàn bộ source lên GitHub.
2. Tạo backend trên Render bằng `render.yaml`.
3. Lấy URL backend Render và tạo GitHub Actions variable `VITE_API_BASE_URL`.
4. Bật GitHub Pages với nguồn `GitHub Actions`.
5. Workflow `.github/workflows/deploy-pages.yml` tự build và xuất bản frontend.

## Cấu trúc chính

```text
frontend/                         Giao diện SvelteKit
backend/                          API và AI engine viết bằng Go
frontend/static/backgrounds/      Ảnh nền đã tối ưu WebP
.github/workflows/                Workflow GitHub Pages
render.yaml                       Cấu hình Render Blueprint
HUONG_DAN_BUILD_DEPLOY.md          Hướng dẫn build và deploy tiếng Việt
```

## API online room

```text
POST /api/online/create
POST /api/online/{code}/join
GET  /api/online/{code}
POST /api/online/{code}/move
```

Phòng online hiện được giữ trong RAM của backend và tự hết hạn sau 45 phút không hoạt động. Khi Render khởi động lại hoặc deploy lại, các phòng đang mở sẽ mất; đây là thiết kế phù hợp cho bản đơn giản, không đăng nhập.

## Yêu cầu môi trường phát triển

- Node.js 20 trở lên
- Go theo phiên bản ghi trong `backend/go.mod`
- npm

## Kiểm tra thủ công đề nghị trước khi public

- Chơi AI cấp 1 và cấp 5 trên điện thoại.
- Tạo phòng bằng máy tính, tham gia bằng điện thoại qua mạng 4G/5G.
- Kiểm tra đường dẫn GitHub Pages sau khi đổi tên repository.
- Kiểm tra thời gian phản hồi đầu tiên của Render sau thời gian không sử dụng.
