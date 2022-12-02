package embedme

type Language string

var (
	PLAIN_TEXT  []Language = []Language{"txt", "embedme"}
	TYPESCRIPT             = []Language{"ts"}
	JAVASCRIPT             = []Language{"js"}
	REASON                 = []Language{"re"}
	SCSS                   = []Language{"scss"}
	RUST                   = []Language{"rust"}
	JAVA                   = []Language{"java"}
	CPP                    = []Language{"cpp"}
	C                      = []Language{"c"}
	HTML                   = []Language{"html"}
	XML                    = []Language{"xml"}
	MARKDOWN               = []Language{"md"}
	TOML                   = []Language{"toml"}
	YAML                   = []Language{"yaml"}
	JSON                   = []Language{"json"}
	JSON_5                 = []Language{"json5"}
  PYTHON                 = []Language{"py", "python"}
	BASH                   = []Language{"bash"}
	SHELL                  = []Language{"sh", "shell"}
	GOLANG                 = []Language{"go", "golang"}
	OBJECTIVE_C            = []Language{"objectivec"}
	PHP                    = []Language{"php"}
	C_SHARP                = []Language{"cs"}
	SWIFT                  = []Language{"swift"}
	RUBY                   = []Language{"rb"}
	KOTLIN                 = []Language{"kotlin"}
	SCALA                  = []Language{"scala"}
	CRYSTAL                = []Language{"cr"}
	PLANT_UML              = []Language{"puml"}
	MERMAID                = []Language{"mermaid"}
	CMAKE                  = []Language{"cmake"}
	PROTOBUF               = []Language{"proto"}
	SQL                    = []Language{"sql"}
	HASKELL                = []Language{"hs"}
	ARDUINO                = []Language{"ino"}
	JSX                    = []Language{"jsx"}
	TSX                    = []Language{"tsx"}
)

type CommentType uint32

const (
	COMMENT_NONE CommentType = iota
	COMMENT_DOUBLE_SLASH
	COMMENT_XML
	COMMENT_HASH
	COMMENT_SINGLE_QUOTE
	COMMENT_DOUBLE_PERCENT
	COMMENT_DOUBLE_HYPHENS
)

var (
	LanguageComments = map[CommentType][]Language{
		COMMENT_NONE: concat(JSON),
		COMMENT_DOUBLE_SLASH: concat(
			PLAIN_TEXT, // this is a lie, but we gotta pick something
			C,
			TYPESCRIPT,
			REASON,
			JAVASCRIPT,
			RUST,
			CPP,
			JAVA,
			GOLANG,
			OBJECTIVE_C,
			SCSS,
			PHP,
			C_SHARP,
			SWIFT,
			KOTLIN,
			SCALA,
			JSON_5,
			PROTOBUF,
			ARDUINO,
			JSX,
			TSX,
		),
		COMMENT_XML: concat(
			HTML,
			MARKDOWN,
			XML,
		),
		COMMENT_HASH: concat(
			PYTHON,
			BASH,
			SHELL,
			YAML,
			TOML,
			RUBY,
			CRYSTAL,
			CMAKE,
		),
		COMMENT_SINGLE_QUOTE: concat(
			PLANT_UML,
		),
		COMMENT_DOUBLE_PERCENT: concat(
			MERMAID,
		),
		COMMENT_DOUBLE_HYPHENS: concat(
			SQL,
			HASKELL,
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
