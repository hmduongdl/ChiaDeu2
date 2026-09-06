package main

import (
	"fmt"
	"testing"
	"time"

	"chiadeu/pkg/expense"
	"chiadeu/pkg/member"
	"chiadeu/pkg/service"
)

// Example_demoScenario minh họa kịch bản 3 thành viên An, Bình, Châu và in các phép duyệt BST (NLR, LNR, LRN)
func Example_demoScenario() {
	// 1. Tạo BST thành viên
	bst := member.NewBST()
	_ = bst.Insert("TV01", "An")
	_ = bst.Insert("TV02", "Bình")
	_ = bst.Insert("TV03", "Châu")

	// In các kết quả duyệt cây BST
	fmt.Println("--- DUYỆT BST THÀNH VIÊN ---")

	fmt.Print("1. Duyệt NLR (Tiền tự): ")
	nlr := bst.TraverseNLR()
	for i, n := range nlr {
		if i > 0 {
			fmt.Print(" -> ")
		}
		fmt.Printf("%s (%s)", n.MaTV, n.HoTen)
	}
	fmt.Println()

	fmt.Print("2. Duyệt LNR (Trung tự - Tăng dần theo MaTV): ")
	lnr := bst.TraverseLNR()
	for i, n := range lnr {
		if i > 0 {
			fmt.Print(" -> ")
		}
		fmt.Printf("%s (%s)", n.MaTV, n.HoTen)
	}
	fmt.Println()

	fmt.Print("3. Duyệt LRN (Hậu tự): ")
	lrn := bst.TraverseLRN()
	for i, n := range lrn {
		if i > 0 {
			fmt.Print(" -> ")
		}
		fmt.Printf("%s (%s)", n.MaTV, n.HoTen)
	}
	fmt.Println()

	// 2. Tạo khoản chi tái hiện đề bài:
	// An ứng 300.000đ (chịu 100.000đ, Bình 120.000đ, Châu 80.000đ)
	exp1 := expense.Expense{
		ID:          "EXP-1",
		PaidBy:      "TV01",
		AmountMinor: 300000,
		SplitType:   expense.SplitCustom,
		Status:      expense.StatusActive,
		Description: "An ứng tiền ăn tối",
		CreatedAt:   time.Now(),
	}
	splits1 := []expense.ExpenseSplit{
		{ExpenseID: "EXP-1", UserID: "TV01", ShareMinor: 100000},
		{ExpenseID: "EXP-1", UserID: "TV02", ShareMinor: 120000},
		{ExpenseID: "EXP-1", UserID: "TV03", ShareMinor: 80000},
	}

	// Châu ứng 60.000đ (Châu chịu thêm 60.000đ cho bản thân)
	exp2 := expense.Expense{
		ID:          "EXP-2",
		PaidBy:      "TV03",
		AmountMinor: 60000,
		SplitType:   expense.SplitCustom,
		Status:      expense.StatusActive,
		Description: "Châu ứng mua đồ dùng",
		CreatedAt:   time.Now(),
	}
	splits2 := []expense.ExpenseSplit{
		{ExpenseID: "EXP-2", UserID: "TV03", ShareMinor: 60000},
	}

	expenses := []expense.Expense{exp1, exp2}
	splits := append(splits1, splits2...)

	// 3. Tính số dư ròng
	balances, err := service.CalculateNetBalances(expenses, splits)
	if err != nil {
		fmt.Printf("Lỗi tính balance: %v\n", err)
		return
	}

	// 4. Chạy SimplifyDebts
	settlements, err := service.SimplifyDebts(balances)
	if err != nil {
		fmt.Printf("Lỗi SimplifyDebts: %v\n", err)
		return
	}

	fmt.Println("\n--- KẾT QUẢ QUYẾT TOÁN SIMPLIFY DEBTS (2 MAX-HEAPS) ---")
	for _, s := range settlements {
		fromNode, _ := bst.Search(s.FromUserID)
		toNode, _ := bst.Search(s.ToUserID)
		fmt.Printf("[%s (%s)] trả [%s (%s)]: %d VND\n",
			fromNode.HoTen, s.FromUserID, toNode.HoTen, s.ToUserID, s.AmountMinor)
	}

	// Output:
	// --- DUYỆT BST THÀNH VIÊN ---
	// 1. Duyệt NLR (Tiền tự): TV01 (An) -> TV02 (Bình) -> TV03 (Châu)
	// 2. Duyệt LNR (Trung tự - Tăng dần theo MaTV): TV01 (An) -> TV02 (Bình) -> TV03 (Châu)
	// 3. Duyệt LRN (Hậu tự): TV03 (Châu) -> TV02 (Bình) -> TV01 (An)
	//
	// --- KẾT QUẢ QUYẾT TOÁN SIMPLIFY DEBTS (2 MAX-HEAPS) ---
	// [Bình (TV02)] trả [An (TV01)]: 120000 VND
	// [Châu (TV03)] trả [An (TV01)]: 80000 VND
}

func TestExampleScenarioRun(t *testing.T) {
	Example_demoScenario()
}
