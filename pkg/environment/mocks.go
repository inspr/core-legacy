package environment

import "os"

// SetMockEnv - sets the environment variables to specific values
func SetMockEnv() {
	channelsValues := "chan;chan1;chan2;chan3"
	os.Setenv("INSPR_LBSIDECAR_PORT", "8888")
	os.Setenv("INSPR_INPUT_CHANNELS", channelsValues)
	os.Setenv("INSPR_OUTPUT_CHANNELS", channelsValues)
	os.Setenv("INSPR_UNIX_SOCKET", "socket_addr")
	os.Setenv("INSPR_SIDECAR_IMAGE", "mock_sidecar_image")
	os.Setenv("INSPR_APP_CTX", "mock.dapp.context")
	os.Setenv("INSPR_ENV", "mock_env")
	os.Setenv("INSPR_APP_ID", "testappid1")
	os.Setenv("INSPR_LBSIDECAR_WRITE_PORT", "3001")
	os.Setenv("INSPR_LBSIDECAR_READ_PORT", "3002")
	os.Setenv("INSPR_SIDECAR_TEST_WRITE_PORT", "1234")
	os.Setenv("INSPR_SIDECAR_TEST_READ_PORT", "4321")
	os.Setenv("INSPR_SIDECAR_TEST_ADDR", "someAddr")
}

// UnsetMockEnv - removes the values of the environment variables
func UnsetMockEnv() {
	os.Unsetenv("INSPR_INPUT_CHANNELS")
	os.Unsetenv("INSPR_LBSIDECAR_PORT")
	os.Unsetenv("INSPR_OUTPUT_CHANNELS")
	os.Unsetenv("INSPR_UNIX_SOCKET")
	os.Unsetenv("INSPR_SIDECAR_IMAGE")
	os.Unsetenv("INSPR_APP_CTX")
	os.Unsetenv("INSPR_ENV")
	os.Unsetenv("INSPR_APP_ID")
	os.Unsetenv("INSPR_LBSIDECAR_WRITE_PORT")
	os.Unsetenv("INSPR_LBSIDECAR_READ_PORT")
	os.Unsetenv("INSPR_LBSIDECAR_PORT")
}
