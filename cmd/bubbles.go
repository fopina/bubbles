package cmd

import (
	"fmt"
	"os"
	"sync"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/fopina/bubbles/internal/bubbles"
)


func run() int {
	var window *sdl.Window
	var renderer *sdl.Renderer
	var mousePos bubbles.MousePos
	var bubs [NumBubbles]*bubbles.Bubble
	var runningMutex sync.Mutex
	var err error

	sdl.Do(func() {
		window, err = sdl.CreateWindow(
			WindowTitle, sdl.WINDOWPOS_UNDEFINED,
			sdl.WINDOWPOS_UNDEFINED, WindowWidth,
			WindowHeight, sdl.WINDOW_OPENGL,
		)
		//window.SetFullscreen(sdl.WINDOW_FULLSCREEN)
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
		return 1
	}
	defer func() {
		sdl.Do(func() {
			window.Destroy()
		})
	}()

	sdl.Do(func() {
		renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	})
	if err != nil {
		fmt.Fprint(os.Stderr, "Failed to create renderer: %s\n", err)
		return 2
	}
	defer func() {
		sdl.Do(func() {
			renderer.Destroy()
		})
	}()

	sdl.Do(func() {
		renderer.Clear()
	})

	for i := range bubs {
		bubs[i] = bubbles.NewBubble(WindowWidth, WindowHeight)
	}

	running := true
	for running {
		sdl.Do(func() {
			for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
				switch event.(type) {
					case *sdl.KeyboardEvent:
						ev2, _ := event.(*sdl.KeyboardEvent)
						if ev2.Keysym.Sym == 27 {
							runningMutex.Lock()
							running = false
							runningMutex.Unlock()
						}
					case *sdl.MouseMotionEvent:
						ev2, _ := event.(*sdl.MouseMotionEvent)
						mousePos.X = ev2.X
						mousePos.Y = ev2.Y
					case *sdl.QuitEvent:
						runningMutex.Lock()
						running = false
						runningMutex.Unlock()
					}
			}

			renderer.Clear()
			renderer.SetDrawColor(0, 0, 0, 0x20)
			renderer.FillRect(&sdl.Rect{0, 0, WindowWidth, WindowHeight})
		})

		wg := sync.WaitGroup{}
		for i := range bubs {
			wg.Add(1)
			go func(i int) {
				bubs[i].Update(&mousePos)
				bubs[i].Render(renderer)
				wg.Done()
			}(i)
		}
		wg.Wait()

		sdl.Do(func() {
			renderer.Present()
			sdl.Delay(2000 / FrameRate)
		})
	}

	return 0
}

func Run() {
	// os.Exit(..) must run AFTER sdl.Main(..) below; so keep track of exit
	// status manually outside the closure passed into sdl.Main(..) below
	var exitcode int
	sdl.Main(func() {
		exitcode = run()
	})
	// os.Exit(..) must run here! If run in sdl.Main(..) above, it will cause
	// premature quitting of sdl.Main(..) function; resource cleaning deferred
	// calls/closing of channels may never run
	os.Exit(exitcode)
}
