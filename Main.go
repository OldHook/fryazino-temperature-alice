package main

import (
	"fmt"
	"github.com/azzzak/alice"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var temperature = -999
var dataTime = time.Now()
var warnLog = log.New(os.Stdout, "WARNING\t", log.Ldate|log.Ltime)

func getPageContent() string {
	resp, err := http.Get("https://fireras.su/thermo.dat/gauge.wml")
	if err != nil {
		warnLog.Println("Error getting fireras.su:" + err.Error())
		return ""
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		warnLog.Println("Error getting fireras.su page body:" + err.Error())
		return ""
	}
	return string(body)
}

func getTemperatureFromText(text string) string {
	tmp := strings.ToLower(text)
	a := strings.Split(tmp, "<b>")
	a = strings.Split(a[1], "</b>")
	a = strings.Split(a[0], "<br/>")
	tmp = strings.TrimSpace(a[2])
	//log.Println(tmp)
	return tmp
}

func refreshData() {
	text := getPageContent()
	if text == "" {
		return
	}
	//fmt.Println(text)
	text = getTemperatureFromText(text)
	tempFloat, err := strconv.ParseFloat(text, 64)
	if err != nil {
		warnLog.Println("Error getting temperature from text:" + err.Error())
	} else {
		temperature = int(math.Round(tempFloat))
		dataTime = time.Now()
	}
}

func infiniteLoop() {
	for {
		refreshData()
		time.Sleep(time.Minute * 10)
	}
}

func main() {
	log.Println("START")
	go infiniteLoop()

	updates := alice.ListenForWebhook("/hook", alice.Debug(true))
	go http.ListenAndServe(":3000", nil)

	updates.Loop(func(k alice.Kit) *alice.Response {
		_, resp := k.Init()

		dataAge := int(math.Round(time.Now().Sub(dataTime).Minutes())) // Сколько минут назад последний раз получали данные
		var respText string
		if temperature == -999 || dataAge > 60 {
			respText = "Информация отсутствует. Попробуйте спросить позже."
		} else if dataAge > 15 {
			respText = fmt.Sprintf("%d %s назад температура была ", dataAge, alice.Plural(dataAge, "минута", "минуты", "минут"))
			respText += fmt.Sprintf("%d %s", temperature, alice.Plural(temperature, "градус", "градуса", "градусов"))
		} else {
			respText = fmt.Sprintf("%d %s", temperature, alice.Plural(temperature, "градус", "градуса", "градусов"))
		}
		return resp.Text(respText).EndSession()
	})
}
