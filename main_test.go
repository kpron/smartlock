package main

import "errors"
import "net/http"
import "net/http/httptest"
import "net/url"
import "bytes"
import "testing"
import "time"

func TestCheckError(t *testing.T) {
	checkError(errors.New("Test error"))
}

func TestLock(t *testing.T) {
	err := lock()
	if err != nil {
		t.Fail()
	}
}

func TestUnlock(t *testing.T) {
	err := unlock()
	if err != nil {
		t.Fail()
	}
}

func TestToggle(t *testing.T) {
	testCases := []struct {
		method string
		token  string
		status int
		answer string
	}{
		{"POST", "HACK!", http.StatusOK, "wrong"},
		{"GET", "shit in get", http.StatusMethodNotAllowed, ""},
		{"POST", "test", http.StatusOK, "lolked"},
		{"POST", "test", http.StatusOK, "unlolked"},
	}
	for _, tc := range testCases {
		data := url.Values{}
		data.Set("key", tc.token)
		buffer := new(bytes.Buffer)
		buffer.WriteString(data.Encode())
		req, err := http.NewRequest(tc.method, "/", buffer)
		req.Header.Set("content-type", "application/x-www-form-urlencoded")
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := func(w http.ResponseWriter, r *http.Request) {
			toggle(w, r, "test")
		}

		handler(rr, req)

		if status := rr.Code; status != tc.status {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
		expected := tc.answer

		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), expected)
		}
		time.Sleep(time.Second)
	}
}
