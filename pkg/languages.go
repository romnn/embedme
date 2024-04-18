package embedme

// LanguageID describes a programming language
type LanguageID string

var (
	// PlainText file extensions
	PlainText = []LanguageID{"txt", "embedme"}
	// None contains no text file extensions
	None = []LanguageID{""}
	// Typescript file extensions
	Typescript = []LanguageID{"ts"}
	// Javascript file extensions
	Javascript = []LanguageID{"js"}
	// Reason file extensions
	Reason = []LanguageID{"re"}
	// Scss file extensions
	Scss = []LanguageID{"scss"}
	// Rust file extensions
	Rust = []LanguageID{"rust"}
	// Java file extensions
	Java = []LanguageID{"java"}
	// Cpp file extensions
	Cpp = []LanguageID{"cpp"}
	// C file extensions
	C = []LanguageID{"c"}
	// HTML file extensions
	HTML = []LanguageID{"html"}
	// XML file extensions
	XML = []LanguageID{"xml"}
	// Markdown file extensions
	Markdown = []LanguageID{"md"}
	// TOML file extensions
	TOML = []LanguageID{"toml"}
	// YAML file extensions
	YAML = []LanguageID{"yaml"}
	// JSON file extensions
	JSON = []LanguageID{"json"}
	// JSON5 file extensions
	JSON5 = []LanguageID{"json5"}
	// Python file extensions
	Python = []LanguageID{"py", "python"}
	// Bash file extensions
	Bash = []LanguageID{"bash"}
	// Shell file extensions
	Shell = []LanguageID{"sh", "shell"}
	// Golang file extensions
	Golang = []LanguageID{"go", "golang"}
	// ObjectiveC file extensions
	ObjectiveC = []LanguageID{"objectivec", "objc"}
	// PHP file extensions
	PHP = []LanguageID{"php"}
	// CSharp file extensions
	CSharp = []LanguageID{"cs"}
	// Swift file extensions
	Swift = []LanguageID{"swift"}
	// Ruby file extensions
	Ruby = []LanguageID{"rb"}
	// Kotlin file extensions
	Kotlin = []LanguageID{"kotlin"}
	// Scala file extensions
	Scala = []LanguageID{"scala"}
	// Crystal file extensions
	Crystal = []LanguageID{"cr"}
	// PlantUML file extensions
	PlantUML = []LanguageID{"puml"}
	// Mermaid file extensions
	Mermaid = []LanguageID{"mermaid"}
	// Cmake file extensions
	Cmake = []LanguageID{"cmake"}
	// Protobuf file extensions
	Protobuf = []LanguageID{"proto"}
	// SQL file extensions
	SQL = []LanguageID{"sql"}
	// Haskell file extensions
	Haskell = []LanguageID{"hs"}
	// Arduino file extensions
	Arduino = []LanguageID{"ino"}
	// Jsx file extensions
	Jsx = []LanguageID{"jsx"}
	// Tsx file extensions
	Tsx = []LanguageID{"tsx"}
)

// CommentType describes the type of comment used by a programming language
type CommentType uint32

const (
	// CommentNone refers to no comment
	CommentNone CommentType = iota
	// CommentDoubleSlash refers to "//" comment
	CommentDoubleSlash
	// CommentXML refers to "<!-- ... -->" comment
	CommentXML
	// CommentHash refers to "#" comment
	CommentHash
	// CommentSingleQuote refers to "'" comment
	CommentSingleQuote
	// CommentDoublePercent refers to "%%" comment
	CommentDoublePercent
	// CommentDoubleHyphens refers to "--" comment
	CommentDoubleHyphens
)

var (
	// LanguageComments ...
	LanguageComments = map[CommentType][]LanguageID{
		CommentNone: concat(JSON),
		CommentDoubleSlash: concat(
			None,      // we define no-language to use double slash
			PlainText, // we define plaintext to use double slash
			C,
			Typescript,
			Reason,
			Javascript,
			Rust,
			Cpp,
			Java,
			Golang,
			ObjectiveC,
			Scss,
			PHP,
			CSharp,
			Swift,
			Kotlin,
			Scala,
			JSON5,
			Protobuf,
			Arduino,
			Jsx,
			Tsx,
		),
		CommentXML: concat(
			HTML,
			Markdown,
			XML,
		),
		CommentHash: concat(
			Python,
			Bash,
			Shell,
			YAML,
			TOML,
			Ruby,
			Crystal,
			Cmake,
		),
		CommentSingleQuote: concat(
			PlantUML,
		),
		CommentDoublePercent: concat(
			Mermaid,
		),
		CommentDoubleHyphens: concat(
			SQL,
			Haskell,
		),
	}
	// SupportedLanguages ...
	SupportedLanguages = buildSupportedLanguages(LanguageComments)
	// CommentForLanguage ...
	CommentForLanguage = buildCommentForLanguage(LanguageComments)
)

func buildSupportedLanguages(mapping map[CommentType][]LanguageID) []string {
	var supported []string
	for _, languages := range mapping {
		for _, lang := range languages {
			supported = append(supported, string(lang))
		}
	}
	return supported
}

func buildCommentForLanguage(mapping map[CommentType][]LanguageID) map[LanguageID]CommentType {
	result := make(map[LanguageID]CommentType)
	for comment, languages := range mapping {
		for _, lang := range languages {
			result[lang] = comment
		}
	}
	return result
}

func concat[T any](slices ...[]T) []T {
	var totalLen int

	for _, s := range slices {
		totalLen += len(s)
	}

	result := make([]T, totalLen)

	var i int

	for _, s := range slices {
		i += copy(result[i:], s)
	}

	return result
}
