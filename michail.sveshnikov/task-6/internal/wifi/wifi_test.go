package wifi_test

import (
	"errors"
	"net"
	"testing"

	wifi "github.com/mdlayher/wifi"
	"github.com/stretchr/testify/require"
	mywifi "github.com/svesh3000/task-6/internal/wifi"
)

//go:generate mockery --name=WiFiHandle --testonly --quiet --outpkg wifi_test --output .

var errNoWifi = errors.New("no wifi")

func parseMAC(t *testing.T, macStr string) net.HardwareAddr {
	t.Helper()

	hwAddr, err := net.ParseMAC(macStr)
	require.NoError(t, err)

	return hwAddr
}

func TestGetAddresses(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

	mockWifi := NewWiFiHandle(t)

	mockWifi.On("Interfaces").Return(nil, errNoWifi)

	service := mywifi.New(mockWifi)
	addrs, err := service.GetAddresses()

	require.Error(t, err)
	require.ErrorContains(t, err, "getting interfaces")
	require.Nil(t, addrs)

	mockWifi.AssertExpectations(t)
}

func TestGetNames(t *testing.T) {
	t.Parallel()

	mockWifi := NewWiFiHandle(t)

	interfaces := []*wifi.Interface{
		{Name: "one", HardwareAddr: parseMAC(t, "00:11:22:33:44:55")},
		{Name: "two", HardwareAddr: parseMAC(t, "aa:bb:cc:dd:ee:ff")},
	}

	mockWifi.On("Interfaces").Return(interfaces, nil)

	service := mywifi.New(mockWifi)
	names, err := service.GetNames()

	require.NoError(t, err)
	require.Len(t, names, 2)
	require.Equal(t, "one", names[0])
	require.Equal(t, "two", names[1])

	mockWifi.AssertExpectations(t)
}

func TestGetNamesError(t *testing.T) {
	t.Parallel()

	mockWifi := NewWiFiHandle(t)

	mockWifi.On("Interfaces").Return(nil, errNoWifi)

	service := mywifi.New(mockWifi)
	names, err := service.GetNames()

	require.Error(t, err)
	require.ErrorContains(t, err, "getting interfaces")
	require.Nil(t, names)

	mockWifi.AssertExpectations(t)
}
