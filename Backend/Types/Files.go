package Types

var ExtensionsToMimeType map[string]string = map[string]string{

	".html": "text/html",
	".css":  "text/css",
	".js":   "text/javascript",

	".json": "application/json",
	".xml":  "application/xml",

	".txt": "text/plain",
	".csv": "text/csv",

	".png":  "image/png",
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".gif":  "image/gif",
	".bmp":  "image/bmp",
	".webp": "image/webp",
	".svg":  "image/svg+xml",

	".pdf": "application/pdf",
	".zip": "application/zip",
	".rar": "application/vnd.rar",

	".mp3": "audio/mpeg",
	".wav": "audio/wav",

	".mp4": "video/mp4",
	".avi": "video/x-msvideo",
	".mov": "video/quicktime",
	".mkv": "video/x-matroska",

	".exe": "application/vnd.microsoft.portable-executable",
	".dll": "application/vnd.microsoft.portable-executable",
	".sh":  "application/x-sh",
	".bat": "application/x-msdos-program",

	".go":   "text/x-go",
	".py":   "text/x-python",
	".java": "text/x-java-source",
	".c":    "text/x-c",
	".cpp":  "text/x-c++",
	".h":    "text/x-c",
	".hpp":  "text/x-c++",
	".rb":   "text/x-ruby",

	".php": "application/x-httpd-php",

	".swift": "text/x-swift",
	".rs":    "text/x-rust",

	".ts":  "application/typescript",
	".tsx": "application/typescript",

	".jsx": "text/jsx",
	".vue": "text/x-vue",
	".md":  "text/markdown",
}

var MimeTypeToIcon map[string]string = map[string]string{

	"text/html":       "document-text-outline",
	"text/css":        "color-palette-outline",
	"text/javascript": "code-slash-outline",

	"application/json": "document-outline",
	"application/xml":  "document-outline",

	"text/plain": "document-text-outline",
	"text/csv":   "document-text-outline",

	"image/png":     "image-outline",
	"image/jpeg":    "image-outline",
	"image/gif":     "image-outline",
	"image/bmp":     "image-outline",
	"image/webp":    "image-outline",
	"image/svg+xml": "image-outline",

	"application/pdf": "document-outline",

	"application/zip":     "archive-outline",
	"application/vnd.rar": "archive-outline",

	"audio/mpeg": "musical-notes-outline",
	"audio/wav":  "musical-notes-outline",

	"video/mp4":        "videocam-outline",
	"video/x-msvideo":  "videocam-outline",
	"video/quicktime":  "videocam-outline",
	"video/x-matroska": "videocam-outline",

	"application/vnd.microsoft.portable-executable": "terminal-outline",

	"application/x-sh":            "terminal-outline",
	"application/x-msdos-program": "terminal-outline",

	"text/x-go":               "code-slash-outline",
	"text/x-python":           "code-slash-outline",
	"text/x-java-source":      "code-slash-outline",
	"text/x-c":                "code-slash-outline",
	"text/x-c++":              "code-slash-outline",
	"text/x-ruby":             "code-slash-outline",
	"application/x-httpd-php": "code-slash-outline",
	"text/x-swift":            "code-slash-outline",
	"text/x-rust":             "code-slash-outline",
	"application/typescript":  "code-slash-outline",
	"text/jsx":                "code-slash-outline",
	"text/x-vue":              "code-slash-outline",
	"text/markdown":           "document-text-outline",
}

var MimeTypeToReadableName map[string]string = map[string]string{

	"text/html":       "HTML Code",
	"text/css":        "CSS Code",
	"text/javascript": "JavaScript Code",

	"application/json": "JSON Source Code",
	"application/xml":  "XML Source Code",
	"text/plain":       "Plain Text",
	"text/csv":         "CSV Document",

	"image/png":     "PNG Image",
	"image/jpeg":    "JPEG Image",
	"image/gif":     "GIF Image",
	"image/bmp":     "BMP Image",
	"image/webp":    "WebP Image",
	"image/svg+xml": "SVG Image",

	"application/pdf": "PDF Document",

	"application/zip":     "ZIP Archive",
	"application/vnd.rar": "RAR Archive",

	"audio/mpeg": "MP3 Audio",
	"audio/wav":  "WAV Audio",

	"video/mp4":        "MP4 Video",
	"video/x-msvideo":  "AVI Video",
	"video/quicktime":  "MOV Video",
	"video/x-matroska": "MKV Video",

	"application/vnd.microsoft.portable-executable": "Executable File",

	"application/x-sh":            "Shell Script",
	"application/x-msdos-program": "Batch File",
	"text/x-go":                   "Go Source Code",
	"text/x-python":               "Python Source Code",
	"text/x-java-source":          "Java Source Code",
	"text/x-c":                    "C Source Code",
	"text/x-c++":                  "C++ Source Code",
	"text/x-ruby":                 "Ruby Source Code",
	"application/x-httpd-php":     "PHP Source Code",
	"text/x-swift":                "Swift Source Code",
	"text/x-rust":                 "Rust Source Code",

	"application/typescript": "TypeScript Code",
	"text/jsx":               "JSX Code",
	"text/x-vue":             "Vue.js Code",
	"text/markdown":          "Markdown Source",
}
