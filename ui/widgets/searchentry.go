package widgets

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/dweymouth/supersonic/ui/util"
)

// SearchEntry is a search entry widget that will issue a search command
// (aka call OnSearched) when a short span of time has elapsed since
// the user typed into the widget.
type SearchEntry struct {
	widget.Entry

	OnSearched func(string)
}

func NewSearchEntry() *SearchEntry {
	sf := &SearchEntry{}
	sf.ExtendBaseWidget(sf)
	sf.Init()
	return sf
}

// For use only by extending widgets
func (sf *SearchEntry) Init() {
	sf.PlaceHolder = lang.L("Search")
	sf.ActionItem = NewClearTextButton(func() {
		sf.SetText("")
	})
	debounceFunc := util.NewDebouncer(200*time.Millisecond, func() {
		sf.sendSearch(sf.Entry.Text)
	})
	sf.Entry.OnChanged = func(_ string) {
		if sf.updateActionButton() {
			sf.ActionItem.Refresh()
		}
		debounceFunc()
	}
}

func (s *SearchEntry) TypedKey(e *fyne.KeyEvent) {
	if e.Name == fyne.KeyEscape {
		s.SetText("")
		fyne.CurrentApp().Driver().CanvasForObject(s).Unfocus()
		return
	}
	s.Entry.TypedKey(e)
}

func (s *SearchEntry) Refresh() {
	s.updateActionButton()
	s.Entry.Refresh()
}

func (s *SearchEntry) MinSize() fyne.Size {
	return fyne.NewSize(200, s.Entry.MinSize().Height)
}

func (s *SearchEntry) updateActionButton() bool {
	btn := s.ActionItem.(*searchActionButton)
	oldResouce := btn.Resource
	if s.Text == "" {
		btn.Resource = theme.SearchIcon()
	} else {
		btn.Resource = theme.ContentClearIcon()
	}
	return oldResouce != btn.Resource
}

var _ fyne.Tappable = (*searchActionButton)(nil)

type searchActionButton struct {
	widget.BaseWidget

	Resource fyne.Resource

	icon *widget.Icon

	OnTapped func()
}

func (s *SearchEntry) sendSearch(text string) {
	if s.OnSearched != nil {
		s.OnSearched(text)
	}
}

func NewClearTextButton(onTapped func()) *searchActionButton {
	c := &searchActionButton{OnTapped: onTapped}
	c.ExtendBaseWidget(c)
	c.Resource = theme.SearchIcon()
	c.icon = widget.NewIcon(c.Resource)
	return c
}

func (c *searchActionButton) Tapped(*fyne.PointEvent) {
	if c.OnTapped != nil {
		c.OnTapped()
	}
}

func (c *searchActionButton) MinSize() fyne.Size {
	th := c.Theme()
	return fyne.NewSquareSize(th.Size(theme.SizeNameInlineIcon) + th.Size(theme.SizeNameInnerPadding)*2)
}

func (c *searchActionButton) Refresh() {
	c.icon.Resource = c.Resource
	c.BaseWidget.Refresh()
}

func (c *searchActionButton) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewCenter(c.icon))
}
