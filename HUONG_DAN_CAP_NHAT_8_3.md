# CẬP NHẬT CARO PASTEL V8.3

## Nội dung sửa chính
- Phòng online mới tạo ở trạng thái chờ: tổng giờ và đồng hồ từng lượt không chạy khi người chơi thứ hai chưa vào.
- Đồng hồ bắt đầu đúng thời điểm người chơi thứ hai tham gia phòng thành công.
- Thêm giới hạn suy nghĩ cho từng lượt: 10, 15, 20, 30, 45, 60 hoặc 90 giây/lượt.
- Nếu người đang đi dùng hết giới hạn của một lượt thì thua ngay, kể cả tổng thời gian vẫn còn.
- Tổng thời gian vẫn bị trừ bình thường trong lúc suy nghĩ; sau nước hợp lệ vẫn được cộng số giây theo cấu hình.
- Mỗi thanh người chơi hiển thị riêng “Tổng giờ” và “Nước này”.
- Màn hình kết quả ghi rõ thua vì hết tổng giờ hay hết giới hạn một lượt.
- Áp dụng đồng bộ cho chơi online, hai người cùng máy và chơi với AI.

## Ví dụ cách tính
Cấu hình: `7 min/bên · cộng 5 giây/nước · tối đa 30 giây/lượt`.
- Mỗi bên bắt đầu với 7 min tổng giờ.
- Trong lượt của mình, cả tổng giờ và đồng hồ 30 giây cùng chạy.
- Đánh hợp lệ sau 12 giây: tổng giờ bị trừ 12 giây rồi cộng lại 5 giây.
- Không đánh trong 30 giây: thua ngay, dù tổng giờ còn nhiều.
- Phòng online chưa đủ hai người: cả hai đồng hồ đứng yên.

## Cách cập nhật lên GitHub
1. Giải nén file ZIP V8.3.
2. Mở thư mục `caro-pastel-online-v8.3`.
3. Upload toàn bộ nội dung bên trong lên repository `caro-pastel-online`, ghi đè các file trùng tên.
4. Commit trực tiếp vào nhánh `main` với nội dung: `Fix online clock and add per-move limit v8.3`.
5. Chờ GitHub Actions chạy xong cả `build` và `deploy`.
6. Vì V8.3 có sửa backend, Render cũng phải tự deploy lại. Vào Render → `caro-pastel-api` → Events và chờ trạng thái `Live`.
7. Trên điện thoại đóng tab cũ rồi mở lại. Nếu vẫn thấy bản cũ, xóa dữ liệu trang `lamhoailinh.github.io` trong Safari.

## Kiểm tra sau cập nhật
1. Tạo phòng online nhưng chưa cho người thứ hai vào: tổng giờ phải đứng nguyên.
2. Cho người thứ hai mở link: đồng hồ của quân Đỏ bắt đầu chạy.
3. Chọn giới hạn 10 giây/lượt và không đánh: sau khoảng 10 giây quân đang đi phải thua.
4. Tạo phòng mới, đánh trước khi hết giới hạn: đồng hồ “Nước này” phải trở về đủ số giây cho người kế tiếp.
