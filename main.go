package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"time"
	//"sort"
	//"runtime/debug"
	"os"
	"regexp"

	"gopkg.in/igm/sockjs-go.v2/sockjs"
	//"strconv"
	"strings"

	//_ "github.com/denisenkom/go-mssqldb"

	_ "net/http/pprof"
)

//go:generate go run scripts/include.go
// mit diesem script das Verzeichnis ./web ausgelesen und damit die ./webContent.go datei erstellt!!!
// im Verzeichnis mit der main.go "go generate" aufrufen!!!

var UpTime time.Time

var Connections int = 0
var AllePings int = 0

var PublishTickertime time.Duration = 100 // default 100

var (
	host          = flag.String("host", ":8080", "specifies Host and Port.")
	ips           = flag.String("ips", "^.*?$", "specifies IP regex")
	debuglevel    = flag.Bool("debug", false, "enable debugging")
	audioPath     = flag.String("audiopath", "/home/pi/Music", "path to the Musik files")
	videoPath     = flag.String("videopath", "/home/pi/Videos", "path to the Video files")
	InternetRadio = flag.String("radiolink", "http://rockabilly-radio.net/", "internet Radio URL")
)

var videoData MISData

type BrowserParameter struct {
	dummy1 string
	dummy2 string
}

type BrowserMessageCommand struct {
	Typ       string           `xml:"Typ"`
	Parameter BrowserParameter `xml:"Parameter"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func isValueInList(value string, list []string) int {
	i := 0
	for _, v := range list {
		if v == value {
			return i
		}
		i++
	}
	return -1
}
func VideoList(path string) []os.FileInfo {
	var tmpVideoList []os.FileInfo
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Test\n")
	for _, file := range files {
		if file.IsDir() {
			VideoList(path + "/" + file.Name())
		} else {
			tmpVideoList = append(tmpVideoList, file)
		}

	}
	return tmpVideoList
}

func handleWebside(w http.ResponseWriter, req *http.Request) {
	match, _ := regexp.MatchString(*ips, strings.Split(req.RemoteAddr, ":")[0])
	if match {
		URL := "web" + req.URL.String()
		if URL == "web/" {
			URL = "web/index.html"
		}
		if isValueInList(URL, WebContentURL) > -1 {
			WebContentPos := isValueInList(URL, WebContentURL)
			WebContentTmpbase64 := WebContent[WebContentPos]
			WebContentTmp, _ := base64.StdEncoding.DecodeString(WebContentTmpbase64)
			mType := mime.TypeByExtension(filepath.Ext(URL))
			w.Header().Set("Content-Type", mType)
			fmt.Fprint(w, string(WebContentTmp[:]))
		}
	} else {
		w.WriteHeader(403)
	}
}
func setNameIndex(nameIndex []NameIndex, Name string, Index int) []NameIndex {
	for i, n := range nameIndex {
		if n.Name == Name {
			nameIndex[i].Index = append(nameIndex[i].Index, Index)
			return nameIndex
		}
	}
	return append(nameIndex, NameIndex{Name, []int{Index}})
}

func getMediaData() MediaData {
	var mediaData MediaData
	files := VideoList(*videoPath)
	for _, f := range files {
		if f.Name()[len(f.Name())-4:] == ".mp4" {

			splitFileName := strings.Split(f.Name()[0:len(f.Name())-4], " - ")
			splitFileName = append(splitFileName, "k.a.", "k.a.", "k.a.")

			mediaData.VideoData.Mediainfos = append(mediaData.VideoData.Mediainfos, MediaInfo{f.Name(), splitFileName[0], splitFileName[1], splitFileName[2]})
			mediaData.VideoData.Interpreten = setNameIndex(mediaData.VideoData.Interpreten, splitFileName[0], len(mediaData.VideoData.Mediainfos))
			mediaData.VideoData.Stile = setNameIndex(mediaData.VideoData.Stile, splitFileName[2], len(mediaData.VideoData.Mediainfos))
		}
	}
	return mediaData
}

func SendData2Browser(session sockjs.Session, message BrowserMessageCommand, EigenePings int) {
	//fmt.Println("SendCommand message: ", message)
	var userBytes []byte
	userBytes = nil
	switch message.Typ {
	case "getMediaData":
		var browserData MediaData
		browserData = getMediaData()
		browserData.XMLName.Local = "setMediaData"
		userBytes, _ = json.Marshal(browserData)
	default:
		fmt.Println("Unbekannter Anfrage vom Browser: ", message.Typ)
	}
	if userBytes != nil {
		session.Send(string(userBytes))
	}
}

func EventHander(session sockjs.Session) {
	//log.Println("Sockjs Session established.")
	Connections++
	var EigenePings int = 0

	log.Printf("%v Session(s) verbunden!\n", Connections)
	var closedSession = make(chan struct{})
	go func() {
		for {
			select {
			case <-closedSession:
				log.Println("Sockjs Session closed!")
				return
			}
		}
	}()
	for {
		if msg, err := session.Recv(); err == nil {
			message := BrowserMessageCommand{}
			err := json.Unmarshal([]byte(msg), &message)
			if err != nil {
				log.Println("\nError Browser Msg: ", err)
			} else {
				EigenePings++
				SendData2Browser(session, message, EigenePings)
			}
			continue
		}
		break
	}
	close(closedSession)
	Connections--
	log.Printf("%v Session(s) verbunden!\n", Connections)
}

func main() {
	// ******  params  ********
	flag.Parse()

	//Conn_Bind()
	//userAuthentication()
	// file(web)server starten
	// diese Zeile auskomentieren und die nächste einblenden, wenn man das verzeichnis ./web verwenden möchte
	//http.Handle("/", http.HandlerFunc(handleWebside))

	http.Handle(fmt.Sprintf("/"), http.StripPrefix("/", http.FileServer(http.Dir("./web"))))
	http.Handle(fmt.Sprintf("/videos/"), http.StripPrefix("/videos/", http.FileServer(http.Dir(*videoPath))))
	http.Handle(fmt.Sprintf("/audios/"), http.StripPrefix("/audios/", http.FileServer(http.Dir(*audioPath))))

	log.Println("Server started")
	log.Println("Please Browse to " + *host)

	// session erstellen!
	handler := sockjs.NewHandler(fmt.Sprintf("/ws/serverdata"), sockjs.DefaultOptions, func(session sockjs.Session) {
		EventHander(session)
		session.Close(1, "not longer in Use.")
	})
	http.Handle(fmt.Sprintf("/ws/serverdata/"), handler)
	http.ListenAndServe(*host, nil)

}
