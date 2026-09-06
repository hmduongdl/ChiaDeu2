package member

import (
	"errors"
	"fmt"
	"strings"
)

// BST quản lý Cây Nhị Phân Tìm Kiếm các thành viên
type BST struct {
	Root *Node
}

// NewBST khởi tạo một cây BST rỗng
func NewBST() *BST {
	return &BST{Root: nil}
}

// Insert chèn một thành viên mới vào cây BST dựa trên MaTV
// Trả về lỗi nếu MaTV đã tồn tại hoặc MaTV rỗng
func (b *BST) Insert(maTV, hoTen string) error {
	maTV = strings.TrimSpace(maTV)
	hoTen = strings.TrimSpace(hoTen)

	if maTV == "" {
		return errors.New("MaTV không được để rỗng")
	}
	if hoTen == "" {
		return errors.New("họ tên không được để rỗng")
	}

	var err error
	b.Root, err = insertNode(b.Root, maTV, hoTen)
	return err
}

func insertNode(node *Node, maTV, hoTen string) (*Node, error) {
	if node == nil {
		return NewNode(maTV, hoTen), nil
	}

	if maTV < node.MaTV {
		var err error
		node.Left, err = insertNode(node.Left, maTV, hoTen)
		return node, err
	} else if maTV > node.MaTV {
		var err error
		node.Right, err = insertNode(node.Right, maTV, hoTen)
		return node, err
	} else {
		// Khóa MaTV trùng lập
		return node, fmt.Errorf("MaTV '%s' đã tồn tại trong hệ thống", maTV)
	}
}

// Search tìm kiếm thành viên theo MaTV dùng cơ chế BST search
func (b *BST) Search(maTV string) (*Node, bool) {
	maTV = strings.TrimSpace(maTV)
	return searchNode(b.Root, maTV)
}

func searchNode(node *Node, maTV string) (*Node, bool) {
	if node == nil {
		return nil, false
	}

	if maTV == node.MaTV {
		return node, true
	} else if maTV < node.MaTV {
		return searchNode(node.Left, maTV)
	} else {
		return searchNode(node.Right, maTV)
	}
}

// Update cập nhật thông tin họ tên của thành viên mà KHÔNG làm thay đổi khóa MaTV
func (b *BST) Update(maTV, hoTenMoi string) error {
	node, found := b.Search(maTV)
	if !found {
		return fmt.Errorf("không tìm thấy thành viên có MaTV '%s' để cập nhật", maTV)
	}
	hoTenMoi = strings.TrimSpace(hoTenMoi)
	if hoTenMoi == "" {
		return errors.New("họ tên mới không được để rỗng")
	}
	node.HoTen = hoTenMoi
	return nil
}

// Remove xóa node khỏi cây BST theo MaTV, xử lý đúng 3 trường hợp:
// 1. Node lá (0 con)
// 2. Node có 1 con
// 3. Node có 2 con (dùng node thế mạng là nút nhỏ nhất của cây con bên phải)
func (b *BST) Remove(maTV string) error {
	maTV = strings.TrimSpace(maTV)
	var removed bool
	b.Root, removed = removeNode(b.Root, maTV)
	if !removed {
		return fmt.Errorf("không tìm thấy thành viên có MaTV '%s' để xóa", maTV)
	}
	return nil
}

func removeNode(node *Node, maTV string) (*Node, bool) {
	if node == nil {
		return nil, false
	}

	if maTV < node.MaTV {
		var removed bool
		node.Left, removed = removeNode(node.Left, maTV)
		return node, removed
	} else if maTV > node.MaTV {
		var removed bool
		node.Right, removed = removeNode(node.Right, maTV)
		return node, removed
	} else {
		// Tìm thấy node cần xóa (maTV == node.MaTV)

		// TH 1: Node lá (0 con)
		if node.Left == nil && node.Right == nil {
			return nil, true
		}

		// TH 2: Node có 1 con
		if node.Left == nil {
			return node.Right, true
		}
		if node.Right == nil {
			return node.Left, true
		}

		// TH 3: Node có 2 con
		// Tìm node thế mạng: node nhỏ nhất thuộc cây con bên phải (In-order Successor)
		successor := findMinNode(node.Right)

		// Sao chép dữ liệu từ node thế mạng sang node hiện tại (giữ nguyên liên kết con)
		node.MaTV = successor.MaTV
		node.HoTen = successor.HoTen
		node.TongDaUng = successor.TongDaUng
		node.TongPhaiChiu = successor.TongPhaiChiu

		// Xóa node thế mạng khỏi cây con bên phải (chắc chắn rơi vào TH 1 hoặc TH 2)
		node.Right, _ = removeNode(node.Right, successor.MaTV)
		return node, true
	}
}

// findMinNode tìm nút có MaTV nhỏ nhất trong cây (nút trái nhất)
func findMinNode(node *Node) *Node {
	current := node
	for current != nil && current.Left != nil {
		current = current.Left
	}
	return current
}

// TraverseNLR duyệt cây theo thứ tự Tiền Tự (Node -> Left -> Right)
func (b *BST) TraverseNLR() []*Node {
	var result []*Node
	traverseNLRHelper(b.Root, &result)
	return result
}

func traverseNLRHelper(node *Node, result *[]*Node) {
	if node == nil {
		return
	}
	*result = append(*result, node)
	traverseNLRHelper(node.Left, result)
	traverseNLRHelper(node.Right, result)
}

// TraverseLNR duyệt cây theo thứ tự Trung Tự (Left -> Node -> Right)
// Mặc định trả về danh sách được sắp xếp tăng dần theo MaTV
func (b *BST) TraverseLNR() []*Node {
	var result []*Node
	traverseLNRHelper(b.Root, &result)
	return result
}

func traverseLNRHelper(node *Node, result *[]*Node) {
	if node == nil {
		return
	}
	traverseLNRHelper(node.Left, result)
	*result = append(*result, node)
	traverseLNRHelper(node.Right, result)
}

// TraverseLRN duyệt cây theo thứ tự Hậu Tự (Left -> Right -> Node)
func (b *BST) TraverseLRN() []*Node {
	var result []*Node
	traverseLRNHelper(b.Root, &result)
	return result
}

func traverseLRNHelper(node *Node, result *[]*Node) {
	if node == nil {
		return
	}
	traverseLRNHelper(node.Left, result)
	traverseLRNHelper(node.Right, result)
	*result = append(*result, node)
}

// CountNodes đếm tổng số lượng nút trên cây
func (b *BST) CountNodes() int {
	return countNodesHelper(b.Root)
}

func countNodesHelper(node *Node) int {
	if node == nil {
		return 0
	}
	return 1 + countNodesHelper(node.Left) + countNodesHelper(node.Right)
}

// Height tính chiều cao của cây BST
func (b *BST) Height() int {
	return heightHelper(b.Root)
}

func heightHelper(node *Node) int {
	if node == nil {
		return 0
	}
	leftH := heightHelper(node.Left)
	rightH := heightHelper(node.Right)
	if leftH > rightH {
		return leftH + 1
	}
	return rightH + 1
}

// ResetTotals xóa sạch số tiền ứng & phải chịu về 0 cho tất cả nút
func (b *BST) ResetTotals() {
	nodes := b.TraverseLNR()
	for _, n := range nodes {
		n.TongDaUng = 0
		n.TongPhaiChiu = 0
	}
}
