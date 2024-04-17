package embedme

type Language string

var (
	PlainText  []Language = []Language{"txt", "embedme"}
	None                  = []Language{""}
	Typescript            = []Language{"ts"}
	Javascript            = []Language{"js"}
	Reason                = []Language{"re"}
	Scss                  = []Language{"scss"}
	Rust                  = []Language{"rust"}
	Java                  = []Language{"java"}
	Cpp                   = []Language{"cpp"}
	C                     = []Language{"c"}
	HTML                  = []Language{"html"}
	XML                   = []Language{"xml"}
	Markdown              = []Language{"md"}
	TOML                  = []Language{"toml"}
	YAML                  = []Language{"yaml"}
	JSON                  = []Language{"json"}
	JSON5                 = []Language{"json5"}
	Python                = []Language{"py", "python"}
	Bash                  = []Language{"bash"}
	Shell                 = []Language{"sh", "shell"}
	Golang                = []Language{"go", "golang"}
	ObjectiveC            = []Language{"objectivec"}
	PHP                   = []Language{"php"}
	CSharp                = []Language{"cs"}
	Swift                 = []Language{"swift"}
	Ruby                  = []Language{"rb"}
	Kotlin                = []Language{"kotlin"}
	Scala                 = []Language{"scala"}
	Crystal               = []Language{"cr"}
	PlantUML              = []Language{"puml"}
	Mermaid               = []Language{"mermaid"}
	Cmake                 = []Language{"cmake"}
	Protobuf              = []Language{"proto"}
	SQL                   = []Language{"sql"}
	Haskell               = []Language{"hs"}
	Arduino               = []Language{"ino"}
	Jsx                   = []Language{"jsx"}
	Tsx                   = []Language{"tsx"}
)

type CommentType uint32

const (
	CommentNone CommentType = iota
	CommentDoubleSlash
	CommentXML
	CommentHash
	CommentSingleQuote
	CommentDoublePercent
	CommentDoubleHyphens
)

var (
	LanguageComments = map[CommentType][]Language{
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
	SupportedLanguages = buildSupportedLanguages(LanguageComments)
	CommentForLanguage = buildCommentForLanguage(LanguageComments)
)

func buildSupportedLanguages(mapping map[CommentType][]Language) []string {
	var supported []string
	for _, languages := range mapping {
		for _, lang := range languages {
			supported = append(supported, string(lang))
		}
	}
	return supported
}

func buildCommentForLanguage(mapping map[CommentType][]Language) map[Language]CommentType {
	result := make(map[Language]CommentType)
	for comment, languages := range mapping {
		for _, lang := range languages {
			result[lang] = comment
		}
	}
	return result
}

// func CommentForLanguage(lang Language) CommentType {
//   for comment, languages := range LanguageComments {
//     if contains(languages, language) {
//     }
//   }

// func contains[T any](slice []T, item T) bool {
// 	for _, s := range slice {
//     if s == item {
//       return true
//     }
//   }
//   return false
// }

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
