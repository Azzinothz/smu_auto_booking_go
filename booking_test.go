package booking

import (
	"testing"
)

func TestNewBooker(t *testing.T) {
	booker := NewBooker(username, password)
	for _, cookie := range booker.collector.Cookies(baseURL) {
		if cookie.Name == "access_token" && len(cookie.Value) > 0 {
			return
		}
	}
	t.Fatal("Failed to get access token.")
}

func TestBookRoom(t *testing.T) {
	booker := NewBooker(username, password)
	roomsStatus := FetchRoomsStatus("2019-05-16")
	err := booker.BookRoom(roomsStatus[4], "21:00", "21:30", "立项相关内容", "立项相关内容立项相关内容立项相关内容", []string{"201610733001", "201610733002", "201610733003", "201610733004", "201610733005", "201610733006"}, "15021617205")
	if err != nil {
		t.Fatal(err)
	}
}

func TestFetchRoomsStatus(t *testing.T) {
	roomsStatus := FetchRoomsStatus("2019-05-14")
	if len(roomsStatus) != 10 {
		t.FailNow()
	}
}
