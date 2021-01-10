package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	pt "github.com/c-bata/go-prompt"
	tui "github.com/manifoldco/promptui"
	"github.com/sirupsen/logrus"
)

// Datas is alias for map[string]string
type Datas map[string]string

var (
	CacheLog = make(map[string][]string)
)

var (
	// Tui for handle some
	Tui = struct {
		Select func(label string, options ...string) string
		Input  func(label string, suggest Datas) string
	}{
		Select: func(label string, options ...string) string {
			prompt := tui.Select{
				Label:        label,
				Items:        options,
				HideSelected: true,
				Size:         20,
				Searcher: func(s string, ix int) bool {
					return strings.Contains(options[ix], s)
				},
			}
			_, result, err := prompt.Run()
			if err != nil {
				logrus.Error("Tui Error:", err)
				return ""
			}
			return result
		},
		Input: func(label string, suggest Datas) string {
			return pt.Input(label, func(d pt.Document) (s []pt.Suggest) {
				for k, v := range suggest {
					s = append(s, pt.Suggest{
						Text:        k,
						Description: v,
					})
				}
				return pt.FilterFuzzy(s, d.GetWordBeforeCursor(), true)
			})
		},
	}
)

func SaveAsCsv(dir string) {
	for k, v := range CacheLog {
		name := filepath.Join(dir, k) + ".csv"
		// fp, err := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
		fp, err := os.Create(name)
		if err != nil {
			logrus.Error(err)
			continue
		}
		buf := ""
		for _, msg := range v {
			if buf == "" {
				buf += msg
			} else {
				buf += "\n" + msg
			}
		}
		fp.WriteString(buf)
		defer fp.Close()
	}
}

func Log(name string, temp string, args ...interface{}) {
	logrus.Infof(temp, args...)
	msg := fmt.Sprintf(strings.ReplaceAll(temp, "|", ","), args...)
	if msgs, ok := CacheLog[name]; ok {
		// s := []string{}
		// for _, a := range args {
		// 	s = append(s, fmt.Sprintf("\"%s\"", a))
		// }
		// msg := strings.Join(s, ",")
		msgs = append(msgs, msg)
		CacheLog[name] = msgs
	} else {
		// s := []string{}
		// for _, a := range args {
		// 	s = append(s, fmt.Sprintf("\"%s\"", a))
		// }
		// msg := strings.Join(s, ",")
		CacheLog[name] = []string{msg}
	}
}

func GenerateConfIni() string {
	api := Tui.Input("telegram API (https://my.telegram.org to get this )  >>", Datas{})
	apihash := Tui.Input("telegram APIHASH (https://my.telegram.org to get this )  >>", Datas{})

	return fmt.Sprintf(`
[auth]
api = %s
apihash = %s
tddir = td-lib
tddb = db
tddbdir = dbdir

err_log = error.txt
log_level = 4

proxyTp = socks5
proxyIP = 127.0.0.1
proxyPort = 1091

[getchats]
limit = 1000

[getcontacts]
limit = 10000
`, api, apihash)
}
