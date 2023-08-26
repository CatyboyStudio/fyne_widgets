package inspector

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type inspItem struct {
	id     int
	label  string
	data   any
	editor Editor
}

func (o inspItem) IsNil() bool {
	return o.editor == nil
}

var _ (fyne.Widget) = (*Inspector)(nil)

type Inspector struct {
	widget.BaseWidget

	bindid     int
	items      []inspItem
	form       *widget.Form
	toolbar    *fyne.Container
	backButton *widget.Button
	pathText   *widget.Label
	editor     Editor
}

func NewInspector() *Inspector {
	o := &Inspector{
		//
	}
	o.ExtendBaseWidget(o)
	return o
}

func (this *Inspector) onBack() {

}

func (this *Inspector) CreateRenderer() fyne.WidgetRenderer {
	this.backButton = widget.NewButtonWithIcon("", theme.NavigateBackIcon(), this.onBack)
	this.pathText = widget.NewLabel("")
	this.pathText.Wrapping = fyne.TextTruncate
	c1 := container.NewHBox(this.backButton, this.pathText)
	c2 := container.NewPadded(c1)
	this.toolbar = c2
	this.toolbar.Hide()
	this.form = widget.NewForm()
	co := container.NewVBox(c2, this.form)
	return widget.NewSimpleRenderer(co)
}

func (this *Inspector) current() inspItem {
	if len(this.items) > 0 {
		return this.items[len(this.items)-1]
	}
	return inspItem{}
}

func (this *Inspector) pushItem(item inspItem) {
	this.items = append(this.items, item)
}

func (this *Inspector) popItem() {
	c := len(this.items)
	if c > 0 {
		this.items[c-1] = inspItem{}
		this.items = this.items[:c-1]
	}
}

func (this *Inspector) showCurrent() {
	this.form.Items = nil

	item := this.current()
	if item.IsNil() {
		this.toolbar.Hide()
	} else {
		item.editor.CreateInspectorGUI(this.form, item.label)
		this.toolbar.Show()
		if len(this.items) > 1 {
			this.backButton.Enable()
		} else {
			this.backButton.Disable()
		}
	}
	this.form.Refresh()
}

func (this *Inspector) Bind(data any, editorType string) (int, error) {
	ed, err := CreateEditor(data, editorType)
	if err != nil {
		return 0, err
	}
	this.bindid++
	id := this.bindid
	item := inspItem{
		id:     id,
		data:   data,
		editor: ed,
	}
	c := len(this.items)
	for i := 0; i < c; i++ {
		this.popItem()
	}
	this.pushItem(item)
	this.showCurrent()
	return id, nil
}

func (this *Inspector) Unbind(id int) {
	idx := -1
	for i, v := range this.items {
		if v.id == id {
			idx = i
			break
		}
	}
	if idx == -1 {
		return
	}
	for i := len(this.items) - 1; i >= idx; i-- {
		this.popItem()
	}
	this.showCurrent()
}

func (this *Inspector) Push(label string, data any, editorType string) error {
	ed, err := CreateEditor(data, editorType)
	if err != nil {
		return err
	}
	this.bindid++
	id := this.bindid
	item := inspItem{
		id:     id,
		label:  label,
		data:   data,
		editor: ed,
	}
	this.pushItem(item)
	this.showCurrent()
	return nil
}
