package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func parcelServe(prefix string, r *gin.Engine) {
	if gin.Mode() == "release" {

		serveAsset := parcelAssetHandler(prefix + "/dist")
		for _, asset := range AssetNames() {
			r.GET(strings.TrimSuffix(strings.TrimPrefix(asset, prefix+"/dist"), "index.html"), serveAsset)
		}

		// Fall back to last index.html if not found
		r.Use(func(c *gin.Context) {
			path := strings.Trim(c.Request.URL.Path, "/")
			segments := strings.SplitAfter(path, "/")
			if path == "" {
				c.Next()
				return
			}
			fmt.Printf("Old path: %s\n", "/"+strings.Join(segments[:], ""))
			c.Request.URL.Path = "/" + strings.Join(segments[:len(segments)-1], "")
			fmt.Printf("New path: %s\n", c.Request.URL.Path)
			r.HandleContext(c)
		})

	} else {

		parcel := exec.Command("parcel", "index.html")
		parcel.Stdout = os.Stdout
		parcel.Stderr = os.Stderr
		parcel.Dir, _ = filepath.Abs(prefix + "/")
		parcel.Start()
		r.Use(parcelProxy("127.0.0.1:1234"))

	}
}

func parcelAssetHandler(prefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := c.Request.URL.Path
		if strings.HasSuffix(p, "/") {
			p = p + "index.html"
		}
		a := MustAsset(prefix + p)
		var contentType = mime.TypeByExtension(filepath.Ext(p))

		c.DataFromReader(http.StatusOK, int64(len(a)), contentType, bytes.NewReader(a), map[string]string{})
	}
}

func parcelProxy(host string) gin.HandlerFunc {
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
