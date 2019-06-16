package LineMessagingSettings

import (
	"os"
	"strconv"
	"testing"
)

func TestLoadingSettings(t *testing.T) {
	var testSettings = Settings("./examplesettings.json")

	if testSettings.APISecret != "<your-API-secret>" {
		t.Errorf("API Secret not beeing read in properly should be 1 is actualy %s", testSettings.APISecret)
	}

	if testSettings.ChannelAccessToken != "<your-channal-access-token>" {
		t.Errorf("Channel Access Token not beeing read in properly should be 2 is actualy %s", testSettings.ChannelAccessToken)
	}

}

func TestLoadingSettingsViaEnv(t *testing.T) {
	testChannal := "channel"
	testSecret := "secret"
	testPort := 80
	store := StoreTempEnvs(testChannal, testSecret)
	os.Setenv(testSecret, testSecret)
	os.Setenv(testChannal, testChannal)
	os.Setenv("PORT", strconv.Itoa(testPort))
	var testSettings = SettingsEnv()
	if testSettings.ChannelAccessToken != testChannal {
		t.Errorf("ChannelAccessToken not beeing read in properly should be 1 is actualy %s", testSettings.APISecret)
	}
	if testSettings.APISecret != testSecret {
		t.Errorf("API Secret not beeing read in properly should be 1 is actualy %s", testSettings.APISecret)
	}

	if testSettings.Port != testPort {
		t.Errorf("port not beeing read in properly should be 1 is actualy %s", testSettings.APISecret)
	}

	SetEnvs(store)
}

func StoreTempEnvs(envs ...string) map[string]string {
	tempEnv := make(map[string]string)
	for _, env := range envs {
		if value := os.Getenv(env); value != "" {
			tempEnv[env] = value
		}
	}
	return tempEnv
}
func SetEnvs(envs map[string]string) {
	for env, value := range envs {
		os.Setenv(env, value)
	}
}
