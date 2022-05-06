/*

Package pht provides programmer-oriented HTML templating, because
everything else frustrates me.

*/
package pht

import (
	"html/template"
	"io"
	"strings"
)

// FIXME: To get that last level of beautification where we wrap the
// attributes nicely, we need to change from writing to an io.Writer to an
// HTMLWriter, which includes the ability to write things out as words
// directly, so we can write entire attributes out as words for the
// WordWrapper.
//
// In the meantime, this is enough for now.

const (
	// The target width the word wrapping algorithm will shoot for
	TargetWidth = 80

	// The minimum width the word wrapping algorithm will permit
	MinimumWidth = 40
)

var v = struct{}{}

var tagsToNotClose = map[string]struct{}{
	"meta":   v,
	"option": v,
	"input":  v,
}

type Tag struct {
	Name       string
	Attributes map[string]string
	Content    []HTML
	NamedTags  map[string]*Tag
}

func (t *Tag) Attr(name, value string) *Tag {
	if t.Attributes == nil {
		t.Attributes = map[string]string{}
	}
	t.Attributes[name] = value
	return t
}

func (t *Tag) Append(html HTML) HTML {
	t.Content = append(t.Content, html)
	return html
}

func (t *Tag) AppendTag(tag *Tag) *Tag {
	t.Content = append(t.Content, tag)
	return tag
}

func (t *Tag) AppendNamedTag(name string, tag *Tag) *Tag {
	t.Content = append(t.Content, tag)
	if t.NamedTags == nil {
		t.NamedTags = map[string]*Tag{}
	}
	t.NamedTags[name] = tag
	return tag
}

func (t *Tag) AppendMany(html ...HTML) {
	t.Content = append(t.Content, html...)
}

func (t *Tag) AddClass(cls ...string) *Tag {
	toAdd := strings.Join(cls, " ")
	if t.Attributes == nil {
		t.Attributes = map[string]string{}
	}
	prevClass := t.Attributes["class"]
	if prevClass == "" {
		t.Attributes["class"] = toAdd
	} else {
		t.Attributes["class"] += " " + toAdd
	}
	return t
}

func (t *Tag) AddStyle(styles ...string) *Tag {
	toAdd := strings.Join(styles, " ")
	if t.Attributes == nil {
		t.Attributes = map[string]string{}
	}
	prevClass := t.Attributes["style"]
	if prevClass == "" {
		t.Attributes["style"] = toAdd
	} else {
		t.Attributes["style"] += " " + toAdd
	}
	return t
}

func first(strs ...string) string {
	for _, s := range strs {
		if s != "" {
			return s
		}
	}
	return ""
}

func writeHTMLIndent(indent int, w io.Writer) {
	if indent <= 0 {
		return
	}
	b := make([]byte, indent*2)
	for i := 0; i < indent*2; i++ {
		b[i] = 32
	}
	w.Write(b)
}

func (t *Tag) On(name string) *Tag {
	if t.NamedTags == nil {
		return nil
	}
	return t.NamedTags[name]
}

func (t *Tag) Render(indent int, w io.Writer) error {
	writeHTMLIndent(indent, w)

	_, _ = w.Write([]byte("<"))
	_, _ = w.Write([]byte(t.Name))
	if len(t.Attributes) > 0 {
		for key, val := range t.Attributes {
			_, _ = w.Write([]byte(" "))
			_, _ = w.Write([]byte(key))
			_, _ = w.Write([]byte("=\""))
			template.HTMLEscape(w, []byte(val))
			_, _ = w.Write([]byte("\""))
		}
	}
	_, err := w.Write([]byte(">"))

	haveIndentableContent := len(t.Content) > 0

	if haveIndentableContent {
		if indent >= 0 {
			_, _ = w.Write([]byte("\n"))
		}

		for _, html := range t.Content {
			var newIndent = indent + 1
			if indent < 0 {
				newIndent = indent
			}
			err = html.Render(newIndent, w)
			if err != nil {
				return err
			}
		}
	}

	if _, noClose := tagsToNotClose[t.Name]; !noClose {
		if haveIndentableContent {
			writeHTMLIndent(indent, w)
		}
		_, _ = w.Write([]byte("</"))
		_, _ = w.Write([]byte(t.Name))
		_, _ = w.Write([]byte(">"))
	}
	if indent >= 0 {
		_, err = w.Write([]byte("\n"))
	}
	return err
}

// Block content renders everything inline as a paragraph would be
// rendered, using the outer indent for itself but causing everything
// internal to bypass it.
type Block struct {
	Components []HTML
	tags       map[string]*Tag
}

func (b *Block) Append(html ...HTML) {
	b.Components = append(b.Components, html...)
}

func (b *Block) AppendNamed(name string, t *Tag) *Tag {
	if b.tags == nil {
		b.tags = map[string]*Tag{}
	}
	b.tags[name] = t
	b.Components = append(b.Components, t)
	return t
}

func (b *Block) On(s string) *Tag {
	if b.tags != nil {
		return b.tags[s]
	}
	return nil
}

func (b *Block) Render(indent int, w io.Writer) error {
	myWidth := TargetWidth - indent
	if myWidth < MinimumWidth {
		myWidth = MinimumWidth
	}

	wrapper := NewWrapper(w, indent*2, indent*2, myWidth)
	for _, html := range b.Components {
		html.Render(-1, wrapper)
	}
	wrapper.Close()
	_, err := w.Write([]byte("\n"))
	return err
}

type Content struct {
	Content string
}

func (c Content) Render(indent int, w io.Writer) error {
	if indent >= 0 {
		writeHTMLIndent(indent, w)
	}
	template.HTMLEscape(w, []byte(c.Content))
	if indent >= 0 {
		_, _ = w.Write([]byte("\n"))
	}
	return nil
}

func (c Content) On(s string) *Tag {
	return nil
}

type PreEscaped struct {
	Content string
	indent  bool
}

func (pe *PreEscaped) Render(indent int, w io.Writer) error {
	if pe.indent {
		writeHTMLIndent(indent, w)
	}
	_, err := w.Write([]byte(pe.Content))
	return err
}

func (pe *PreEscaped) Indent() *PreEscaped {
	pe.indent = true
	return pe
}

func (pe *PreEscaped) On(s string) *Tag {
	return nil
}

type Section struct{}

func (s Section) Render(indent int, w io.Writer) error {
	_, err := w.Write([]byte("\n"))
	return err
}

func (s Section) On(_ string) *Tag {
	return nil
}

// A Sequence simply sequences some HTML into a single object, without
// introducing any tags.
type Sequence struct {
	html []HTML
	tags map[string]*Tag
}

func (s *Sequence) Append(html ...HTML) {
	s.html = append(s.html, html...)
}

func (s *Sequence) AppendNamedTag(name string, t *Tag) *Tag {
	s.html = append(s.html, t)
	if s.tags == nil {
		s.tags = map[string]*Tag{}
	}
	s.tags[name] = t
	return t
}

func (s *Sequence) AppendTag(t *Tag) *Tag {
	s.html = append(s.html, t)
	return t
}

func (s *Sequence) Render(indent int, w io.Writer) error {
	var err error
	for _, html := range s.html {
		err = html.Render(indent, w)
	}
	return err
}

func (s *Sequence) On(tag string) *Tag {
	if s.tags != nil {
		return s.tags[tag]
	}
	return nil
}

func (s *Sequence) Register(name string, t *Tag) {
	if s.tags == nil {
		s.tags = map[string]*Tag{}
	}
	s.tags[name] = t
}

type ClosureHTML func(indent int, w io.Writer) error

func (chtml ClosureHTML) Render(indent int, w io.Writer) error {
	return chtml(indent, w)
}

func (chtml ClosureHTML) On(_ string) *Tag {
	return nil
}

func GridTextInput(id string,
	placeholder string,
	enterkeyhint string,
	value string,
) HTML {
	return ClosureHTML(func(indent int, w io.Writer) error {
		input := Tag{Name: "input"}
		input.Attr("type", "text").
			Attr("id", id).
			Attr("placeholder", placeholder).
			Attr("list", "service_options").
			Attr("autocorrect", "off").
			Attr("enterkeyhint", first(enterkeyhint, "next")).
			Attr("value", value)

		return input.Render(indent, w)
	})
}

// HTML is something that knows how to render itself to HTML.
type HTML interface {
	Render(indent int, w io.Writer) error

	// This yield up particular internal tags for further
	// modification. It is legal and frequent to return nil, so it
	// should be done carefully. But often you have static guarantees
	// it'll be valid.
	On(string) *Tag
}

// TagReferences is a convenient way to implement non-trivial Inner
// references.
type TagReferences struct {
	HTML
	Tags map[string]*Tag
}

func (tr TagReferences) On(s string) *Tag {
	return tr.Tags[s]
}
