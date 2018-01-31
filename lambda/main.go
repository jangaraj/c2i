package main

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/Jeffail/gabs"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/influxdata/influxdb/client/v2"
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	if string(os.Getenv("DEBUG")) == "true" {
		log.Printf("Processing Lambda request %s\n", request.RequestContext.RequestID)
	}

	if len(request.Body) < 1 {
		return events.APIGatewayProxyResponse{}, errors.New("no name was provided in the HTTP body")
	}

	jsonParsed, err := gabs.ParseJSON([]byte(request.Body))
	if err != nil {
		log.Printf("JSON decode error: " + err.Error())
		return events.APIGatewayProxyResponse{}, errors.New("JSON decode error: " + err.Error())
	}

	// Create a new HTTPClient
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     os.Getenv("INFLUXDB_URL"),
		Username: os.Getenv("INFLUXDB_USERNAME"),
		Password: os.Getenv("INFLUXDB_PASSWORD"),
	})
	if err != nil {
		log.Printf("InfluxDB error: " + err.Error())
		return events.APIGatewayProxyResponse{}, errors.New("InfluxDB error: " + err.Error())
	}

	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  os.Getenv("INFLUXDB_DB_DATA"),
		Precision: "s",
	})
	if err != nil {
		log.Printf("InfluxDB error: " + err.Error())
		return events.APIGatewayProxyResponse{}, errors.New("InfluxDB error: " + err.Error())
	}

	var value float64
	var text string
	var ok bool
	var datestamp time.Time
	tags := map[string]string{}
	fields := map[string]interface{}{}

	// timestamp
	text, ok = jsonParsed.Path("Summary.Timestamp").Data().(string)
	if ok {
		datestamp, err = time.Parse("20060102150405", text[0:14])
		if err != nil {
			log.Printf("InfluxDB error: " + err.Error())
			return events.APIGatewayProxyResponse{}, errors.New("InfluxDB error: " + err.Error())
		}
	}

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
		log.Printf("InfluxDB point error: " + err.Error())
		return events.APIGatewayProxyResponse{}, errors.New("InfluxDB point error: " + err.Error())
	}
	bp.AddPoint(pt)
	if string(os.Getenv("DEBUG")) == "true" {
		log.Printf("DEBUG: InfluxDB test_timing batch point: %s\n", pt)
	}

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
		log.Printf("InfluxDB point error: " + err.Error())
		return events.APIGatewayProxyResponse{}, errors.New("InfluxDB point error: " + err.Error())
	}
	bp.AddPoint(pt)
	if string(os.Getenv("DEBUG")) == "true" {
		log.Printf("DEBUG: InfluxDB test_byte batch point: %s\n", pt)
	}

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
	fields["month"] = datestamp.Format("01")
	fields["year"] = datestamp.Format("2006")
	fields["year-month"] = datestamp.Format("2006-01")
	pt, err = client.NewPoint("test_counter", tags, fields, datestamp)
	if err != nil {
		log.Printf("InfluxDB point error: " + err.Error())
		return events.APIGatewayProxyResponse{}, errors.New("InfluxDB point error: " + err.Error())
	}
	bp.AddPoint(pt)
	if string(os.Getenv("DEBUG")) == "true" {
		log.Printf("DEBUG: InfluxDB test_counter batch point: %s\n", pt)
	}

	// Write the batch
	if err := c.Write(bp); err != nil {
		log.Printf("InfluxDB error: " + err.Error())
		return events.APIGatewayProxyResponse{}, errors.New("InfluxDB error: " + err.Error())
	}

	return events.APIGatewayProxyResponse{
		Body:       "OK",
		StatusCode: 200,
	}, nil

}

func main() {
	lambda.Start(Handler)
}
