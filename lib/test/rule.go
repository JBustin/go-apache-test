package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-apache-test/lib/utils"
)

type Rule struct {
	Name         string
	Host         string
	Pathname     string
	Skip         bool
	Headers      map[string]string
	Expectations Expectations
}

func (r *Rule) Init(others ...Expectations) {
	if r.Name == "" {
		r.Name = fmt.Sprintf("%v%v", r.Host, r.Pathname)
	}

	for _, other := range others {
		r.Expectations.Backend = utils.ExtendStr(r.Expectations.Backend, other.Backend)
		r.Expectations.Host = utils.ExtendStr(r.Expectations.Host, other.Host)
		r.Expectations.Pathname = utils.ExtendStr(r.Expectations.Pathname, other.Pathname)
		r.Expectations.HttpStatusCode = utils.ExtendInt(r.Expectations.HttpStatusCode, other.HttpStatusCode)
	}

	r.Expectations.Host = utils.ExtendStr(r.Expectations.Host, r.Host)
	r.Expectations.Pathname = utils.ExtendStr(r.Expectations.Pathname, r.Pathname)
}

func (r *Rule) Request() (Expectations, error) {
	var result Expectations

	url := fmt.Sprintf("http://%v%v", r.Host, r.Pathname)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return result, err
	}
	defer func() { req.Close = true }()
	for key, value := range r.Headers {
		req.Header.Add(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	if resp.StatusCode >= 400 {
		result.HttpStatusCode = resp.StatusCode
		return result, nil
	}

	err = json.Unmarshal(body, &result)
	result.HttpStatusCode = resp.StatusCode

	return result, err
}

func (r Rule) ShouldBe(result Expectations) error {
	if r.Expectations.HttpStatusCode != result.HttpStatusCode {
		return fmt.Errorf(`
	%v
	Error with Status code
		Expected: %v
		Received: %v`, r.Name, r.Expectations.HttpStatusCode, result.HttpStatusCode)
	}

	if result.HttpStatusCode >= 400 {
		return nil
	}

	if r.Host != result.Host {
		return fmt.Errorf(`
	%v
	Error with Host
		Expected: %v
		Received: %v`, r.Name, r.Expectations.Host, result.Host)
	}

	if r.Expectations.Backend != result.Backend {
		return fmt.Errorf(`
	%v
	Error with Backend
		Expected: %v
		Received: %v`, r.Name, r.Expectations.Backend, result.Backend)
	}

	if r.Expectations.Pathname != result.Pathname {
		return fmt.Errorf(`
	%v
	Error with Pathname
		Expected: %v
		Received: %v`, r.Name, r.Expectations.Pathname, result.Pathname)
	}

	return nil
}

func (r Rule) StrSuccess() string {
	return fmt.Sprintf("\t\t✅ %v", r.Name)
}

func (r Rule) StrError(err error) string {
	return fmt.Sprintf("\t\t❌ %v", r.Name)
}
