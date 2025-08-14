package processor

type FrontMatter struct {
	Title    string `yaml:"title"`
	Template string `yaml:"template"`
}

const (
	FrontMatterDelimiter = "---\n"
	DefaultTemplate      = "index"
)