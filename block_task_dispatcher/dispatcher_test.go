package block_task_dispatcher_test

import (
	"github.com/stretchr/testify/require"
	"solana-program-scanner/block_task_dispatcher"
	"sync"
	"testing"
)

func prepareTest(start uint64, end uint64, chanSize int) (uint64, uint64, *block_task_dispatcher.TaskDispatcher, chan uint64) {
	taskCh := make(chan uint64, chanSize)
	dispatcher := block_task_dispatcher.New(taskCh)
	return start, end, dispatcher, taskCh
}

func TestGetBlockTaskDispatcher_DispatchingTask1(t *testing.T) {
	startSlot, endSlot, d, taskCh := prepareTest(10, 20, 100)

	nextSlot, stopped := d.dispatchTasks(startSlot, endSlot, true)
	require.Equal(t, uint64(21), nextSlot)
	require.Equal(t, false, stopped)
	for slot := range taskCh {
		require.Equal(t, startSlot, slot)
		startSlot++
	}
}

func TestGetBlockTaskDispatcher_DispatchingTask2(t *testing.T) {
	startSlot, endSlot, d, taskCh := prepareTest(10, 20, 5)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for slot := range taskCh {
			require.Equal(t, startSlot, slot)
			startSlot++
		}
	}()

	nextSlot, stopped := d.dispatchTasks(startSlot, endSlot, true)
	require.Equal(t, uint64(21), nextSlot)
	require.Equal(t, false, stopped)

	wg.Wait()
}
