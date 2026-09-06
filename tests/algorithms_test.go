package tests

import (
	"testing"

	"chiadeu/pkg/expense"
	"chiadeu/pkg/member"
	"chiadeu/pkg/service"
)

// TestSplitEqual kiểm tra thuật toán SplitEqual
func TestSplitEqual(t *testing.T) {
	// Case 1: 100k chia cho 3 người ("U3", "U1", "U2") -> base = 33333, remainder = 1
	// Sắp xếp: "U1", "U2", "U3" -> U1 nhận +1 (33334), U2 (33333), U3 (33333)
	memberIDs := []string{"U3", "U1", "U2"}
	splits, err := service.SplitEqual("EXP-1", 100000, memberIDs)
	if err != nil {
		t.Fatalf("SplitEqual thất bại: %v", err)
	}

	if len(splits) != 3 {
		t.Fatalf("Mong muốn 3 splits, nhận được: %d", len(splits))
	}

	if splits[0].UserID != "U1" || splits[0].ShareMinor != 33334 {
		t.Errorf("U1 sai share: %v", splits[0])
	}
	if splits[1].UserID != "U2" || splits[1].ShareMinor != 33333 {
		t.Errorf("U2 sai share: %v", splits[1])
	}
	if splits[2].UserID != "U3" || splits[2].ShareMinor != 33333 {
		t.Errorf("U3 sai share: %v", splits[2])
	}

	// Case 2: Lỗi khi amount <= 0
	_, err = service.SplitEqual("EXP-2", 0, memberIDs)
	if err == nil {
		t.Error("Mong muốn lỗi khi amountMinor <= 0 nhưng không có lỗi")
	}

	// Case 3: Lỗi khi danh sách rỗng
	_, err = service.SplitEqual("EXP-3", 100, []string{})
	if err == nil {
		t.Error("Mong muốn lỗi khi memberIDs rỗng nhưng không có lỗi")
	}
}

// TestCalculateNetBalances kiểm tra thuật toán CalculateNetBalances
func TestCalculateNetBalances(t *testing.T) {
	expenses := []expense.Expense{
		{ID: "E1", PaidBy: "U1", AmountMinor: 300000, Status: expense.StatusActive},
		{ID: "E2", PaidBy: "U2", AmountMinor: 100000, Status: expense.StatusVoided}, // Bỏ qua
	}

	splits := []expense.ExpenseSplit{
		{ExpenseID: "E1", UserID: "U1", ShareMinor: 100000},
		{ExpenseID: "E1", UserID: "U2", ShareMinor: 200000},
		{ExpenseID: "E2", UserID: "U2", ShareMinor: 100000}, // Thuộc E2 bị VOIDED
	}

	balances, err := service.CalculateNetBalances(expenses, splits)
	if err != nil {
		t.Fatalf("CalculateNetBalances thất bại: %v", err)
	}

	// U1: ứng 300k - chịu 100k = +200k
	// U2: ứng 0k - chịu 200k = -200k
	if balances["U1"] != 200000 {
		t.Errorf("U1 balance sai: %d (kỳ vọng 200000)", balances["U1"])
	}
	if balances["U2"] != -200000 {
		t.Errorf("U2 balance sai: %d (kỳ vọng -200000)", balances["U2"])
	}
}

// TestSimplifyDebtsDemoScenario kiểm tra thuật toán SimplifyDebts với bài toán mẫu:
// An (TV01): ứng 300k, chịu 100k -> net = +200k
// Bình (TV02): ứng 0k, chịu 120k -> net = -120k
// Châu (TV03): ứng 60k, chịu 140k -> net = -80k
// Kết quả kỳ vọng: Bình trả An 120k, Châu trả An 80k.
func TestSimplifyDebtsDemoScenario(t *testing.T) {
	balances := map[string]int64{
		"TV01": 200000,
		"TV02": -120000,
		"TV03": -80000,
	}

	settlements, err := service.SimplifyDebts(balances)
	if err != nil {
		t.Fatalf("SimplifyDebts thất bại: %v", err)
	}

	if len(settlements) != 2 {
		t.Fatalf("Mong muốn 2 giao dịch settlement, nhận được: %d", len(settlements))
	}

	// Giao dịch 1: Bình (TV02) -> An (TV01): 120,000 VND
	if settlements[0].FromUserID != "TV02" || settlements[0].ToUserID != "TV01" || settlements[0].AmountMinor != 120000 {
		t.Errorf("Giao dịch 1 sai: %v (Kỳ vọng TV02 -> TV01: 120000)", settlements[0])
	}

	// Giao dịch 2: Châu (TV03) -> An (TV01): 80,000 VND
	if settlements[1].FromUserID != "TV03" || settlements[1].ToUserID != "TV01" || settlements[1].AmountMinor != 80000 {
		t.Errorf("Giao dịch 2 sai: %v (Kỳ vọng TV03 -> TV01: 80000)", settlements[1])
	}
}

// TestBSTOperations kiểm tra các chức năng của BST
func TestBSTOperations(t *testing.T) {
	bst := member.NewBST()

	// 1. Test Insert & Search
	_ = bst.Insert("TV03", "Châu")
	_ = bst.Insert("TV01", "An")
	_ = bst.Insert("TV04", "Dũng")
	_ = bst.Insert("TV02", "Bình")
	_ = bst.Insert("TV05", "Giang")

	if bst.CountNodes() != 5 {
		t.Errorf("Tổng số nút sai: %d (kỳ vọng 5)", bst.CountNodes())
	}

	// 2. Test LNR Traversal (phải tăng dần theo MaTV: TV01, TV02, TV03, TV04, TV05)
	lnr := bst.TraverseLNR()
	expectedOrder := []string{"TV01", "TV02", "TV03", "TV04", "TV05"}
	for i, node := range lnr {
		if node.MaTV != expectedOrder[i] {
			t.Errorf("LNR sai vị trí %d: có '%s', kỳ vọng '%s'", i, node.MaTV, expectedOrder[i])
		}
	}

	// 3. Test Delete Node Case 1: Leaf node ("TV02")
	err := bst.Remove("TV02")
	if err != nil {
		t.Errorf("Xóa node lá TV02 thất bại: %v", err)
	}
	if _, found := bst.Search("TV02"); found {
		t.Error("Vẫn tìm thấy TV02 sau khi xóa")
	}

	// 4. Test Delete Node Case 3: Node 2 children ("TV03" là root có con trái TV01, con phải TV04)
	err = bst.Remove("TV03")
	if err != nil {
		t.Errorf("Xóa node 2 con TV03 thất bại: %v", err)
	}
	if _, found := bst.Search("TV03"); found {
		t.Error("Vẫn tìm thấy TV03 sau khi xóa")
	}

	// Kiểm tra số lượng nút sau khi xóa 2 nút (còn 3 nút: TV01, TV04, TV05)
	if bst.CountNodes() != 3 {
		t.Errorf("Tổng số nút sau xóa sai: %d (kỳ vọng 3)", bst.CountNodes())
	}
}
