package local

import (
	"context"
	"ocpp16/protocol"
)

type LocalActionPlugin struct {
	requestHandlerMap  map[string]protocol.RequestHandler
	responseHandlerMap map[string]protocol.ResponseHandler
}

func NewActionPlugin() *LocalActionPlugin {
	plugin := &LocalActionPlugin{}
	plugin.registerRequestHandler()
	plugin.registerResponseHandler()
	return plugin
}

func (l *LocalActionPlugin) BootNotification(ctx context.Context, request protocol.Request) (protocol.Response, error) {
	return nil, nil
}

func (l *LocalActionPlugin) StatusNotification(ctx context.Context, request protocol.Request) (protocol.Response, error) {
	return nil, nil
}

func (l *LocalActionPlugin) MeterValues(ctx context.Context, request protocol.Request) (protocol.Response, error) {
	return nil, nil
}

func (l *LocalActionPlugin) Authorize(ctx context.Context, request protocol.Request) (protocol.Response, error) {
	return nil, nil
}

func (l *LocalActionPlugin) StartTransaction(ctx context.Context, request protocol.Request) (protocol.Response, error) {
	return nil, nil
}

func (l *LocalActionPlugin) StopTransaction(ctx context.Context, request protocol.Request) (protocol.Response, error) {
	return nil, nil

}

func (l *LocalActionPlugin) ChargingPointOffline(id string) error {
	return nil

}

// firmwareManagement - request
func (l *LocalActionPlugin) FirmwareStatusNotification(ctx context.Context, request protocol.Request) (protocol.Response, error) {
	return nil, nil
}

func (l *LocalActionPlugin) DiagnosticsStatusNotification(ctx context.Context, request protocol.Request) (protocol.Response, error) {
	return nil, nil
}

func (l *LocalActionPlugin) registerRequestHandler() {
	l.requestHandlerMap = map[string]protocol.RequestHandler{
		protocol.BootNotificationName:           protocol.RequestHandler(l.BootNotification),
		protocol.StatusNotificationName:         protocol.RequestHandler(l.StatusNotification),
		protocol.MeterValuesName:                protocol.RequestHandler(l.MeterValues),
		protocol.AuthorizeName:                  protocol.RequestHandler(l.Authorize),
		protocol.StartTransactionName:           protocol.RequestHandler(l.StartTransaction),
		protocol.StopTransactionName:            protocol.RequestHandler(l.StopTransaction),
		protocol.FirmwareStatusNotificationName: protocol.RequestHandler(l.FirmwareStatusNotification),
	}
}

//RequestHandler represent device active request Center
func (l *LocalActionPlugin) RequestHandler(action string) (protocol.RequestHandler, bool) {
	handler, ok := l.requestHandlerMap[action]
	return handler, ok
}

type Reply struct {
	Err error
}

// chargingCore-response
func (l *LocalActionPlugin) ChangeConfigurationResponse(ctx context.Context, res protocol.Response) error {
	return nil
}

func (l *LocalActionPlugin) DataTransferResponse(ctx context.Context, res protocol.Response) error {
	return nil
}

func (l *LocalActionPlugin) RemoteStartTransactionResponse(ctx context.Context, res protocol.Response) error {
	return nil
}

func (l *LocalActionPlugin) ResetResponse(ctx context.Context, res protocol.Response) error {
	return nil
}

func (l *LocalActionPlugin) RemoteStopTransactionResponse(ctx context.Context, res protocol.Response) error {
	return nil

}

func (l *LocalActionPlugin) UnlockConnectorResponse(ctx context.Context, res protocol.Response) error {
	return nil
}

func (l *LocalActionPlugin) GetConfigurationResponse(ctx context.Context, res protocol.Response) error {
	return nil
}

func (l *LocalActionPlugin) CallError(ctx context.Context, res protocol.Response) error {
	return nil
}

// smartCharging - repsonse
func (l *LocalActionPlugin) SetChargingProfileResponse(ctx context.Context, res protocol.Response) error {
	return nil
}

func (l *LocalActionPlugin) ClearChargingProfileResponse(ctx context.Context, res protocol.Response) error {
	return nil
}

func (l *LocalActionPlugin) GetCompositeScheduleResponse(ctx context.Context, res protocol.Response) error {
	return nil
}

// firmwareManagement - response
func (l *LocalActionPlugin) GetDiagnosticsResponse(ctx context.Context, res protocol.Response) error {
	return nil
}

func (l *LocalActionPlugin) UpdateFirmWareResponse(ctx context.Context, res protocol.Response) error {
	return nil
}

//Reservation - response

func (l *LocalActionPlugin) ReserveNowResponse(ctx context.Context, res protocol.Response) error {
	return nil
}

func (l *LocalActionPlugin) CancelReservationResponse(ctx context.Context, res protocol.Response) error {
	return nil
}

//RemoteTrigger -response
func (l *LocalActionPlugin) TriggerMessageResponse(ctx context.Context, res protocol.Response) error {
	return nil
}

//LocalAuthListManagement -response
func (l *LocalActionPlugin) SendLocalListResponse(ctx context.Context, res protocol.Response) error {
	return nil
}

func (l *LocalActionPlugin) GetLocalListVersionResponse(ctx context.Context, res protocol.Response) error {
	return nil
}

//ResponseHandler represent The device reply to the center request
func (l *LocalActionPlugin) ResponseHandler(action string) (protocol.ResponseHandler, bool) {
	handler, ok := l.responseHandlerMap[action]
	return handler, ok
}

func (l *LocalActionPlugin) registerResponseHandler() {
	l.responseHandlerMap = map[string]protocol.ResponseHandler{
		protocol.ChangeConfigurationName:    protocol.ResponseHandler(l.ChangeConfigurationResponse),
		protocol.DataTransferName:           protocol.ResponseHandler(l.DataTransferResponse),
		protocol.RemoteStartTransactionName: protocol.ResponseHandler(l.RemoteStartTransactionResponse),
		protocol.ResetName:                  protocol.ResponseHandler(l.ResetResponse),
		protocol.RemoteStopTransactionName:  protocol.ResponseHandler(l.RemoteStopTransactionResponse),
		protocol.UnlockConnectorName:        protocol.ResponseHandler(l.UnlockConnectorResponse),
		protocol.GetConfigurationName:       protocol.ResponseHandler(l.GetConfigurationResponse),
		protocol.SetChargingProfileName:     protocol.ResponseHandler(l.SetChargingProfileResponse),
		protocol.ClearChargingProfileName:   protocol.ResponseHandler(l.ClearChargingProfileResponse),
		protocol.GetCompositeScheduleName:   protocol.ResponseHandler(l.GetCompositeScheduleResponse),
		protocol.ReserveNowName:             protocol.ResponseHandler(l.ReserveNowResponse),
		protocol.CancelReservationName:      protocol.ResponseHandler(l.CancelReservationResponse),
		protocol.TriggerMessageName:         protocol.ResponseHandler(l.TriggerMessageResponse),
		protocol.SendLocalListName:          protocol.ResponseHandler(l.SendLocalListResponse),
		protocol.GetLocalListVersionName:    protocol.ResponseHandler(l.GetLocalListVersionResponse),
		protocol.GetDiagnosticsName:         protocol.ResponseHandler(l.GetDiagnosticsResponse),
		protocol.UpdateFirmwareName:         protocol.ResponseHandler(l.UpdateFirmWareResponse),
		protocol.CallErrorName:              protocol.ResponseHandler(l.CallError),
	}
}
