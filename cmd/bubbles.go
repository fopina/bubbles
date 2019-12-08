package cmd

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/fopina/bubbles/data"
	"github.com/fopina/bubbles/internal/bubbles"
	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
)

func loadChunk(filename string) (*mix.Chunk, error) {
	file, err := data.Assets.Open(filename)
	if err != nil {
		return nil, err
	}
	p, _ := ioutil.ReadAll(file)
	src, _ := sdl.RWFromMem(p)
	chunk, err := mix.LoadWAVRW(src, true)
	if err != nil {
		return nil, err
	}
	return chunk, nil
}

func run(windowWidth, windowHeight int32, fullscreen bool) int {
	var window *sdl.Window
	var renderer *sdl.Renderer
	var mousePos bubbles.MousePos
	var bubs [NumBubbles]*bubbles.Bubble
	var pops []*mix.Chunk
	var runningMutex sync.Mutex
	var err error

	if err := sdl.Init(sdl.INIT_AUDIO); err != nil {
		fmt.Println(err)
		return 1
	}
	defer sdl.Quit()

	if err := mix.OpenAudio(44100, mix.DEFAULT_FORMAT, 1, 4096); err != nil {
		fmt.Println(err)
		return 1
	}
	defer mix.CloseAudio()

	pops = make([]*mix.Chunk, 4)
	pops[0], err = loadChunk("pop1.ogg")
	if err != nil {
		fmt.Println(err)
		return 1
	}
	defer pops[0].Free()

	pops[1], err = loadChunk("pop2.ogg")
	if err != nil {
		fmt.Println(err)
		return 1
	}
	defer pops[1].Free()

	pops[2], err = loadChunk("pop3.ogg")
	if err != nil {
		fmt.Println(err)
		return 1
	}
	defer pops[2].Free()

	pops[3], err = loadChunk("pop4.ogg")
	if err != nil {
		fmt.Println(err)
		return 1
	}
	defer pops[3].Free()

	sdl.Do(func() {

		window, err = sdl.CreateWindow(
			WindowTitle,
			sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
			windowWidth, windowHeight,
			sdl.WINDOW_OPENGL,
		)
		if fullscreen {
			window.SetFullscreen(sdl.WINDOW_FULLSCREEN)
		}
		sdl.SetCursor(sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_HAND))
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
		bubs[i] = bubbles.NewBubble(windowWidth, windowHeight)
	}

	running := true
	for running {
		sdl.Do(func() {
			renderer.Clear()
			renderer.SetDrawColor(0, 0, 0, 0x20)
			renderer.FillRect(&sdl.Rect{0, 0, windowWidth, windowHeight})
			for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
				switch event.(type) {
				case *sdl.KeyboardEvent:
					ev2, _ := event.(*sdl.KeyboardEvent)
					if ev2.Keysym.Sym == 27 || ev2.Keysym.Sym == 113 {
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
		})

		wg := sync.WaitGroup{}
		for i := range bubs {
			wg.Add(1)
			go func(i int) {
				bubs[i].Update(&mousePos, pops)
				sdl.Do(func() {
					bubs[i].Render(renderer)
				})
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
	var height, width int
	var fullscreen bool
	flag.IntVar(&width, "width", WindowWidth, "window width")
	flag.IntVar(&width, "w", WindowWidth, "window width")
	flag.IntVar(&height, "height", WindowHeight, "window height")
	flag.IntVar(&height, "hh", WindowHeight, "window height")
	flag.BoolVar(&fullscreen, "fullscreen", false, "use fullscreen mode")
	flag.BoolVar(&fullscreen, "fs", false, "use fullscreen mode")

	flag.Parse()

	// os.Exit(..) must run AFTER sdl.Main(..) below; so keep track of exit
	// status manually outside the closure passed into sdl.Main(..) below
	var exitcode int
	sdl.Main(func() {
		exitcode = run(int32(width), int32(height), fullscreen)
	})
	// os.Exit(..) must run here! If run in sdl.Main(..) above, it will cause
	// premature quitting of sdl.Main(..) function; resource cleaning deferred
	// calls/closing of channels may never run
	os.Exit(exitcode)
}
