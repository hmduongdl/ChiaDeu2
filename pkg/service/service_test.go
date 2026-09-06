package service

import (
	"testing"

	"chiadeu/pkg/expense"
	"chiadeu/pkg/member"
)

// TestSplitEqual kiểm tra thuật toán SplitEqual trong package service
func TestSplitEqual(t *testing.T) {
	memberIDs := []string{"TV03", "TV01", "TV02"}
	splits, err := SplitEqual("EXP-01", 100000, memberIDs)
	if err != nil {
		t.Fatalf("SplitEqual lỗi: %v", err)
	}

	if len(splits) != 3 {
		t.Fatalf("Mong muốn 3 splits, nhận được: %d", len(splits))
	}

	// Sắp xếp TV01, TV02, TV03 -> TV01 nhận 33334, TV02 33333, TV03 33333
	if splits[0].UserID != "TV01" || splits[0].ShareMinor != 33334 {
		t.Errorf("TV01 sai share: %v", splits[0])
	}
	if splits[1].UserID != "TV02" || splits[1].ShareMinor != 33333 {
		t.Errorf("TV02 sai share: %v", splits[1])
	}
	if splits[2].UserID != "TV03" || splits[2].ShareMinor != 33333 {
		t.Errorf("TV03 sai share: %v", splits[2])
	}
}

// TestCalculateNetBalances kiểm tra tính số dư ròng
func TestCalculateNetBalances(t *testing.T) {
	expenses := []expense.Expense{
		{ID: "E1", PaidBy: "TV01", AmountMinor: 300000, Status: expense.StatusActive},
		{ID: "E2", PaidBy: "TV02", AmountMinor: 100000, Status: expense.StatusVoided},
	}

	splits := []expense.ExpenseSplit{
		{ExpenseID: "E1", UserID: "TV01", ShareMinor: 100000},
		{ExpenseID: "E1", UserID: "TV02", ShareMinor: 200000},
		{ExpenseID: "E2", UserID: "TV02", ShareMinor: 100000},
	}

	balances, err := CalculateNetBalances(expenses, splits)
	if err != nil {
		t.Fatalf("CalculateNetBalances lỗi: %v", err)
	}

	if balances["TV01"] != 200000 {
		t.Errorf("TV01 balance sai: %d", balances["TV01"])
	}
	if balances["TV02"] != -200000 {
		t.Errorf("TV02 balance sai: %d", balances["TV02"])
	}
}

// TestSimplifyDebts kiểm tra thuật toán 2 Max-Heaps
func TestSimplifyDebts(t *testing.T) {
	balances := map[string]int64{
		"TV01": 200000,
		"TV02": -120000,
		"TV03": -80000,
	}

	settlements, err := SimplifyDebts(balances)
	if err != nil {
		t.Fatalf("SimplifyDebts lỗi: %v", err)
	}

	if len(settlements) != 2 {
		t.Fatalf("Mong muốn 2 settlements, nhận được %d", len(settlements))
	}

	if settlements[0].FromUserID != "TV02" || settlements[0].ToUserID != "TV01" || settlements[0].AmountMinor != 120000 {
		t.Errorf("Settlement 1 sai: %v", settlements[0])
	}
	if settlements[1].FromUserID != "TV03" || settlements[1].ToUserID != "TV01" || settlements[1].AmountMinor != 80000 {
		t.Errorf("Settlement 2 sai: %v", settlements[1])
	}
}

// TestValidateExpense kiểm tra validation khoản chi
func TestValidateExpense(t *testing.T) {
	bst := member.NewBST()
	_ = bst.Insert("TV01", "An")
	_ = bst.Insert("TV02", "Bình")

	splits := []expense.ExpenseSplit{
		{ExpenseID: "E1", UserID: "TV01", ShareMinor: 50000},
		{ExpenseID: "E1", UserID: "TV02", ShareMinor: 50000},
	}

	err := ValidateExpense("TV01", 100000, expense.SplitEqual, splits, bst)
	if err != nil {
		t.Errorf("ValidateExpense thất bại: %v", err)
	}
}
