package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"
)

var (
	SMTP_HOST     = "localhost"
	SMTP_PORT     string
	HTTP_PORT     string
	SMTP_USERNAME = ""
	SMTP_PASSWORD = ""
)

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// we will use a dockerfile to build the image for testing
	resource, err := pool.BuildAndRun("mailhog-test-server", "./Dockerfile", []string{})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		SMTP_PORT = resource.GetPort("1025/tcp")
		HTTP_PORT = resource.GetPort("8025/tcp")
		// ping to ensure that the server is up and running
		_, err := net.Dial("tcp", net.JoinHostPort("localhost", SMTP_PORT))
		return err
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestMail_Send(t *testing.T) {
	mail := Mail{
		host:     SMTP_HOST,
		port:     SMTP_PORT,
		from:     "from@example.com",
		password: "",
		username: "",
	}
	spew.Dump(mail)
	err := mail.Send([]string{"to@example.com"}, "Test Subject", "Sending an automated test email")
	assert.Nil(t, err)

	resp, err := http.Get(fmt.Sprintf("http://localhost:%s/api/v2/messages", HTTP_PORT))

	assert.Nil(t, err)

	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)

	assert.Nil(t, err)
	spew.Dump(string(body))

	// check whether the email was saved successfully

}
