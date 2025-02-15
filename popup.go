package nanogui

import (
	"fmt"
	"github.com/shibukawa/nanovgo"
)

type Popup struct {
	Window
	parentWindow IWindow
	anchorX      int
	anchorY      int
	anchorHeight int
	vScroll      *VScrollPanel
	panel        Widget
}

func NewPopup(parent Widget, parentWindow IWindow) *Popup {
	popup := &Popup{
		parentWindow: parentWindow,
		anchorHeight: 30,
	}
	InitWidget(popup, parent)
	popup.vScroll = NewVScrollPanel(popup)
	popup.panel = NewVScrollPanelChild(popup.vScroll)
	return popup
}

// SetAnchorPosition() sets the anchor position in the parent window; the placement of the popup is relative to it
func (p *Popup) SetAnchorPosition(x, y int) {
	p.anchorX = x
	p.anchorY = y
}

// AnchorPosition() returns the anchor position in the parent window; the placement of the popup is relative to it
func (p *Popup) AnchorPosition() (int, int) {
	return p.anchorX, p.anchorY
}

// SetAnchorHeight() sets the anchor height; this determines the vertical shift relative to the anchor position
func (p *Popup) SetAnchorHeight(h int) {
	p.anchorHeight = h
}

// AnchorHeight() returns the anchor height; this determines the vertical shift relative to the anchor position
func (p *Popup) AnchorHeight() int {
	return p.anchorHeight
}

// SetParentWindow() sets the parent window of the popup
func (p *Popup) SetParentWindow(w *Window) {
	p.parentWindow = w
}

// ParentWindow() returns the parent window of the popup
func (p *Popup) ParentWindow() IWindow {
	return p.parentWindow
}

func (p *Popup) OnPerformLayout(self Widget, ctx *nanovgo.Context) {
	if p.layout != nil || len(p.children) != 1 {
		p.WidgetImplement.OnPerformLayout(self, ctx)
	} else {
		p.children[0].SetPosition(0, 0)
		p.children[0].SetSize(p.WidgetWidth, p.WidgetHeight)
		p.children[0].OnPerformLayout(p.children[0], ctx)
	}
}

func (p *Popup) IsPositionAbsolute() bool {
	return true
}

func (p *Popup) Draw(self Widget, ctx *nanovgo.Context) {
	p.RefreshRelativePlacement()

	if !p.visible {
		return
	}
	ds := float32(p.theme.WindowDropShadowSize)
	cr := float32(p.theme.WindowCornerRadius)

	px := float32(p.WidgetPosX)
	py := float32(p.WidgetPosY)
	pw := float32(p.WidgetWidth)
	ph := float32(p.WidgetHeight)
	ah := float32(p.anchorHeight)

	/* Draw a drop shadow */
	shadowPaint := nanovgo.BoxGradient(px, py, pw, ph, cr*2, ds*2, p.theme.DropShadow, p.theme.Transparent)
	ctx.BeginPath()
	ctx.Rect(px-ds, py-ds, pw+ds*2, ph+ds*2)
	ctx.RoundedRect(px, py, pw, ph, cr)
	ctx.PathWinding(nanovgo.Hole)
	ctx.SetFillPaint(shadowPaint)
	ctx.Fill()

	/* Draw window */
	ctx.BeginPath()
	ctx.RoundedRect(px, py, pw, ph, cr)

	ctx.MoveTo(px-15, py+ah)
	ctx.LineTo(px+1, py+ah-15)
	ctx.LineTo(px+1, py+ah+15)

	ctx.SetFillColor(p.theme.WindowPopup)

	ctx.Fill()

	p.WidgetImplement.Draw(self, ctx)
}

// RefreshRelativePlacement is internal helper function to maintain nested window position values; overridden in \ref Popup
func (p *Popup) RefreshRelativePlacement() {
	p.parentWindow.RefreshRelativePlacement()
	p.visible = p.visible && p.parentWindow.VisibleRecursive()
	x, y := p.parentWindow.Position()
	p.WidgetPosX = x + p.anchorX
	p.WidgetPosY = y + p.anchorY - p.anchorHeight
}

func (p *Popup) FindWindow() IWindow {
	return p
}

func (p *Popup) String() string {
	return p.StringHelper(fmt.Sprintf("Popup(%d)", p.Depth()), "")
}
