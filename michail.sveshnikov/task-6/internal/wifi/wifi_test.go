package wifi_test

import (
	"net"
	"testing"

	wifi "github.com/mdlayher/wifi"
	"github.com/stretchr/testify/require"
	mywifi "github.com/svesh3000/task-6/internal/wifi"
)

func parseMAC(t *testing.T, macStr string) net.HardwareAddr {
	t.Helper()

	hwAddr, err := net.ParseMAC(macStr)
	require.NoError(t, err)

	return hwAddr
}

func TestGetAddresses(t *testing.T) {
	mockWifi := NewWiFiHandle(t)

	iface := &wifi.Interface{
		Name:         "one",
		HardwareAddr: parseMAC(t, "00:11:22:33:44:55"),
	}
	mockWifi.On("Interfaces").Return([]*wifi.Interface{iface}, nil)

	service := mywifi.New(mockWifi)

	addrs, err := service.GetAddresses()

	require.NoError(t, err)
	require.Len(t, addrs, 1)
	require.Equal(t, parseMAC(t, "00:11:22:33:44:55"), addrs[0])

	mockWifi.AssertExpectations(t)
}
