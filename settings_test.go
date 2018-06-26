package main

import (
	"testing"
)

func TestLoadingSettings(t *testing.T) {
	var testSettings = Settings("./examplesettings.json")

	if testSettings.APISecret != "<your-API-secret>" {
		t.Errorf("API Secret not beeing read in properly should be 1 is actualy {0}", testSettings.APISecret)
	}

	if testSettings.ChannelAccessToken != "<your-channal-access-token>" {
		t.Errorf("Channel Access Token not beeing read in properly should be 2 is actualy {0}", testSettings.ChannelAccessToken)
	}

}
