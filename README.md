# Caro Pastel Online V8.3

Game Caro 16×16 giao diện xanh lá pastel, chơi trên điện thoại và máy tính qua GitHub Pages + Render.

## Tính năng chính
- Giao diện responsive, 4 ảnh nền pastel và lưu lựa chọn nền.
- Chơi với AI 5 cấp độ, hai người cùng máy hoặc hai người online bằng mã/link phòng.
- Phòng online không chạy giờ cho đến khi người chơi thứ hai tham gia.
- Đồng hồ kép: tổng thời gian mỗi bên và giới hạn suy nghĩ cho từng lượt.
- Tùy chọn giới hạn mỗi lượt: 10–90 giây; hết giới hạn thì thua dù tổng giờ vẫn còn.
- Cộng giờ sau mỗi nước theo cấu hình như `7 min + 5 giây/nước`.
- Người thứ ba có thể vào xem phòng.
- Luật thắng: đúng 5 quân liên tiếp; 6 quân trở lên hoặc chuỗi bị chặn cả hai đầu không tính thắng.
- Frontend SvelteKit tĩnh triển khai bằng GitHub Pages; backend Go triển khai bằng Render Blueprint.

## Khởi động nhanh trên Windows
Nhấp đúp:
```text
run_local_windows.bat
```
Sau đó mở:
```text
http://localhost:5173
```

## Build kiểm tra
```text
build_all_windows.bat
```
Kết quả:
- Frontend: `frontend/build/`
- Backend Windows: `dist/caro-server.exe`

## Triển khai online
Xem:
```text
HUONG_DAN_BUILD_DEPLOY.md
```
Cập nhật riêng V8.3:
```text
HUONG_DAN_CAP_NHAT_8_3.md
```

## Cấu trúc
```text
frontend/                       Giao diện SvelteKit
backend/                        API, phòng online và AI engine Go
frontend/static/backgrounds/    Ảnh nền WebP
.github/workflows/              Workflow GitHub Pages
render.yaml                     Render Blueprint
```

## Lưu ý vận hành
Phòng online được lưu trong RAM của backend và tự hết hạn sau 45 phút không hoạt động. Khi Render restart hoặc deploy lại, các phòng đang mở sẽ mất.
