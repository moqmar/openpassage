package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func proxy(host string) gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
			return
		}
		c.Request.Body = ioutil.NopCloser(bytes.NewReader(body))
		url := fmt.Sprintf("%s://%s%s", "http", host, c.Request.RequestURI)

		proxyReq, err := http.NewRequest(c.Request.Method, url, bytes.NewReader(body))

		proxyReq.Header = make(http.Header)
		for h, val := range c.Request.Header {
			proxyReq.Header[h] = val
		}

		client := &http.Client{}
		resp, err := client.Do(proxyReq)
		if err != nil {
			http.Error(c.Writer, err.Error(), http.StatusBadGateway)
			return
		}

		c.Status(resp.StatusCode)
		for h, val := range resp.Header {
			c.Writer.Header()[h] = val
		}

		defer resp.Body.Close()
		bodyContent, _ := ioutil.ReadAll(resp.Body)
		c.Writer.Write(bodyContent)

		c.Abort()
	}
}
