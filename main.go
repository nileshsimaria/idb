package main

import (
	"fmt"
	"github.com/influxdata/influxdb/client/v2"
	flag "github.com/spf13/pflag"
	"log"
	"math/rand"
	"time"
)

var (
	sSize  = flag.Int("samples", 100, "Number of points in one write")
	writes = flag.Int("writes", 30, "Total number of InfluxDB write")
)

const (
	MyDB     = "bm"
	username = "influx"
	password = "influxdb"
)

func writePoints(clnt client.Client, sSize int) {
	sampleSize := sSize

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  MyDB,
		Precision: "us",
	})
	if err != nil {
		log.Fatal(err)
	}

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < sampleSize; i++ {
		regions := []string{"us-west1", "us-west2", "us-west3", "us-east1"}
		tags := map[string]string{
			"cpu":    "cpu-total",
			"host":   fmt.Sprintf("host%d", rand.Intn(1000)),
			"region": regions[rand.Intn(len(regions))],
		}

		idle := rand.Float64() * 100.0
		fields := map[string]interface{}{
			"idle": idle,
			"busy": 100.0 - idle,
		}

		pt, err := client.NewPoint(
			"cpu_usage",
			tags,
			fields,
			time.Now(),
		)
		if err != nil {
			log.Fatal(err)
		}
		bp.AddPoint(pt)
	}

	if err := clnt.Write(bp); err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://localhost:8086",
		Username: username,
		Password: password,
	})
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < *writes; i++ {
		writePoints(c, *sSize)
	}

	if err := c.Close(); err != nil {
		log.Fatal(err)
	}
}
