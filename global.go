package nanogui

import (
	"sync"
	"time"

	"github.com/goxjs/gl"
	"github.com/shibukawa/glfw"
)

var mainloopActive bool = false
var startTime time.Time
var debugFlag bool

type Application struct {
	Screen            *Screen
	MainThreadThunker chan func()
	Progress          *ProgressBar
	Shader            *GLShader
	Globals           map[string]string
}

func (a *Application) SetGlobal(s, v string) {
	a.Globals[s] = v
}

func (a *Application) GetGlobal(s string) string {
	v := a.Globals[s]
	return v
}

func Init() {
	err := glfw.Init(gl.ContextWatcher)
	if err != nil {
		panic(err)
	}
	startTime = time.Now()
}

func GetTime() float32 {
	return float32(time.Now().Sub(startTime)/time.Millisecond) * 0.001
}

func MainLoop(app *Application) {
	mainloopActive = true

	var wg sync.WaitGroup

	/* If there are no mouse/keyboard events, try to refresh the
	view roughly every 50 ms; this is to support animations
	such as progress bars while keeping the system load
	reasonably low */
	wg.Add(1)
	go func() {
		for mainloopActive {
			time.Sleep(100 * time.Millisecond)
			glfw.PostEmptyEvent()
		}
		wg.Done()
	}()
	for mainloopActive {
		haveActiveScreen := false
		for x := len(app.MainThreadThunker); len(app.MainThreadThunker) > 0; x = x + len(app.MainThreadThunker) {
			//fmt.Println("Processing", x, "events")
			f := <-app.MainThreadThunker
			f()
			//fmt.Println(" ", x-1, "events remaining")
		}
		for _, screen := range nanoguiScreens {
			if !screen.Visible() {
				continue
			} else if screen.GLFWWindow().ShouldClose() {
				screen.SetVisible(false)
				continue
			}

			//screen.DebugPrint()

			screen.PerformLayout()

			screen.DrawAll()
			haveActiveScreen = true
		}
		if !haveActiveScreen {
			mainloopActive = false
			break
		}
		glfw.WaitEvents()
	}

	wg.Wait()
}

func SetDebug(d bool) {
	debugFlag = d
}

func InitWidget(child, parent Widget) {
	//w.cursor = Arrow
	if parent != nil {
		parent.AddChild(parent, child)
		child.SetTheme(parent.Theme())
	}
	child.SetVisible(true)
	child.SetEnabled(true)
	child.SetFontSize(-1)
}
