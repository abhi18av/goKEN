package main


type goKEN{

TalkMeta Meta
TalkTranscriptPage TranscriptPage
TalkVideoPage VideoPage

}
type Meta struct {
	DateProcessed []string
	type ItemsUpdated struct{
   	TranscriptPage
      	VideoPage
	}
}

type TranscriptPage struct {
	AvailableTranscripts []string
	DatePosted           string
	LocalTitle           string
	Rated                string
	TalkTranscript       []string {
		LocalTalkTitle string
		Paragraphs     []string
	}
	TimeStamps []string
}

type VideoPage struct {
	AvailableSubtitlesCount string   `json:"AvailableSubtitlesCount"`
	Speaker                 string   `json:"Speaker"`
	Duration                string   `json:"Duration"`
	TimeFilmed              string   `json:"TimeFilmed"`
	TalkViewsCount          string   `json:"TalkViewsCount"`
	TalkTopicsList          []string `json:"TalkTopicsList"`
	TalkCommentsCount       string   `json:"TalkCommentsCount"`
}
