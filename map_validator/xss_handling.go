package map_validator

import (
	"regexp"
)

func isPossibleXSS(input string) bool {
	xssPatterns := []string{
		"(?i)<script.*?>.*?</script>",
		"(?i)<iframe.*?>.*?</iframe>",
		"(?i)<object.*?>.*?</object>",
		"(?i)<embed.*?>.*?</embed>",
		"(?i)<svg.*?>.*?</svg>",
		"(?i)<img.*?src.*?=.*?>",
		"(?i)<a.*?href=\"javascript:.*?>.*?</a>",
		"(?i)onerror",
		"(?i)onload",
		"(?i)alert\\(",
		"[\\s\"'`;\\/0-9]\\b(alert|prompt|confirm|fetch|eval|new Function)\\s*\\(",
		"(?i)style\\s*=.*?expression|url\\s*\\(",
		"(?i)style\\s*=['\"].*?behavior:\\s*url",
		"(?i)<style.*?>.*?</style>",
		"(?i)@import\\s+['\"].*?;",
		"(?i)<meta.*?refresh.*?>",
		"(?i)<link.*?rel=['\"]?stylesheet['\"]?.*?>",
		"(?i)(base64|eval|atob|btoa|decodeURIComponent|encodeURIComponent)\\(",
		"(?i)\\.innerhtml\\s*=\\s*['\"].*?<script",
		"(?i)\\.outerhtml\\s*=\\s*['\"].*?<script",
		"(?i)\\.addEventListener\\s*\\(",
		"(?i)<canvas.*?>.*?</canvas>",
		"(?i)\\.getContext\\s*\\(\\s*['\"]webgl['\"]\\s*\\)",
		"(?i)\\.shaderSource\\s*\\(.*?,",
		"(?i)<template.*?>.*?</template>",
		"(?i)String\\.fromCharCode\\s*\\(.*?\\)",
		"(?i)eval\\s*\\(['\"].*?['\"]\\s*\\)",
		"(?i)this\\.alert\\s*\\(['\"].*?['\"]\\)",
		"(?i)window\\.name",
		"(?i)#<script>eval\\(name\\)</script>",
		"(?i)<div style=\".*?:\\s*expression\\(.*?\\);?\">",
		"(?i)<div style=\".*?:\\s*url\\(.*?\\);?\">",
		"(?i)<div style=\".*?:\\s*attr\\(.*?\\);?\">",
		"(?i)<div style=\".*?:\\s*behavior\\(.*?\\);?\">",
		"(?i)<div style=\".*?:\\s*expression\\(.*?\\);?\">",
		"(?i)<div style=\".*?:\\s*moz-binding\\(.*?\\);?\">",
		"(?i)<div style=\".*?:\\s*moz-xxx\\(.*?\\);?\">",
		"(?i)<div style=\".*?:\\s*webkit-xxx\\(.*?\\);?\">",
	}

	for _, pattern := range xssPatterns {
		match, err := regexp.MatchString(pattern, input)
		if err != nil {
			return true
		}
		if match {
			return true
		}
	}

	return false
}
