# CARO PASTEL ONLINE 8.0

## Phần đã nâng cấp
- Đổi toàn bộ giao diện sang xanh lá pastel, tối ưu cảm ứng và màn hình nhỏ.
- Tích hợp 4 hình nền người dùng cung cấp, chuyển sang WebP để giảm dung lượng tải.
- Thêm chọn hình nền và lưu lựa chọn trên từng thiết bị.
- Giữ chế độ chơi AI, hai người cùng máy và bổ sung phòng online bằng mã 6 ký tự/link chia sẻ.
- Thêm quyền Đỏ, Xanh và người xem; máy chủ kiểm tra token và lượt đi.
- Đồng bộ phòng gần thời gian thực bằng polling 0,9 giây.
- Thêm chống gửi hai nước liên tiếp do nhấp kép hoặc hai yêu cầu đồng thời.
- Đồng hồ online lấy thời gian thực từ máy chủ và xử lý hết giờ trên máy chủ.
- Chuyển frontend sang SvelteKit adapter-static để triển khai GitHub Pages.
- Thêm workflow GitHub Actions, `render.yaml`, script chạy/build Windows và hướng dẫn tiếng Việt.
- Chuyển nơi lưu lịch sử ván từ SQLite sang JSON để backend production không cần CGO.
- Giảm bộ nhớ AI mặc định để phù hợp máy chủ tài nguyên thấp.

## Kiểm tra đã thực hiện
- `npm run check`: 0 lỗi, 0 cảnh báo Svelte.
- Build frontend với `BASE_PATH` kiểu GitHub project site: thành công.
- Build backend Go production: thành công.
- API health, tạo ván AI, đánh với AI, tạo/vào phòng online: thành công.
- Hai yêu cầu nước đi đồng thời của cùng người: một yêu cầu 200, một yêu cầu 409 đúng thiết kế.
- Đồng hồ phòng online giảm theo thời gian máy chủ.

## Giới hạn hiện tại
- Phòng online lưu trong RAM nên mất khi backend restart/deploy.
- Đồng bộ online có độ trễ khoảng 0,9 giây, không phải WebSocket thời gian thực tuyệt đối.
- Bản này không có tài khoản, xếp hạng, chat hoặc lưu phòng lâu dài.
