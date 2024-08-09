package config

import (
	"os"
	"reflect"
	"testing"

	"github.com/Sadere/ya-metrics/internal/common"
	"github.com/stretchr/testify/assert"
)

func TestConfigFromArgsEnv(t *testing.T) {
	defCfg := defaultConfig()

	type want struct {
		conf Config
		err  bool
	}
	var tests = []struct {
		name string
		args []string
		env  map[string]string
		want want
	}{
		{
			name: "address from arg",
			args: []string{"-a", "localhost:1111"},
			want: want{
				conf: Config{
					ServerAddress: common.NetAddress{
						Host: "localhost",
						Port: 1111,
					},
					PollInterval:   DefaultPollInterval,
					ReportInterval: DefaultReportInterval,
					RateLimit:      DefaultRateLimit,
				},
				err: false,
			},
		},
		{
			name: "address from env",
			args: []string{"-a", "localhost:1111"},
			env: map[string]string{
				"ADDRESS": "localhost:2222",
			},
			want: want{
				conf: Config{
					ServerAddress: common.NetAddress{
						Host: "localhost",
						Port: 2222,
					},
					PollInterval:   DefaultPollInterval,
					ReportInterval: DefaultReportInterval,
					RateLimit:      DefaultRateLimit,
				},
				err: false,
			},
		},
		{
			name: "invalid address",
			env: map[string]string{
				"ADDRESS": "nop",
			},
			want: want{
				conf: defCfg,
				err:  true,
			},
		},
		{
			name: "poll interval from arg",
			args: []string{"-p", "999"},
			want: want{
				conf: Config{
					ServerAddress: common.NetAddress{
						Host: "localhost",
						Port: 8080,
					},
					PollInterval:   999,
					ReportInterval: DefaultReportInterval,
					RateLimit:      DefaultRateLimit,
				},
				err: false,
			},
		},
		{
			name: "poll interval from env",
			args: []string{"-p", "999"},
			env: map[string]string{
				"POLL_INTERVAL": "333",
			},
			want: want{
				conf: Config{
					ServerAddress: common.NetAddress{
						Host: "localhost",
						Port: 8080,
					},
					PollInterval:   333,
					ReportInterval: DefaultReportInterval,
					RateLimit:      DefaultRateLimit,
				},
				err: false,
			},
		},
		{
			name: "malformed poll interval from env",
			args: []string{"-p", "999"},
			env: map[string]string{
				"POLL_INTERVAL": "nop",
			},
			want: want{
				conf: Config{
					ServerAddress: common.NetAddress{
						Host: "localhost",
						Port: 8080,
					},
					PollInterval:   DefaultPollInterval,
					ReportInterval: DefaultReportInterval,
					RateLimit:      DefaultRateLimit,
				},
				err: false,
			},
		},
		{
			name: "report interval from arg",
			args: []string{"-r", "999"},
			want: want{
				conf: Config{
					ServerAddress: common.NetAddress{
						Host: "localhost",
						Port: 8080,
					},
					PollInterval:   DefaultPollInterval,
					ReportInterval: 999,
					RateLimit:      DefaultRateLimit,
				},
				err: false,
			},
		},
		{
			name: "report interval from env",
			args: []string{"-r", "999"},
			env: map[string]string{
				"REPORT_INTERVAL": "333",
			},
			want: want{
				conf: Config{
					ServerAddress: common.NetAddress{
						Host: "localhost",
						Port: 8080,
					},
					PollInterval:   DefaultPollInterval,
					ReportInterval: 333,
					RateLimit:      DefaultRateLimit,
				},
				err: false,
			},
		},
		{
			name: "malformed report interval from env",
			args: []string{"-r", "999"},
			env: map[string]string{
				"REPORT_INTERVAL": "nop",
			},
			want: want{
				conf: Config{
					ServerAddress: common.NetAddress{
						Host: "localhost",
						Port: 8080,
					},
					PollInterval:   DefaultPollInterval,
					ReportInterval: DefaultReportInterval,
					RateLimit:      DefaultRateLimit,
				},
				err: false,
			},
		},
		{
			name: "rate limit from arg",
			args: []string{"-l", "11"},
			want: want{
				conf: Config{
					ServerAddress: common.NetAddress{
						Host: "localhost",
						Port: 8080,
					},
					PollInterval:   DefaultPollInterval,
					ReportInterval: DefaultReportInterval,
					RateLimit:      11,
				},
				err: false,
			},
		},
		{
			name: "rate limit from env",
			args: []string{"-l", "11"},
			env: map[string]string{
				"RATE_LIMIT": "22",
			},
			want: want{
				conf: Config{
					ServerAddress: common.NetAddress{
						Host: "localhost",
						Port: 8080,
					},
					PollInterval:   DefaultPollInterval,
					ReportInterval: DefaultReportInterval,
					RateLimit:      22,
				},
				err: false,
			},
		},
		{
			name: "malformed rate limit from env",
			args: []string{"-l", "11"},
			env: map[string]string{
				"RATE_LIMIT": "nop",
			},
			want: want{
				conf: Config{
					ServerAddress: common.NetAddress{
						Host: "localhost",
						Port: 8080,
					},
					PollInterval:   DefaultPollInterval,
					ReportInterval: DefaultReportInterval,
					RateLimit:      DefaultRateLimit,
				},
				err: false,
			},
		},
		{
			name: "hash key from arg",
			args: []string{"-k", "test_key"},
			want: want{
				conf: Config{
					ServerAddress: common.NetAddress{
						Host: "localhost",
						Port: 8080,
					},
					PollInterval:   DefaultPollInterval,
					ReportInterval: DefaultReportInterval,
					RateLimit:      DefaultRateLimit,
					HashKey:        "test_key",
				},
				err: false,
			},
		},
		{
			name: "hash key from env",
			args: []string{"-k", "test_key"},
			env: map[string]string{
				"KEY": "env_key",
			},
			want: want{
				conf: Config{
					ServerAddress: common.NetAddress{
						Host: "localhost",
						Port: 8080,
					},
					PollInterval:   DefaultPollInterval,
					ReportInterval: DefaultReportInterval,
					RateLimit:      DefaultRateLimit,
					HashKey:        "env_key",
				},
				err: false,
			},
		},
		{
			name: "pub key from arg",
			args: []string{"-crypto-key", "./arg-crypto.key"},
			want: want{
				conf: Config{
					ServerAddress: common.NetAddress{
						Host: "localhost",
						Port: 8080,
					},
					PollInterval:   DefaultPollInterval,
					ReportInterval: DefaultReportInterval,
					RateLimit:      DefaultRateLimit,
					PubKeyFilePath: "./arg-crypto.key",
				},
				err: false,
			},
		},
		{
			name: "pub key from env",
			args: []string{"-crypto-key", "./arg-crypto.key"},
			env: map[string]string{
				"CRYPTO_KEY": "./env-crypto.key",
			},
			want: want{
				conf: Config{
					ServerAddress: common.NetAddress{
						Host: "localhost",
						Port: 8080,
					},
					PollInterval:   DefaultPollInterval,
					ReportInterval: DefaultReportInterval,
					RateLimit:      DefaultRateLimit,
					PubKeyFilePath: "./env-crypto.key",
				},
				err: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key, value := range tt.env {
				t.Setenv(key, value)
			}

			conf, err := NewConfig(tt.args)

			if tt.want.err {
				assert.Error(t, err)
				t.SkipNow()
			} else {
				assert.NoError(t, err)
			}

			if !reflect.DeepEqual(conf, tt.want.conf) {
				t.Errorf("conf got %+v, want %+v", conf, tt.want.conf)
			}

			// Remove envs for next tests
			for key := range tt.env {
				t.Setenv(key, "")
			}
		})
	}
}

func TestConfigFromFile(t *testing.T) {
	// Создаем временный файл с конфигурацией
	f, err := os.CreateTemp("", "config-test")
	if err != nil {
		t.Errorf("failed to create temp file: %s", err)
	}

	// Удаляем файл после прохождения теста
	defer func() {
		err = os.Remove(f.Name())

		if err != nil {
			t.Errorf("failed to remove temp file: %s", err)
		}
	}()

	// Пишем в тестовый файл
	testCfgContent := `
	{
		"address": "localhost:5555",
		"report_interval": 11,
		"poll_interval": 22,
		"crypto_key": "pub.pem"
	}
	`

	_, err = f.WriteString(testCfgContent)
	if err != nil {
		t.Errorf("failed to write to temp file: %s", err)
	}

	err = f.Close()
	if err != nil {
		t.Errorf("failed to close temp file: %s", err)
	}

	args := []string{"-c", f.Name()}

	conf, err := NewConfig(args)

	assert.NoError(t, err)

	expectedCfg := Config{
		ServerAddress: common.NetAddress{
			Host: "localhost",
			Port: 5555,
		},
		PollInterval:   22,
		ReportInterval: 11,
		PubKeyFilePath: "pub.pem",
	}

	if !reflect.DeepEqual(conf, expectedCfg) {
		t.Errorf("conf got %+v, want %+v", conf, expectedCfg)
	}
}
