package actions

// JWT struct
type JWT struct {
	Secret string
}

// Proxy struct
type Proxy struct {
	OutboundEndpoint string
	KeepOrigin       bool
}

// ResponseTransform struct
type ResponseTransform struct {
	ModifyStatus int
	SetHeaders   map[string]string
	ReformatBody map[string]string
}

// Transform struct
type Transform struct {
	Append    map[string]string
	Duplicate map[string]string
}

// RequestTransform struct
type RequestTransform struct {
	Body  *Transform
	Query *Transform
}

// RateLimit struct
type RateLimit struct {
	Max      int
	WindowMs int
}
