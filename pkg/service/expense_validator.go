package service

import (
	"errors"
	"fmt"
	"strings"

	"chiadeu/pkg/expense"
	"chiadeu/pkg/member"
)

// ExpenseValidator cung cấp các phương thức kiểm tra tính hợp lệ của khoản chi
type ExpenseValidator struct{}

// ValidateExpense kiểm tra 5 điều kiện bắt buộc trước khi tạo hoặc chốt expense:
// 1. AmountMinor > 0
// 2. PaidBy không rỗng và phải tồn tại trong BST thành viên
// 3. SplitType hợp lệ (EQUAL, PERCENT, WEIGHT, CUSTOM)
// 4. Mọi ShareMinor >= 0 và UserID thuộc split phải tồn tại trong BST thành viên
// 5. Tổng tất cả ShareMinor == AmountMinor
func ValidateExpense(paidBy string, amountMinor int64, splitType expense.SplitType, splits []expense.ExpenseSplit, bst *member.BST) error {
	paidBy = strings.TrimSpace(paidBy)
	if paidBy == "" {
		return errors.New("mã người trả (PaidBy) không được để rỗng")
	}

	// Kiểm tra nguời trả có tồn tại trong BST không
	if bst != nil {
		if _, found := bst.Search(paidBy); !found {
			return fmt.Errorf("người trả money (MaTV: '%s') không tồn tại trong danh sách thành viên", paidBy)
		}
	}

	if amountMinor <= 0 {
		return fmt.Errorf("số tiền khoản chi phải > 0, giá trị hiện tại: %d", amountMinor)
	}

	switch splitType {
	case expense.SplitEqual, expense.SplitPercent, expense.SplitWeight, expense.SplitCustom:
		// Hợp lệ
	default:
		return fmt.Errorf("kiểu chia '%s' không hợp lệ (phải là EQUAL, PERCENT, WEIGHT, CUSTOM)", splitType)
	}

	if len(splits) == 0 {
		return errors.New("danh sách người tham gia chia tiền không được rỗng")
	}

	var sumShare int64 = 0
	userSeen := make(map[string]bool)

	for i, s := range splits {
		sUserID := strings.TrimSpace(s.UserID)
		if sUserID == "" {
			return fmt.Errorf("split thứ %d: MaTV không được để rỗng", i+1)
		}

		if userSeen[sUserID] {
			return fmt.Errorf("thành viên '%s' bị lặp lại nhiều lần trong cùng khoản chi", sUserID)
		}
		userSeen[sUserID] = true

		if bst != nil {
			if _, found := bst.Search(sUserID); !found {
				return fmt.Errorf("thành viên chia tiền (MaTV: '%s') không tồn tại trong danh sách thành viên", sUserID)
			}
		}

		if s.ShareMinor < 0 {
			return fmt.Errorf("số tiền chia của thành viên '%s' không được âm (%d)", sUserID, s.ShareMinor)
		}

		sumShare += s.ShareMinor
	}

	if sumShare != amountMinor {
		return fmt.Errorf("tổng các phần chia (%d VND) không bằng tổng số tiền khoản chi (%d VND)", sumShare, amountMinor)
	}

	return nil
}
