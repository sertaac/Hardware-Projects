package main

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

/**************************************/
/*                                    */
/*     Arcade Project v1.0 - Main     */
/*     Frutiger Aero + Y2K Edition    */
/*           Programmed by            */
/*            Sertaç Ataç             */
/*            02.01.2026              */
/*                                    */
/**************************************/

/**************************************************/
/*                                                */
/*            SCREEN CONFIGURATION                */
/*                                                */
/**************************************************/

const (
	SCREEN_WIDTH  = 1280
	SCREEN_HEIGHT = 720
)

/**************************************************/
/*                                                */
/*       NEON GREEN FRUTIGER AERO PALETTE         */
/*                                                */
/**************************************************/

var (
	/*       Primary colors - Y2K Neon Green vibes       */
	NeonGreen       = rl.Color{R: 57, G: 255, B: 20, A: 255}
	NeonGreenDark   = rl.Color{R: 30, G: 180, B: 10, A: 255}
	NeonGreenGlow   = rl.Color{R: 57, G: 255, B: 20, A: 100}
	NeonGreenBright = rl.Color{R: 100, G: 255, B: 100, A: 255}

	/*  Frutiger Aero backgrounds - glossy, clean style  */
	AeroBackground = rl.Color{R: 15, G: 25, B: 35, A: 255}
	AeroDarkPanel  = rl.Color{R: 20, G: 35, B: 50, A: 230}
	AeroGloss      = rl.Color{R: 255, G: 255, B: 255, A: 50}
	AeroShadow     = rl.Color{R: 0, G: 0, B: 0, A: 150}

	/*                   Accent colors                   */
	AccentCyan   = rl.Color{R: 0, G: 255, B: 255, A: 255}
	AccentPurple = rl.Color{R: 180, G: 100, B: 255, A: 255}
	AccentOrange = rl.Color{R: 255, G: 165, B: 0, A: 255}
	AccentPink   = rl.Color{R: 255, G: 100, B: 200, A: 255}

	/*                    Text colors                    */
	TextWhite = rl.Color{R: 255, G: 255, B: 255, A: 255}
	TextGray  = rl.Color{R: 180, G: 180, B: 180, A: 255}

	/*              Global animation time                */
	globalTime float32 = 0
)

/**************************************************/
/*                                                */
/*               GAME ITEM STRUCTURE              */
/*                                                */
/**************************************************/

type GameItem struct {
	Title       string
	Description string
	Color       rl.Color
	IconText    string
	IconShape   int
}

/**************************************************/
/*                                                */
/*        BLADE NAVIGATION (XBOX 360 STYLE)       */
/*                                                */
/**************************************************/

type BladeNav struct {
	Categories     []string
	Icons          []string
	CurrentBlade   int
	TargetBlade    int
	TransitionTime float32
	BladeWidth     float32
}

func NewBladeNav() *BladeNav {
	return &BladeNav{
		Categories:   []string{"GAMES", "MEDIA", "SETTINGS", "NETWORK"},
		Icons:        []string{">", "||", "**", "@@"},
		CurrentBlade: 0,
		TargetBlade:  0,
		BladeWidth:   280,
	}
}

func (b *BladeNav) Update() {
	if b.TransitionTime < 1.0 {
		b.TransitionTime += rl.GetFrameTime() * 3.0
		if b.TransitionTime > 1.0 {
			b.TransitionTime = 1.0
			b.CurrentBlade = b.TargetBlade
		}
	}
}

func (b *BladeNav) NextBlade() {
	if b.TargetBlade < len(b.Categories)-1 && b.TransitionTime >= 1.0 {
		b.TargetBlade++
		b.TransitionTime = 0.0
	}
}

func (b *BladeNav) PrevBlade() {
	if b.TargetBlade > 0 && b.TransitionTime >= 1.0 {
		b.TargetBlade--
		b.TransitionTime = 0.0
	}
}

func (b *BladeNav) Draw() {
	easeT := easeOutCubic(b.TransitionTime)
	currentOffset := float32(b.CurrentBlade) * b.BladeWidth
	targetOffset := float32(b.TargetBlade) * b.BladeWidth
	offset := currentOffset + (targetOffset-currentOffset)*easeT

	/*   Center the blades on screen   */
	totalWidth := float32(len(b.Categories)) * b.BladeWidth
	startX := (SCREEN_WIDTH - totalWidth) / 2

	for i := 0; i < len(b.Categories); i++ {
		x := startX + float32(i)*b.BladeWidth - offset + b.BladeWidth*float32(b.CurrentBlade)
		isActive := i == b.TargetBlade
		drawAeroBlade(int32(x), 85, int32(b.BladeWidth-15), 50, isActive, b.Icons[i], b.Categories[i])
	}
}

/**************************************************/
/*                                                */
/*         GAME GRID (NINTENDO STYLE)             */
/*                                                */
/**************************************************/

type GameGrid struct {
	Games         []GameItem
	SelectedIndex int
	HoverScale    float32
	ScrollOffset  float32
	TargetScroll  float32
	bounceTime    float32
}

func NewGameGrid() *GameGrid {
	return &GameGrid{
		Games: []GameItem{
			{Title: "RETRO RACER", Description: "High-speed arcade racing", Color: NeonGreen, IconText: "RR", IconShape: 0},
			{Title: "PIXEL QUEST", Description: "8-bit adventure game", Color: AccentCyan, IconText: "PQ", IconShape: 1},
			{Title: "NEON BLASTER", Description: "Shoot-em-up action", Color: AccentPurple, IconText: "NB", IconShape: 2},
			{Title: "SYNTH BEATS", Description: "Rhythm game", Color: AccentPink, IconText: "SB", IconShape: 3},
			{Title: "CYBER MAZE", Description: "Puzzle challenge", Color: NeonGreenBright, IconText: "CM", IconShape: 0},
			{Title: "STAR FORCE", Description: "Space shooter", Color: AccentOrange, IconText: "SF", IconShape: 1},
		},
		SelectedIndex: 0,
		HoverScale:    1.0,
	}
}

func (g *GameGrid) Update() {
	g.HoverScale = lerp(g.HoverScale, 1.12, rl.GetFrameTime()*8.0)
	g.ScrollOffset = lerp(g.ScrollOffset, g.TargetScroll, rl.GetFrameTime()*6.0)
	g.bounceTime += rl.GetFrameTime()
}

func (g *GameGrid) SelectNext() {
	if g.SelectedIndex < len(g.Games)-1 {
		g.SelectedIndex++
		g.HoverScale = 1.0
		g.bounceTime = 0
	}
}

func (g *GameGrid) SelectPrev() {
	if g.SelectedIndex > 0 {
		g.SelectedIndex--
		g.HoverScale = 1.0
		g.bounceTime = 0
	}
}

func (g *GameGrid) Draw() {
	cardWidth := int32(160)
	cardHeight := int32(220)
	spacing := int32(20)

	/*   Center the game grid on screen   */
	totalWidth := int32(len(g.Games))*(cardWidth+spacing) - spacing
	startX := (SCREEN_WIDTH - totalWidth) / 2
	startY := int32(190)

	for i, game := range g.Games {
		x := startX + int32(i)*(cardWidth+spacing)
		y := startY
		isSelected := i == g.SelectedIndex

		scale := float32(1.0)
		bounceOffset := float32(0)
		if isSelected {
			scale = g.HoverScale
			/*   Subtle floating animation   */
			bounceOffset = float32(math.Sin(float64(g.bounceTime)*3)) * 4
		}

		drawGameCard(x, y-int32(bounceOffset), cardWidth, cardHeight, scale, isSelected, game)
	}
}

/**************************************************/
/*                                                */
/*       FRUTIGER AERO DECORATIVE ELEMENTS        */
/*         Water droplets, bubbles, shine         */
/*                                                */
/**************************************************/

type AeroDecorations struct {
	bubbles   []Bubble
	droplets  []WaterDroplet
	scanlineY float32
	glowPulse float32
}

type Bubble struct {
	x, y   float32
	radius float32
	speed  float32
	alpha  uint8
}

type WaterDroplet struct {
	x, y     float32
	size     float32
	alpha    uint8
	lifeTime float32
	maxLife  float32
}

func NewAeroDecorations() *AeroDecorations {
	dec := &AeroDecorations{
		bubbles:  make([]Bubble, 20),
		droplets: make([]WaterDroplet, 8),
	}

	/*   Initialize floating bubbles   */
	for i := range dec.bubbles {
		dec.bubbles[i] = Bubble{
			x:      float32(rl.GetRandomValue(0, SCREEN_WIDTH)),
			y:      float32(rl.GetRandomValue(0, SCREEN_HEIGHT)),
			radius: float32(rl.GetRandomValue(4, 25)),
			speed:  float32(rl.GetRandomValue(20, 70)),
			alpha:  uint8(rl.GetRandomValue(25, 70)),
		}
	}

	/*   Initialize water droplets   */
	for i := range dec.droplets {
		dec.droplets[i] = WaterDroplet{
			x:       float32(rl.GetRandomValue(100, SCREEN_WIDTH-100)),
			y:       float32(rl.GetRandomValue(100, SCREEN_HEIGHT-100)),
			size:    float32(rl.GetRandomValue(15, 40)),
			alpha:   uint8(rl.GetRandomValue(30, 80)),
			maxLife: float32(rl.GetRandomValue(3, 8)),
		}
	}

	return dec
}

func (d *AeroDecorations) Update() {
	dt := rl.GetFrameTime()

	/*   Update bubbles - float upward   */
	for i := range d.bubbles {
		d.bubbles[i].y -= d.bubbles[i].speed * dt
		d.bubbles[i].x += float32(math.Sin(float64(globalTime+float32(i)))) * 0.5

		if d.bubbles[i].y < -d.bubbles[i].radius {
			d.bubbles[i].y = SCREEN_HEIGHT + d.bubbles[i].radius
			d.bubbles[i].x = float32(rl.GetRandomValue(0, SCREEN_WIDTH))
		}
	}

	/*   Update water droplets   */
	for i := range d.droplets {
		d.droplets[i].lifeTime += dt
		if d.droplets[i].lifeTime > d.droplets[i].maxLife {
			d.droplets[i].lifeTime = 0
			d.droplets[i].x = float32(rl.GetRandomValue(100, SCREEN_WIDTH-100))
			d.droplets[i].y = float32(rl.GetRandomValue(100, SCREEN_HEIGHT-100))
			d.droplets[i].size = float32(rl.GetRandomValue(15, 40))
		}
	}

	d.scanlineY += 180 * dt
	if d.scanlineY > SCREEN_HEIGHT {
		d.scanlineY = 0
	}

	d.glowPulse += dt * 2.0
}

func (d *AeroDecorations) Draw() {
	/*   Draw floating bubbles with glass effect   */
	for _, b := range d.bubbles {
		/*   Bubble body   */
		color := rl.Color{R: NeonGreen.R, G: NeonGreen.G, B: NeonGreen.B, A: b.alpha}
		rl.DrawCircle(int32(b.x), int32(b.y), b.radius, color)

		/*   Glossy highlight (Frutiger Aero signature)   */
		highlightX := b.x - b.radius*0.35
		highlightY := b.y - b.radius*0.35
		rl.DrawCircle(int32(highlightX), int32(highlightY), b.radius*0.35, rl.Color{R: 255, G: 255, B: 255, A: 100})

		/*   Small secondary highlight   */
		rl.DrawCircle(int32(b.x+b.radius*0.2), int32(b.y+b.radius*0.3), b.radius*0.15, rl.Color{R: 255, G: 255, B: 255, A: 40})
	}

	/*   Draw water droplets (Frutiger Aero style)   */
	for _, drop := range d.droplets {
		lifeRatio := drop.lifeTime / drop.maxLife
		fadeAlpha := uint8(float32(drop.alpha) * (1 - lifeRatio*0.5))

		/*   Droplet body - ellipse shape   */
		rl.DrawEllipse(int32(drop.x), int32(drop.y), drop.size*0.7, drop.size, rl.Color{R: 180, G: 220, B: 255, A: fadeAlpha / 2})

		/*   Inner glow   */
		rl.DrawEllipse(int32(drop.x), int32(drop.y), drop.size*0.5, drop.size*0.7, rl.Color{R: 200, G: 240, B: 255, A: fadeAlpha / 3})

		/*   Highlight   */
		rl.DrawCircle(int32(drop.x-drop.size*0.2), int32(drop.y-drop.size*0.4), drop.size*0.2, rl.Color{R: 255, G: 255, B: 255, A: fadeAlpha})
	}

	/*   Scanline effect   */
	scanAlpha := uint8(10)
	rl.DrawRectangle(0, int32(d.scanlineY), SCREEN_WIDTH, 2, rl.Color{R: 255, G: 255, B: 255, A: scanAlpha})

	/*   Bottom glow   */
	glowIntensity := uint8(20 + 12*float32(math.Sin(float64(d.glowPulse))))
	rl.DrawRectangleGradientV(0, SCREEN_HEIGHT-70, SCREEN_WIDTH, 70,
		rl.Color{R: 0, G: 0, B: 0, A: 0},
		rl.Color{R: NeonGreen.R, G: NeonGreen.G, B: NeonGreen.B, A: glowIntensity})
}

/**************************************************/
/*                                                */
/*      DRAWING HELPERS - FRUTIGER AERO STYLE     */
/*                                                */
/**************************************************/

func drawAeroBlade(x, y, width, height int32, isActive bool, icon, text string) {
	panelColor := AeroDarkPanel
	if isActive {
		panelColor = rl.Color{R: 25, G: 50, B: 35, A: 245}
	}

	/*   Shadow   */
	rl.DrawRectangle(x+3, y+3, width, height, AeroShadow)

	/*   Main panel   */
	rl.DrawRectangle(x, y, width, height, panelColor)

	/*   Glass shine effect (Frutiger Aero)   */
	rl.DrawRectangleGradientV(x, y, width, height/2, AeroGloss, rl.Color{A: 0})
	rl.DrawRectangleGradientH(x, y, width/3, height, rl.Color{R: 255, G: 255, B: 255, A: 20}, rl.Color{A: 0})

	if isActive {
		/*   Neon glow border   */
		rl.DrawRectangleLines(x-2, y-2, width+4, height+4, NeonGreenGlow)
		rl.DrawRectangleLines(x-1, y-1, width+2, height+2, NeonGreenGlow)
		rl.DrawRectangleLines(x, y, width, height, NeonGreen)

		/*   Icon   */
		rl.DrawText(icon, x+12, y+14, 22, NeonGreen)
		rl.DrawText(text, x+38, y+15, 20, NeonGreen)
	} else {
		rl.DrawRectangleLines(x, y, width, height, rl.Color{R: 60, G: 70, B: 80, A: 180})
		rl.DrawText(icon, x+12, y+14, 22, TextGray)
		rl.DrawText(text, x+38, y+15, 20, TextGray)
	}
}

func drawGameCard(x, y, width, height int32, scale float32, isSelected bool, game GameItem) {
	scaledWidth := int32(float32(width) * scale)
	scaledHeight := int32(float32(height) * scale)
	offsetX := (scaledWidth - width) / 2
	offsetY := (scaledHeight - height) / 2

	drawX := x - offsetX
	drawY := y - offsetY

	/*   Multi-layer glow for selected   */
	if isSelected {
		for g := int32(4); g >= 1; g-- {
			glowAlpha := uint8(20 * g)
			rl.DrawRectangle(drawX-g*5, drawY-g*5,
				scaledWidth+g*10, scaledHeight+g*10,
				rl.Color{R: game.Color.R, G: game.Color.G, B: game.Color.B, A: glowAlpha})
		}
	}

	/*   Shadow   */
	rl.DrawRectangle(drawX+5, drawY+5, scaledWidth, scaledHeight, AeroShadow)

	/*   Card background   */
	rl.DrawRectangleGradientV(drawX, drawY, scaledWidth, scaledHeight,
		rl.Color{R: 35, G: 45, B: 55, A: 255},
		rl.Color{R: 20, G: 28, B: 38, A: 255})

	/*   Top accent bar   */
	rl.DrawRectangle(drawX, drawY, scaledWidth, 5, game.Color)

	/*   Glass shine (Frutiger Aero)   */
	rl.DrawRectangleGradientV(drawX, drawY, scaledWidth, scaledHeight/3, AeroGloss, rl.Color{A: 0})
	rl.DrawRectangleGradientH(drawX, drawY, scaledWidth/3, scaledHeight/2, rl.Color{R: 255, G: 255, B: 255, A: 25}, rl.Color{A: 0})

	/*   Icon area   */
	iconCenterX := drawX + scaledWidth/2
	iconCenterY := drawY + scaledHeight/3 + 5
	iconRadius := float32(scaledWidth) / 4

	/*   Icon background with shape   */
	rl.DrawCircle(iconCenterX, iconCenterY, iconRadius+3,
		rl.Color{R: game.Color.R, G: game.Color.G, B: game.Color.B, A: 180})
	rl.DrawCircle(iconCenterX, iconCenterY, iconRadius,
		rl.Color{R: 25, G: 35, B: 45, A: 255})

	/*   Inner shape decoration   */
	switch game.IconShape {
	case 0:
		/*   Triangle   */
		rl.DrawTriangle(
			rl.Vector2{X: float32(iconCenterX), Y: float32(iconCenterY) - iconRadius*0.5},
			rl.Vector2{X: float32(iconCenterX) - iconRadius*0.4, Y: float32(iconCenterY) + iconRadius*0.3},
			rl.Vector2{X: float32(iconCenterX) + iconRadius*0.4, Y: float32(iconCenterY) + iconRadius*0.3},
			rl.Color{R: game.Color.R, G: game.Color.G, B: game.Color.B, A: 80})
	case 1:
		/*   Diamond   */
		rl.DrawRectanglePro(
			rl.Rectangle{X: float32(iconCenterX), Y: float32(iconCenterY), Width: iconRadius * 0.7, Height: iconRadius * 0.7},
			rl.Vector2{X: iconRadius * 0.35, Y: iconRadius * 0.35}, 45,
			rl.Color{R: game.Color.R, G: game.Color.G, B: game.Color.B, A: 80})
	case 2:
		/*   Ring   */
		rl.DrawRing(rl.Vector2{X: float32(iconCenterX), Y: float32(iconCenterY)},
			iconRadius*0.3, iconRadius*0.5, 0, 360, 20,
			rl.Color{R: game.Color.R, G: game.Color.G, B: game.Color.B, A: 80})
	case 3:
		/*   Star points   */
		for j := 0; j < 4; j++ {
			angle := float64(j) * math.Pi / 2
			px := float32(iconCenterX) + iconRadius*0.4*float32(math.Cos(angle))
			py := float32(iconCenterY) + iconRadius*0.4*float32(math.Sin(angle))
			rl.DrawCircle(int32(px), int32(py), iconRadius*0.12,
				rl.Color{R: game.Color.R, G: game.Color.G, B: game.Color.B, A: 100})
		}
	}

	/*   Icon letter   */
	textWidth := rl.MeasureText(game.IconText, 22)
	rl.DrawText(game.IconText, iconCenterX-textWidth/2, iconCenterY-10, 22, game.Color)

	/*   Icon highlight   */
	rl.DrawCircle(iconCenterX-int32(iconRadius*0.4), iconCenterY-int32(iconRadius*0.4), iconRadius*0.2, rl.Color{R: 255, G: 255, B: 255, A: 60})

	/*   Title   */
	titleY := drawY + scaledHeight - 55
	rl.DrawText(game.Title, drawX+10, titleY, 14, TextWhite)

	/*   Description   */
	descY := titleY + 18
	rl.DrawText(game.Description, drawX+10, descY, 10, TextGray)

	/*   Border   */
	borderColor := rl.Color{R: 50, G: 60, B: 70, A: 255}
	if isSelected {
		borderColor = game.Color
	}
	rl.DrawRectangleLines(drawX, drawY, scaledWidth, scaledHeight, borderColor)
}

func drawHeader() {
	rl.DrawRectangleGradientV(0, 0, SCREEN_WIDTH, 80,
		rl.Color{R: 8, G: 15, B: 25, A: 255},
		AeroBackground)

	/*   Glass shine   */
	rl.DrawRectangleGradientV(0, 0, SCREEN_WIDTH, 40, AeroGloss, rl.Color{A: 0})

	/*   Title centered   */
	title := "ARCADE PROJECT"
	titleWidth := rl.MeasureText(title, 30)
	titleX := (SCREEN_WIDTH - titleWidth) / 2

	rl.DrawText(title, titleX+2, 17, 30, rl.Color{R: 57, G: 255, B: 20, A: 40})
	rl.DrawText(title, titleX+1, 16, 30, NeonGreenGlow)
	rl.DrawText(title, titleX, 15, 30, NeonGreen)

	/*   Subtitle   */
	subtitle := "FRUTIGER AERO | Y2K EDITION"
	subWidth := rl.MeasureText(subtitle, 10)
	rl.DrawText(subtitle, (SCREEN_WIDTH-subWidth)/2, 48, 10, TextGray)

	/*   Date right side   */
	rl.DrawText("02.01.2026", SCREEN_WIDTH-110, 20, 16, AccentCyan)

	/*   Decorative line   */
	rl.DrawRectangle(100, 75, SCREEN_WIDTH-200, 1, rl.Color{R: NeonGreen.R, G: NeonGreen.G, B: NeonGreen.B, A: 50})
}

func drawFooter() {
	footerY := int32(SCREEN_HEIGHT - 40)

	rl.DrawRectangleGradientV(0, footerY, SCREEN_WIDTH, 40,
		rl.Color{A: 0}, rl.Color{R: 8, G: 15, B: 25, A: 220})

	/*   Decorative line   */
	rl.DrawRectangle(100, footerY+5, SCREEN_WIDTH-200, 1, rl.Color{R: NeonGreen.R, G: NeonGreen.G, B: NeonGreen.B, A: 50})

	/*   Control hints with arrow symbols   */
	controls := "<< >>  NAVIGATE    ^v  SELECT    [ENTER] LAUNCH    [ESC] EXIT"
	ctrlWidth := rl.MeasureText(controls, 12)
	rl.DrawText(controls, (SCREEN_WIDTH-ctrlWidth)/2, footerY+16, 12, TextGray)

	/*   FPS and version   */
	fps := fmt.Sprintf("%d FPS", rl.GetFPS())
	rl.DrawText(fps, SCREEN_WIDTH-100, footerY+16, 12, AccentCyan)
	rl.DrawText("v1.0", 50, footerY+16, 12, NeonGreenDark)
}

/**************************************************/
/*                                                */
/*             UTILITY FUNCTIONS                  */
/*                                                */
/**************************************************/

func lerp(a, b, t float32) float32 {
	return a + (b-a)*t
}

func easeOutCubic(t float32) float32 {
	return 1 - float32(math.Pow(float64(1-t), 3))
}

/**************************************************/
/*                                                */
/*              MAIN APPLICATION                  */
/*                                                */
/**************************************************/

func main() {
	rl.InitWindow(SCREEN_WIDTH, SCREEN_HEIGHT, "Arcade Project v1.0 - Frutiger Aero Y2K Edition")
	rl.SetTargetFPS(60)

	bladeNav := NewBladeNav()
	gameGrid := NewGameGrid()
	decorations := NewAeroDecorations()

	for !rl.WindowShouldClose() {
		globalTime += rl.GetFrameTime()

		/*   Input handling   */
		if rl.IsKeyPressed(rl.KeyRight) {
			bladeNav.NextBlade()
		}
		if rl.IsKeyPressed(rl.KeyLeft) {
			bladeNav.PrevBlade()
		}

		if bladeNav.TargetBlade == 0 {
			if rl.IsKeyPressed(rl.KeyDown) || rl.IsKeyPressed(rl.KeyS) {
				gameGrid.SelectNext()
			}
			if rl.IsKeyPressed(rl.KeyUp) || rl.IsKeyPressed(rl.KeyW) {
				gameGrid.SelectPrev()
			}
		}

		bladeNav.Update()
		gameGrid.Update()
		decorations.Update()

		/*   Drawing   */
		rl.BeginDrawing()

		rl.DrawRectangleGradientV(0, 0, SCREEN_WIDTH, SCREEN_HEIGHT,
			AeroBackground,
			rl.Color{R: 5, G: 12, B: 22, A: 255})

		decorations.Draw()
		drawHeader()
		bladeNav.Draw()

		if bladeNav.TargetBlade == 0 {
			gameGrid.Draw()
		} else {
			content := ""
			switch bladeNav.TargetBlade {
			case 1:
				content = "MEDIA CENTER - Coming Soon"
			case 2:
				content = "SYSTEM SETTINGS"
			case 3:
				content = "NETWORK STATUS"
			}
			contentWidth := rl.MeasureText(content, 26)
			rl.DrawText(content, (SCREEN_WIDTH-contentWidth)/2, 280, 26, TextGray)
		}

		drawFooter()

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
