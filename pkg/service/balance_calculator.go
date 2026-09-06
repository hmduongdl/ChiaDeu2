package service

import (
	"fmt"

	"chiadeu/pkg/expense"
)

// CalculateNetBalances tính toán số dư ròng của từng thành viên dựa trên danh sách khoản chi và các phần chia
// GIỮ NGUYÊN CHÍNH XÁC LOGIC THUẬT TOÁN GỐC:
// 1. Bỏ qua expense có Status == VOIDED
// 2. Với mỗi expense: kiểm tra tổng splits của nó == AmountMinor, nếu sai -> trả lỗi
// 3. balances[expense.PaidBy] += expense.AmountMinor
// 4. balances[split.UserID] -= split.ShareMinor
func CalculateNetBalances(expenses []expense.Expense, splits []expense.ExpenseSplit) (map[string]int64, error) {
	balances := make(map[string]int64)

	// Gom nhóm splits theo ExpenseID để kiểm tra nhanh
	splitMap := make(map[string][]expense.ExpenseSplit)
	for _, s := range splits {
		splitMap[s.ExpenseID] = append(splitMap[s.ExpenseID], s)
	}

	for _, exp := range expenses {
		// Bỏ qua expense bị VOIDED
		if exp.Status == expense.StatusVoided {
			continue
		}

		expSplits, exists := splitMap[exp.ID]
		if !exists || len(expSplits) == 0 {
			return nil, fmt.Errorf("khoản chi '%s' không có thông tin chi tiết các phần chia (splits)", exp.ID)
		}

		var totalSplitAmount int64 = 0
		for _, s := range expSplits {
			totalSplitAmount += s.ShareMinor
		}

		// Validation: Tổng các splits của khoản chi phải đúng bằng AmountMinor
		if totalSplitAmount != exp.AmountMinor {
			return nil, fmt.Errorf("lỗi dữ liệu khoản chi '%s': tổng các phần chia (%d VND) khác AmountMinor (%d VND)",
				exp.ID, totalSplitAmount, exp.AmountMinor)
		}

		// Tăng số dư ròng cho người trả tiền
		balances[exp.PaidBy] += exp.AmountMinor

		// Trừ số dư ròng cho các người chịu tiền
		for _, s := range expSplits {
			balances[s.UserID] -= s.ShareMinor
		}
	}

	return balances, nil
}
