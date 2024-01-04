package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"runtime"
	"sync"
	"time"
)

var (
	mu             sync.Mutex
	prevGoroutines int
)

const (
	warningThreshold = 0.2
	checkInterval    = 300 * time.Millisecond
)

func monitorGoroutines(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(checkInterval):
			currentGoroutines := runtime.NumGoroutine()

			mu.Lock()
			defer mu.Unlock()

			if prevGoroutines != 0 {
				change := currentGoroutines - prevGoroutines
				changePercentage := float64(change) / float64(prevGoroutines)
				if changePercentage > warningThreshold {
					fmt.Printf("⚠️ Предупреждение: Количество горутин увеличилось более чем на 20%%!\n")
				} else if changePercentage < -warningThreshold {
					fmt.Printf("⚠️ Предупреждение: Количество горутин уменьшилось более чем на 20%%!\n")
				}
			}

			fmt.Printf("Текущее количество горутин: %d\n", currentGoroutines)
			prevGoroutines = currentGoroutines
		}
	}
}

func main() {
	g, ctx := errgroup.WithContext(context.Background())

	// Мониторинг горутин
	go func() {
		monitorGoroutines(ctx)
	}()

	// Имитация активной работы приложения с созданием горутин
	for i := 0; i < 64; i++ {
		g.Go(func() error {
			time.Sleep(5 * time.Second)
			return nil
		})
		time.Sleep(80 * time.Millisecond)
	}

	// Ожидание завершения всех горутин
	if err := g.Wait(); err != nil {
		fmt.Println("Ошибка:", err)
	}
}
