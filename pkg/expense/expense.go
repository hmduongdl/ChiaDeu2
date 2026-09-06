package expense

import (
	"fmt"
	"time"
)

// SplitType kiểu chia tiền
type SplitType string

const (
	SplitEqual   SplitType = "EQUAL"   // Chia đều
	SplitPercent SplitType = "PERCENT" // Chia theo tỷ lệ phần trăm
	SplitWeight  SplitType = "WEIGHT"  // Chia theo trọng số phần
	SplitCustom  SplitType = "CUSTOM"  // Chia theo số tiền tùy ý
)

// ExpenseStatus trạng thái khoản chi
type ExpenseStatus string

const (
	StatusActive ExpenseStatus = "ACTIVE" // Khoản chi có hiệu lực
	StatusVoided ExpenseStatus = "VOIDED" // Khoản chi đã hủy
)

// Expense struct đại diện cho một khoản chi
type Expense struct {
	ID          string        `json:"id"`
	PaidBy      string        `json:"paid_by"`      // MaTV của người ứng tiền
	AmountMinor int64         `json:"amount_minor"` // Tổng số tiền (đơn vị minor/VNĐ)
	SplitType   SplitType     `json:"split_type"`   // EQUAL, PERCENT, WEIGHT, CUSTOM
	Status      ExpenseStatus `json:"status"`       // ACTIVE hoặc VOIDED
	Description string        `json:"description"`  // Mô tả khoản chi
	CreatedAt   time.Time     `json:"created_at"`   // Thời gian tạo
}

func (e Expense) String() string {
	return fmt.Sprintf("ID: %-8s | Người trả: %-8s | Tiền: %10d VND | Loại: %-7s | Trạng thái: %-6s | Mô tả: %s",
		e.ID, e.PaidBy, e.AmountMinor, e.SplitType, e.Status, e.Description)
}
