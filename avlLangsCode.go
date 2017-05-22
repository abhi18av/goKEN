package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
)

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

func genTranscriptURLs(langCodes map[string]string, avaiLableLanguages []string, videoURL string) []string {

	langBaseURL := "/transcript?language="

	var urls []string

	for key, _ := range langCodes {

		aLangCode := langCodes[key]
		newURL := videoURL + langBaseURL + aLangCode
		//fmt.Println(x)
		urls = append(urls, newURL)
	}
	//fmt.Println(len(urls))

	return urls
}
func WiFi() {
	// Make a get request
	rs, err := http.Get("https://google.com")
	// Process response
	if err != nil {
		color.Red("WiFi OFF")
		//panic("Not connected to the net") // More idiomatic way would be to print the error and die unless it's a serious error

		// Learn about exit status in Golang
		os.Exit(1)
	}

	defer rs.Body.Close()

}

func main() {

	WiFi()

	videoURL := "https://www.ted.com/talks/ken_robinson_says_schools_kill_creativity"
	transcriptEnURL := videoURL + "/transcript?language=en"

	transcriptPage, _ := goquery.NewDocument(transcriptEnURL)

	// Using append here to add to the array-field
	avaiLableLanguages := transcriptAvailableTranscripts(transcriptPage)

	urls := genTranscriptURLs(langCodes, avaiLableLanguages, videoURL)

	fmt.Println(urls)

	wg.Add(len(urls) + 1)
	//fmt.Println(urls)

	var wg sync.WaitGroup
	// @@@@@@@@@@
	// Page Common
	transcriptEnURL := videoURL + "/transcript?language=en"

	var transcriptPageCommon TranscriptPage
	transcriptPage, _ := goquery.NewDocument(url)

	// Using append here to add to the array-field
	//transcriptPageInstance.TalkTranscript = append(transcriptPageInstance.TalkTranscript, transcript)

	go func(url string) {
		defer wg.Done()
		transcriptAvailableTranscripts(transcriptPage)
	}(transcriptEnURL)

	// @@@@@@@@@@@@
	wg.Wait()
}

// transcriptPage

// OUTPUT
// Do schools kill creativity?
func transcriptLocalTalkTitle(doc *goquery.Document) string {
	title := doc.Find(".m5").Contents().Text()
	//fmt.Println(strings.Split(title, "\n")[2])
	return strings.Split(title, "\n")[2]
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
