package settlement

import (
	"fmt"
	"time"

	"chiadeu/pkg/expense"
)

// BatchStatus trạng thái kỳ quyết toán
type BatchStatus string

const (
	BatchOpen      BatchStatus = "OPEN"      // Kỳ quyết toán đang mở, tiếp tục nhận khoản chi
	BatchClosed    BatchStatus = "CLOSED"    // Kỳ quyết toán đã chốt
	BatchCancelled BatchStatus = "CANCELLED" // Kỳ quyết toán bị hủy
)

// SettlementBatch quản lý 1 kỳ chốt sổ quyết toán
type SettlementBatch struct {
	ID          string            `json:"id"`
	Status      BatchStatus       `json:"status"`
	Expenses    []expense.Expense `json:"expenses"`
	Splits      []expense.ExpenseSplit `json:"splits"`
	Settlements []Settlement      `json:"settlements"`
	CreatedAt   time.Time         `json:"created_at"`
	ClosedAt    *time.Time        `json:"closed_at,omitempty"`
}

func (b SettlementBatch) String() string {
	closedStr := "N/A"
	if b.ClosedAt != nil {
		closedStr = b.ClosedAt.Format("2006-01-02 15:04:05")
	}
	return fmt.Sprintf("Batch ID: %-8s | Trang thai: %-7s | So Khoan Chi: %d | So Giao Dich No: %d | Dong ky: %s",
		b.ID, b.Status, len(b.Expenses), len(b.Settlements), closedStr)
}
