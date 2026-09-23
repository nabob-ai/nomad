package dialog

// DefaultEnglishKeywords 返回默认英文关键词矩阵
func DefaultEnglishKeywords() KeywordMatrix {
	return KeywordMatrix{
		// General preferences
		CategoryPreference: {
			"want", "prefer", "like", "hope", "wish",
			"would like", "I'd like", "I want",
			"looking for", "need",
		},

		// Constraints
		CategoryConstraint: {
			"must", "have to", "need to",
			"don't", "do not", "never", "avoid",
			"shouldn't", "should not", "must not",
			"please don't", "no", "without",
		},

		// Style preferences
		CategoryStyle: {
			"style", "tone", "voice", "feel",
			"formal", "casual", "professional", "friendly",
			"humorous", "serious", "conversational",
			"concise", "detailed", "brief",
		},

		// Audience preferences
		CategoryAudience: {
			"reader", "audience", "target",
			"user", "customer", "viewer",
			"for whom", "intended for",
		},

		// Length preferences
		CategoryLength: {
			"words", "length", "characters",
			"short", "brief", "concise",
			"long", "detailed", "comprehensive",
			"keep it", "limit to", "around",
		},

		// Format preferences
		CategoryFormat: {
			"format", "layout", "structure",
			"list", "bullet", "numbered",
			"paragraph", "heading", "section",
			"markdown", "plain text", "html",
		},
	}
}

// TechnicalWritingKeywords 技术写作关键词
func TechnicalWritingKeywords() KeywordMatrix {
	base := DefaultEnglishKeywords()

	// Add technical writing specific keywords
	base[CategoryStyle] = append(base[CategoryStyle],
		"technical", "documentation",
		"tutorial", "guide", "reference",
		"step-by-step", "how-to",
	)

	base[CategoryFormat] = append(base[CategoryFormat],
		"code block", "snippet", "example",
		"api", "endpoint", "parameter",
	)

	return base
}

// MarketingContentKeywords 营销内容关键词
func MarketingContentKeywords() KeywordMatrix {
	base := DefaultEnglishKeywords()

	// Add marketing specific keywords
	base[CategoryStyle] = append(base[CategoryStyle],
		"persuasive", "engaging", "compelling",
		"call to action", "cta",
		"benefits", "features",
	)

	base[CategoryAudience] = append(base[CategoryAudience],
		"prospect", "lead", "buyer",
		"decision maker", "stakeholder",
	)

	return base
}
