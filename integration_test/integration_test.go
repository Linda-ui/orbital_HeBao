package integration_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

func TestIntegrationGateway(t *testing.T) {

	paths, err := filepath.Glob(filepath.Join("testdata", "**", "*.input"))
	if err != nil {
		t.Fatal(err)
	}

	var url string
	for _, path := range paths {

		directory, filename := filepath.Split(path)

		switch directory {
		case "testdata/gateway/":
			url = "http://localhost:8080/"
		case "testdata/echo/":
			url = "http://localhost:8080/gateway/echo/EchoMethod"
		case "testdata/sum/":
			url = "http://localhost:8080/gateway/sum/SumMethod"
		case "testdata/noService/":
			url = "http://localhost:8080/gateway/echoXXX/EchoMethod"
		case "testdata/noServiceMethod/":
			url = "http://localhost:8080/gateway/echo/EchoMethodXXX"
		default:
			url = "http://localhost:8080/gateway"
		}

		testname := filename[:len(filename)-len(filepath.Ext(path))]

		t.Run(testname, func(t *testing.T) {
			payload, err := os.ReadFile(path)
			if err != nil {
				t.Fatal("error reading source file:", err)
			}

			var resp *http.Response
			if directory == "testdata/gateway/" {
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
