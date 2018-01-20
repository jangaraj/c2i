package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Jeffail/gabs"
	"github.com/influxdata/influxdb/client/v2"
)

var (
	Info  *log.Logger
	Error *log.Logger
)

func Init(
	infoHandle io.Writer,
	errorHandle io.Writer) {

	Info = log.New(infoHandle,
		"",
		log.Ldate|log.Ltime)

	Error = log.New(errorHandle,
		"",
		log.Ldate|log.Ltime)
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	Info.Printf("Processing request: %v %v %v %v %v",
		r.RemoteAddr, r.Method, r.Proto, r.Host, r.URL)

	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	defer r.Body.Close()
	newStr := buf.String()
	jsonParsed, err := gabs.ParseJSON([]byte(newStr))
	if err != nil {
		Error.Printf("JSON decode error: " + err.Error())
		http.Error(w, "JSON decode error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	//Info.Printf("Parsed struct: %s\n", jsonParsed)

	fmt.Fprintf(w, dataInflux(jsonParsed, w))
}

func main() {
	Init(os.Stdout, os.Stderr)
	Info.Printf("Starting c2i")
	// TODO hide password values
        /*
	for _, pair := range os.Environ() {
		Info.Printf("Env variable: %s\n", pair)
	}
        */
	http.HandleFunc("/data/", dataHandler)
	http.ListenAndServe(":"+os.Getenv("APP_PORT"), nil)
}

func dataInflux(jsonParsed *gabs.Container, w http.ResponseWriter) string {
	var datestamp = time.Now()

	// Create a new HTTPClient
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     os.Getenv("INFLUXDB_URL"),
		Username: os.Getenv("INFLUXDB_USERNAME"),
		Password: os.Getenv("INFLUXDB_PASSWORD"),
	})
	if err != nil {
		Error.Printf("InfluxDB error: " + err.Error())
		http.Error(w, "InfluxDB error: "+err.Error(), http.StatusInternalServerError)
		return "ERROR"
	}

	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  os.Getenv("INFLUXDB_DB_DATA"),
		Precision: "s",
	})
	if err != nil {
		Error.Printf("InfluxDB error: " + err.Error())
		http.Error(w, "InfluxDB error: "+err.Error(), http.StatusInternalServerError)
		return "ERROR"

	}

	var value float64
	var text string
	var ok bool
	tags := map[string]string{}
	fields := map[string]interface{}{}

	// POINT test_timing
	text, ok = jsonParsed.Path("TestDetail.Name").Data().(string)
	if ok {
		tags["testname"] = text
	}
	text, ok = jsonParsed.Path("NodeName").Data().(string)
	if ok {
		tags["nodename"] = text
	}
	value, ok = jsonParsed.Path("Summary.Timing.Total").Data().(float64)
	if ok {
		fields["total"] = value
	}
	value, ok = jsonParsed.Path("Summary.Timing.Dns").Data().(float64)
	if ok {
		fields["dns"] = value
	}
	value, ok = jsonParsed.Path("Summary.Timing.Wait").Data().(float64)
	if ok {
		fields["wait"] = value
	}
	value, ok = jsonParsed.Path("Summary.Timing.Connect").Data().(float64)
	if ok {
		fields["connect"] = value
	}
	value, ok = jsonParsed.Path("Summary.Timing.Send").Data().(float64)
	if ok {
		fields["send"] = value
	}
	value, ok = jsonParsed.Path("Summary.Timing.Ssl").Data().(float64)
	if ok {
		fields["ssl"] = value
	}
	value, ok = jsonParsed.Path("Summary.Timing.Wire").Data().(float64)
	if ok {
		fields["wire"] = value
	}
	value, ok = jsonParsed.Path("Summary.Timing.Client").Data().(float64)
	if ok {
		fields["client"] = value
	}
	value, ok = jsonParsed.Path("Summary.Timing.DocumentComplete").Data().(float64)
	if ok {
		fields["doccomplete"] = value
	} else {
		fields["doccomplete"] = float64(0)
	}
	value, ok = jsonParsed.Path("Summary.Timing.renderStart").Data().(float64)
	if ok {
		fields["renderstart"] = value
	} else {
		fields["renderstart"] = float64(0)
	}
	value, ok = jsonParsed.Path("Summary.Timing.domLoad").Data().(float64)
	if ok {
		fields["domload"] = value
	} else {
		fields["domload"] = float64(0)
	}
	pt, err := client.NewPoint("test_timing", tags, fields, datestamp)
	if err != nil {
		Error.Printf("InfluxDB point error: " + err.Error())
		http.Error(w, "InfluxDB point error: "+err.Error(), http.StatusInternalServerError)
		return "ERROR"
	}
	bp.AddPoint(pt)

	// POINT test_byte
	fields = map[string]interface{}{}
	value, ok = jsonParsed.Path("Summary.Byte.Response.TotalContent").Data().(float64)
	if ok {
		fields["totalcontent"] = value
	} else {
		fields["totalcontent"] = float64(0)
	}
	value, ok = jsonParsed.Path("Summary.Byte.Response.Image").Data().(float64)
	if ok {
		fields["image"] = value
	} else {
		fields["image"] = float64(0)
	}
	value, ok = jsonParsed.Path("Summary.Byte.Response.Script").Data().(float64)
	if ok {
		fields["script"] = value
	} else {
		fields["script"] = float64(0)
	}
	value, ok = jsonParsed.Path("Summary.Byte.Response.Css").Data().(float64)
	if ok {
		fields["css"] = value
	} else {
		fields["css"] = float64(0)
	}
	value, ok = jsonParsed.Path("Summary.Byte.Response.Html").Data().(float64)
	if ok {
		fields["html"] = value
	} else {
		fields["html"] = float64(0)
	}
	pt, err = client.NewPoint("test_byte", tags, fields, datestamp)
	if err != nil {
		Error.Printf("InfluxDB point error: " + err.Error())
		http.Error(w, "InfluxDB point error: "+err.Error(), http.StatusInternalServerError)
		return "ERROR"
	}
	bp.AddPoint(pt)

	// POINT test_counter
	fields = map[string]interface{}{}
	value, ok = jsonParsed.Path("Summary.Counter.Hosts").Data().(float64)
	if ok {
		fields["hosts"] = value
	}
	value, ok = jsonParsed.Path("Summary.Counter.Requests").Data().(float64)
	if ok {
		fields["requests"] = value
	}
	value, ok = jsonParsed.Path("Summary.Counter.FailedRequests").Data().(float64)
	if ok {
		fields["failedrequests"] = value
	}
	value, ok = jsonParsed.Path("Summary.Counter.JsFailures").Data().(float64)
	if ok {
		fields["jsfailures"] = value
	} else {
		fields["jsfailures"] = float64(0)
	}
	if jsonParsed.ExistsP("Summary.Error") == true && jsonParsed.ExistsP("Summary.Error.Code") == true {
		fields["availability"] = float64(0)
	} else {
		fields["availability"] = float64(100)
	}
	pt, err = client.NewPoint("test_counter", tags, fields, datestamp)
	if err != nil {
		Error.Printf("InfluxDB point error: " + err.Error())
		http.Error(w, "InfluxDB point error: "+err.Error(), http.StatusInternalServerError)
		return "ERROR"
	}
	bp.AddPoint(pt)

	// Write the batch
	if err := c.Write(bp); err != nil {
		Error.Printf("InfluxDB error: " + err.Error())
		http.Error(w, "InfluxDB error: "+err.Error(), http.StatusInternalServerError)
		return "ERROR"
	}

	return "OK"
}
