package main

import (
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"
)

var host = "127.0.0.1"
var port = "8002"
var apiport = "9999"

func TestScacheGet(t *testing.T) {
	cases := []struct {
		key, exp string
		expCode  int
	}{
		{"Tom", "630", 200},
		{"Jack", "589", 200},
		{"Sam", "567", 200},
	}

	for _, c := range cases {
		url := "http://" + host + ":" + port + "/scache/scores/" + c.key
		req, _ := http.NewRequest("GET", url, nil)
		response, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		t.Log("response status code is", response.StatusCode)
		body, _ := io.ReadAll(response.Body)
		t.Log("response body is ", string(body))
		if response.StatusCode != c.expCode {
			t.Fatal("unexpected code :", response.StatusCode, "exp:", c.expCode)
		}
		if string(body) != c.exp {
			t.Fatal("unepxected body:", string(body), "exp:", c.exp)
		}
	}
}

func TestScacheConcurrentGet(t *testing.T) {
	cases := []struct {
		key, exp string
		expCode  int
	}{
		{"Tom", "630", 200},
		{"Tom", "630", 200},
		{"Tom", "630", 200},
		{"Tom", "630", 200},
		{"Tom", "630", 200},
		{"Tom", "630", 200},
		{"Tom", "630", 200},
		{"Tom", "630", 200},
		{"Tom", "630", 200},
		{"Tom", "630", 200},
		{"Tom", "630", 200},
		{"Tom", "630", 200},
		{"Tom", "630", 200},
	}

	numCase := len(cases)
	var wg sync.WaitGroup
	wg.Add(numCase)

	for _, c := range cases {
		url := "http://" + host + ":" + port + "/scache/scores/" + c.key
		req, _ := http.NewRequest("GET", url, nil)
		go func() {
			defer wg.Done()
			response, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			t.Log("response status code is", response.StatusCode)
			body, _ := io.ReadAll(response.Body)
			t.Log("response body is ", string(body))
			if response.StatusCode != c.expCode {
				t.Fatal("unexpected code :", response.StatusCode, "exp:", c.expCode)
			}
			if string(body) != c.exp {
				t.Fatal("unepxected body:", string(body), "exp:", c.exp)
			}
		}()
	}
	wg.Wait()
}

func TestScacheGetLoad(t *testing.T) {
	cases := []struct {
		key, exp string
		expCode  int
	}{
		{"Tom", "630", 200},
		{"Jack", "589", 200},
		{"Sam", "567", 200},
	}

	for _, c := range cases {
		url := "http://" + host + ":" + port + "/scache/scores/" + c.key
		req, _ := http.NewRequest("GET", url, nil)
		req, _ = http.NewRequest("GET", url, nil)
		response, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		t.Log("response status code is", response.StatusCode)
		body, _ := io.ReadAll(response.Body)
		t.Log("response body is ", string(body))
		if response.StatusCode != c.expCode {
			t.Fatal("unexpected code :", response.StatusCode, "exp:", c.expCode)
		}
		if string(body) != c.exp {
			t.Fatal("unepxected body:", string(body), "exp:", c.exp)
		}
		if LoadCounts[c.key] > 1 {
			t.Errorf("cache miss:key %s loaded %d times", c.key, LoadCounts[c.key])
		}
	}

}

func TestScacheGetWithUnknownCache(t *testing.T) {
	url := "http://" + host + ":" + port + "/unkowncache/scores/TOM"
	req, _ := http.NewRequest("GET", url, nil)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("response status code is", response.StatusCode)
	body, _ := io.ReadAll(response.Body)
	t.Log("response body is ", string(body))
	if response.StatusCode != 404 {
		t.Fatal("unexpected code :", response.StatusCode, "exp:", 404)
	}
	if !strings.Contains(string(body), "no such cache server: /unkowncache/scores/TOM") {
		t.Fatal("unepxected body:", string(body), "exp:", "no such cache server: /unkowncache/scores/TOM")
	}
}

func TestScacheGetWithUnknownGroup(t *testing.T) {
	url := "http://" + host + ":" + port + "/scache/unkowngroup/TOM"
	req, _ := http.NewRequest("GET", url, nil)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("response status code is", response.StatusCode)
	body, _ := io.ReadAll(response.Body)
	t.Log("response body is ", string(body))
	if response.StatusCode != 404 {
		t.Fatal("unexpected code :", response.StatusCode, "exp:", 404)
	}
	if !strings.Contains(string(body), "no such group: unkowngroup") {
		t.Fatal("unepxected body:", string(body), "exp:", "no such group: unkowngroup")
	}
}

func TestScacheGetWithUnknownKey(t *testing.T) {
	url := "http://" + host + ":" + port + "/scache/scores/unkey"
	req, _ := http.NewRequest("GET", url, nil)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("response status code is", response.StatusCode)
	body, _ := io.ReadAll(response.Body)
	t.Log("response body is ", string(body))
	if response.StatusCode != 404 {
		t.Fatal("unexpected code :", response.StatusCode, "exp:", 404)
	}
	if !strings.Contains(string(body), "unkey not exist") {
		t.Fatal("unepxected body:", string(body), "exp:", "unkey not exist")
	}
}

func TestScacheGetAPI(t *testing.T) {
	cases := []struct {
		key, exp string
		expCode  int
	}{
		{"Tom", "630", 200},
		{"Jack", "589", 200},
		{"Sam", "567", 200},
	}

	for _, c := range cases {
		url := "http://" + host + ":" + apiport + "/api?key=" + c.key
		t.Log("request GET: " + url)
		req, _ := http.NewRequest("GET", url, nil)
		response, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		t.Log("response status code is", response.StatusCode)
		body, _ := io.ReadAll(response.Body)
		t.Log("response body is ", string(body))
		if response.StatusCode != c.expCode {
			t.Fatal("unexpected code :", response.StatusCode, "exp:", c.expCode)
		}
		if string(body) != c.exp {
			t.Fatal("unepxected body:", string(body), "exp:", c.exp)
		}
	}
}
