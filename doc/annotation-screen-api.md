# API For Annotation Section

## General

- Thông tin của mỗi response bao gồm:
    - "code": mã lỗi của response -> luôn có.
    - "message": thông báo của response -> luôn có.
    - "data": dữ liệu kèm theo, nếu có.

- Authorization: Bearer + ' ' + token
    - Ví dụ: `Bearer pKsdxxxprnHubASdc0`

## Load tools and labels 
- Mô tả: Lấy danh sách các tools và label của project. Vì vậy, cần cung cấp thông tin project và thông tin người dùng cho server.
- Path: `<host>/api/management/v1/tools`
- Query:
    - `projectId`: number
- Authorization: Bearer
- Sample: `https://www.example.com/api/management/v1/tools?project=1904040007`
- Response:
```json
{
    "code": 1,
    "message": "Success",
    "data": {
        "rectangle": ["car", "person"],
        "polygon": ["girl", "truck", "ball"],
        "line": ["dog tail"]
    }
}
```
