package main

import (
	"io"
	"net/http"
	"os"
	"strings"
	"sync"

	"strconv"

	"encoding/json"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
	"github.com/imdario/mergo"
)

type TedTalk struct {
	TalkVideoPage      VideoPage      `json:"TalkVideoPage"`
	TalkTranscriptPage TranscriptPage `json:"TalkTranscriptPage"`
}

type VideoPage struct {
	TalkURL                 string   `json:"VideoURL"`
	AvailableSubtitlesCount string   `json:"AvailableSubtitlesCount"`
	Speaker                 string   `json:"Speaker"`
	Duration                string   `json:"Duration"`
	TimeFilmed              string   `json:"TimeFilmed"`
	TalkViewsCount          string   `json:"TalkViewsCount"`
	TalkTopicsList          []string `json:"TalkTopicsList"`
	TalkCommentsCount       string   `json:"TalkCommentsCount"`
}

var langCodes = map[string]string{
	"Afrikaans":             "af",
	"Albanian":              "sq",
	"Algerian Arabic":       "arq",
	"Amharic":               "am",
	"Arabic":                "ar",
	"Armenian":              "hy",
	"Assamese":              "as",
	"Asturian":              "ast",
	"Azerbaijani":           "az",
	"Basque":                "eu",
	"Belarusian":            "be",
	"Bengali":               "bn",
	"Bislama":               "bi",
	"Bosnian":               "bs",
	"Bulgarian":             "bg",
	"Burmese":               "my",
	"Catalan":               "ca",
	"Cebuano":               "ceb",
	"Chinese, Simplified":   "zh-cn",
	"Chinese, Traditional":  "zh-tw",
	"Chinese, Yue":          "zh",
	"Creole, Haitian":       "ht",
	"Croatian":              "hr",
	"Czech":                 "cs",
	"Danish":                "da",
	"Dutch":                 "nl",
	"Dzongkha":              "dz",
	"English":               "en",
	"Esperanto":             "eo",
	"Estonian":              "et",
	"Filipino":              "fil",
	"Finnish":               "fi",
	"French":                "fr",
	"French (Canada)":       "fr-ca",
	"Galician":              "gl",
	"Georgian":              "ka",
	"German":                "de",
	"Greek":                 "el",
	"Gujarati":              "gu",
	"Hakha Chin":            "cnh",
	"Hausa":                 "ha",
	"Hebrew":                "he",
	"Hindi":                 "hi",
	"Hungarian":             "hu",
	"Hupa":                  "hup",
	"Icelandic":             "is",
	"Igbo":                  "ig",
	"Indonesian":            "id",
	"Ingush":                "inh",
	"Irish":                 "ga",
	"Italian":               "it",
	"Japanese":              "ja",
	"Kannada":               "kn",
	"Kazakh":                "kk",
	"Khmer":                 "km",
	"Klingon":               "tlh",
	"Korean":                "ko",
	"Kurdish":               "ku",
	"Kyrgyz":                "ky",
	"Lao":                   "lo",
	"Latgalian":             "ltg",
	"Latin":                 "la",
	"Latvian":               "lv",
	"Lithuanian":            "lt",
	"Luxembourgish":         "lb",
	"Macedo":                "rup",
	"Macedonian":            "mk",
	"Malagasy":              "mg",
	"Malay":                 "ms",
	"Malayalam":             "ml",
	"Maltese":               "mt",
	"Marathi":               "mr",
	"Mauritian Creole":      "mfe",
	"Mongolian":             "mn",
	"Montenegrin":           "srp",
	"Nepali":                "ne",
	"Norwegian Bokmal":      "nb",
	"Norwegian Nynorsk":     "nn",
	"Occitan":               "oc",
	"Pashto":                "ps",
	"Persian":               "fa",
	"Polish":                "pl",
	"Portuguese":            "pt",
	"Portuguese, Brazilian": "pt-br",
	"Punjabi":               "pa",
	"Romanian":              "ro",
	"Russian":               "ru",
	"Rusyn":                 "ry",
	"Serbian":               "sr",
	"Serbo-Croatian":        "sh",
	"Silesian":              "szl",
	"Sinhala":               "si",
	"Slovak":                "sk",
	"Slovenian":             "sl",
	"Somali":                "so",
	"Spanish":               "es",
	"Swahili":               "sw",
	"Swedish":               "sv",
	"Swedish Chef":          "art-x-bork",
	"Tagalog":               "tl",
	"Tajik":                 "tg",
	"Tamil":                 "ta",
	"Tatar":                 "tt",
	"Telugu":                "te",
	"Thai":                  "th",
	"Tibetan":               "bo",
	"Turkish":               "tr",
	"Turkmen":               "tk",
	"Ukrainian":             "uk",
	"Urdu":                  "ur",
	"Uyghur":                "ug",
	"Uzbek":                 "uz",
	"Vietnamese":            "vi",
}

func genTranscriptURLs(langCodes map[string]string, availableLanguages []string, videoURL string) []string {

	langBaseURL := "/transcript?language="

	var urls []string

	for _, lang := range availableLanguages {

		newURL := videoURL + langBaseURL + langCodes[lang]
		//fmt.Println(x)
		urls = append(urls, newURL)
	}
	//numOfURLs := len(urls)
	//fmt.Println("generated URLs : ", numOfURLs)

	return urls
}

type talkTranscript struct {
	LocalTalkTitle              string   `json:"LocalTalkTitle"`
	Paragraphs                  []string `json:"Paragraphs"`
	TimeStamps                  []string `json:"TimeStamps"`
	TalkTranscriptAndTimeStamps []string `json:"TalkTranscriptAndTimeStamps"`
}

type TranscriptPage struct {
	AvailableTranscripts []string                  `json:"AvailableTranscripts"`
	DatePosted           string                    `json:"DatePosted"`
	Rated                string                    `json:"Rated"`
	ImageURL             string                    `json:"ImageURL"`
	TalkTranscript       map[string]talkTranscript `json:"TalkTranscript"`
}

func main() {

	// Add logger and stubs for better debugging
	checkInternet()

	videoURL := "https://www.ted.com/talks/ken_robinson_says_schools_kill_creativity"
	//videoURL := "https://www.ted.com/talks/elon_musk_the_future_we_re_building_and_boring"

	// We are knowingly making sync. calls to the main Video page and
	// in case we find there are One or more subtitle lanuguages we make
	// more async. requests
	var videoPageInfo VideoPage
	videoPageInfo = videoFetchInfo(videoURL)

	// Checking if there are any subtitles at all
	// In case there are, we send a default query to fetch the list of available languages
	numOfSubtitles, e1 := strconv.ParseInt(videoPageInfo.AvailableSubtitlesCount, 10, 64)

	checkErr(e1, "ERROR: in parsing available subtitiles")

	// This function will cause the program to EXIT if there are no subtitles
	// Else we continue to fill the basic page
	exitIfNoSubtitlesExist(numOfSubtitles)

	transcriptEnURL := videoURL + "/transcript?language=en"

	// Since we've already made the request to default lang transcript
	// we fill in the common details into a transcript info struct
	transcriptPageCommonInfo := transcriptFetchCommonInfo(transcriptEnURL)

	urls := genTranscriptURLs(langCodes, transcriptPageCommonInfo.AvailableTranscripts, videoURL)
	//fmt.Println(transcriptCommonInfo.AvailableTranscripts)

	// @@@@@@@@@@
	// Page UnCommon

	//var transcriptS []talkTranscript

	langSpecificMap := make(map[string]talkTranscript)

	var wg sync.WaitGroup

	numOfURLs := len(urls)
	//fmt.Println(numOfURLs)
	wg.Add(numOfURLs)

	for _, url := range urls {

		go func(url string) {
			defer wg.Done()
			//color.Green(url)
			x, langName := transcriptFetchUncommonInfo(url)
			langSpecificMap[langName] = x
			//transcriptS = append(transcriptS, x)
		}(url)

	}

	wg.Wait()

	//writeJSON(videoPageInfo)

	//fmt.Println(transcriptS)

	// STUB for the actual construction of the complete talk struct

	var transcriptPageUnCommonInfo TranscriptPage
	transcriptPageUnCommonInfo.TalkTranscript = langSpecificMap

	transcriptPageCompleteInfo := transcriptPageCommonInfo
	mergo.Merge(&transcriptPageCompleteInfo, transcriptPageUnCommonInfo)
	//writeJSON(transcriptPageCompleteInfo)

	//temp1, _ := json.Marshal(transcriptPageCompleteInfo)
	//fmt.Println(string(temp1))

	var tedTalk TedTalk
	tedTalk.TalkVideoPage = videoPageInfo
	tedTalk.TalkTranscriptPage = transcriptPageCompleteInfo
	//mergo.Merge(&tedTalk, transcriptPageCompleteInfo)
	//fmt.Println(tedTalk)
	//temp2, _ := json.Marshal(tedTalk)
	//fmt.Println(string(temp2))
	writeJSON(tedTalk)
} // end of main()

func writeJSON(aStruct TedTalk) {

	temp1, e1 := json.Marshal(aStruct)
	checkErr(e1, "ERROR: unable to marshal the talk struct in writeJSON function")

	//fmt.Println(string(temp1))
	htmlSplit := strings.Split(aStruct.TalkVideoPage.TalkURL, "/")
	talkName := htmlSplit[len(htmlSplit)-1]

	fileName := "./" + talkName + ".json"

	f, e2 := os.Create(fileName)

	checkErr(e2, "ERROR: unable to create a file on disk.")

	f.Write(temp1)
	defer f.Close()
}

func checkErr(e error, errInfo string) {
	if e != nil {
		panic(errInfo)
	}
}

func checkInternet() {
	// Make a GET request
	result, err := http.Get("https://google.com")
	// Process response
	if err != nil {
		color.Red("We're OFF-Line!")
		//panic("Not connected to the net") // More idiomatic way would be to print the error and die unless it's a serious error

		// Learn about exit status in Golang
		os.Exit(1)
	}

	defer result.Body.Close()

}

func exitIfNoSubtitlesExist(numOfSubtitles int64) {
	if numOfSubtitles < 1 {
		color.Red("No subtitles available yet")
		os.Exit(1)
	}
}

func videoFetchInfo(url string) VideoPage {

	videoPage, e1 := goquery.NewDocument(url)

	checkErr(e1, "ERROR: unable to fetch the video page in videoFetchFunction")

	videoPageInstance := VideoPage{
		TalkURL:                 videoTalkURL(url),
		AvailableSubtitlesCount: videoAvailableSubtitlesCount(videoPage),
		Speaker:                 videoSpeaker(videoPage),
		Duration:                videoDuration(videoPage),
		TimeFilmed:              videoTimeFilmed(videoPage),
		TalkViewsCount:          videoTalkViewsCount(videoPage),
		TalkTopicsList:          videoTalkTopicsList(videoPage),
		TalkCommentsCount:       videoTalkCommentsCount(videoPage),
	}

	return videoPageInstance
}

func transcriptFetchCommonInfo(url string) TranscriptPage {
	transcriptPage, e1 := goquery.NewDocument(url)
	checkErr(e1, "ERROR: unable to fetch the common transcript page info.")

	transcriptPageInstance := TranscriptPage{

		AvailableTranscripts: transcriptAvailableTranscripts(transcriptPage),
		DatePosted:           transcriptDatePosted(transcriptPage),
		Rated:                transcriptRated(transcriptPage),
		ImageURL:             transcriptGetImage(transcriptPage, url),
	}
	return transcriptPageInstance
}

func transcriptFetchUncommonInfo(url string) (talkTranscript, string) {

	//fmt.Println(url)
	transcriptPage, e1 := goquery.NewDocument(url)
	//fmt.Println(transcriptLocalTalkTitle(transcriptPage))
	checkErr(e1, "ERROR: unable to fetch the Uncommon transcript page info.")

	transcript := talkTranscript{

		LocalTalkTitle:              transcriptLocalTalkTitle(transcriptPage),
		Paragraphs:                  transcriptTalkTranscript(transcriptPage),
		TimeStamps:                  transcriptTimeStamps(transcriptPage),
		TalkTranscriptAndTimeStamps: transcriptTalkTranscriptAndTimeStamps(transcriptPage),
	}

	langName := strings.Split(url, "=")[1]
	return transcript, langName
}

// @@@@@@@@@@@@@@@@@@
// VIDEO PAGE

func videoAvailableSubtitlesCount(doc *goquery.Document) string {

	subtitles := doc.Find(".player-hero__meta__link").Contents().Text()
	//fmt.Println(subtitles)

	//for _, x := range strings.Split(subtitles, "\n") {
	//fmt.Println(x)
	//println("~~~~~~")
	//}

	y := strings.Split(subtitles, "\n")
	z := strings.Split(y[3], " ")[0]
	// In case I need an INT
	//numOfSubtitles, _ := strconv.ParseInt(z, 10, 32)
	numOfSubtitles := z
	return numOfSubtitles
}

func videoSpeaker(doc *goquery.Document) string {
	speaker := doc.Find(".talk-speaker__name").Contents().Text()
	//fmt.Println(speaker)
	speaker = strings.Trim(speaker, "\n")
	return speaker
}

/*
// This is now taken from the transcripts page
func title(doc *goquery.Document) {
	title := doc.Find(".player-hero__title__content").Contents().Text()
	fmt.Println(title)
}
*/

func videoDuration(doc *goquery.Document) string {

	duration := doc.Find(".player-hero__meta").Contents().Text()
	//fmt.Println(duration)

	//for _, x := range strings.Split(duration, "\n") {
	//	fmt.Println(x)
	//	println("~~~~~~")
	//}

	x := strings.Split(duration, "\n")
	//fmt.Println(x[6])
	return x[6]

}

// TimeFilmed : Time at which the talk was filmed
func videoTimeFilmed(doc *goquery.Document) string {

	timeFilmed := doc.Find(".player-hero__meta").Contents().Text()

	//	fmt.Println(time_filmed)

	y := strings.Split(timeFilmed, "\n")
	//fmt.Println(y[11])
	return y[11]
}

func videoTalkViewsCount(doc *goquery.Document) string {

	talkViewsCount := doc.Find("#sharing-count").Contents().Text()
	//	fmt.Println(talkViewsCount)

	a := strings.Split(talkViewsCount, "\n")
	b := strings.TrimSpace(a[2])
	//fmt.Println(b)
	return b

}

func videoTalkTopicsList(doc *goquery.Document) []string {

	talkTopics := doc.Find(".talk-topics__list").Contents().Text()

	c := strings.Split(talkTopics, "\n")
	var topics []string
	for i := 3; i < len(c); i++ {
		//fmt.Println(c[i])
		if c[i] == "" {

		} else {
			topics = append(topics, c[i])
		}
	}
	return topics
}

func videoTalkCommentsCount(doc *goquery.Document) string {

	talkCommentsCount := doc.Find(".h11").Contents().Text()
	//fmt.Println(talkCommentsCount)
	d := strings.Split(talkCommentsCount, " ")
	//fmt.Println(d[0])
	return strings.TrimLeft(d[0], "\n")
}

func videoTalkURL(videoURL string) string {
	return videoURL
}

// TRANSCRIPT PAGE

// OUTPUT
// Do schools kill creativity?
func transcriptLocalTalkTitle(doc *goquery.Document) string {
	title := doc.Find(".m5").Contents().Text()
	//fmt.Println(strings.Split(title, "\n")[2])
	return strings.Split(title, "\n")[2]
}

// OUTPUT
// ["0:11","0:15","0:16","0:23","0:29","0:56","1:11","1:15","1:18","1:21","1:35","1:37","1:38","1:41","2:23","2:56","3:09","3:11","3:16","3:18","3:20","3:22","3:25","3:27","3:30","3:56","4:07","4:12","4:14","4:20","4:21","4:25","4:27","5:09","5:21","6:05","6:21","6:31","6:33","6:54","7:01","7:03","7:10","7:11","7:15","7:21","7:22","7:24","7:28","7:29","7:34","7:57","7:58","8:10","8:18","8:21","8:27","9:10","9:13","9:22","9:48","9:51","10:21","10:27","10:30","10:36","10:45","10:48","10:54","10:56","11:00","11:02","11:18","11:41","12:06","12:23","12:56","13:33","13:56","13:58","14:18","14:25","14:26","14:28","14:43","14:50","15:46","15:48","15:50","15:53","16:26","16:50","17:32","17:39","18:13","18:33","19:04","19:05"]
func transcriptTimeStamps(doc *goquery.Document) []string {
	times := doc.Find(".talk-transcript__para__time").Contents().Text()
	var times1 []string

	for _, time := range strings.Split(times, " ") {
		if time == "" {

		} else {

			//fmt.Println(time)
			times1 = append(times1, strings.TrimRight(time, " "))

		}
	}

	//fmt.Println(times1)
	var times2 []string
	for _, time := range strings.Split(times1[len(times1)-1], "\n\n") {

		times2 = append(times2, time)
	}

	//fmt.Println(times2)
	var timestamps []string

	for i := 1; i < len(times1)-1; i++ {

		timestamps = append(timestamps, times1[i])

	}

	for i := 0; i < len(times2); i++ {

		timestamps = append(timestamps, times2[i])
	}

	//fmt.Println(timestamps)

	for i := 0; i < len(timestamps); i++ {

		timestamps[i] = strings.Trim(timestamps[i], "\n")
	}

	//x, _ := json.Marshal(timestamps)
	//fmt.Println(string(x))

	return timestamps
}

// OUTPUT
// Seperate hunks of the Textual string
func transcriptTalkTranscript(doc *goquery.Document) []string {
	texts := doc.Find(".talk-transcript__para__text").Contents().Text()
	var para []string
	for _, text := range strings.Split(texts, "  ") {

		//fmt.Println(text)
		para = append(para, text)
	}

	var lines []string
	for _, para := range strings.Split(texts, "\n\n") {

		//fmt.Println(text)
		lines = append(lines, para)
	}

	return para
	//return lines
}

// OUTPUT
// The entire text chunk
func transcriptTalkTranscriptAndTimeStamps(doc *goquery.Document) []string {

	texts := doc.Find(".talk-transcript__para").Contents().Text()
	var para []string
	for _, text := range strings.Split(texts, "  ") {

		//fmt.Println(text)
		para = append(para, text)
	}

	var lines []string
	for _, para := range strings.Split(texts, "\n\n") {

		//fmt.Println(text)
		lines = append(lines, para)
	}

	return para
	//return lines
}

// OUTPUT
// [Afrikaans Albanian Arabic Armenian Azerbaijani Basque Belarusian Bengali Bulgarian Catalan Chinese, Simplified Chinese, Traditional Croatian Czech Danish Dutch English Esperanto Estonian Filipino Finnish French French (Canada) Galician Georgian German Greek Hebrew Hungarian Indonesian Ingush Italian Japanese Korean Lao Latvian Lithuanian Macedonian Marathi Mongolian Nepali Norwegian Bokmal Persian Polish Portuguese Portuguese, Brazilian Romanian Russian Serbian Slovak Slovenian Spanish Swedish Thai Turkish Ukrainian Urdu Uzbek Vietnamese]
// This should return an array of strings => ["langs"]
func transcriptAvailableTranscripts(doc *goquery.Document) []string {

	var langsList []string

	langs := doc.Find(".talk-transcript__language").Contents().Text()

	//	fmt.Println(langs)
	langsSeparated := strings.Split(langs, "\n")

	for i := 1; i < len(langsSeparated)-1; i++ {
		//fmt.Println(i, ":", langsSeparated[i])
		langsList = append(langsList, langsSeparated[i])
	}

	return langsList
}

// OUTPUT
// Jun 2006
func transcriptDatePosted(doc *goquery.Document) string {
	posted := doc.Find(".meta__item").Contents().Text()
	p := strings.Split(posted, "\n")
	//fmt.Println(p[3])
	return (p[3])

}

// OUTPUT
// Inspiring, Funny
func transcriptRated(doc *goquery.Document) string {

	rated := doc.Find(".meta__row").Contents().Text()

	r := strings.Split(rated, "\n")
	//fmt.Println(r[3])
	return r[3]
	/*
	   rx := strings.Split(r[3], ",")

	   	for _, x := range rx{
	   		append(ls,x)
	   	}
	*/

	//println(len(rx))
	//println(r[0])
	//println(r[1])
	//return(p[3])
}

func transcriptGetImage(doc *goquery.Document, videoURL string) string {

	imageURL, _ := doc.Find(".thumb__image").Attr("src")

	response, e1 := http.Get(imageURL)

	checkErr(e1, "ERROR: unable to fetch the image from the transcript page")

	defer response.Body.Close()

	//open a file for writing
	htmlSplit := strings.Split(videoURL, "/")
	talkName := htmlSplit[len(htmlSplit)-2]

	// Establish a file name
	fileName := "./" + talkName + ".jpg"

	f, e2 := os.Create(fileName)

	checkErr(e2, "ERROR: unable to create a file on disk in transcriptGetImage()")

	defer f.Close()

	// Use io.Copy to just dump the response body to the file. This supports huge files
	_, e3 := io.Copy(f, response.Body)
	checkErr(e3, "ERROR: unable to copy the body of response in transcriptGetImage()")

	return imageURL
}
