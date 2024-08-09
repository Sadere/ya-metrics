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
					ServeGRPC: false,
					Address: common.NetAddress{
						Host: "localhost",
						Port: 1111,
					},
					LogLevel:        "fatal",
					StoreInterval:   DefaultStoreInterval,
					FileStoragePath: DefaultFileStoragePath,
					Restore:         true,
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
					ServeGRPC: false,
					Address: common.NetAddress{
						Host: "localhost",
						Port: 2222,
					},
					LogLevel:        "fatal",
					StoreInterval:   DefaultStoreInterval,
					FileStoragePath: DefaultFileStoragePath,
					Restore:         true,
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
			name: "store interval from arg",
			args: []string{"-i", "111"},
			want: want{
				conf: Config{
					ServeGRPC: false,
					Address: common.NetAddress{
						Host: "localhost",
						Port: 8080,
					},
					LogLevel:        "fatal",
					StoreInterval:   111,
					FileStoragePath: DefaultFileStoragePath,
					Restore:         true,
				},
				err: false,
			},
		},
		{
			name: "store interval from env",
			args: []string{"-i", "111"},
			env: map[string]string{
				"STORE_INTERVAL": "222",
			},
			want: want{
				conf: Config{
					ServeGRPC: false,
					Address: common.NetAddress{
						Host: "localhost",
						Port: 8080,
					},
					LogLevel:        "fatal",
					StoreInterval:   222,
					FileStoragePath: DefaultFileStoragePath,
					Restore:         true,
				},
				err: false,
			},
		},
		{
			name: "malformed store interval from env",
			args: []string{"-i", "111"},
			env: map[string]string{
				"STORE_INTERVAL": "nop",
			},
			want: want{
				conf: Config{
					ServeGRPC: false,
					Address: common.NetAddress{
						Host: "localhost",
						Port: 8080,
					},
					LogLevel:        "fatal",
					StoreInterval:   DefaultStoreInterval,
					FileStoragePath: DefaultFileStoragePath,
					Restore:         true,
				},
				err: false,
			},
		},
		{
			name: "file path from arg",
			args: []string{"-f", "arg-file.db"},
			want: want{
				conf: Config{
					ServeGRPC: false,
					Address: common.NetAddress{
						Host: "localhost",
						Port: 8080,
					},
					LogLevel:        "fatal",
					StoreInterval:   DefaultStoreInterval,
					FileStoragePath: "arg-file.db",
					Restore:         true,
				},
				err: false,
			},
		},
		{
			name: "file path from env",
			args: []string{"-f", "arg-file.db"},
			env: map[string]string{
				"FILE_STORAGE_PATH": "env-file.db",
			},
			want: want{
				conf: Config{
					ServeGRPC: false,
					Address: common.NetAddress{
						Host: "localhost",
						Port: 8080,
					},
					LogLevel:        "fatal",
					StoreInterval:   DefaultStoreInterval,
					FileStoragePath: "env-file.db",
					Restore:         true,
				},
				err: false,
			},
		},
		{
			name: "restore from arg",
			args: []string{"-r", "true"},
			want: want{
				conf: defCfg,
				err:  false,
			},
		},
		{
			name: "restore from env",
			args: []string{"-r", "false"},
			env: map[string]string{
				"RESTORE": "true",
			},
			want: want{
				conf: defCfg,
				err:  false,
			},
		},
		{
			name: "DSN from arg",
			args: []string{"-d", "DSN_ARG"},
			want: want{
				conf: Config{
					ServeGRPC: false,
					Address: common.NetAddress{
						Host: "localhost",
						Port: 8080,
					},
					LogLevel:        "fatal",
					StoreInterval:   DefaultStoreInterval,
					FileStoragePath: DefaultFileStoragePath,
					PostgresDSN:     "DSN_ARG",
					Restore:         true,
				},
				err: false,
			},
		},
		{
			name: "DSN from env",
			args: []string{"-d", "DSN_ARG"},
			env: map[string]string{
				"DATABASE_DSN": "DSN_ENV",
			},
			want: want{
				conf: Config{
					ServeGRPC: false,
					Address: common.NetAddress{
						Host: "localhost",
						Port: 8080,
					},
					LogLevel:        "fatal",
					StoreInterval:   DefaultStoreInterval,
					FileStoragePath: DefaultFileStoragePath,
					PostgresDSN:     "DSN_ENV",
					Restore:         true,
				},
				err: false,
			},
		},
		{
			name: "hash key from arg",
			args: []string{"-k", "hash-key-arg"},
			want: want{
				conf: Config{
					ServeGRPC: false,
					Address: common.NetAddress{
						Host: "localhost",
						Port: 8080,
					},
					LogLevel:        "fatal",
					StoreInterval:   DefaultStoreInterval,
					FileStoragePath: DefaultFileStoragePath,
					HashKey:         "hash-key-arg",
					Restore:         true,
				},
				err: false,
			},
		},
		{
			name: "hash key from env",
			args: []string{"-k", "hash-key-arg"},
			env: map[string]string{
				"KEY": "hash-key-env",
			},
			want: want{
				conf: Config{
					ServeGRPC: false,
					Address: common.NetAddress{
						Host: "localhost",
						Port: 8080,
					},
					LogLevel:        "fatal",
					StoreInterval:   DefaultStoreInterval,
					FileStoragePath: DefaultFileStoragePath,
					HashKey:         "hash-key-env",
					Restore:         true,
				},
				err: false,
			},
		},
		{
			name: "priv key from arg",
			args: []string{"-crypto-key", "arg.pem"},
			want: want{
				conf: Config{
					ServeGRPC: false,
					Address: common.NetAddress{
						Host: "localhost",
						Port: 8080,
					},
					LogLevel:        "fatal",
					StoreInterval:   DefaultStoreInterval,
					FileStoragePath: DefaultFileStoragePath,
					PrivateKeyPath:  "arg.pem",
					Restore:         true,
				},
				err: false,
			},
		},
		{
			name: "priv key from env",
			args: []string{"-crypto-key", "arg.pem"},
			env: map[string]string{
				"CRYPTO_KEY": "env.pem",
			},
			want: want{
				conf: Config{
					ServeGRPC: false,
					Address: common.NetAddress{
						Host: "localhost",
						Port: 8080,
					},
					LogLevel:        "fatal",
					StoreInterval:   DefaultStoreInterval,
					FileStoragePath: DefaultFileStoragePath,
					PrivateKeyPath:  "env.pem",
					Restore:         true,
				},
				err: false,
			},
		},
		{
			name: "trusted net from arg",
			args: []string{"-t", "10.0.0.0/24"},
			want: want{
				conf: Config{
					ServeGRPC: false,
					Address: common.NetAddress{
						Host: "localhost",
						Port: 8080,
					},
					LogLevel:        "fatal",
					StoreInterval:   DefaultStoreInterval,
					FileStoragePath: DefaultFileStoragePath,
					TrustedSubnet:   "10.0.0.0/24",
					Restore:         true,
				},
				err: false,
			},
		},
		{
			name: "trusted net from env",
			args: []string{"-t", "10.0.0.0/24"},
			env: map[string]string{
				"TRUSTED_SUBNET": "11.0.0.0/24",
			},
			want: want{
				conf: Config{
					ServeGRPC: false,
					Address: common.NetAddress{
						Host: "localhost",
						Port: 8080,
					},
					LogLevel:        "fatal",
					StoreInterval:   DefaultStoreInterval,
					FileStoragePath: DefaultFileStoragePath,
					TrustedSubnet:   "11.0.0.0/24",
					Restore:         true,
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
		"address": "localhost:1336",
		"restore": true,
		"store_interval": 1,
		"store_file": "/path/to/file.db",
		"database_dsn": "",
		"crypto_key": "/path/to/key.pem"
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
		ServeGRPC: false,
		Address: common.NetAddress{
			Host: "localhost",
			Port: 1336,
		},
		StoreInterval:   1,
		FileStoragePath: "/path/to/file.db",
		Restore:         true,
		PrivateKeyPath:  "/path/to/key.pem",
	}

	if !reflect.DeepEqual(conf, expectedCfg) {
		t.Errorf("conf got %+v, want %+v", conf, expectedCfg)
	}
}
