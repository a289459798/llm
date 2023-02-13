// Code generated by goctl. DO NOT EDIT.
package types

import "net/http"

type LoginRequest struct {
}

type LoginResponse struct {
}

type InfoRequest struct {
}

type InfoResponse struct {
}

type ReportRequest struct {
	Content string `json:"content"`
}

type ReportResponse struct {
	Data string `json:"data"`
}

type ImageRequest struct {
	Content string `json:"content"`
}

type ImageResponse struct {
	Url string `json:"url"`
}

type WatermarkRequest struct {
	Url      string  `json:"url"`
	Content  string  `json:"content"`
	Position string  `json:"position"`
	Opacity  float32 `json:"opacity"`
	FontSize uint    `json:"fontSize"`
	Color    string  `json:"color"`
}

type StreamResponse struct {
	Writer http.ResponseWriter `json:"writer"`
}
