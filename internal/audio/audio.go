package audio

import (
	"log"
	"sync"

	"wsjtx-web/internal/hub"

	"github.com/gen2brain/malgo"
	"github.com/gorilla/websocket"
)

// Manager captures PCM from a malgo device and broadcasts binary frames on the hub.
type Manager struct {
	ctx         *malgo.AllocatedContext
	device      *malgo.Device
	hub         *hub.Hub
	currentID   string
	mu          sync.Mutex
	isCapturing bool
}

// New initializes the malgo context. Call Start to begin capture.
func New(h *hub.Hub) (*Manager, error) {
	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, func(message string) {
		log.Printf("LOG <%v>\n", message)
	})
	if err != nil {
		return nil, err
	}

	return &Manager{
		ctx: ctx,
		hub: h,
	}, nil
}

// Start begins (or restarts) capture on the given device ID. Empty ID uses the default device.
func (am *Manager) Start(deviceID string) error {
	am.mu.Lock()
	defer am.mu.Unlock()

	if am.isCapturing {
		am.device.Uninit()
		am.isCapturing = false
	}

	deviceConfig := malgo.DefaultDeviceConfig(malgo.Capture)
	deviceConfig.Capture.Format = malgo.FormatS16
	deviceConfig.Capture.Channels = 1
	deviceConfig.SampleRate = 12000
	deviceConfig.Alsa.NoMMap = 1

	if deviceID != "" {
		infos, err := am.ctx.Devices(malgo.Capture)
		if err == nil {
			for _, info := range infos {
				if info.ID.String() == deviceID {
					deviceConfig.Capture.DeviceID = info.ID.Pointer()
					break
				}
			}
		}
	}

	onRecvFrames := func(pOutputSample, pInputSamples []byte, framecount uint32) {
		am.hub.Broadcast(hub.BroadcastMessage{
			Type:    websocket.BinaryMessage,
			Payload: pInputSamples,
		})
	}

	captureCallbacks := malgo.DeviceCallbacks{
		Data: onRecvFrames,
	}

	device, err := malgo.InitDevice(am.ctx.Context, deviceConfig, captureCallbacks)
	if err != nil {
		return err
	}

	if err := device.Start(); err != nil {
		device.Uninit()
		return err
	}

	am.device = device
	am.currentID = deviceID
	am.isCapturing = true
	log.Printf("Started audio capture on device: %s", deviceID)
	return nil
}

// ListDevices returns capture devices as id/name maps.
func (am *Manager) ListDevices() ([]map[string]string, error) {
	infos, err := am.ctx.Devices(malgo.Capture)
	if err != nil {
		return nil, err
	}

	var devices []map[string]string
	for _, info := range infos {
		devices = append(devices, map[string]string{
			"id":   info.ID.String(),
			"name": info.Name(),
		})
	}
	return devices, nil
}
