package service

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"chiadeu/pkg/expense"
	"chiadeu/pkg/member"
	"chiadeu/pkg/settlement"
)

// MemberData struct lưu thông tin thành viên dạng phẳng để serialize JSON hoặc đệm đọc file
type MemberData struct {
	MaTV         string `json:"ma_tv"`
	HoTen        string `json:"ho_ten"`
	TongDaUng    int64  `json:"tong_da_ung"`
	TongPhaiChiu int64  `json:"tong_phai_chiu"`
}

// GroupStorageData struct đại diện cho toàn bộ dữ liệu ứng dụng cần lưu file JSON
type GroupStorageData struct {
	Members          []MemberData                 `json:"members"`
	Expenses         []expense.Expense            `json:"expenses"`
	Splits           []expense.ExpenseSplit       `json:"splits"`
	Batches          []settlement.SettlementBatch `json:"batches"`
	ActiveBatch      *settlement.SettlementBatch  `json:"active_batch,omitempty"`
	ActiveSettlement []settlement.Settlement      `json:"active_settlements,omitempty"`
}

// SaveToFile lưu toàn bộ dữ liệu ra file JSON
func SaveToFile(filePath string, bst *member.BST, expenses []expense.Expense, splits []expense.ExpenseSplit, batches []settlement.SettlementBatch, activeBatch *settlement.SettlementBatch, activeSettlement []settlement.Settlement) error {
	nodes := bst.TraverseLNR()
	memberDataList := make([]MemberData, len(nodes))
	for i, n := range nodes {
		memberDataList[i] = MemberData{
			MaTV:         n.MaTV,
			HoTen:        n.HoTen,
			TongDaUng:    n.TongDaUng,
			TongPhaiChiu: n.TongPhaiChiu,
		}
	}

	data := GroupStorageData{
		Members:          memberDataList,
		Expenses:         expenses,
		Splits:           splits,
		Batches:          batches,
		ActiveBatch:      activeBatch,
		ActiveSettlement: activeSettlement,
	}

	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("lỗi chuyển đổi dữ liệu sang JSON: %w", err)
	}

	err = os.WriteFile(filePath, bytes, 0644)
	if err != nil {
		return fmt.Errorf("lỗi ghi file '%s': %w", filePath, err)
	}

	return nil
}

// LoadFromFile đọc dữ liệu từ file JSON và tái tạo BST thành viên
func LoadFromFile(filePath string) (*member.BST, []expense.Expense, []expense.ExpenseSplit, []settlement.SettlementBatch, *settlement.SettlementBatch, []settlement.Settlement, error) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, fmt.Errorf("lỗi đọc file '%s': %w", filePath, err)
	}

	var data GroupStorageData
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, fmt.Errorf("lỗi giải mã JSON từ file '%s': %w", filePath, err)
	}

	bst := member.NewBST()
	for _, m := range data.Members {
		err := bst.Insert(m.MaTV, m.HoTen)
		if err != nil {
			// Bỏ qua nếu MaTV trùng lập
			continue
		}
		if node, found := bst.Search(m.MaTV); found {
			node.TongDaUng = m.TongDaUng
			node.TongPhaiChiu = m.TongPhaiChiu
		}
	}

	return bst, data.Expenses, data.Splits, data.Batches, data.ActiveBatch, data.ActiveSettlement, nil
}

// ImportMembersFromTxt đọc danh sách thành viên từ file TXT (định dạng: MaTV,HoTen hoặc MaTV|HoTen mỗi dòng)
func ImportMembersFromTxt(filePath string) ([]MemberData, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("không thể mở file txt '%s': %w", filePath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var members []MemberData

	lineNo := 0
	for scanner.Scan() {
		lineNo++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") {
			continue // Bỏ qua dòng trống hoặc comment
		}

		var parts []string
		if strings.Contains(line, ",") {
			parts = strings.Split(line, ",")
		} else if strings.Contains(line, "|") {
			parts = strings.Split(line, "|")
		} else if strings.Contains(line, "\t") {
			parts = strings.Split(line, "\t")
		} else {
			parts = strings.Fields(line)
		}

		if len(parts) < 2 {
			continue
		}

		maTV := strings.TrimSpace(parts[0])
		hoTen := strings.TrimSpace(strings.Join(parts[1:], " "))

		members = append(members, MemberData{
			MaTV:  maTV,
			HoTen: hoTen,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("lỗi đọc dòng từ file '%s': %w", filePath, err)
	}

	return members, nil
}
