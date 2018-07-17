package funcs

import (
	"fmt"
	"github.com/soniah/gosnmp"
	"log"
	"strconv"
	"strings"
	"time"
)

type QueryExecuter struct {
	Interal *gosnmp.GoSNMP
}

func GetQuerier(ip string, community string, version string, timeout int) *QueryExecuter {
	var querier QueryExecuter
	ip_port := strings.Split(ip, ":")
	port := uint16(161)
	if len(ip_port) == 2 {
		if data, err := strconv.ParseUint(ip_port[1], 10, 16); err == nil {
			port = uint16(data)
		}
	}
	querier.Interal = &gosnmp.GoSNMP{
		Target:    ip,
		Port:      port,
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

	result, err := querier.Interal.Get(oids)
	if err != nil {
		log.Println("Get oid value failed, err is: ", err, "oid is:", oid)
	}

	data := result.Variables[0]
	switch data.Type {
	case gosnmp.OctetString:
		return fmt.Sprintf("%s", string(data.Value.([]byte))), nil
	default:
		ret := fmt.Sprintf("%d", gosnmp.ToBigInt(data.Value))
		if i, err := strconv.ParseInt(ret, 10, 64); err == nil && i != 0 {
			return i, nil
		}
	}
	return int64(0), nil
}

func (querier *QueryExecuter) GetBulkMetricValue(oid string) (map[string]interface{}, error) {
	results, err := querier.Interal.BulkWalkAll(oid)
	data := make(map[string]interface{})
	if err != nil {
		log.Println(" builk oid value failed, err is ", err, "oid is:", oid)
		return data, err
	}

	for _, variable := range results {
		keys := strings.Split(variable.Name, ".")
		key := keys[len(keys)-1]
		switch variable.Type {
		case gosnmp.OctetString:
			data[key] = fmt.Sprintf("%s", string(variable.Value.([]byte)))
		default:
			ret := fmt.Sprintf("%d", gosnmp.ToBigInt(variable.Value))
			if i, err := strconv.ParseInt(ret, 10, 64); err == nil {
				data[key] = i
			}
		}
	}
	return data, nil
}
