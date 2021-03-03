package utils

import (
	"os"
	"reflect"
	"testing"

	"github.com/spf13/viper"
)

func Test_initViperConfig(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "sets_default_values",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitViperConfig()
			scope := viper.Get(configScope)
			if scope != defaultValues[configScope] {
				t.Errorf("viper's scope, expected %v, got %v",
					scope,
					defaultValues[configScope])
			}

			ip := viper.Get(configServerIP)
			if ip != defaultValues[configServerIP] {
				t.Errorf("viper's scope, expected %v, got %v",
					ip,
					defaultValues[configServerIP])
			}
		})
	}
}

func Test_readViperConfig(t *testing.T) {
	type args struct {
		baseDir string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "working_read",
			args:    args{baseDir: "./test"},
			wantErr: false,
		},
		{
			name:    "not_working_read",
			args:    args{baseDir: "1/2/3/"},
			wantErr: true,
		},
	}

	// sets defaults values in ./test/.inspr/config
	setupViperTest()
	// sets the
	viper.SetConfigFile("./test/.inspr/config")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ReadViperConfig(tt.args.baseDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("readViperConfig() error = %v, wantErr %v",
					err,
					tt.wantErr)
			}
		})
	}
	os.Remove("./test/.inspr/config")
	os.Remove("./test/.inspr")
}

func Test_changeViperValues(t *testing.T) {
	type args struct {
		key   string
		value interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "changing_scope",
			args: args{
				key:   configScope,
				value: "new_scope",
			},
			wantErr: false,
		},
		{
			name: "changing_IP",
			args: args{
				key:   configServerIP,
				value: "XXX.YYY.ZZZ.0",
			},
			wantErr: false,
		},
		{
			name: "error_writing",
			args: args{
				key:   configServerIP,
				value: nil,
			},
			wantErr: true,
		},
	}

	// read mock config values
	setupViperTest()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				viper.SetConfigFile("/etc/")
			}

			// reads the current values of the viper config
			ReadViperConfig("./test")

			err := ChangeViperValues(tt.args.key, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("changeViperValues() error = %v, wantErr %v",
					err,
					tt.wantErr)
			}

			got := viper.Get(tt.args.key)
			if !reflect.DeepEqual(got, tt.args.value) && !tt.wantErr {
				t.Errorf("viper.Get(key) got = %v, want %v",
					got,
					tt.args.value)
			}
		})
	}
	os.Remove("./test/.inspr/config")
	os.Remove("./test/.inspr")
}

func Test_createViperConfig(t *testing.T) {
	type args struct {
		folderPath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "local_viper_config",
			args:    args{folderPath: "./test/config"},
			wantErr: false,
		},
		{
			name:    "error_folder_location",
			args:    args{folderPath: "/1//2/3/4/5"},
			wantErr: true,
		},
	}

	// read mock config values
	setupViperTest()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := createViperConfig(tt.args.folderPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("createViperConfig() error = %v, wantErr %v",
					err,
					tt.wantErr)
			}
		})
	}
	os.Remove("./test/config")
}

func Test_createInsprConfigFolder(t *testing.T) {
	type args struct {
		folderPath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "create_folder",
			args:    args{folderPath: "./test/inspr"},
			wantErr: false,
		},
		{
			name:    "create_folder",
			args:    args{folderPath: "1/inspr"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := createInsprConfigFolder(tt.args.folderPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("createInsprConfigFolder() error = %v, wantErr %v",
					err,
					tt.wantErr)
			}
		})
	}
	os.Remove("./test/inspr")
}

func Test_existingKeys(t *testing.T) {
	tests := []struct {
		name string
		want []string
	}{
		{
			name: "all_keys",
			want: []string{"extra", "scope", "serverip"},
		},
	}

	// read mock config values
	setupViperTest()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExistingKeys()

			// doesn't return in a specific order, so we use maps to compare
			receivedValues := make(map[string]bool)
			for _, k := range got {
				receivedValues[k] = true
			}

			// checking values
			for _, k := range tt.want {
				if receivedValues[k] == false {
					t.Errorf("existingKeys() => %v doesn't exist but is expected",
						k)
				}
			}
		})
	}
}

/// test utils functions

func setupViperTest() {
	// specifies the path in which the config file present
	viper.AddConfigPath("./test/")
	viper.SetConfigName("viper_config")
	viper.SetConfigType("yaml")

	// contains defaults values to be used in others functions
	viper.ReadInConfig()
}
