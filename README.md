# ỨNG DỤNG "CHIA ĐỀU" - QUẢN LÝ CHI TIÊU NHÓM & QUYẾT TOÁN TỐI ƯU (GOLANG CLI)

> **Dự án CLI độc lập bằng ngôn ngữ Go (Golang)** áp dụng chuẩn thiết kế **OOP trong Go** (`struct` + `method` với `pointer receiver`, encapsulations qua exported/unexported identifiers, phân tách package theo trách nhiệm). 
> 
> Ứng dụng tích hợp cấu trúc dữ liệu **Cây Nhị Phân Tìm Kiếm (BST)** quản lý thành viên và **2 Max-Heaps (`container/heap`)** để rút gọn chuỗi giao dịch nợ tối ưu với độ phức tạp $O(n \log n)$.

---

## 📌 MỤC LỤC
1. [Tính Năng Chính](#-tính-năng-chính)
2. [Cấu Trúc Thư Mục & Packages](#-cấu-trúc-thư-mục--packages)
3. [Cấu Trúc Dữ Liệu & Thuật Toán Trọng Tâm](#-cấu-trúc-dữ-liệu--thuật-toán-trọng-tâm)
   - [1. Cây Nhị Phân Tìm Kiếm (BST)](#1-cây-nhị-phân-tìm-kiếm-bst)
   - [2. Thuật Toán SplitEqual](#2-thuật-toán-splitequal)
   - [3. Thuật Toán CalculateNetBalances](#3-thuật-toán-calculatenetbalances)
   - [4. Thuật Toán SimplifyDebts (2 Max-Heaps)](#4-thuật-toán-simplifydebts-2-max-heaps)
4. [Hướng Dẫn Biên Dịch & Chạy Chương Trình](#-hướng-dẫn-biên-dịch--chạy-chương-trình)
5. [Hướng Dẫn Chạy Test Với File `.txt`](#-hướng-dẫn-chạy-test-với-file-txt)
6. [Kịch Bản Test Mẫu (An, Bình, Châu)](#-kịch-bản-test-mẫu-an-bình-châu)

---

## 🚀 TÍNH NĂNG CHÍNH

- **Quản lý Thành viên bằng BST**:
  - Nhập thành viên từ bàn phím hoặc đọc file `.txt` / `.json`.
  - Thêm 1 thành viên (`BST.Insert`), Tìm kiếm (`BST.Search`), Sửa tên (`BST.Update`).
  - Hủy thành viên (`BST.Remove` xử lý đủ 3 trường hợp: node lá, node 1 con, node 2 con).
  - Duyệt cây theo 3 thứ tự: **NLR** (Tiền tự), **LNR** (Trung tự - tăng dần theo MaTV), **LRN** (Hậu tự).
  - Sắp xếp linh hoạt theo Tên, Tổng đã ứng, Số dư ròng.
- **Quản lý Chi tiêu & Kiểm tra Validation**:
  - Hỗ trợ 4 kiểu chia: `EQUAL` (Chia đều), `PERCENT` (Theo %), `WEIGHT` (Trọng số), `CUSTOM` (Số tiền tùy ý).
  - Ràng buộc kiểm tra: `AmountMinor > 0`, `PaidBy` và `UserID` tham gia bắt buộc phải tồn tại trong BST (`BST.Search`), tổng các phần chia bằng đúng tổng số tiền.
- **Rút gọn nợ bằng 2 Max-Heaps (`container/heap`)**:
  - Rút gọn số lượng giao dịch chuyển tiền giữa các thành viên về mức tối thiểu.
  - Vòng đời trạng thái Settlement: `PENDING` -> `AWAITING_CONFIRMATION` -> `PAID` / `CANCELLED`.
- **Lưu trữ & Đọc File**:
  - Đọc danh sách thành viên từ file văn bản dạng `.txt`.
  - Lưu và đọc toàn bộ trạng thái nhóm ra file cấu trúc `.json`.

---

## 📁 CẤU TRÚC THƯ MỤC & PACKAGES

```
ChiaDeu2/
├── go.mod                     # Go module definitions
├── main.go                    # Entrypoint chính của ứng dụng CLI
├── example_test.go            # Test case minh họa bài toán mẫu & duyệt BST
├── chiadeu                    # File thực thi binary sau khi build
├── data/                      # Thư mục lưu dữ liệu test & json
│   ├── members.txt            # File TXT test nhập 5 thành viên
│   ├── members_large.txt      # File TXT test nhập 10 thành viên
│   ├── input_expenses.txt     # Hướng dẫn kịch bản nhập dữ liệu
│   └── chiadeu.json           # File JSON lưu vết dữ liệu nhóm
├── pkg/                       # Các package nghiệp vụ độc lập
│   ├── member/                # Quản lý BST thành viên
│   │   ├── node.go            # Struct Node & các helper
│   │   └── bst.go             # Struct BST, Insert, Search, Remove (3 TH), Duyệt NLR/LNR/LRN
│   ├── expense/               # Struct khoản chi & các enum
│   │   ├── expense.go         # Struct Expense & SplitType, ExpenseStatus
│   │   └── split.go           # Struct ExpenseSplit
│   ├── settlement/            # Struct quyết toán & đợt chốt sổ
│   │   ├── settlement.go      # Struct Settlement & SettlementStatus
│   │   └── batch.go           # Struct SettlementBatch & BatchStatus
│   └── service/               # Pure functions phục vụ tính toán & kiểm thử
│       ├── expense_calculator.go  # Thuật toán SplitEqual, Percent, Weight, Custom
│       ├── balance_calculator.go  # Thuật toán CalculateNetBalances
│       ├── settlement_resolver.go # Thuật toán SimplifyDebts (2 Max-Heaps)
│       ├── expense_validator.go   # Ràng buộc ValidateExpense
│       ├── storage_service.go     # Đọc/ghi file JSON & import TXT
│       └── service_test.go        # Unit tests cho package service
├── internal/
│   └── cli/
│       └── menu.go            # Giao diện CLI tương tác theo chuẩn khung đề bài
└── tests/
    └── algorithms_test.go     # Bộ kiểm thử thuật toán và các thao tác BST
```

---

## 🛠 CẤU TRÚC DỮ LIỆU & THUẬT TOÁN TRỌNG TÂM

### 1. Cây Nhị Phân Tìm Kiếm (BST)
Quản lý các nút `Node` với khóa chính là `MaTV` (chuỗi duy nhất do người dùng nhập):
- **Xóa Node (`BST.Remove`)**:
  - **Trường hợp 1 (0 con - Node lá)**: Xóa nút, trả về `nil`.
  - **Trường hợp 2 (1 con)**: Thay thế nút bằng nút con duy nhất của nó.
  - **Trường hợp 3 (2 con)**: Tìm nút thế mạng là nút nhỏ nhất thuộc cây con bên phải (*In-order Successor*), sao chép dữ liệu thế mạng sang nút hiện tại, sau đó đệ quy xóa nút thế mạng.

### 2. Thuật toán `SplitEqual`
```go
base := amountMinor / count
remainder := amountMinor % count
// Mỗi người nhận base; Sắp xếp theo MaTV tăng dần; Cộng +1 cho remainder người đầu tiên.
```

### 3. Thuật toán `CalculateNetBalances`
```go
// Bỏ qua expense có Status == VOIDED
balances[expense.PaidBy] += expense.AmountMinor
balances[split.UserID] -= split.ShareMinor
```

### 4. Thuật toán `SimplifyDebts` (2 Max-Heaps)
Cài đặt `heap.Interface` của package `container/heap`:
```go
type BalanceHeap []*HeapItem

func (h BalanceHeap) Less(i, j int) bool {
    if h[i].Amount != h[j].Amount {
        return h[i].Amount > h[j].Amount // Ưu tiên số tiền lớn hơn (Max-Heap)
    }
    return h[i].MaTV < h[j].MaTV         // Ưu tiên MaTV nhỏ hơn (Deterministic tie-breaker)
}
```
Lặp rút gọn nợ giữa phần tử top của `creditors` và `debtors` cho đến khi 1 trong 2 heap rỗng.

---

## 💻 HƯỚNG DẪN BIÊN DỊCH & CHẠY CHƯƠNG TRÌNH

### 1. Chạy tất cả Unit Tests
```bash
go test -v ./...
```

### 2. Biên dịch ra file Binary
```bash
go build -o chiadeu main.go
```

### 3. Khởi chạy ứng dụng CLI
```bash
./chiadeu
# Hoặc chạy trực tiếp bằng:
go run main.go
```

---

## 📝 HƯỚNG DẪN CHẠY TEST VỚI FILE `.TXT`

Ứng dụng hỗ trợ đọc trực tiếp file danh sách thành viên `.txt` dạng `MaTV, Họ và tên` (hoặc phân cách bởi dấu phẩy, dấu gạch đứng `|`, tab):

1. **Khởi chạy ứng dụng**: `./chiadeu`
2. **Chọn menu `1.2`** (Đọc từ file):
   - Nhập đường dẫn: `data/members.txt` (đã tạo sẵn 5 thành viên `TV01` đến `TV05`)
   - Hoặc nhập: `data/members_large.txt` (đã tạo sẵn 10 thành viên)
3. **Chọn menu `2.2`** để kiểm tra cây BST xuất theo thứ tự LNR tăng dần.

---

## 🎯 KỊCH BẢN TEST MẪU (AN, BÌNH, CHÂU)

Bạn có thể chạy nhanh kịch bản kiểm thử bằng 2 cách:

### Cách 1: Gõ chữ `demo` trong Menu chính của ứng dụng CLI
Hệ thống sẽ tự động khởi tạo 3 thành viên:
- **An (TV01)**: Ứng 300.000đ, chịu 100.000đ $\rightarrow$ Số dư ròng: **+200.000đ** (Chủ nợ)
- **Bình (TV02)**: Ứng 0đ, chịu 120.000đ $\rightarrow$ Số dư ròng: **-120.000đ** (Con nợ)
- **Châu (TV03)**: Ứng 60.000đ, chịu 140.000đ $\rightarrow$ Số dư ròng: **-80.000đ** (Con nợ)

#### Kết quả hiển thị:
```text
--- DUYỆT BST THÀNH VIÊN DEMO ---
1. Duyệt NLR (Tiền tự): TV01 (An) -> TV02 (Bình) -> TV03 (Châu)
2. Duyệt LNR (Trung tự): TV01 (An) -> TV02 (Bình) -> TV03 (Châu)
3. Duyệt LRN (Hậu tự): TV03 (Châu) -> TV02 (Bình) -> TV01 (An)

==================================================================================
                   KẾT QUẢ QUYẾT TOÁN TỐI ƯU (EXPECTED RESULT)
==================================================================================
  Giao dịch 1: [Bình (TV02)]  ===>  trả  ===>  [An (TV01)] :  120000 VND
  Giao dịch 2: [Châu (TV03)]  ===>  trả  ===>  [An (TV01)] :  80000 VND
==================================================================================
```

### Cách 2: Chạy kiểm thử tự động với `example_test.go`
```bash
go test -v -run TestExampleScenarioRun
```
