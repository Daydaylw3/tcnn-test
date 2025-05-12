package biz

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var introduce []string

var cgoSleep = 3

var cgoSleepRandom = 2

func init() {
	if err := filepath.Walk("biz/wiki", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if filepath.Base(path) == "Readme.md" {
			return nil
		}
		file, err := readFile(path)
		if err != nil {
			log.Printf("[WARN]read file(%s) err: %v", path, err)
			return nil
		}
		introduce = append(introduce, file)
		return nil
	}); err != nil {
		log.Printf("[WARN]init wiki file err: %v", err)
	}
	if s, _ := strconv.Atoi(os.Getenv("CGO_SLEEP")); s >= 0 {
		cgoSleep = s
	}
	log.Printf("Cgo Sleep: %d", cgoSleep)
	if s, _ := strconv.Atoi(os.Getenv("CGO_SLEEP_RANDOM")); s >= 0 {
		cgoSleepRandom = s
	}
	log.Printf("Cgo Sleep Random: %d", cgoSleepRandom)
}

func readFile(filepath string) (string, error) {
	file, err := os.ReadFile(filepath)
	if err != nil {
		return "", fmt.Errorf("read file(%s) err: %w", filepath, err)
	}
	return string(file), nil
}

func Decode(ctx context.Context) error {
	var files []map[string]string
	for _, intro := range introduce {
		if err := ctx.Err(); err != nil {
			return err
		}
		file, err := parseIntroFile(intro)
		if err != nil {
			log.Printf("parse introduce file(%s) err: %v", intro, err)
			continue
		}
		files = append(files, file)
	}
	var wikis []wikiIntroduce
	for _, f := range files {
		if err := ctx.Err(); err != nil {
			return err
		}
		wikis = append(wikis, decodeToWikiIntroduce(f))
	}
	var uppers [][]string
	for _, w := range wikis {
		upper, err := callFunc(ctx, w)
		if errors.Is(err, ctx.Err()) {
			return err
		}
		uppers = append(uppers, upper)
	}
	return nil
}

var introRegexp = regexp.MustCompile(`===([^=]+)===`)

func parseIntroFile(input string) (map[string]string, error) {
	result := make(map[string]string)
	matches := introRegexp.FindAllStringIndex(input, -1)
	for i, match := range matches {
		start, end := match[0], match[1]
		// 确定内容的开始位置
		contentStart := end
		contentEnd := len(input) // 到文档末尾
		if i+1 < len(matches) {
			contentEnd = matches[i+1][0] // 下一个标题的位置
		}
		// 截取并存储标题和内容
		result[strings.TrimSpace(input[start+3:end-3])] = strings.TrimSpace(input[contentStart:contentEnd])
	}
	return result, nil
}

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

func decodeToWikiIntroduce(file map[string]string) wikiIntroduce {
	def := defaultIntro{
		introduce: file["Introduce"],
		life:      file["Life"],
		method:    "DefaultMethod",
		sleep:     cgoSleep + random.Intn(cgoSleepRandom+1),
	}
	if p, ok := file["Performances"]; ok {
		def.method = "Performances"
		return shakespeare{
			defaultIntro: def,
			performances: p,
		}
	}
	return def
}

func callFunc(ctx context.Context, w wikiIntroduce) ([]string, error) {
	intro, err := uppercaseByC(ctx, w.Introduce())
	if err != nil {
		return nil, err
	}
	life, err := uppercaseByC(ctx, w.Life())
	if err != nil {
		return nil, err
	}
	vw := reflect.ValueOf(w)
	method := vw.MethodByName(w.Method())
	if !method.IsValid() {
		log.Printf("method(%s) for %T not found", w.Method(), w)
		return []string{intro, life}, nil
	}
	c, err := uppercaseByC(ctx, method.Call(getCallIn())[0].String(), w.Sleep())
	if err != nil {
		return nil, err
	}
	return []string{
		intro, life, c,
	}, nil
}

func getCallIn() []reflect.Value {
	argSel := []interface{}{
		1,
		1.2,
		int64(3),
		"string",
		true,
		struct{}{},
		[]string{},
		map[string]string{},
		make(chan int),
	}
	nums := random.Intn(120)
	var args []interface{}
	for i := 0; i < nums; i++ {
		args = append(args, argSel[i%len(argSel)])
	}
	return append([]reflect.Value{reflect.ValueOf(1)}, reflect.ValueOf(args))
}

type wikiIntroduce interface {
	Introduce() string
	Life() string
	Method() string
	Sleep() int
}

type defaultIntro struct {
	introduce string
	life      string
	method    string
	sleep     int
}

func (d defaultIntro) Introduce() string {
	return d.introduce
}

func (d defaultIntro) Life() string {
	return d.life
}

func (d defaultIntro) Method() string {
	if d.method == "" {
		return "DefaultMethod"
	}
	return d.method
}

func (d defaultIntro) Sleep() int {
	return d.sleep
}

func (_ defaultIntro) DefaultMethod(_ int) string {
	return "Nothing"
}

type shakespeare struct {
	defaultIntro
	performances string
}

func (s shakespeare) Performances(_ int, _ ...interface{}) string {
	return s.performances
}
