package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/armakuni/circleci-workflow-dashboard/circleci"
	"github.com/armakuni/circleci-workflow-dashboard/dashboard"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
)

type Dashboard struct {
	RefreshInterval   int
	Now               string
	DashboardMonitors dashboard.Monitors
}

func updateDashboard(circleCIClient *circleci.Client, filter *circleci.Filter, ticker *time.Ticker, c *cache.Cache, animateBuildError bool) {
	for ; true; <-ticker.C {
		dashboardMonitors, err := dashboard.Build(circleCIClient, filter, animateBuildError)
		c.Set("dashErr", err, cache.NoExpiration)
		c.Set("dashboardMonitors", dashboardMonitors, cache.NoExpiration)
		c.Set("now", time.Now().Format("2006-01-02 15:04:05 -0700"), cache.NoExpiration)
	}
}

func getCachedDashboard(c *cache.Cache, refreshInterval int) (Dashboard, error) {
	if err, found := c.Get("dashErr"); err != nil && found {
		return Dashboard{}, err.(error)
	}
	now, found := c.Get("now")
	if !found {
		return Dashboard{}, fmt.Errorf("Could not find cached dashboard data")
	}
	dashboardMonitors, found := c.Get("dashboardMonitors")
	if !found {
		return Dashboard{}, fmt.Errorf("Could not find cached dashboard data")
	}
	return Dashboard{
		DashboardMonitors: dashboardMonitors.(dashboard.Monitors),
		Now:               now.(string),
		RefreshInterval:   refreshInterval,
	}, nil
}

func setup(refreshInterval int) (*time.Ticker, *cache.Cache) {
	ticker := time.NewTicker(time.Duration(refreshInterval) * time.Second)
	cacher := cache.New(5*time.Minute, 5*time.Minute)
	return ticker, cacher
}

func getConfig() (*circleci.Config, *circleci.Filter, error) {
	apiToken := os.Getenv("CIRCLECI_TOKEN")
	apiURL := os.Getenv("CIRCLECI_API_URL")
	jobsURL := os.Getenv("CIRCLECI_JOBS_URL")
	config := &circleci.Config{
		APIToken: apiToken,
		APIURL:   apiURL,
		JobsURL:  jobsURL,
	}
	filterJson := os.Getenv("DASHBOARD_FILTER")
	if filterJson == "" {
		filterJson = "null"
	}
	var filter circleci.Filter
	if err := json.Unmarshal([]byte(filterJson), &filter); err != nil {
		return nil, nil, fmt.Errorf("Error loading dashboard filter: %v", err.Error())
	}
	return config, &filter, nil
}

func getRefershInterval() int {
	refreshInterval := os.Getenv("REFRESH_INTERVAL")
	if refreshInterval == "" {
		return 30
	}
	refreshInt, err := strconv.Atoi(refreshInterval)
	if err != nil {
		fmt.Println("REFRESH_INTERVAL must be an int")
		os.Exit(1)
	}
	return refreshInt
}

func main() {
	var animateBuildError = true
	refreshInterval := getRefershInterval()
	if os.Getenv("ANIMATED_BUILD_ERROR") == "false" {
		animateBuildError = false
	}
	config, filter, err := getConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	circleCIClient, err := circleci.NewClient(config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ticker, cacher := setup(refreshInterval)
	defer ticker.Stop()
	go updateDashboard(circleCIClient, filter, ticker, cacher, animateBuildError)
	r := gin.Default()
	r.LoadHTMLGlob("templates/*.tmpl")
	r.GET("/", func(c *gin.Context) {
		dashboard, err := getCachedDashboard(cacher, refreshInterval)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		c.HTML(200, "dashboard.tmpl", dashboard)
	})
	r.Static("/assets", "./assets")
	r.Run() // listen and serve on 0.0.0.0:8080
}
