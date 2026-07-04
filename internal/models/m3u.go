package models

type Playlist struct {
	Channels []Channel
}

type Channel struct {
	Length     int64  `json:"length"`
	TvgID      string `json:"tvg_id,omitempty"`
	TvgName    string `json:"tvg_name,omitempty"`
	TvgLogo    string `json:"tvg_logo,omitempty"`
	GroupTitle string `json:"group_title,omitempty"`
	Name       string `json:"name"`
	Uri        string `json:"uri"`
}
