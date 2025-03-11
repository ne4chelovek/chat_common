package closer

import (
	"log"
	"os"
	"os/signal"
	"sync"
)

// globalCloser - глобальный экземпляр Closer, который используется для управления
// закрытием всех зарегистрированных функций.
var globalCloser = New()

// Add - добавляет одну или несколько функций в глобальный Closer.
// Эти функции будут выполнены при вызове CloseAll.
func Add(f ...func() error) {
	globalCloser.Add(f...)
}

// Wait - блокирует выполнение до тех пор, пока не будет вызван метод CloseAll.
// Это позволяет дождаться завершения всех зарегистрированных функций.
func Wait() {
	globalCloser.Wait()
}

// CloseAll - вызывает выполнение всех зарегистрированных функций в глобальном Closer.
// После выполнения всех функций, Wait больше не блокирует выполнение.
func CloseAll() {
	globalCloser.CloseAll()
}

// Closer - структура, которая управляет набором функций, которые должны быть выполнены
// при завершении работы. Она обеспечивает безопасное добавление и выполнение функций.
type Closer struct {
	mu    sync.Mutex       // мьютекс для защиты доступа к списку функций
	once  sync.Once        // гарантирует, что CloseAll будет вызван только один раз
	done  chan struct{}    // канал для сигнализации о завершении работы
	funcs []func() error   // список функций, которые нужно выполнить при закрытии
}

// New - создает новый экземпляр Closer. Если переданы сигналы, то Closer будет
// автоматически вызывать CloseAll при получении одного из этих сигналов.
func New(sig ...os.Signal) *Closer {
	c := &Closer{done: make(chan struct{})}
	if len(sig) > 0 {
		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, sig...)
			<-ch
			signal.Stop(ch)
			c.CloseAll()
		}()
	}
	return c
}

// Add - добавляет одну или несколько функций в список функций, которые будут
// выполнены при вызове CloseAll.
func (c *Closer) Add(f ...func() error) {
	c.mu.Lock()
	c.funcs = append(c.funcs, f...)
	c.mu.Unlock()
}

// Wait - блокирует выполнение до тех пор, пока не будет вызван метод CloseAll.
// Это позволяет дождаться завершения всех зарегистрированных функций.
func (c *Closer) Wait() {
	<-c.done
}

// CloseAll - выполняет все зарегистрированные функции и закрывает канал done.
// Метод гарантирует, что все функции будут выполнены только один раз, даже если
// CloseAll вызывается несколько раз.
func (c *Closer) CloseAll() {
	c.once.Do(func() {
		defer close(c.done) // гарантирует, что канал done будет закрыт после выполнения всех функций

		c.mu.Lock()
		funcs := c.funcs
		c.funcs = nil // очищаем список функций, чтобы избежать повторного выполнения
		c.mu.Unlock()

		// Создаем канал для сбора ошибок, которые могут вернуть функции
		errs := make(chan error, len(funcs))
		for _, f := range funcs {
			go func(f func() error) {
				errs <- f() // выполняем функцию и отправляем ошибку в канал
			}(f)
		}

		// Ожидаем завершения всех функций и логируем ошибки, если они есть
		for i := 0; i < cap(errs); i++ {
			if err := <-errs; err != nil {
				log.Println("error returned from Closer:", err)
			}
		}
	})
}
