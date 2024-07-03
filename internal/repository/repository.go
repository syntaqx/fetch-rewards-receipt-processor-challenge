package repository

import (
	"sync"

	"github.com/syntaqx/fetch-rewards-receipt-processor-challenge/internal/model"
)

type ReceiptRepository interface {
	SaveReceipt(id string, receipt model.Receipt, points int64)
	GetPoints(id string) (int64, bool)
}

type receiptRepository struct {
	sync.RWMutex
	receipts map[string]model.Receipt
	points   map[string]int64
}

func NewReceiptRepository() ReceiptRepository {
	return &receiptRepository{
		receipts: make(map[string]model.Receipt),
		points:   make(map[string]int64),
	}
}

func (r *receiptRepository) SaveReceipt(id string, receipt model.Receipt, points int64) {
	r.Lock()
	defer r.Unlock()
	r.receipts[id] = receipt
	r.points[id] = points
}

func (r *receiptRepository) GetPoints(id string) (int64, bool) {
	r.RLock()
	defer r.RUnlock()
	points, exists := r.points[id]
	return points, exists
}
