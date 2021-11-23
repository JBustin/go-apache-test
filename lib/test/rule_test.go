package test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_RuleInit(t *testing.T) {
	r := Rule{
		Name:     "product",
		Host:     "www.mydomain.com",
		Pathname: "product/1",
		Skip:     false,
		Headers:  map[string]string{},
		Expectations: Expectations{
			Backend:        "",
			Host:           "",
			Pathname:       "",
			HttpStatusCode: 0,
			Body:           "",
		},
	}

	r.Init(Expectations{
		Backend:        "newstack.com",
		Host:           "",
		Pathname:       "",
		HttpStatusCode: 200,
		Body:           "",
	}, Expectations{
		Backend:        "",
		Host:           "",
		Pathname:       "",
		HttpStatusCode: 0,
		Body:           "",
	})

	assert.Equal(t, "newstack.com", r.Expectations.Backend, "should populate the backend")
	assert.Equal(t, "www.mydomain.com", r.Expectations.Host, "should populate the host")
	assert.Equal(t, "product/1", r.Expectations.Pathname, "should populate the pathname")
	assert.Equal(t, 200, r.Expectations.HttpStatusCode, "should populate the status code")
}

func Test_RuleShould(t *testing.T) {
	r := Rule{
		Name:     "product",
		Host:     "www.mydomain.com",
		Pathname: "product/1",
		Skip:     false,
		Headers:  map[string]string{},
		Expectations: Expectations{
			Backend:        "newstack.com",
			Host:           "www.mydomain.com",
			Pathname:       "product/1",
			HttpStatusCode: 200,
			Body:           "",
		},
	}

	err := r.ShouldBe(Expectations{
		Backend:        "newstack.com",
		Host:           "www.mydomain.com",
		Pathname:       "product/1",
		HttpStatusCode: 200,
		Body:           "",
	})

	assert.Equal(t, nil, err, "should valid the expectations")

	err = r.ShouldBe(Expectations{
		Backend:        "legacy.com",
		Host:           "www.mydomain.com",
		Pathname:       "product/1",
		HttpStatusCode: 200,
		Body:           "",
	})

	assert.Equal(t, "\n\tproduct\n\tError with Backend\n\t\tExpected: newstack.com\n\t\tReceived: legacy.com", fmt.Sprintf("%v", err), "should not valid the expectations")

	err = r.ShouldBe(Expectations{
		Backend:        "newstack.com",
		Host:           "www.mydomain.com",
		Pathname:       "product/1",
		HttpStatusCode: 301,
		Body:           "",
	})

	assert.Equal(t, "\n\tproduct\n\tError with Status code\n\t\tExpected: 200\n\t\tReceived: 301", fmt.Sprintf("%v", err), "should not valid the expectations")
}
