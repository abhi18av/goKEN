package main

import "fmt"

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

func main() {

	videoURL := "https://www.ted.com/talks/elon_musk_the_future_we_re_building_and_boring"
	// [Afrikaans Albanian Arabic Armenian Azerbaijani Basque Belarusian Bengali Bulgarian Catalan Chinese, Simplified Chinese, Traditional Croatian Czech Danish Dutch English Esperanto Estonian Filipino Finnish French French (Canada) Galician Georgian German Greek Hebrew Hungarian Indonesian Ingush Italian Japanese Korean Lao Latvian Lithuanian Macedonian Marathi Mongolian Nepali Norwegian Bokmal Persian Polish Portuguese Portuguese, Brazilian Romanian Russian Serbian Slovak Slovenian Spanish Swedish Thai Turkish Ukrainian Urdu Uzbek Vietnamese]
	availableLanguages := []string{"Chinese, Simplified", "Russian", "Korean", "German", "Spanish", "English"}
	urls := genTranscriptURLs(langCodes, availableLanguages, videoURL)

	numOfURLs := len(urls)
	fmt.Println("generated URLs : ", numOfURLs)

	fmt.Println(urls)
}
