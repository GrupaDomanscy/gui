package gui

import (
	"errors"
	"runtime"

	"domanscy.group/gui/components"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type AppRoutine = func() (stopWindowRendering bool)

type App struct {
	title       string
	rootElement components.Component

	fontStore map[string]rl.Font
	routine   AppRoutine

	windowSize rl.Vector2
}

func newApp(title string, initialSize rl.Vector2, root components.Component, routine AppRoutine) *App {
	// This lock os thread thing protects from crashing in tests when loading fonts.
	// I don't know exactly why it works, but it works.
	runtime.LockOSThread()

	app := &App{
		title:       title,
		rootElement: root,
		fontStore:   map[string]rl.Font{},
		routine:     routine,
		windowSize:  initialSize,
	}

	rl.InitWindow(int32(initialSize.X), int32(initialSize.Y), app.title)

	rl.SetWindowState(rl.FlagWindowResizable)

	for !rl.IsWindowReady() {
		// ...
	}

	return app
}

func (app *App) run() {
	app.rootElement.CalculateSize(app.getFont, rl.Vector2{X: float32(rl.GetRenderWidth()), Y: float32(rl.GetRenderHeight())})

	for !rl.WindowShouldClose() {
		// window resize event thingies
		newWindowSize := rlGetWindowSize()

		if !rl.Vector2Equals(newWindowSize, app.windowSize) {
			oldWindowSize := app.windowSize
			app.windowSize = newWindowSize

			app.rootElement.PropagateEvent("gui:window-resize", oldWindowSize, app.windowSize)
			app.rootElement.CalculateSize(app.getFont, app.windowSize)
		}

		// the rest
		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)

		app.rootElement.Render(app.getFont)

		if app.routine != nil {
			exitFlag := app.routine()

			rl.EndDrawing()

			if exitFlag {
				break
			}
		} else {
			rl.EndDrawing()
		}
	}

	app.unloadAllFonts()
	rl.CloseWindow()

	// This lock os thread thing protects from crashing in tests when loading fonts.
	// I don't know exactly why it works, but it works.
	runtime.UnlockOSThread()
}

func rlGetWindowSize() rl.Vector2 {
	return rl.Vector2{
		X: float32(rl.GetRenderWidth()),
		Y: float32(rl.GetRenderHeight()),
	}
}

var ErrFontDoesNotExist = errors.New("font does not exist")

func (app *App) loadFont(fontName string, fontFilePath string) {
	if _, err := app.getFont(fontName); err == nil {
		return
	}

	font := rl.LoadFontEx(fontFilePath, 32, nil, 1024)

	app.fontStore[fontName] = font
}

func (app *App) unloadAllFonts() {
	for _, font := range app.fontStore {
		rl.UnloadFont(font)
	}
}

func (app *App) getFont(fontName string) (rl.Font, error) {
	font, ok := app.fontStore[fontName]

	if !ok {
		return rl.Font{}, ErrFontDoesNotExist
	}

	return font, nil
}

type AppBuilder struct {
	title       string
	initialSize rl.Vector2
	fontsToLoad map[string]string
	rootElement components.Component
	appRoutine  AppRoutine
}

func BuildApp() *AppBuilder {
	return &AppBuilder{
		title:       "",
		initialSize: rl.Vector2Zero(),
		fontsToLoad: map[string]string{},
		rootElement: nil,
		appRoutine:  nil,
	}
}

func (builder *AppBuilder) WithTitle(title string) *AppBuilder {
	builder.title = title
	return builder
}

func (builder *AppBuilder) WithInitialSize(x int, y int) *AppBuilder {
	builder.initialSize.X = float32(x)
	builder.initialSize.Y = float32(y)
	return builder
}

func (builder *AppBuilder) WithFont(fontName string, fontPath string) *AppBuilder {
	builder.fontsToLoad[fontName] = fontPath
	return builder
}

func (builder *AppBuilder) WithRootElement(rootElement components.Component) *AppBuilder {
	builder.rootElement = rootElement
	return builder
}

func (builder *AppBuilder) WithAppRoutine(routine AppRoutine) *AppBuilder {
	builder.appRoutine = routine
	return builder
}

func (builder *AppBuilder) Run() {
	app := newApp(builder.title, builder.initialSize, builder.rootElement, builder.appRoutine)

	for fontName, fontPath := range builder.fontsToLoad {
		app.loadFont(fontName, fontPath)
	}

	app.run()
}
