package dialog

// DefaultChineseKeywords 返回默认中文关键词矩阵
func DefaultChineseKeywords() KeywordMatrix {
	return KeywordMatrix{
		// 一般偏好
		CategoryPreference: {
			"想要", "希望", "喜欢", "偏好", "倾向", "更喜欢",
			"觉得", "认为", "我想", "我要",
		},

		// 约束条件
		CategoryConstraint: {
			"必须", "一定要", "不要", "禁止", "不能", "避免", "不允许",
			"请勿", "严禁", "切勿", "务必",
		},

		// 风格偏好
		CategoryStyle: {
			"风格", "语气", "调性", "口吻", "感觉",
			"正式", "轻松", "幽默", "严肃", "专业",
			"口语化", "书面语", "简洁", "详细",
		},

		// 受众偏好
		CategoryAudience: {
			"读者", "受众", "目标人群", "给谁看",
			"用户", "客户", "观众",
		},

		// 长度偏好
		CategoryLength: {
			"字数", "长度", "篇幅", "多少字",
			"简短", "简洁", "详细", "展开",
			"不要太长", "控制在",
		},

		// 格式偏好
		CategoryFormat: {
			"格式", "排版", "布局",
			"列表", "段落", "标题", "编号",
			"markdown", "纯文本",
		},
	}
}

// WritingStyleKeywords 写作风格关键词（扩展）
func WritingStyleKeywords() KeywordMatrix {
	base := DefaultChineseKeywords()

	// 添加写作特定关键词
	base[CategoryStyle] = append(base[CategoryStyle],
		"文风", "笔触", "叙述方式",
		"第一人称", "第三人称",
		"对话多", "描写多", "议论多",
		"故事性", "说明性", "抒情",
	)

	base[CategoryPreference] = append(base[CategoryPreference],
		"主题", "题材", "方向",
		"重点", "核心", "关键",
		"突出", "强调", "侧重",
	)

	return base
}

// ContentPlatformKeywords 内容平台关键词
func ContentPlatformKeywords() KeywordMatrix {
	base := DefaultChineseKeywords()

	// 添加平台特定关键词
	platformKeywords := []string{
		"微信", "公众号", "小红书", "抖音", "知乎",
		"微博", "B站", "头条", "百家号",
		"豆瓣", "简书", "掘金", "CSDN",
	}

	base[CategoryFormat] = append(base[CategoryFormat], platformKeywords...)

	return base
}
