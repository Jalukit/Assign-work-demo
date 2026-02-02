# Assign Work Demo (React + Go + MongoDB)

ฟีเจอร์:
- Register/Login (JWT)
- Role: user / admin
- Admin สร้างงาน (capacity)
- User เห็นงานทั้งหมด + จำนวนที่เหลือ + รายชื่อคนที่ลงทะเบียน
- User ลงทะเบียน/ยกเลิกได้

## วิธีรันด้วย Docker Compose
1) ไปที่โฟลเดอร์ backend แล้ว copy env:
   cp backend/.env.example backend/.env
   (แก้ JWT_SECRET ได้)
2) Run:
   docker compose up

เปิด:
- Frontend: http://localhost:5173
- Backend:  http://localhost:8080/health

## หมายเหตุด้านความปลอดภัย (เดโม่)
- หน้า Register อนุญาตเลือก role ได้ (สะดวกทดสอบ) ของจริงอย่าทำแบบนี้
- Endpoint registerTask ใช้ check capacity แบบ 2-step เพื่อความอ่านง่าย
  ของจริงควรทำ update แบบ atomic + validation ที่แน่นกว่า
