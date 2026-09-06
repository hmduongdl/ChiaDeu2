package service

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"chiadeu/pkg/expense"
)

// SplitEqual chia đều khoản chi cho danh sách thành viên memberIDs
// GIỮ NGUYÊN CHÍNH XÁC LOGIC THUẬT TOÁN GỐC:
// 1. base := amountMinor / count, remainder := amountMinor % count
// 2. Mỗi thành viên nhận base
// 3. sort danh sách theo MaTV tăng dần
// 4. Cộng thêm 1 đơn vị cho "remainder" người đầu tiên sau khi sort
// 5. Trả lỗi nếu amountMinor <= 0 hoặc danh sách thành viên rỗng
func SplitEqual(expenseID string, amountMinor int64, memberIDs []string) ([]expense.ExpenseSplit, error) {
	if amountMinor <= 0 {
		return nil, fmt.Errorf("amountMinor phải > 0, giá trị nhận được: %d", amountMinor)
	}

	count := len(memberIDs)
	if count == 0 {
		return nil, errors.New("danh sách thành viên tham gia chia tiền rỗng")
	}

	// Sao chép và sắp xếp danh sách MaTV tăng dần
	sortedIDs := make([]string, count)
	for i, id := range memberIDs {
		sortedIDs[i] = strings.TrimSpace(id)
	}
	sort.Strings(sortedIDs)

	base := amountMinor / int64(count)
	remainder := amountMinor % int64(count)

	splits := make([]expense.ExpenseSplit, count)
	for i, maTV := range sortedIDs {
		share := base
		if int64(i) < remainder {
			share++
		}
		splits[i] = expense.ExpenseSplit{
			ExpenseID:  expenseID,
			UserID:     maTV,
			ShareMinor: share,
		}
	}

	return splits, nil
}

// SplitPercent chia tiền theo phần trăm từng người (tổng percents phải = 100)
func SplitPercent(expenseID string, amountMinor int64, percentMap map[string]float64) ([]expense.ExpenseSplit, error) {
	if amountMinor <= 0 {
		return nil, errors.New("amountMinor phải > 0")
	}
	if len(percentMap) == 0 {
		return nil, errors.New("danh sách phần trăm rỗng")
	}

	var sumPercent float64 = 0
	memberIDs := make([]string, 0, len(percentMap))
	for id, pct := range percentMap {
		if pct < 0 {
			return nil, fmt.Errorf("tỷ lệ %% của thành viên '%s' không được âm", id)
		}
		sumPercent += pct
		memberIDs = append(memberIDs, id)
	}

	if sumPercent < 99.99 || sumPercent > 100.01 {
		return nil, fmt.Errorf("tổng tỷ lệ phần trăm phải bằng 100%% (hiện tại: %.2f%%)", sumPercent)
	}

	sort.Strings(memberIDs)

	var calculatedSum int64 = 0
	splits := make([]expense.ExpenseSplit, len(memberIDs))

	for i, id := range memberIDs {
		pct := percentMap[id]
		share := int64((float64(amountMinor) * pct) / 100.0)
		splits[i] = expense.ExpenseSplit{
			ExpenseID:  expenseID,
			UserID:     id,
			ShareMinor: share,
		}
		calculatedSum += share
	}

	// Phân bổ phần lẻ chênh lệch (nếu có do làm tròn số nguyên) cho người đầu tiên
	diff := amountMinor - calculatedSum
	if diff != 0 && len(splits) > 0 {
		splits[0].ShareMinor += diff
	}

	return splits, nil
}

// SplitWeight chia tiền theo trọng số (ví dụ: người 1 phần, người 2 phần...)
func SplitWeight(expenseID string, amountMinor int64, weightMap map[string]int) ([]expense.ExpenseSplit, error) {
	if amountMinor <= 0 {
		return nil, errors.New("amountMinor phải > 0")
	}
	if len(weightMap) == 0 {
		return nil, errors.New("danh sách trọng số rỗng")
	}

	var totalWeight int = 0
	memberIDs := make([]string, 0, len(weightMap))
	for id, w := range weightMap {
		if w <= 0 {
			return nil, fmt.Errorf("trọng số của thành viên '%s' phải > 0", id)
		}
		totalWeight += w
		memberIDs = append(memberIDs, id)
	}

	sort.Strings(memberIDs)

	var calculatedSum int64 = 0
	splits := make([]expense.ExpenseSplit, len(memberIDs))

	for i, id := range memberIDs {
		w := weightMap[id]
		share := (amountMinor * int64(w)) / int64(totalWeight)
		splits[i] = expense.ExpenseSplit{
			ExpenseID:  expenseID,
			UserID:     id,
			ShareMinor: share,
		}
		calculatedSum += share
	}

	// Phân bổ phần dư làm tròn cho người đầu tiên
	diff := amountMinor - calculatedSum
	if diff != 0 && len(splits) > 0 {
		splits[0].ShareMinor += diff
	}

	return splits, nil
}

// SplitCustom nhận trực tiếp mảng số tiền tùy ý do người dùng nhập
func SplitCustom(expenseID string, customMap map[string]int64) ([]expense.ExpenseSplit, error) {
	if len(customMap) == 0 {
		return nil, errors.New("danh sách chia tiền tùy ý rỗng")
	}

	memberIDs := make([]string, 0, len(customMap))
	for id := range customMap {
		memberIDs = append(memberIDs, id)
	}
	sort.Strings(memberIDs)

	splits := make([]expense.ExpenseSplit, len(memberIDs))
	for i, id := range memberIDs {
		splits[i] = expense.ExpenseSplit{
			ExpenseID:  expenseID,
			UserID:     id,
			ShareMinor: customMap[id],
		}
	}

	return splits, nil
}
