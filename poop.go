package main

import (
	"fmt"
	"github.com/tarm/serial"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const longPoop = int64(10)

var lastPoop int64

func main() {
	device := getEnv("DEVICE", "/dev/ttyUSB0")
	port := getEnv("PORT", "8080")

	go serialRead(device)
	http.HandleFunc("/", gotPoop)
	log.Fatal(http.ListenAndServe(
		fmt.Sprintf(":%s", port), nil))
}

func serialRead(device string) {
	c := &serial.Config{Name: device, Baud: 9600}
	s, err := serial.OpenPort(c)
	if err != nil {
		lastPoop = -1
		return
	}
	defer s.Close()

	lastPoop = time.Now().Unix()
	poopStart := lastPoop

	for true {
		buf := make([]byte, 4)
		_, err := s.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		// Wait because this might be a long pooper who sits in  the darkness for _longPoop_ seconds
		if canIPoop() {
			writeTimestamp("poop.log", fmt.Sprintf("%d %d\n", poopStart, lastPoop))
			poopStart = time.Now().Unix()
		}
		lastPoop = time.Now().Unix()
		// write lastPoop
	}
}

func writeTimestamp(fileName, text string) {
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(text); err != nil {
		panic(err)
	}
}

func readLogs(fileName string) string {
	f, err := os.OpenFile(fileName, os.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	logs, err := ioutil.ReadFile("poop.log")
	if err != nil {
		panic(err)
	}
	return string(logs)
}

func gotPoop(w http.ResponseWriter, r *http.Request) {
	if lastPoop == -1 {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if "cani/last" == r.URL.Path[1:] {
		fmt.Fprintf(w, "Last Poop was %d seconds ago", time.Now().Unix()-lastPoop)
		return
	}

	if "cani/logs" == r.URL.Path[1:] {
		fmt.Fprintf(w, "%s", readLogs("poop.log"))
		return
	}

	if "cani/" == r.URL.Path[1:] {
		if canIPoop() {
			fmt.Fprintf(w, "%d", time.Now().Unix()-lastPoop)
			return
		}
		http.Error(w, http.StatusText(http.StatusGone), http.StatusGone)
		return
	}

	imgSrc := "https://i.imgur.com/QhWI4Mg.gif"
	decision := "Nope"

	if canIPoop() {
		imgSrc = "https://i.imgur.com/l3Bs46c.jpg"
		decision = "Yes you can!"
	}

	respTpl, err := template.New("resp").Parse(htmlTmpl)
	if err != nil {
		log.Fatal("Template f up ", err)
	}
	respTpl.Execute(w, struct {
		ImgSrc   string
		Decision string
	}{imgSrc, decision})
}

func canIPoop() bool {
	return time.Now().Unix()-lastPoop > longPoop
}

const htmlTmpl = `
<head>
	<title>Can I Poop?</title>
	<meta property="og:title" content="Can I Poop?" />
	<meta property="og:description" content="{{.Decision}}" />
	<meta property="og:image" content="{{.ImgSrc}}" />
</head>
<script>
// reload every 5 seconds
setTimeout(function(){
   window.location.reload(1);
}, 5000);
</script>
<body>
	<img src='{{.ImgSrc}}' />
</body>
`

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
