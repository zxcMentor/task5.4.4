package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"runtime"
	"time"
)

const (
	warningThreshold = 0.2
	checkInterval    = 300 * time.Millisecond
)

func monitorGoroutines(ctx context.Context) {
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()
	prevGoroutines := runtime.NumGoroutine()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			currentGoroutines := runtime.NumGoroutine()
			change := float64(currentGoroutines-prevGoroutines) / float64(prevGoroutines)
			if change > warningThreshold || change < -warningThreshold {
				fmt.Printf("⚠️ Предупреждение: Количество горутин изменилось более чем на 20%%!\n")
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
		fmt.Println("Ошибка :", err)
	}

}
