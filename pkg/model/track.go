package model

type Track struct {
	Title    string
	Position uint32
}

type TrackResult struct {
	Track
	Artist   string
	Album    string
	CoverUrl string
}
