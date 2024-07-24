package gui

import (
	"errors"

	"domanscy.group/gui/components"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type App struct {
	title       string
	rootElement components.Component

	fontStore map[string]rl.Font
}

func newApp(title string, initialSize rl.Vector2, root components.Component) *App {
	app := &App{
		title:       title,
		rootElement: root,
		fontStore:   map[string]rl.Font{},
	}

	rl.InitWindow(int32(initialSize.X), int32(initialSize.Y), app.title)

	return app
}

func (app *App) run() {
	app.rootElement.CalculateSize(app.getFont, rl.Vector2{X: float32(rl.GetRenderWidth()), Y: float32(rl.GetRenderHeight())})

	defer rl.CloseWindow()

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)

		app.rootElement.Render(app.getFont)

		rl.EndDrawing()
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
}

func BuildApp() *AppBuilder {
	return &AppBuilder{
		title:       "",
		initialSize: rl.Vector2Zero(),
		fontsToLoad: map[string]string{},
		rootElement: nil,
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

func (builder *AppBuilder) Run() {
	app := newApp(builder.title, builder.initialSize, builder.rootElement)

	for fontName, fontPath := range builder.fontsToLoad {
		app.loadFont(fontName, fontPath)
	}

	app.run()
}
