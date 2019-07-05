package main

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"time"
)

var (
	optBind string

	optSiteFQDN  string
	optSiteName  string
	optGithubOrg string
)

var (
	Packages map[string]*Package
)

func init() {
	if optBind = os.Getenv("BIND"); len(optBind) == 0 {
		optBind = ":5000"
	}
	if optSiteFQDN = os.Getenv("SITE_FQDN"); len(optSiteFQDN) == 0 {
		optSiteFQDN = "go.guoyk.net"
	}
	if optSiteName = os.Getenv("SITE_NAME"); len(optSiteName) == 0 {
		optSiteName = "Go - Guo Y.K."
	}
	if optGithubOrg = os.Getenv("GITHUB_ORG"); len(optGithubOrg) == 0 {
		optGithubOrg = "go-guoyk"
	}
}

func exit(err *error) {
	if *err != nil {
		os.Exit(1)
	}
}

func routineUpdate() {
	for {
		update()
		time.Sleep(time.Minute)
	}
}

func update() {
	var err error
	var resp *http.Response
	if resp, err = http.Get(fmt.Sprintf("https://api.github.com/orgs/%s/repos?per_page=100", optGithubOrg)); err != nil {
		log.Println("failed to get api: " + err.Error())
		return
	}
	defer resp.Body.Close()
	var body []byte
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		log.Println("failed to read api: " + err.Error())
		return
	}
	// unmarshal json
	var pkgs []*Package
	if err = json.Unmarshal(body, &pkgs); err != nil {
		log.Println("failed to decode api: " + err.Error())
		return
	}
	// build map
	m := make(map[string]*Package)
	for _, pkg := range pkgs {
		m[pkg.Name] = pkg
	}
	Packages = m
	return
}

func routeIndex(ctx echo.Context) error {
	var keys []string
	for k := range Packages {
		keys = append(keys, k)
	}
	sort.Sort(sort.StringSlice(keys))
	var values []*Package
	for _, k := range keys {
		values = append(values, Packages[k])
	}
	return ctx.Render(http.StatusOK, "index", map[string]interface{}{
		"SiteFQDN": optSiteFQDN,
		"SiteName": optSiteName,
		"Packages": values,
	})
}

func routePackage(ctx echo.Context) error {
	pkg := Packages[ctx.Param("name")]
	if pkg == nil {
		return ctx.String(http.StatusNotFound, "not found")
	}
	return ctx.Render(http.StatusOK, "package", map[string]interface{}{
		"SiteFQDN":    optSiteFQDN,
		"SiteName":    optSiteName,
		"Name":        pkg.Name,
		"HTMLURL":     pkg.HTMLURL,
		"CloneURL":    pkg.CloneURL,
		"Description": pkg.Description,
	})
}

func main() {
	var err error
	defer exit(&err)

	go routineUpdate()

	e := echo.New()

	e.HideBanner = true
	e.HidePort = true
	e.Renderer = DefaultTemplate

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	e.GET("/", routeIndex)
	e.GET("/:name", routePackage)
	e.GET("/:name/*", routePackage)

	err = e.Start(optBind)
}
