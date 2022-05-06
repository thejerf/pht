package pht

// T is a shortcut function to set up a tag with no attributes.
func T(name string) *Tag {
	return &Tag{name, nil, nil, nil}
}

// TA is a shortcut function to set up a tag with possible initial
// attributes.
func TA(name string, atts map[string]string) *Tag {
	return &Tag{name, atts, nil, nil}
}

// A is a shortcut type for declaring attribute maps.
type A map[string]string

func C(content string) Content {
	return Content{content}
}

func S() Section {
	return Section{}
}

func Seq() *Sequence {
	return &Sequence{}
}

func PE(s string) *PreEscaped {
	return &PreEscaped{s, false}
}

func B(html ...HTML) *Block {
	return &Block{Components: html}
}
