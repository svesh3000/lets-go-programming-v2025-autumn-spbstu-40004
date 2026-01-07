package wifi_test

import (
	"fmt"
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

	interfaces := []*wifi.Interface{
		{Name: "one", HardwareAddr: parseMAC(t, "00:11:22:33:44:55")},
		{Name: "two", HardwareAddr: parseMAC(t, "aa:bb:cc:dd:ee:ff")},
	}

	mockWifi.On("Interfaces").Return(interfaces, nil)

	service := mywifi.New(mockWifi)
	addrs, err := service.GetAddresses()

	require.NoError(t, err)
	require.Len(t, addrs, 2)
	require.Equal(t, parseMAC(t, "00:11:22:33:44:55"), addrs[0])
	require.Equal(t, parseMAC(t, "aa:bb:cc:dd:ee:ff"), addrs[1])

	mockWifi.AssertExpectations(t)
}

func TestGetAddressesError(t *testing.T) {
	mockWifi := NewWiFiHandle(t)

	mockWifi.On("Interfaces").Return(nil, fmt.Errorf("no wifi"))

	service := mywifi.New(mockWifi)
	addrs, err := service.GetAddresses()

	require.Error(t, err)
	require.ErrorContains(t, err, "getting interfaces")
	require.Nil(t, addrs)

	mockWifi.AssertExpectations(t)
}
