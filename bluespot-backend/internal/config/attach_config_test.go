package config

import (
	"strings"
	"testing"

	"github.com/spf13/viper"
)

func TestAttachConfigFileSettings(t *testing.T) {
	previous := GlobalConfig
	t.Cleanup(func() { GlobalConfig = previous })
	for _, fromEnv := range []bool{false, true} {
		name := "yaml"
		if fromEnv {
			name = "environment"
		}
		t.Run(name, func(t *testing.T) {
			loader := viper.New()
			if fromEnv {
				t.Setenv("BLUESPOT_ATTACH_UPLOAD_CONFIG_FILE_DIR_PATH", "/tmp/config")
				t.Setenv("BLUESPOT_ATTACH_VIEW_CONFIG_FILE_BASE_URL", "https://example.test/config")
				t.Setenv("BLUESPOT_ATTACH_UPLOAD_CONFIG_FILE_SECRET", " secret ")
				bindEnvVars(loader)
			} else {
				loader.SetConfigType("yaml")
				if err := loader.ReadConfig(strings.NewReader(`attach:
  upload_config_file_dir_path: /tmp/config
  view_config_file_base_url: https://example.test/config
  upload_config_file_secret: " secret "
`)); err != nil {
					t.Fatal(err)
				}
			}
			GlobalConfig = &Config{}
			if err := loader.Unmarshal(GlobalConfig); err != nil {
				t.Fatal(err)
			}
			if GetAttachConfigFileUploadDirPath() != "/tmp/config" || GetAttachConfigFileViewBaseURL() != "https://example.test/config" || GetAttachConfigFileUploadSecret() != " secret " {
				t.Fatalf("unexpected config: %+v", GlobalConfig.Attach)
			}
		})
	}
	GlobalConfig = nil
	if GetAttachConfigFileUploadDirPath() != "" || GetAttachConfigFileViewBaseURL() != "" || GetAttachConfigFileUploadSecret() != "" {
		t.Fatal("nil config should return empty values")
	}
}
