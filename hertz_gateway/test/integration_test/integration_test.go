package integrationtest

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/Linda-ui/orbital_HeBao/hertz_gateway/config"
)

func TestIntegrationGateway(t *testing.T) {

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to set cwd to integration test directory")
	}

	paths, err := filepath.Glob(filepath.Join(path.Dir(filename), "testdata", "**", "*.input"))
	if err != nil {
		t.Fatal(err)
	}

	var url string
	for _, path := range paths {

		directory, filename := filepath.Split(path)
		gatewayAddr := "http://" + config.ServiceAddr

		switch filepath.Base(directory) {
		case "gateway":
			url = gatewayAddr
		case "echo":
			url = gatewayAddr + "/gateway/echo/EchoMethod"
		case "sum":
			url = gatewayAddr + "/gateway/sum/SumMethod"
		case "noService":
			url = gatewayAddr + "/gateway/echoXXX/EchoMethod"
		case "noServiceMethod":
			url = gatewayAddr + "/gateway/echo/EchoMethodXXX"
		default:
			url = gatewayAddr + "/gateway"
		}

		// removing the file extension (.input) to obtain the test name
		testname := filename[:len(filename)-len(filepath.Ext(filename))]

		t.Run(testname, func(t *testing.T) {
			payload, err := os.ReadFile(path)
			if err != nil {
				t.Fatal("error reading source file:", err)
			}

			var resp *http.Response
			if filepath.Base(directory) == "gateway" {
				resp, err = http.Get(url)
				if err != nil {
					t.Fatalf("Failed to make request: %v", err)
				}
			} else {
				resp, err = http.Post(url, "application/json", bytes.NewBuffer(payload))
				if err != nil {
					t.Fatalf("Failed to make request: %v", err)
				}
			}

			defer resp.Body.Close()

			responseBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to read response body: %v", err)
			}

			goldenfile := filepath.Join(directory, testname+".golden")
			expectedResponse, err := os.ReadFile(goldenfile)
			if err != nil {
				t.Fatal("error reading golden file:", err)
			}

			if !bytes.Equal(responseBody, expectedResponse) {
				t.Errorf("Response body does not match expected value.\nExpected: %s\nActual: %s", string(expectedResponse), string(responseBody))
			}
		})
	}
}
