package websocket

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	randn "math/rand"
	"net/url"
	local "ocpp16/plugin/passive/local"
	"ocpp16/protocol"

	// registry "ocpp16/registry/rpcx"
	"github.com/gorilla/websocket"
	//"github.com/sirupsen/logrus"
	// "github.com/smallnest/rpcx/client"
	"net/http"
	// "ocpp16/logwriter"
	"sync"
	"testing"
	"time"
)

//go test -timeout=30m -v
var mx sync.Mutex
var r = randn.New(randn.NewSource(time.Now().Unix()))
var ws_addr = flag.String("ws_addr", "127.0.0.1:8090", "websocket service address")

func RandString(len int) string {
	mx.Lock()
	defer mx.Unlock()
	bytes := make([]byte, len, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}

var fnBootNotificationRequest = func() protocol.BootNotificationRequest {
	return protocol.BootNotificationRequest{
		ChargePointVendor:       "qinglianyun",
		ChargePointModel:        "lihuaye",
		ChargePointSerialNumber: RandString(15),
		ChargeBoxSerialNumber:   RandString(15),
		FirmwareVersion:         RandString(15),
		Iccid:                   RandString(15),
		Imsi:                    RandString(15),
		MeterType:               RandString(15),
		MeterSerialNumber:       RandString(15),
	}
}

var fnStatusNotificationRequest = func() protocol.StatusNotificationRequest {
	return protocol.StatusNotificationRequest{ //valid request
		ConnectorId:     15,
		ErrorCode:       "ConnectorLockFailure",
		Info:            RandString(40),
		Status:          "Available",
		Timestamp:       time.Now().Format(protocol.ISO8601),
		VendorId:        RandString(240),
		VendorErrorCode: RandString(40),
	}
	// return protocol.StatusNotificationRequest{ //invalid request
	// 	ConnectorId:     1,
	// 	ErrorCode:       "ConnectorLockFailure",
	// 	Info:            RandString(40),
	// 	Status:          "Available",
	// 	Timestamp:       time.Now().Format(protocol.ISO8601),
	// 	VendorId:        RandString(260),
	// 	VendorErrorCode: RandString(55),
	// }
}

var fnAuthorizeRequest = func() protocol.AuthorizeRequest {
	return protocol.AuthorizeRequest{
		IdTag: "qinglianyun",
	}
}

var fnMeterValueRequest = func() protocol.MeterValuesRequest {
	var meterValueReq = protocol.MeterValuesRequest{
		ConnectorId:   10,
		TransactionId: 81771202930020352,
	}
	var meterValue = protocol.MeterValue{
		Timestamp: time.Now().Format(time.RFC3339),
	}
	var sampledValue = protocol.SampledValue{
		Value:   "100",
		Context: "",
		// Context:   "Interruption.Begin",
		Format: "Raw",
		// Measurand: "Energy.Active.Export.Register",
		Measurand: "",
		// Phase:     "L1",
		Phase: "",
		// Location:  "Cable",
		Location: "",
		Unit:     "Wh",
	}
	meterValue.SampledValue = append(meterValue.SampledValue, sampledValue)
	meterValueReq.MeterValue = append(meterValueReq.MeterValue, meterValue)
	return meterValueReq
}

var fnStartTransactionRequest = func() protocol.StartTransactionRequest {
	return protocol.StartTransactionRequest{
		ConnectorId:   10,
		IdTag:         "qinglianyun",
		MeterStart:    10,
		ReservationId: 10,
		Timestamp:     time.Now().Format(protocol.ISO8601),
	}
}
var fnStopTransactionRequest = func() protocol.StopTransactionRequest {
	var meterValue = protocol.MeterValue{
		Timestamp: time.Now().Format(protocol.ISO8601),
	}
	var sampledValue = protocol.SampledValue{
		Value:     RandString(10),
		Context:   "Interruption.Begin",
		Format:    "Raw",
		Measurand: "Energy.Active.Export.Register",
		Phase:     "L1",
		Location:  "Cable",
		Unit:      "Wh",
	}
	meterValue.SampledValue = append(meterValue.SampledValue, sampledValue)
	return protocol.StopTransactionRequest{
		IdTag:           "qinglianyun",
		MeterStop:       100,
		Timestamp:       time.Now().Format(protocol.ISO8601),
		TransactionId:   80030155044556800,
		Reason:          "EmergencyStop",
		TransactionData: []protocol.MeterValue{meterValue},
	}
}

func clientHandler(ctx context.Context, t *testing.T, d *dispatcher) {
	flag.Parse()
	// name, id := RandString(5), RandString(5)
	name, id := "qinglianyun", "lihuaye"
	path := fmt.Sprintf("/ocpp/%v/%v", name, id)
	u := url.URL{Scheme: "ws", Host: "182.92.132.15:8090", Path: path}
	dialer := websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 45 * time.Second,
		Subprotocols:     []string{"ocpp1.5", "ocpp1.6"},
	}
	c, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		t.Fatal("dial:", err)
	}
	defer c.Close()
	ch := make(chan string, 10)
	defer close(ch)
	// queue := NewQueue()
	var waitgroup sync.WaitGroup
	waitgroup.Add(1)
	var mtx sync.Mutex
	go func() {
		for range time.Tick(time.Second * 10) {
			if err = c.WriteMessage(websocket.PingMessage, []byte("ping")); err != nil {
				t.Error(err)
			}
		}
	}()
	go func() { //test for center request
		defer waitgroup.Done()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				// call := &protocol.Call{
				// 	MessageTypeID: protocol.CALL,
				// 	UniqueID:      RandString(7),
				// 	Action:        "BootNotification",
				// 	Request:       fnBootNotificationRequest(),
				// }
				// queue.Push(call.UniqueID)
				// if err := d.appendRequest(context.Background(), fmt.Sprintf("%v-%v", name, id), call); err != nil {
				// 	return
				// }
				// time.Sleep(time.Second * time.Duration(randn.Intn(3)) / 5)
				time.Sleep(time.Second * 10000)
			}
		}
	}()
	// waitgroup.Add(1)
	// go func() {
	// 	defer waitgroup.Done()
	// 	for {
	// 		select {
	// 		case <-ctx.Done():
	// 			return
	// 		case res_uniqueid := <-ch:
	// 			rep_uniqueid, _ := queue.Pop()
	// 			// next_uniqueid, _ := queue.Peek()
	// 			// t.Logf("ws_id(%v), res_uniqueid(%v),rep_uniqueid(%v),queue remain(%v), next_uniqueid(%v)", fmt.Sprintf("%v-%v", name, id), res_uniqueid, rep_uniqueid, queue.Len(), next_uniqueid)
	// 			if res_uniqueid != rep_uniqueid {
	// 				t.Errorf("ws_id(%v), res_uniqueid(%v) != rep_uniqueid(%v)", fmt.Sprintf("%v-%v", name, id), res_uniqueid, rep_uniqueid)
	// 			}
	// 		}
	// 	}
	// }()
	waitgroup.Add(1)
	go func() {
		defer waitgroup.Done()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				_, message, err := c.ReadMessage()
				if err != nil {
					return
				}
				fields, err := parseMessage(message)
				if err != nil {
					return
				}
				switch fields[0].(float64) {
				case float64(protocol.CALL):
					go func() {
						uniqueid := fields[1].(string)
						callResult := &protocol.CallResult{
							MessageTypeID: protocol.CALL_RESULT,
							UniqueID:      uniqueid,
							Response: &protocol.BootNotificationResponse{
								CurrentTime: time.Now().Format(time.RFC3339),
								Interval:    10,
								Status:      "Accepted",
							},
						}
						callResultMsg, err := json.Marshal(callResult)
						if err != nil {
							return
						}
						time.Sleep(time.Second * time.Duration(randn.Intn(3)) / 10)
						t.Logf("test for center call: recv msg(%+v), resp_msg(%+v)", string(message), string(callResultMsg))
						mtx.Lock()
						// err = c.WriteMessage(websocket.TextMessage, callResultMsg)
						mtx.Unlock()
						if err != nil {
							return
						}
						// ch <- callResult.UniqueID
					}()
				case float64(protocol.CALL_RESULT), float64(protocol.CALL_ERROR):
					t.Logf("test for client call: recv msg(%v), ", string(message))
				default:
					t.Log(string(message))
				}

			}
		}
	}()
	//test for client call
	waitgroup.Add(1)
	go func() {
		defer waitgroup.Done()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				var action = "BootNotification"
				call := &protocol.Call{
					MessageTypeID: protocol.CALL,
					UniqueID:      RandString(7),
					Action:        action,
				}
				switch action {
				case "StatusNotification":
					call.Request = fnStatusNotificationRequest()
				case "Authorize":
					call.Request = fnAuthorizeRequest()
				case "BootNotification":
					call.Request = fnBootNotificationRequest()
				case "MeterValues":
					call.Request = fnMeterValueRequest()
					t.Logf("%+v", call.Request)
				case "StartTransaction":
					call.Request = fnStartTransactionRequest()
				case "StopTransaction":
					call.Request = fnStopTransactionRequest()
				default:
				}
				callMsg, err := json.Marshal(call)
				if err != nil {
					t.Error(err)
					return
				}
				mtx.Lock()
				err = c.WriteMessage(websocket.TextMessage, callMsg)
				mtx.Unlock()
				if err != nil {
					t.Error(err)
					return
				}
				time.Sleep(time.Second * 100)
			}
		}
	}()
	waitgroup.Wait()
	t.Logf("(%v) grace exit gorutine", path)
}

func WsHandler(t *testing.T, waitGroup *sync.WaitGroup) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*36000)
	lg := initLogger()
	SetLogger(lg)
	server := NewDefaultServer()
	plugin := local.NewActionPlugin()
	server.RegisterActionPlugin(plugin)
	go func() {
		server.Serve(*ws_addr, "/ocpp/:name/:id")
	}()
	for i := 0; i < 1; i++ { //numbers of client
		time.Sleep(time.Second / 10)
		go func() {
			clientHandler(ctx, t, server.dispatcher)
		}()
	}
	select {
	case <-ctx.Done():
		time.Sleep(time.Second * 50)
		cancel()
	}
	waitGroup.Done()
}
