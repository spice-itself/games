package main

import (
    "bufio"
    "fmt"
    "os"
    "os/exec"
    "runtime"
    "sync"
    "time"
)

// Геймстат
type GameState struct {
    x, y     int
    maxX     int
    maxY     int
    running  bool
    mu       sync.RWMutex
}

// Функция контроля персонажа
func backend(state *GameState, updates chan<- struct{}) {
    reader := bufio.NewReader(os.Stdin)
    for state.isRunning() {
        char, _, err := reader.ReadRune()
        if err != nil {
            fmt.Println("Ошибка ввода:", err)
            return
        }

        state.mu.Lock()
        switch char {
        case 'w', 'W':
            state.y--
            if state.y < 0 {
                state.y = state.maxY
            }
        case 's', 'S':
            state.y++
            if state.y > state.maxY {
                state.y = 0
            }
        case 'd', 'D':
            state.x++
            if state.x > state.maxX {
                state.x = 0
            }
        case 'a', 'A':
            state.x--
            if state.x < 0 {
                state.x = state.maxX
            }
        case 'q', 'Q':
            state.running = false
            state.mu.Unlock()
            return
        }
        state.mu.Unlock()
        
        clear_screen()
        updates <- struct{}{} // Сигнал для очищения (или обновления хз) экрана
    }
}

// Функция отрисовки
func frontend(state *GameState, updates <-chan struct{}) {
    for state.isRunning() {
        select {
        case <-updates:
            draw_game(state)
        case <-time.After(16 * time.Millisecond): // ~60 FPS
            draw_game(state)
        }
    }
}

func draw_game(state *GameState) {
    state.mu.RLock()
    defer state.mu.RUnlock()

    // Отрисовка поля
    for i := 0; i <= state.maxY; i++ {
        for j := 0; j <= state.maxX; j++ {
            if i == state.y && j == state.x {
                fmt.Printf("[] ")
            } else {
                fmt.Printf("   ")
            }
        }
        fmt.Println()
    }
    fmt.Println("\nWASD для передвижения. Нажмите на Q чтоб выйти. Энтер чтобы обновлять экран ибо я рукожоп.")
}

// Очистка экрана с обработкой ошибок
func clear_screen() {
    var cmd *exec.Cmd
    if runtime.GOOS == "windows" {
        cmd = exec.Command("cmd", "/c", "cls")
    } else {
        cmd = exec.Command("clear")
    }
    cmd.Stdout = os.Stdout
    if err := cmd.Run(); err != nil {
        fmt.Println("Ошибка очистки экрана:", err)
    }
}

// Проверка
func (s *GameState) isRunning() bool {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.running
}

func main() {
    // Инициализация
    state := &GameState{
        x:       10,
        y:       10,
        maxX:    20,
        maxY:    20,
        running: true,
    }

    updates := make(chan struct{}, 1) // Буферизация

    fmt.Println("Игра началась! WASD для передвижения. Нажмите на Q чтоб выйти. Энтер чтобы обновлять экран ибо я рукожоп.\nПофиксите сами, мне поебать. Мои нервы уже в пределе")
    
    // Запуск горутин
    go backend(state, updates)
    go frontend(state, updates)

    // Ожидаем завершения
    for state.isRunning() {
        time.Sleep(100 * time.Millisecond)
    }
}