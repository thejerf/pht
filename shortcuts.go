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

// C is a shortcut for creating conten.
func C(content string) Content {
	return Content{content}
}

// S is a shortcut for a section.
func S() Section {
	return Section{}
}

// Seq is a shortcut for a sequence.
func Seq(html ...HTML) *Sequence {
	return &Sequence{html: html}
}

// PE is a shortcut for creating a pre-escaped section.
func PE(s string) *PreEscaped {
	return &PreEscaped{s, false}
}

// B is a shortcut for creating a block.
func B(html ...HTML) *Block {
	return &Block{Components: html}
}
