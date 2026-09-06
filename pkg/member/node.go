package member

import "fmt"

// Node đại diện cho một nút trong Cây Nhị Phân Tìm Kiếm (BST)
type Node struct {
	MaTV         string // Khóa tìm kiếm (MaTV - unique string)
	HoTen        string // Họ và tên thành viên
	TongDaUng    int64  // Tổng số tiền đã ứng ra trả (đơn vị minor/VNĐ)
	TongPhaiChiu int64  // Tổng số tiền phải chịu (đơn vị minor/VNĐ)
	Left         *Node  // Cây con bên trái (MaTV nhỏ hơn)
	Right        *Node  // Cây con bên phải (MaTV lớn hơn)
}

// NewNode khởi tạo một nút mới
func NewNode(maTV, hoTen string) *Node {
	return &Node{
		MaTV:         maTV,
		HoTen:        hoTen,
		TongDaUng:    0,
		TongPhaiChiu: 0,
		Left:         nil,
		Right:        nil,
	}
}

// NetBalance tính số dư ròng của thành viên (Tổng đã ứng - Tổng phải chịu)
// Kết quả > 0: Chủ nợ (được nhận lại tiền)
// Kết quả < 0: Con nợ (phải trả thêm tiền)
func (n *Node) NetBalance() int64 {
	if n == nil {
		return 0
	}
	return n.TongDaUng - n.TongPhaiChiu
}

// String trả về chuỗi thông tin tóm tắt của Node
func (n *Node) String() string {
	if n == nil {
		return "<nil>"
	}
	return fmt.Sprintf("MaTV: %-8s | HoTen: %-20s | Ung: %10d VND | Chiu: %10d VND | Net: %10d VND",
		n.MaTV, n.HoTen, n.TongDaUng, n.TongPhaiChiu, n.NetBalance())
}
