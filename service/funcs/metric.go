package funcs

import (
	"errors"
	_ "fmt"
	"github.com/soniah/gosnmp"
	"log"
	"strings"
	"time"
)

type QueryExecuter struct {
	Interal *gosnmp.GoSNMP
}

func GetQuerier(ip string, port int, community string, version string, timeout int) *QueryExecuter {
	var querier QueryExecuter
	querier.Interal = &gosnmp.GoSNMP{
		Target:    ip,
		Port:      uint16(port),
		Community: community,
		Version:   gosnmp.Version2c,
		Timeout:   time.Duration(timeout) * time.Second,
		Logger:    nil,
		Retries:   5,
	}
	err := querier.Interal.Connect()
	if err != nil {
		log.Println("querier connect err.")
	}
	return &querier
}

func (querier *QueryExecuter) Close() {
	querier.Interal.Conn.Close()
}

func (querier *QueryExecuter) GetMetricValue(oid string) (interface{}, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(" Recovered in GetMetricValue", r, "oid is:", oid)
		}
	}()

	oids := []string{oid}

	for {
		result, err := querier.Interal.GetNext(oids)
		if err != nil {
			log.Println("Get oid value failed, err is: ", err, "oid is:", oid)
			break
		}

		data := result.Variables[0]
		if !strings.Contains(data.Name, oid) {
			break
		}
		if data.Value != nil {
			return data.Value, nil
		}
	}
	return nil, errors.New("Empty metric value for " + oid)
}

func (querier *QueryExecuter) GetBulkMetricValue(oid string) (interface{}, error) {
	results, err := querier.Interal.BulkWalkAll(oid)
	if err != nil {
		log.Println(" builk oid value failed, err is ", err, "oid is:", oid)
		return nil, err
	}

	for _, variable := range results {
		return variable.Value, nil
	}
	return nil, errors.New("Empty bulk metric value for " + oid)
}
