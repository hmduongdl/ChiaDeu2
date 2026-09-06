package service

import (
	"container/heap"
	"fmt"
	"time"

	"chiadeu/pkg/settlement"
)

// HeapItem biểu diễn một phần tử trong Hàng Đợi Ưu Tiên (Binary Heap)
type HeapItem struct {
	MaTV   string // Mã thành viên (MaTV)
	Amount int64  // Số tiền (dương cho cả chủ nợ và con nợ)
}

// BalanceHeap triển khai heap.Interface để tạo Hàng Đợi Nhị Phân (Binary Heap)
// Đặt tên riêng creditorHeap và debtorHeap như yêu cầu đề bài
type BalanceHeap []*HeapItem

func (h BalanceHeap) Len() int { return len(h) }

// Swap đổi chỗ 2 phần tử trong heap
func (h BalanceHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

// Less định nghĩa thứ tự ưu tiên của MAX-HEAP:
// 1. Ưu tiên số tiền (Amount) lớn hơn đứng trước (Max-Heap)
// 2. Nếu số tiền bằng nhau, ưu tiên MaTV nhỏ hơn theo thứ tự từ điển (Deterministic tie-breaker)
func (h BalanceHeap) Less(i, j int) bool {
	if h[i].Amount != h[j].Amount {
		return h[i].Amount > h[j].Amount // Số tiền lớn hơn có độ ưu tiên cao hơn
	}
	return h[i].MaTV < h[j].MaTV // MaTV nhỏ hơn có độ ưu tiên cao hơn
}

// Push thêm phần tử vào heap
func (h *BalanceHeap) Push(x interface{}) {
	item := x.(*HeapItem)
	*h = append(*h, item)
}

// Pop lấy phần tử ở cuối mảng sau khi heap đã điều chỉnh vị trí gốc
func (h *BalanceHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // Xóa con trỏ để GC thu hồi bộ nhớ
	*h = old[0 : n-1]
	return item
}

// Định nghĩa 2 type alias để phân biệt rõ ràng creditorHeap và debtorHeap theo yêu cầu
type creditorHeap = BalanceHeap
type debtorHeap = BalanceHeap

// SimplifyDebts thực hiện thuật toán tối ưu hóa nợ (Greedy + 2 Max-Heaps)
// GIỮ NGUYÊN CHÍNH XÁC LOGIC THUẬT TOÁN GỐC:
// 1. Kiểm tra tổng tất cả balance phải = 0, sai thì trả lỗi
// 2. Tạo 2 MAX-HEAP: creditors (chứa balance > 0) và debtors (chứa balance < 0 với giá trị tuyệt đối)
// 3. Trong Less(i, j): ưu tiên amount lớn hơn; nếu bằng thì ưu tiên MaTV nhỏ hơn
// 4. Vòng lặp khi cả 2 heap không rỗng:
//    a. creditor := heap.Pop(creditors), debtor := heap.Pop(debtors)
//    b. amount := min(creditor.Amount, debtor.Amount)
//    c. Nếu amount <= 0 thì break
//    d. Tạo Settlement{FromUserID: debtor.MaTV, ToUserID: creditor.MaTV, AmountMinor: amount, Status: PENDING}
//    e. Trừ amount khỏi creditor.Amount và debtor.Amount
//    f. Nếu còn dư (> 0) thì heap.Push lại vào đúng heap của nó
func SimplifyDebts(balances map[string]int64) ([]settlement.Settlement, error) {
	// Bước 1: Kiểm tra tổng tất cả số dư phải đúng bằng 0
	var sumBalance int64 = 0
	for _, b := range balances {
		sumBalance += b
	}

	if sumBalance != 0 {
		return nil, fmt.Errorf("tổng số dư ròng của tất cả thành viên phải bằng 0 (tổng hiện tại: %d VND)", sumBalance)
	}

	// Bước 2: Đưa các thành viên vào 2 MAX-HEAP tương ứng
	creditors := &creditorHeap{}
	debtors := &debtorHeap{}

	heap.Init(creditors)
	heap.Init(debtors)

	for maTV, bal := range balances {
		if bal > 0 {
			// Chủ nợ: giữ nguyên giá trị dương
			heap.Push(creditors, &HeapItem{MaTV: maTV, Amount: bal})
		} else if bal < 0 {
			// Con nợ: lấy giá trị tuyệt đối (chuyển âm sang dương)
			heap.Push(debtors, &HeapItem{MaTV: maTV, Amount: -bal})
		}
	}

	settlements := make([]settlement.Settlement, 0)
	settlementCount := 1
	now := time.Now()

	// Bước 3: Vòng lặp khớp nợ giữa phần tử lớn nhất của Creditors và Debtors
	for creditors.Len() > 0 && debtors.Len() > 0 {
		// a. Pop phần tử có ưu tiên cao nhất từ cả 2 max-heap
		creditor := heap.Pop(creditors).(*HeapItem)
		debtor := heap.Pop(debtors).(*HeapItem)

		// b. Lấy số tiền giao dịch nhỏ nhất giữa 2 khoản
		amount := creditor.Amount
		if debtor.Amount < amount {
			amount = debtor.Amount
		}

		// c. Nếu số tiền <= 0 thì dừng vòng lặp
		if amount <= 0 {
			break
		}

		// d. Tạo giao dịch Settlement với trạng thái PENDING
		sID := fmt.Sprintf("STL-%03d", settlementCount)
		settlementCount++

		st := settlement.Settlement{
			ID:          sID,
			FromUserID:  debtor.MaTV,   // Con nợ trả tiền
			ToUserID:    creditor.MaTV, // Chủ nợ nhận tiền
			AmountMinor: amount,
			Status:      settlement.StatusPending,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		settlements = append(settlements, st)

		// e. Trừ số tiền giao dịch khỏi cả 2 phía
		creditor.Amount -= amount
		debtor.Amount -= amount

		// f. Nếu còn dư (> 0) thì push lại vào đúng heap của nó
		if creditor.Amount > 0 {
			heap.Push(creditors, creditor)
		}
		if debtor.Amount > 0 {
			heap.Push(debtors, debtor)
		}
	}

	return settlements, nil
}
