package settlement

import (
	"fmt"
	"time"
)

// SettlementStatus trạng thái giao dịch thanh toán nợ
type SettlementStatus string

const (
	StatusPending              SettlementStatus = "PENDING"               // Mới tạo, con nợ chưa chuyển tiền
	StatusAwaitingConfirmation SettlementStatus = "AWAITING_CONFIRMATION" // Con nợ báo đã chuyển, chờ chủ nợ xác nhận
	StatusPaid                 SettlementStatus = "PAID"                  // Chủ nợ đã xác nhận thành công
	StatusCancelled            SettlementStatus = "CANCELLED"             // Giao dịch bị hủy
)

// Settlement đại diện cho 1 giao dịch trả nợ giữa 2 người
type Settlement struct {
	ID          string           `json:"id"`
	FromUserID  string           `json:"from_user_id"` // MaTV con nợ (người trả)
	ToUserID    string           `json:"to_user_id"`   // MaTV chủ nợ (người nhận)
	AmountMinor int64            `json:"amount_minor"` // Số tiền thanh toán (minor/VNĐ)
	Status      SettlementStatus `json:"status"`       // PENDING, AWAITING_CONFIRMATION, PAID, CANCELLED
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

func (s Settlement) String() string {
	return fmt.Sprintf("ID: %-8s | [%s] tra [%s]: %10d VND | Trang thai: %s",
		s.ID, s.FromUserID, s.ToUserID, s.AmountMinor, s.Status)
}
