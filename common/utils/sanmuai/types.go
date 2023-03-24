package sanmuai

type ImageCreate struct {
	Prompt         string `json:"prompt,omitempty"`
	N              int    `json:"n,omitempty"`
	Size           string `json:"size,omitempty"`
	ResponseFormat string `json:"response_format,omitempty"`
}

type ImageRepair struct {
	Image string `json:"image"`
}

type Image2Text struct {
	Image string `json:"image"`
}
