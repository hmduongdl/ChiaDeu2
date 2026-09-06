package expense

import "fmt"

// ExpenseSplit đại diện cho số tiền chia cho 1 thành viên trong một khoản chi
type ExpenseSplit struct {
	ExpenseID  string `json:"expense_id"`  // Mã khoản chi liên kết
	UserID     string `json:"user_id"`     // MaTV thành viên chịu tiền
	ShareMinor int64  `json:"share_minor"` // Số tiền phải chịu (minor/VNĐ)
}

func (s ExpenseSplit) String() string {
	return fmt.Sprintf("ExpenseID: %s | MaTV: %s | Share: %d VND", s.ExpenseID, s.UserID, s.ShareMinor)
}
