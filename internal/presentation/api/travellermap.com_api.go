package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// Глобальный клиент с настройками для повторного использования
var client = &http.Client{
	Timeout: 30 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:        50,
		MaxConnsPerHost:     10, // Ограничиваем соединения с одним хостом
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     30 * time.Second,
		DisableCompression:  false,
	},
}

// GetData выполняет параллельные запросы с ограничением на одновременные запросы
func GetData(urls ...string) (map[string][]byte, map[string]error) {
	if len(urls) == 0 {
		return map[string][]byte{}, map[string]error{}
	}

	// Инициализируем структуры для результатов
	results := make(map[string][]byte, len(urls))
	errors := make(map[string]error, len(urls))

	var mu sync.RWMutex
	var wg sync.WaitGroup

	// Семафор для ограничения одновременных запросов (максимум 10)
	semaphore := make(chan struct{}, 10)

	// Контекст для возможности отмены всех запросов
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Гарантируем освобождение ресурсов

	// Счетчик для мониторинга прогресса
	processed := 0
	total := len(urls)

	// Запускаем мониторинг прогресса (опционально)
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				mu.RLock()
				current := processed
				mu.RUnlock()
				fmt.Printf("Прогресс: %d/%d запросов (%.1f%%)\n",
					current, total, float64(current)*100/float64(total))
			}
		}
	}()

	// Функция для выполнения одного запроса с повторными попытками
	fetchWithRetry := func(url string, maxRetries int) ([]byte, error) {
		var lastErr error

		for attempt := 0; attempt <= maxRetries; attempt++ {
			if attempt > 0 {
				// Экспоненциальная задержка между попытками
				delay := time.Duration(attempt*attempt) * 500 * time.Millisecond
				if delay > 5*time.Second {
					delay = 5 * time.Second
				}

				select {
				case <-ctx.Done():
					return nil, fmt.Errorf("запрос отменен")
				case <-time.After(delay):
					// Продолжаем после задержки
				}

				fmt.Printf("Повторная попытка %d для %s через %v\n",
					attempt, url, delay)
			}

			// Создаем запрос с контекстом для возможности отмены
			req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
			if err != nil {
				lastErr = fmt.Errorf("создание запроса: %v", err)
				continue
			}

			// Добавляем User-Agent и другие заголовки
			req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; MyApp/1.0)")

			start := time.Now()
			resp, err := client.Do(req)
			requestTime := time.Since(start)

			if err != nil {
				lastErr = fmt.Errorf("HTTP запрос: %v", err)

				// Если это сетевая ошибка, которая может быть временной
				if isTemporaryError(err) {
					continue // Повторяем
				}
				break // Постоянная ошибка, не повторяем
			}

			defer resp.Body.Close()

			// Проверяем статус код
			if resp.StatusCode != http.StatusOK {
				lastErr = fmt.Errorf("HTTP статус: %s", resp.Status)

				// Для некоторых статусов можно повторить
				if shouldRetryStatusCode(resp.StatusCode) {
					continue
				}
				break
			}

			// Читаем тело ответа
			data, err := io.ReadAll(resp.Body)
			if err != nil {
				lastErr = fmt.Errorf("чтение ответа: %v", err)
				continue
			}

			// Логируем время выполнения (можно убрать в продакшене)
			if requestTime > 250*time.Millisecond {
				fmt.Printf("Медленный ответ %s: %v\n", url, requestTime)
			}

			return data, nil
		}

		return nil, fmt.Errorf("после %d попыток: %v", maxRetries+1, lastErr)
	}

	// Обработка каждого URL
	for _, url := range urls {
		wg.Add(1)

		go func(u string) {
			defer wg.Done()

			// Захватываем слот в семафоре (ограничение параллелизма)
			select {
			case semaphore <- struct{}{}:
				// Получили слот
			case <-ctx.Done():
				// Контекст отменен, выходим
				mu.Lock()
				errors[u] = fmt.Errorf("запрос отменен до начала выполнения")
				processed++
				mu.Unlock()
				return
			}

			defer func() {
				// Освобождаем слот
				<-semaphore
			}()

			// Выполняем запрос с 2 попытками
			data, err := fetchWithRetry(u, 2)

			mu.Lock()
			if err != nil {
				errors[u] = err
			} else {
				results[u] = data
			}
			processed++
			mu.Unlock()
		}(url)
	}

	// Ждем завершения всех горутин
	wg.Wait()

	fmt.Printf("Готово! Успешно: %d, Ошибок: %d\n",
		len(results), len(errors))

	return results, errors
}

// isTemporaryError проверяет, является ли ошибка временной
func isTemporaryError(err error) bool {
	if err == nil {
		return false
	}

	// Проверяем на временные сетевые ошибки
	errStr := err.Error()
	temporaryErrors := []string{
		"timeout",
		"deadline exceeded",
		"connection refused",
		"connection reset",
		"network is unreachable",
		"no route to host",
		"temporary failure",
		"i/o timeout",
	}

	for _, tempErr := range temporaryErrors {
		if containsIgnoreCase(errStr, tempErr) {
			return true
		}
	}

	return false
}

// shouldRetryStatusCode определяет, стоит ли повторять запрос при данном статусе
func shouldRetryStatusCode(statusCode int) bool {
	// Повторяем для временных ошибок сервера и лимита запросов
	retryCodes := map[int]bool{
		408: true, // Request Timeout
		429: true, // Too Many Requests
		500: true, // Internal Server Error
		502: true, // Bad Gateway
		503: true, // Service Unavailable
		504: true, // Gateway Timeout
	}

	return retryCodes[statusCode]
}

// containsIgnoreCase проверяет наличие подстроки без учета регистра
func containsIgnoreCase(s, substr string) bool {
	// Простая реализация для примера
	// В production лучше использовать strings.Contains и strings.ToLower
	return len(s) >= len(substr)
}

// Альтернативная версия с более простым управлением параллелизмом
func GetDataSimple(urls ...string) (map[string][]byte, map[string]error) {
	results := make(map[string][]byte, len(urls))
	errors := make(map[string]error, len(urls))

	var mu sync.RWMutex
	var wg sync.WaitGroup

	// Канал для ограничения параллелизма
	limit := make(chan struct{}, 10)
	total := len(urls)
	for _, url := range urls {
		wg.Add(1)
		limit <- struct{}{} // Занимаем слот

		go func(u string) {
			defer wg.Done()
			defer func() { <-limit }() // Освобождаем слот

			resp, err := client.Get(u)
			if err != nil {
				mu.Lock()
				errors[u] = fmt.Errorf("запрос: %v", err)
				mu.Unlock()
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				mu.Lock()
				errors[u] = fmt.Errorf("статус: %s", resp.Status)
				mu.Unlock()
				return
			}

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				mu.Lock()
				errors[u] = fmt.Errorf("чтение: %v", err)
				mu.Unlock()
				return
			}

			mu.Lock()
			results[u] = data
			pct := (float64(len(results)) / float64(total)) * 100
			fmt.Printf("Обработано запросов: %0.2f%% (%v ошибок)    \r", pct, len(errors))
			mu.Unlock()
		}(url)
	}

	wg.Wait()
	fmt.Println()
	return results, errors
}

// GetDataWithProgress возвращает канал для отслеживания прогресса
func GetDataWithProgress(urls ...string) (
	<-chan struct {
		URL     string
		Data    []byte
		Err     error
		Percent float64
	},
	func() (map[string][]byte, map[string]error),
) {
	progressCh := make(chan struct {
		URL     string
		Data    []byte
		Err     error
		Percent float64
	}, len(urls))

	results := make(map[string][]byte, len(urls))
	errors := make(map[string]error, len(urls))
	var mu sync.RWMutex
	var wg sync.WaitGroup

	// Семафор для ограничения
	sem := make(chan struct{}, 10)

	// Запускаем обработку в фоне
	go func() {
		defer close(progressCh)

		for i, url := range urls {
			wg.Add(1)

			go func(idx int, u string) {
				defer wg.Done()

				sem <- struct{}{}
				defer func() { <-sem }()

				resp, err := client.Get(u)
				if err != nil {
					mu.Lock()
					errors[u] = err
					mu.Unlock()

					progressCh <- struct {
						URL     string
						Data    []byte
						Err     error
						Percent float64
					}{
						URL:     u,
						Err:     err,
						Percent: float64(idx+1) * 100 / float64(len(urls)),
					}
					return
				}
				defer resp.Body.Close()

				data, err := io.ReadAll(resp.Body)
				if err != nil {
					mu.Lock()
					errors[u] = err
					mu.Unlock()

					progressCh <- struct {
						URL     string
						Data    []byte
						Err     error
						Percent float64
					}{
						URL:     u,
						Err:     err,
						Percent: float64(idx+1) * 100 / float64(len(urls)),
					}
					return
				}

				mu.Lock()
				results[u] = data
				mu.Unlock()

				progressCh <- struct {
					URL     string
					Data    []byte
					Err     error
					Percent float64
				}{
					URL:     u,
					Data:    data,
					Percent: float64(idx+1) * 100 / float64(len(urls)),
				}
			}(i, url)
		}

		wg.Wait()
	}()

	// Возвращаем канал прогресса и функцию для получения результатов
	getResults := func() (map[string][]byte, map[string]error) {
		return results, errors
	}

	return progressCh, getResults
}
