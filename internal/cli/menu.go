package cli

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"chiadeu/pkg/expense"
	"chiadeu/pkg/member"
	"chiadeu/pkg/service"
	"chiadeu/pkg/settlement"
)

type CLIApp struct {
	bst               *member.BST
	expenses          []expense.Expense
	splits            []expense.ExpenseSplit
	batches           []settlement.SettlementBatch
	activeSettlements []settlement.Settlement
	reader            *bufio.Reader
}

func NewCLIApp() *CLIApp {
	return &CLIApp{
		bst:               member.NewBST(),
		expenses:          make([]expense.Expense, 0),
		splits:            make([]expense.ExpenseSplit, 0),
		batches:           make([]settlement.SettlementBatch, 0),
		activeSettlements: make([]settlement.Settlement, 0),
		reader:            bufio.NewReader(os.Stdin),
	}
}

func (app *CLIApp) Run() {
	for {
		app.printMainMenu()
		choice := app.readString(">> Chọn nhóm chức năng (1-5) hoặc chức năng con (VD: 1.1, 4.3): ")
		switch choice {
		// Nhóm 1: Nhập / Đọc danh sách thành viên
		case "1":
			app.handleGroup1()
		case "1.1":
			app.inputMembersFromKeyboard()
		case "1.2":
			app.readFromFile()
		case "1.3":
			app.addOneMember()

		// Nhóm 2: Xuất danh sách thành viên
		case "2":
			app.handleGroup2()
		case "2.1":
			app.traverseNLR()
		case "2.2":
			app.traverseLNR()
		case "2.3":
			app.traverseLRN()

		// Nhóm 3: Sắp xếp / Chỉnh sửa / Tìm kiếm / Hủy
		case "3":
			app.handleGroup3()
		case "3.1":
			app.customSortMembers()
		case "3.2":
			app.searchMember()
		case "3.3":
			app.updateMember()
		case "3.4":
			app.deleteMember()

		// Nhóm 4: Quản lý chi tiêu & Thống kê
		case "4":
			app.handleGroup4()
		case "4.1":
			app.recordExpense()
		case "4.2":
			app.viewBalances()
		case "4.3":
			app.finalizeSettlementBatch()
		case "4.4":
			app.showStatistics()

		// Nhóm 5: Chức năng khác
		case "5":
			app.handleGroup5()
		case "5.1":
			app.saveToFile()
		case "5.2":
			fmt.Println("\n=======================================================")
			fmt.Println("  Cảm ơn bạn đã sử dụng ứng dụng CHIA ĐỀU! Tạm biệt!")
			fmt.Println("=======================================================")
			return

		// Chức năng Demo bổ sung
		case "demo", "DEMO":
			app.handleRunDemoScenario()

		default:
			fmt.Println("\n[!] Lựa chọn không hợp lệ. Vui lòng chọn lại (VD: 1.1, 2.2, 4.3, 5.2...).")
		}
	}
}

func (app *CLIApp) printMainMenu() {
	fmt.Println("\n==========================================================================")
	fmt.Println("         QUẢN LÝ THÀNH VIÊN & CHI TIÊU NHÓM (BST & MAX-HEAP)              ")
	fmt.Println("==========================================================================")
	fmt.Println("1. Nhập / Đọc danh sách thành viên")
	fmt.Println("   1.1 Nhập từ bàn phím")
	fmt.Println("   1.2 Đọc từ file JSON")
	fmt.Println("   1.3 Thêm 1 thành viên (insert vào BST theo MaTV)")
	fmt.Println("2. Xuất danh sách thành viên")
	fmt.Println("   2.1 Duyệt NLR (Node-Left-Right / tiền tự)")
	fmt.Println("   2.2 Duyệt LNR (Left-Node-Right / trung tự — tăng dần theo MaTV)")
	fmt.Println("   2.3 Duyệt LRN (Left-Right-Node / hậu tự)")
	fmt.Println("3. Sắp xếp / Chỉnh sửa / Tìm kiếm / Hủy")
	fmt.Println("   3.1 Sắp xếp danh sách (theo tên, theo tổng đã chi, theo số dư...)")
	fmt.Println("   3.2 Tìm kiếm thành viên theo MaTV (dùng đúng cơ chế BST.Search)")
	fmt.Println("   3.3 Sửa thông tin thành viên (BST.Update, không đổi khóa)")
	fmt.Println("   3.4 Hủy thành viên (BST.Remove — xử lý đủ 3 trường hợp)")
	fmt.Println("4. Quản lý chi tiêu & Thống kê")
	fmt.Println("   4.1 Ghi khoản chi (người trả & người chia đều kiểm tra qua BST.Search)")
	fmt.Println("   4.2 Xem số dư ròng từng thành viên (CalculateNetBalances)")
	fmt.Println("   4.3 Chốt kỳ quyết toán (SimplifyDebts bằng 2 max-heap)")
	fmt.Println("   4.4 Thống kê: tổng chi, người chi nhiều nhất, số người nợ/nhận")
	fmt.Println("5. Chức năng khác")
	fmt.Println("   5.1 Lưu danh sách ra file JSON")
	fmt.Println("   5.2 Thoát chương trình")
	fmt.Println("--------------------------------------------------------------------------")
	fmt.Println(" (Mẹo: gõ 'demo' để chạy nhanh kịch bản An, Bình, Châu)")
	fmt.Println("--------------------------------------------------------------------------")
}

// -----------------------------------------------------------------------------
// NHÓM 1: NHẬP / ĐỌC DANH SÁCH THÀNH VIÊN
// -----------------------------------------------------------------------------
func (app *CLIApp) handleGroup1() {
	fmt.Println("\n--- [1] NHẬP / ĐỌC DANH SÁCH THÀNH VIÊN ---")
	fmt.Println("1.1 Nhập từ bàn phím")
	fmt.Println("1.2 Đọc từ file JSON")
	fmt.Println("1.3 Thêm 1 thành viên (insert vào BST theo MaTV)")
	sub := app.readString(">> Chọn (1.1 - 1.3): ")
	switch sub {
	case "1.1", "1":
		app.inputMembersFromKeyboard()
	case "1.2", "2":
		app.readFromFile()
	case "1.3", "3":
		app.addOneMember()
	default:
		fmt.Println("[!] Lựa chọn không hợp lệ.")
	}
}

func (app *CLIApp) inputMembersFromKeyboard() {
	countStr := app.readString(">> Nhập số lượng thành viên muốn thêm: ")
	count, err := strconv.Atoi(countStr)
	if err != nil || count <= 0 {
		fmt.Println("[!] Số lượng không hợp lệ.")
		return
	}

	successCount := 0
	for i := 1; i <= count; i++ {
		fmt.Printf("\n--- Thành viên %d/%d ---\n", i, count)
		maTV := app.readString("Mã thành viên (MaTV): ")
		hoTen := app.readString("Họ và tên: ")

		err := app.bst.Insert(maTV, hoTen)
		if err != nil {
			fmt.Printf("[X] Lỗi thêm thành viên: %v\n", err)
		} else {
			fmt.Printf("[V] Đã chèn nút '%s' (%s) vào BST thành công!\n", maTV, hoTen)
			successCount++
		}
	}
	fmt.Printf("\n[+] Đã thêm thành công %d/%d thành viên vào BST.\n", successCount, count)
}

func (app *CLIApp) readFromFile() {
	filePath := app.readString("Nhập đường dẫn file đọc (mặc định: data/members.txt hoặc data/chiadeu.json): ")
	filePath = strings.TrimSpace(filePath)
	if filePath == "" {
		filePath = "data/members.txt"
	}

	if strings.HasSuffix(strings.ToLower(filePath), ".txt") {
		members, err := service.ImportMembersFromTxt(filePath)
		if err != nil {
			fmt.Printf("[X] Lỗi đọc file TXT: %v\n", err)
			return
		}
		successCount := 0
		for _, m := range members {
			err := app.bst.Insert(m.MaTV, m.HoTen)
			if err == nil {
				successCount++
			}
		}
		fmt.Printf("[V] Đã đọc file TXT '%s' thành công! Chèn thành công %d/%d nút vào cây BST thành viên.\n",
			filePath, successCount, len(members))
	} else {
		bst, exp, splits, batches, _, settlements, err := service.LoadFromFile(filePath)
		if err != nil {
			fmt.Printf("[X] Lỗi đọc file JSON: %v\n", err)
		} else {
			app.bst = bst
			app.expenses = exp
			app.splits = splits
			app.batches = batches
			app.activeSettlements = settlements
			fmt.Printf("[V] Đã đọc dữ liệu thành công từ file JSON '%s'! Nạp được %d nút BST thành viên.\n", filePath, app.bst.CountNodes())
		}
	}
}

func (app *CLIApp) addOneMember() {
	fmt.Println("\n--- THÊM 1 THÀNH VIÊN VÀO BST (1.3) ---")
	maTV := app.readString("Mã thành viên (MaTV): ")
	hoTen := app.readString("Họ và tên: ")

	err := app.bst.Insert(maTV, hoTen)
	if err != nil {
		fmt.Printf("[X] Lỗi: %v\n", err)
	} else {
		fmt.Printf("[V] Đã chèn nút '%s' (%s) vào cây BST thành công!\n", maTV, hoTen)
	}
}

// -----------------------------------------------------------------------------
// NHÓM 2: XUẤT DANH SÁCH THÀNH VIÊN (TRAVERSALS)
// -----------------------------------------------------------------------------
func (app *CLIApp) handleGroup2() {
	fmt.Println("\n--- [2] XUẤT DANH SÁCH THÀNH VIÊN ---")
	fmt.Println("2.1 Duyệt NLR (Node-Left-Right / tiền tự)")
	fmt.Println("2.2 Duyệt LNR (Left-Node-Right / trung tự — tăng dần theo MaTV)")
	fmt.Println("2.3 Duyệt LRN (Left-Right-Node / hậu tự)")
	sub := app.readString(">> Chọn (2.1 - 2.3): ")
	switch sub {
	case "2.1", "1":
		app.traverseNLR()
	case "2.2", "2":
		app.traverseLNR()
	case "2.3", "3":
		app.traverseLRN()
	default:
		fmt.Println("[!] Lựa chọn không hợp lệ.")
	}
}

func (app *CLIApp) traverseNLR() {
	fmt.Println("\n--- DUYỆT NLR (TIỀN TỰ) ---")
	nodes := app.bst.TraverseNLR()
	app.printNodeList(nodes)
}

func (app *CLIApp) traverseLNR() {
	fmt.Println("\n--- DUYỆT LNR (TRUNG TỰ - TĂNG DẦN THEO MATV) ---")
	nodes := app.bst.TraverseLNR()
	app.printNodeList(nodes)
}

func (app *CLIApp) traverseLRN() {
	fmt.Println("\n--- DUYỆT LRN (HẬU TỰ) ---")
	nodes := app.bst.TraverseLRN()
	app.printNodeList(nodes)
}

func (app *CLIApp) printNodeList(nodes []*member.Node) {
	if len(nodes) == 0 {
		fmt.Println(" (Danh sách BST rỗng)")
		return
	}
	fmt.Println("--------------------------------------------------------------------------------")
	for i, n := range nodes {
		fmt.Printf("%2d. %s\n", i+1, n.String())
	}
	fmt.Println("--------------------------------------------------------------------------------")
}

// -----------------------------------------------------------------------------
// NHÓM 3: SẮP XẾP / CHỈNH SỬA / TÌM KIẾM / HỦY
// -----------------------------------------------------------------------------
func (app *CLIApp) handleGroup3() {
	fmt.Println("\n--- [3] SẮP XẾP / CHỈNH SỬA / TÌM KIẾM / HỦY ---")
	fmt.Println("3.1 Sắp xếp danh sách (theo tên, theo tổng đã chi, theo số dư...)")
	fmt.Println("3.2 Tìm kiếm thành viên theo MaTV (dùng đúng cơ chế BST.Search)")
	fmt.Println("3.3 Sửa thông tin thành viên (BST.Update, không đổi khóa)")
	fmt.Println("3.4 Hủy thành viên (BST.Remove — xử lý đủ 3 trường hợp)")
	sub := app.readString(">> Chọn (3.1 - 3.4): ")
	switch sub {
	case "3.1", "1":
		app.customSortMembers()
	case "3.2", "2":
		app.searchMember()
	case "3.3", "3":
		app.updateMember()
	case "3.4", "4":
		app.deleteMember()
	default:
		fmt.Println("[!] Lựa chọn không hợp lệ.")
	}
}

func (app *CLIApp) customSortMembers() {
	nodes := app.bst.TraverseLNR()
	if len(nodes) == 0 {
		fmt.Println("[!] Cây BST rỗng.")
		return
	}

	fmt.Println("\n--- SẮP XẾP DANH SÁCH THÀNH VIÊN TÙY CHỌN (3.1) ---")
	fmt.Println("a) Theo Họ Tên (A-Z)")
	fmt.Println("b) Theo Tổng đã ứng tiền (Giảm dần)")
	fmt.Println("c) Theo Số dư ròng (Chủ nợ nhiều nhất -> Con nợ)")

	sub := app.readString(">> Chọn (a/b/c): ")
	switch strings.ToLower(sub) {
	case "a", "1":
		sort.Slice(nodes, func(i, j int) bool {
			return nodes[i].HoTen < nodes[j].HoTen
		})
		fmt.Println("\n--- DANH SÁCH THEO HỌ TÊN (A-Z) ---")
	case "b", "2":
		sort.Slice(nodes, func(i, j int) bool {
			return nodes[i].TongDaUng > nodes[j].TongDaUng
		})
		fmt.Println("\n--- DANH SÁCH THEO TỔNG ĐÃ ỨNG (GIẢM DẦN) ---")
	case "c", "3":
		sort.Slice(nodes, func(i, j int) bool {
			return nodes[i].NetBalance() > nodes[j].NetBalance()
		})
		fmt.Println("\n--- DANH SÁCH THEO SỐ DƯ RÒNG (CHỦ NỢ -> CON NỢ) ---")
	default:
		fmt.Println("[!] Lựa chọn không hợp lệ.")
		return
	}
	app.printNodeList(nodes)
}

func (app *CLIApp) searchMember() {
	fmt.Println("\n--- TÌM KIẾM THÀNH VIÊN THEO MATV (BST.SEARCH - 3.2) ---")
	maTV := app.readString("Nhập MaTV cần tìm: ")

	node, found := app.bst.Search(maTV)
	if found {
		fmt.Println("\n[V] Đã tìm thấy nút thành viên trong BST:")
		fmt.Println("--------------------------------------------------------------------------------")
		fmt.Println(node.String())
		fmt.Println("--------------------------------------------------------------------------------")
	} else {
		fmt.Printf("[X] Không tìm thấy thành viên nào có MaTV '%s' trong BST.\n", maTV)
	}
}

func (app *CLIApp) updateMember() {
	fmt.Println("\n--- SỬA THÔNG TIN THÀNH VIÊN (BST.UPDATE - 3.3) ---")
	maTV := app.readString("Nhập MaTV của thành viên cần sửa: ")
	hoTenMoi := app.readString("Nhập Họ tên mới: ")

	err := app.bst.Update(maTV, hoTenMoi)
	if err != nil {
		fmt.Printf("[X] Lỗi: %v\n", err)
	} else {
		fmt.Printf("[V] Đã cập nhật họ tên mới cho thành viên '%s' thành công!\n", maTV)
	}
}

func (app *CLIApp) deleteMember() {
	fmt.Println("\n--- HỦY THÀNH VIÊN (BST.REMOVE - 3.4) ---")
	maTV := app.readString("Nhập MaTV của thành viên cần xóa: ")

	node, found := app.bst.Search(maTV)
	if !found {
		fmt.Printf("[X] Không tìm thấy thành viên có MaTV '%s' để xóa.\n", maTV)
		return
	}

	// Giải thích 3 trường hợp xóa node cho người dùng minh bạch
	if node.Left == nil && node.Right == nil {
		fmt.Println("[*] Phân tích BST: Node cần xóa là NODE LÁ (0 con).")
	} else if node.Left == nil || node.Right == nil {
		fmt.Println("[*] Phân tích BST: Node cần xóa có 1 CHILD (1 con).")
	} else {
		fmt.Println("[*] Phân tích BST: Node cần xóa có 2 CHILDREN (2 con) -> Dùng node thế mạng (nút nhỏ nhất cây con phải).")
	}

	err := app.bst.Remove(maTV)
	if err != nil {
		fmt.Printf("[X] Lỗi xóa: %v\n", err)
	} else {
		fmt.Printf("[V] Đã xóa thành công thành viên '%s' khỏi BST!\n", maTV)
	}
}

// -----------------------------------------------------------------------------
// NHÓM 4: QUẢN LÝ CHI TIÊU & THỐNG KÊ
// -----------------------------------------------------------------------------
func (app *CLIApp) handleGroup4() {
	fmt.Println("\n--- [4] QUẢN LÝ CHI TIÊU & THỐNG KÊ ---")
	fmt.Println("4.1 Ghi khoản chi (người trả & người chia đều phải kiểm tra trong BST)")
	fmt.Println("4.2 Xem số dư ròng từng thành viên (CalculateNetBalances)")
	fmt.Println("4.3 Chốt kỳ quyết toán (SimplifyDebts bằng 2 max-heap)")
	fmt.Println("4.4 Thống kê: tổng chi, người chi nhiều nhất, số người đang nợ/được nhận")
	sub := app.readString(">> Chọn (4.1 - 4.4): ")
	switch sub {
	case "4.1", "1":
		app.recordExpense()
	case "4.2", "2":
		app.viewBalances()
	case "4.3", "3":
		app.finalizeSettlementBatch()
	case "4.4", "4":
		app.showStatistics()
	default:
		fmt.Println("[!] Lựa chọn không hợp lệ.")
	}
}

func (app *CLIApp) recordExpense() {
	if app.bst.CountNodes() == 0 {
		fmt.Println("[!] Cây BST rỗng. Vui lòng thêm thành viên trước!")
		return
	}

	fmt.Println("\n--- GHI NHẬN KHOẢN CHI MỚI (4.1) ---")
	paidBy := app.readString("Mã người trả tiền (PaidBy - MaTV): ")
	payerNode, found := app.bst.Search(paidBy)
	if !found {
		fmt.Printf("[X] Không tìm thấy thành viên '%s' trong cây BST!\n", paidBy)
		return
	}

	desc := app.readString("Mô tả khoản chi: ")

	amountStr := app.readString("Tổng số tiền (VND): ")
	amount, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil || amount <= 0 {
		fmt.Println("[!] Số tiền phải là số nguyên > 0.")
		return
	}

	fmt.Println("\nChọn kiểu chia tiền:")
	fmt.Println("1. EQUAL   (Chia đều)")
	fmt.Println("2. PERCENT (Chia theo tỷ lệ phần trăm %)")
	fmt.Println("3. WEIGHT  (Chia theo trọng số phần)")
	fmt.Println("4. CUSTOM  (Chia theo số tiền nhập tùy ý)")

	splitChoice := app.readString(">> Chọn kiểu chia (1-4): ")
	expID := fmt.Sprintf("EXP-%03d", len(app.expenses)+1)

	var splitType expense.SplitType
	var splits []expense.ExpenseSplit

	switch splitChoice {
	case "1":
		splitType = expense.SplitEqual
		membersStr := app.readString("Nhập các MaTV chia tiền (cách nhau bởi dấu phẩy, hoặc nhấn Enter để chọn TẤT CẢ): ")
		var memberIDs []string
		if strings.TrimSpace(membersStr) == "" {
			nodes := app.bst.TraverseLNR()
			for _, n := range nodes {
				memberIDs = append(memberIDs, n.MaTV)
			}
		} else {
			parts := strings.Split(membersStr, ",")
			for _, p := range parts {
				p = strings.TrimSpace(p)
				if p != "" {
					memberIDs = append(memberIDs, p)
				}
			}
		}

		splits, err = service.SplitEqual(expID, amount, memberIDs)

	case "2":
		splitType = expense.SplitPercent
		nodes := app.bst.TraverseLNR()
		fmt.Println("\nNhập % cho từng thành viên (tổng phải = 100%):")
		pctMap := make(map[string]float64)
		for _, n := range nodes {
			pStr := app.readString(fmt.Sprintf(" Tỷ lệ %% của %s (%s): ", n.MaTV, n.HoTen))
			pVal, pErr := strconv.ParseFloat(pStr, 64)
			if pErr == nil && pVal > 0 {
				pctMap[n.MaTV] = pVal
			}
		}
		splits, err = service.SplitPercent(expID, amount, pctMap)

	case "3":
		splitType = expense.SplitWeight
		nodes := app.bst.TraverseLNR()
		fmt.Println("\nNhập trọng số phần cho từng thành viên:")
		wMap := make(map[string]int)
		for _, n := range nodes {
			wStr := app.readString(fmt.Sprintf(" Trọng số của %s (%s): ", n.MaTV, n.HoTen))
			wVal, wErr := strconv.Atoi(wStr)
			if wErr == nil && wVal > 0 {
				wMap[n.MaTV] = wVal
			}
		}
		splits, err = service.SplitWeight(expID, amount, wMap)

	case "4":
		splitType = expense.SplitCustom
		nodes := app.bst.TraverseLNR()
		fmt.Println("\nNhập số tiền chịu cho từng thành viên:")
		cMap := make(map[string]int64)
		for _, n := range nodes {
			cStr := app.readString(fmt.Sprintf(" Số tiền của %s (%s): ", n.MaTV, n.HoTen))
			cVal, cErr := strconv.ParseInt(cStr, 10, 64)
			if cErr == nil && cVal > 0 {
				cMap[n.MaTV] = cVal
			}
		}
		splits, err = service.SplitCustom(expID, cMap)

	default:
		fmt.Println("[!] Kiểu chia không hợp lệ.")
		return
	}

	if err != nil {
		fmt.Printf("[X] Lỗi tính chia tiền: %v\n", err)
		return
	}

	// Validate khoản chi và kiểm tra sự tồn tại trong BST
	valErr := service.ValidateExpense(paidBy, amount, splitType, splits, app.bst)
	if valErr != nil {
		fmt.Printf("[X] Lỗi Validate khoản chi: %v\n", valErr)
		return
	}

	exp := expense.Expense{
		ID:          expID,
		PaidBy:      paidBy,
		AmountMinor: amount,
		SplitType:   splitType,
		Status:      expense.StatusActive,
		Description: desc,
		CreatedAt:   time.Now(),
	}

	app.expenses = append(app.expenses, exp)
	app.splits = append(app.splits, splits...)

	payerNode.TongDaUng += amount
	for _, s := range splits {
		if sNode, found := app.bst.Search(s.UserID); found {
			sNode.TongPhaiChiu += s.ShareMinor
		}
	}

	fmt.Printf("\n[V] Đã ghi nhận khoản chi '%s' thành công!\n", expID)
}

func (app *CLIApp) viewBalances() {
	fmt.Println("\n--- SỐ DƯ RÒNG TỪNG THÀNH VIÊN (CalculateNetBalances - 4.2) ---")
	balances, err := service.CalculateNetBalances(app.expenses, app.splits)
	if err != nil {
		fmt.Printf("[X] Lỗi tính số dư: %v\n", err)
		return
	}

	nodes := app.bst.TraverseLNR()
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Printf("%-10s | %-20s | %-15s | %-15s | %-15s\n", "MaTV", "Họ tên", "Đã ứng (VND)", "Phải chịu (VND)", "Số dư ròng (VND)")
	fmt.Println("--------------------------------------------------------------------------------")

	for _, n := range nodes {
		bal := balances[n.MaTV]
		statusStr := ""
		if bal > 0 {
			statusStr = "(Chủ nợ)"
		} else if bal < 0 {
			statusStr = "(Con nợ)"
		} else {
			statusStr = "(Hòa vốn)"
		}

		fmt.Printf("%-10s | %-20s | %15d | %15d | %15d %s\n",
			n.MaTV, n.HoTen, n.TongDaUng, n.TongPhaiChiu, bal, statusStr)
	}
	fmt.Println("--------------------------------------------------------------------------------")
}

func (app *CLIApp) finalizeSettlementBatch() {
	fmt.Println("\n=======================================================")
	fmt.Println(" 4.3 CHỐT KỲ QUYẾT TOÁN - SIMPLIFY DEBTS (2 MAX-HEAPS)")
	fmt.Println("=======================================================")

	if len(app.expenses) == 0 {
		fmt.Println("[!] Không có khoản chi nào để quyết toán.")
		return
	}

	balances, err := service.CalculateNetBalances(app.expenses, app.splits)
	if err != nil {
		fmt.Printf("[X] Lỗi tính số dư: %v\n", err)
		return
	}

	settlements, err := service.SimplifyDebts(balances)
	if err != nil {
		fmt.Printf("[X] Lỗi SimplifyDebts: %v\n", err)
		return
	}

	app.activeSettlements = settlements
	batchID := fmt.Sprintf("BATCH-%03d", len(app.batches)+1)
	now := time.Now()

	batch := settlement.SettlementBatch{
		ID:          batchID,
		Status:      settlement.BatchClosed,
		Expenses:    app.expenses,
		Splits:      app.splits,
		Settlements: settlements,
		CreatedAt:   now,
		ClosedAt:    &now,
	}

	app.batches = append(app.batches, batch)

	fmt.Printf("\n[V] ĐÃ CHỐT KỲ QUYẾT TOÁN '%s' THÀNH CÔNG!\n", batchID)
	fmt.Println("\n--- KẾT QUẢ GIAO DỊCH QUYẾT TOÁN TỐI ƯU ---")
	for i, s := range settlements {
		fromNode, _ := app.bst.Search(s.FromUserID)
		toNode, _ := app.bst.Search(s.ToUserID)

		fromName := s.FromUserID
		if fromNode != nil {
			fromName = fmt.Sprintf("%s (%s)", fromNode.HoTen, s.FromUserID)
		}
		toName := s.ToUserID
		if toNode != nil {
			toName = fmt.Sprintf("%s (%s)", toNode.HoTen, s.ToUserID)
		}

		fmt.Printf("  %d. [%s]  ===>  trả  ===>  [%s] : %12d VND\n",
			i+1, fromName, toName, s.AmountMinor)
	}
}

func (app *CLIApp) showStatistics() {
	fmt.Println("\n=======================================================")
	fmt.Println("                 4.4 THỐNG KÊ CHI TIÊU                ")
	fmt.Println("=======================================================")

	var totalSpent int64 = 0
	for _, e := range app.expenses {
		if e.Status == expense.StatusActive {
			totalSpent += e.AmountMinor
		}
	}

	nodes := app.bst.TraverseLNR()
	var topSpender *member.Node
	debtorCount := 0
	creditorCount := 0

	for _, n := range nodes {
		if topSpender == nil || n.TongDaUng > topSpender.TongDaUng {
			topSpender = n
		}
		if n.NetBalance() < 0 {
			debtorCount++
		} else if n.NetBalance() > 0 {
			creditorCount++
		}
	}

	fmt.Printf("1. Tổng tiền nhóm đã chi: %d VND\n", totalSpent)
	if topSpender != nil {
		fmt.Printf("2. Người ứng tiền nhiều nhất: %s (%s) với %d VND\n",
			topSpender.HoTen, topSpender.MaTV, topSpender.TongDaUng)
	}
	fmt.Printf("3. Số người đang nợ: %d | Số người được nhận (chủ nợ): %d\n", debtorCount, creditorCount)
	fmt.Println("=======================================================")
}

// -----------------------------------------------------------------------------
// NHÓM 5: CHỨC NĂNG KHÁC
// -----------------------------------------------------------------------------
func (app *CLIApp) handleGroup5() {
	fmt.Println("\n--- [5] CHỨC NĂNG KHÁC ---")
	fmt.Println("5.1 Lưu danh sách ra file JSON")
	fmt.Println("5.2 Thoát chương trình")
	sub := app.readString(">> Chọn (5.1 - 5.2): ")
	switch sub {
	case "5.1", "1":
		app.saveToFile()
	case "5.2", "2":
		fmt.Println("\nTạm biệt!")
		os.Exit(0)
	default:
		fmt.Println("[!] Lựa chọn không hợp lệ.")
	}
}

func (app *CLIApp) saveToFile() {
	filePath := app.readString("Nhập đường dẫn file lưu (mặc định: data/chiadeu.json): ")
	if strings.TrimSpace(filePath) == "" {
		filePath = "data/chiadeu.json"
	}
	err := service.SaveToFile(filePath, app.bst, app.expenses, app.splits, app.batches, nil, app.activeSettlements)
	if err != nil {
		fmt.Printf("[X] Lỗi lưu file: %v\n", err)
	} else {
		fmt.Printf("[V] Đã lưu dữ liệu thành công ra file '%s'!\n", filePath)
	}
}

func (app *CLIApp) handleRunDemoScenario() {
	fmt.Println("\n==================================================================================")
	fmt.Println("      KỊCH BẢN DEMO TÁI HIỆN BÀI TOÁN MẪU (AN, BÌNH, CHÂU)")
	fmt.Println("==================================================================================")

	demoBST := member.NewBST()
	demoBST.Insert("TV01", "An")
	demoBST.Insert("TV02", "Bình")
	demoBST.Insert("TV03", "Châu")

	fmt.Println("\n--- DUYỆT BST THÀNH VIÊN DEMO ---")
	fmt.Print("1. Duyệt NLR (Tiền tự): ")
	for i, n := range demoBST.TraverseNLR() {
		if i > 0 {
			fmt.Print(" -> ")
		}
		fmt.Printf("%s (%s)", n.MaTV, n.HoTen)
	}
	fmt.Println()

	fmt.Print("2. Duyệt LNR (Trung tự): ")
	for i, n := range demoBST.TraverseLNR() {
		if i > 0 {
			fmt.Print(" -> ")
		}
		fmt.Printf("%s (%s)", n.MaTV, n.HoTen)
	}
	fmt.Println()

	fmt.Print("3. Duyệt LRN (Hậu tự): ")
	for i, n := range demoBST.TraverseLRN() {
		if i > 0 {
			fmt.Print(" -> ")
		}
		fmt.Printf("%s (%s)", n.MaTV, n.HoTen)
	}
	fmt.Println()

	exp1 := expense.Expense{
		ID:          "EXP-DEMO-1",
		PaidBy:      "TV01",
		AmountMinor: 300000,
		SplitType:   expense.SplitCustom,
		Status:      expense.StatusActive,
		Description: "An ứng ăn tối",
		CreatedAt:   time.Now(),
	}
	splits1 := []expense.ExpenseSplit{
		{ExpenseID: "EXP-DEMO-1", UserID: "TV01", ShareMinor: 100000},
		{ExpenseID: "EXP-DEMO-1", UserID: "TV02", ShareMinor: 120000},
		{ExpenseID: "EXP-DEMO-1", UserID: "TV03", ShareMinor: 80000},
	}

	exp2 := expense.Expense{
		ID:          "EXP-DEMO-2",
		PaidBy:      "TV03",
		AmountMinor: 60000,
		SplitType:   expense.SplitCustom,
		Status:      expense.StatusActive,
		Description: "Châu ứng mua đồ linh tinh",
		CreatedAt:   time.Now(),
	}
	splits2 := []expense.ExpenseSplit{
		{ExpenseID: "EXP-DEMO-2", UserID: "TV03", ShareMinor: 60000},
	}

	demoExpenses := []expense.Expense{exp1, exp2}
	demoSplits := append(splits1, splits2...)

	balances, err := service.CalculateNetBalances(demoExpenses, demoSplits)
	if err != nil {
		fmt.Printf("Lỗi tính số dư demo: %v\n", err)
		return
	}

	settlements, err := service.SimplifyDebts(balances)
	if err != nil {
		fmt.Printf("Lỗi SimplifyDebts demo: %v\n", err)
		return
	}

	fmt.Println("\n==================================================================================")
	fmt.Println("                   KẾT QUẢ QUYẾT TOÁN TỐI ƯU (EXPECTED RESULT)")
	fmt.Println("==================================================================================")
	for i, s := range settlements {
		fromNode, _ := demoBST.Search(s.FromUserID)
		toNode, _ := demoBST.Search(s.ToUserID)
		fmt.Printf("  Giao dịch %d: [%s (%s)]  ===>  trả  ===>  [%s (%s)] :  %d VND\n",
			i+1, fromNode.HoTen, s.FromUserID, toNode.HoTen, s.ToUserID, s.AmountMinor)
	}
	fmt.Println("==================================================================================")
}

func (app *CLIApp) readString(prompt string) string {
	fmt.Print(prompt)
	text, _ := app.reader.ReadString('\n')
	return strings.TrimSpace(text)
}
