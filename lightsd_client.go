package main

import (
	"context"
	"io"
	"net"
	"time"

	"github.com/mafredri/cdp/rpcc"
)

type args struct {
	Target string `json:"target"`
}

type LightsClient struct {
	conn *rpcc.Conn
}

func NewLightsClient(sockettype, addr string) (*LightsClient, error) {
	conn, err := rpcc.Dial(addr, rpcc.WithDialer(func(ctx context.Context, addr string) (io.ReadWriteCloser, error) {
		conn, err := net.Dial(sockettype, addr)
		if err != nil {
			return nil, err
		}
		return conn, nil
	}))
	if err != nil {
		return nil, err
	}
	return &LightsClient{
		conn: conn,
	}, nil
}

type GetLightStateRet struct {
	Lifx struct {
		Addr    string `json:"addr"`
		Gateway struct {
			Site    string `json:"site"`
			Url     string `json:"url"`
			Latency int    `json:"latency"`
		} `json:"gateway"`
		MCU struct {
			FirmwareVersion string `json:"firmware_version"`
		} `json:"mcu"`
		Wifi struct {
			FirmwareVersion string `json:"firmware_version"`
		} `json:"wifi"`
	} `json:"_lifx"`
	Model  string        `json:"_model"`
	Vendor string        `json:"_vendor"`
	HSBK   []float64     `json:"hsbk"`
	Power  bool          `json:"power"`
	Label  string        `json:"label"`
	Tags   []interface{} `json:"tags"`
}

func (lc *LightsClient) GetLightState(target string) (result []GetLightStateRet, err error) {
	args := []interface{}{target}
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	err = rpcc.Invoke(ctx, "get_light_state", args, &result, lc.conn)
	if err != nil {
		return nil, err
	}
	return result, err
}

func (lc *LightsClient) SetLightFromHSBK(target string, h, s, b float64, k int, t time.Duration) (result interface{}, err error) {
	t = t / time.Millisecond
	args := []interface{}{&target, &h, &s, &b, &k, &t}
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	err = rpcc.Invoke(ctx, "set_light_from_hsbk", args, &result, lc.conn)
	if err != nil {
		return "", err
	}
	return result, err
}

func (lc *LightsClient) PowerToggle(target string) (result interface{}, err error) {
	args := &args{
		Target: "*",
	}
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	err = rpcc.Invoke(ctx, "power_toggle", args, &result, lc.conn)
	if err != nil {
		return "", err
	}
	return result, err
}

func (lc *LightsClient) Close() {
	lc.conn.Close()
}
