package main

import "testing"

func TestFeedEntry_Filename(t *testing.T) {
	tests := []struct {
		url  string
		want string
	}{
		{
			url:  "http://fw.ota.homesmart.ikea.net/global/GW1.0/01.16.026/bin/10005777-TRADFRI-control-outlet-2.0.024.ota.ota.signed",
			want: "10005777-TRADFRI-control-outlet-2.0.024.ota.ota.signed",
		},
	}
	for _, tt := range tests {
		fe := FeedEntry{
			BinaryURL: tt.url,
		}
		if got := fe.Filename(); got != tt.want {
			t.Errorf("FeedEntry.Filename() = %v, want %v", got, tt.want)
		}
	}
}
