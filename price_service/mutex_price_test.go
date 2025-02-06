package price_service_test

import (
	"github.com/shopspring/decimal"
	"solana-program-scanner/price_service"
	"sync"
	"testing"
	"time"
)

func TestMutexPrice(t *testing.T) {
	price100 := decimal.NewFromFloat(100)
	price200 := decimal.NewFromFloat(200)
	mp := price_service.NewMutexPrice()
	mp.Set(price100)

	if !mp.Get().Equal(price100) {
		t.Errorf("期望初始价格为100.0，实际获得：%v", mp.Get())
	}

	mp.Set(price200)
	if !mp.Get().Equal(price200) {
		t.Errorf("期望更新后价格为200.0，实际获得：%v", mp.Get())
	}

	// 测试并发安全性
	var wg sync.WaitGroup
	iterations := 1000

	wg.Add(2)
	// 并发写入
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			mp.Set(decimal.NewFromInt(int64(i)))
			time.Sleep(time.Microsecond)
		}
	}()

	// 并发读取
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			_ = mp.Get()
			time.Sleep(time.Microsecond)
		}
	}()

	wg.Wait()
}
