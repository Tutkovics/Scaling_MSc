package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/octago/sflags/gen/gflag"
)

// example: main.go -name Frontend -delay 9 -port 9090 -cpu 90 -memory 900 -endpoint-url /read -endpoint-cpu 99 -endpoint-delay 98 -endpoint-url /index -endpoint-cpu 22 -endpoint-delay 202
type config struct {
	Name           string   `flag:"name" desc:"Server/service name"`
	InitDelay      uint     `flag:"delay" desc:"Delay after start up [ms]"`
	Port           uint     `flag:"port" desc:"Open port to listen"`
	CPUusage       uint     `flag:"cpu" desc:"CPU usage in idle time [mCPU]"`
	MemoryUsage    uint     `flag:"memory" desc:"Memory usage in idle time [kB]"`
	Endpoints      []string `flag:"endpoint-url" desc:"Endpoints to listen"`
	EndpointsCPU   []uint   `flag:"endpoint-cpu" desc:"CPU usage for the endpoints"`
	EndpointsDelay []uint   `flag:"endpoint-delay" desc:"Delay for each endpoint [ms]"`
	EndpointsCall  []string `flag:"endpoint-call" desc:"If the endpoint need to call other service"`
}

func readConfigParameters() *config {
	// Set default parameters
	c := &config{
		Name:        "Service-#ID",
		InitDelay:   0,
		Port:        8080,
		CPUusage:    50,
		MemoryUsage: 64,
		Endpoints: []string{
			"/",
			"/health",
		},
		EndpointsCPU: []uint{
			200,
			10,
		},
		EndpointsDelay: []uint{
			30,
			0,
		},
	}

	err := gflag.ParseToDef(c)
	if err != nil {
		log.Fatalf("[READ_PARAMS]\terr: %v", err)
	}
	flag.Parse()

	// Check given paramters
	fmt.Printf("[READ_PARAMS]\tParameters OK: %t\n", c.check())

	return c
}

func (c *config) check() bool {
	if len(c.Endpoints) == len(c.EndpointsCPU) &&
		len(c.Endpoints) == len(c.EndpointsDelay) &&
		len(c.Endpoints) == len(c.EndpointsCall) {
		return true
	}

	return false
}

func main() {
	var cfg = readConfigParameters()

	fmt.Printf("[MAIN]\t\tConfig values:\t%+v\n", cfg)

	var addr string = ":" + fmt.Sprint(cfg.Port)

	// start webserver
	for i, endpoint := range cfg.Endpoints {
		fmt.Printf("%d --> %s --> \n", i, endpoint)
		http.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("[REQUEST-INCOME] %s --> %s\n", r.URL, r.URL.Path)
			start := time.Now()
			//foundEndpoint := false

			for k, endp := range cfg.Endpoints {
				if endp == r.URL.Path {
					// Sleep not relevant here
					// time.Sleep(time.Duration(cfg.EndpointsDelay[k]) * time.Millisecond)
					fmt.Printf("[REQUEST]\t%s\n", endp)
					fmt.Fprintf(w, "<h1>Hello from: %s!</h1><hl>\n", cfg.Name)
					fmt.Fprintln(w, "<h3>Config</h3><ul>")
					fmt.Fprintf(w, "<li>Config values: %t</li>\n", cfg.check())
					fmt.Fprintf(w, "<li>Endpoint: %s</li>\n", endp)
					fmt.Fprintf(w, "<li>CPU usage: %d</li>\n", cfg.EndpointsCPU[k])
					fmt.Fprintf(w, "<li>Delay time: %d</li>\n", cfg.EndpointsDelay[k])
					fmt.Fprintf(w, "<li>Call out: %s</li>\n", cfg.EndpointsCall[k])
					//fmt.Fprintf(w, "\t(%s)\n", strings.Split(cfg.EndpointsCall[k], ";"))

					// Check for empty/only whitespace string
					call := strings.ReplaceAll(cfg.EndpointsCall[k], "'", "")
					length := len(strings.TrimSpace(call))

					if length != 0 {
						fmt.Printf("Callout: '%s' - '%s' --> len: %d", endpoint, cfg.EndpointsCall[k], length)
						fmt.Fprintf(w, "<ul>")
						for i, callOut := range strings.Split(cfg.EndpointsCall[k], "__") {
							callOut = strings.ReplaceAll(callOut, "'", "")
							fmt.Printf("[CALL_OUT]\t#no%d --> %s\n", i, callOut)
							url := "http://" + callOut
							resp, err := http.Get(url)

							if err != nil {
								fmt.Fprintf(w, "<li>%d: <b>%s</b>: Oops, something went wrong</li>\n", i, callOut)
							} else {
								fmt.Fprintf(w, "<li>%d: <b>%s</b>: %s</li>\n", i, callOut, resp.Status)
							}

						}
						fmt.Fprintf(w, "</ul>")
					}
					fmt.Fprintln(w, "</ul>")

					// Generate CPU usage
					// Create waitgroup to wait all calculations done
					var waitgroup sync.WaitGroup
					waitgroup.Add(int(cfg.EndpointsCPU[k]))
					for i := 0; i < int(cfg.EndpointsCPU[k]); i++ {
						go algo(600, &waitgroup)
					}
					waitgroup.Wait()

					// After CPU calcualation wait if the delay time not passed
					waitTime := (time.Duration(cfg.EndpointsDelay[k]) * time.Millisecond) - time.Now().Sub(start)

					if waitTime > 0 {
						time.Sleep(waitTime)
					} else {
						fmt.Fprintf(w, "<p>CPU calculation took more time than delay time (%s)</p>\n", waitTime)
					}

					// Give more information about request/response
					fmt.Fprintf(w, "<h3>Info</h3>\n<ul>\n")
					fmt.Fprintf(w, "<li>Time: %s</li>\n", time.Now())
					fmt.Fprintf(w, "<li>Method: %s</li>\n", r.Method)
					fmt.Fprintf(w, "<li>URL: %s</li>\n", r.URL)
					fmt.Fprintf(w, "<li>RemoteAddr: %s</li>\n", r.RemoteAddr)
					fmt.Fprintf(w, "<li>Host: %s</li>\n", r.Host)
					fmt.Fprintln(w, "</ul>")
				}

			}

			// send response time
			fmt.Fprintf(w, "\nResponse time: %s\n", time.Now().Sub(start))
		})
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// log
		fmt.Printf("[REQUEST-INCOME] '%s' --> '%s'\n", r.URL, r.URL.Path)
		fmt.Printf("[MAIN]\t\tConfig values:\t%+v\n", cfg)
		fmt.Printf("%+v", r)
		// response
		fmt.Fprintf(w, "<h1>'/' or 404 page</h1>\n")
		fmt.Fprintf(w, "%+v", r)
	})

	if err := http.ListenAndServe(addr, nil); err != nil {

		log.Fatal(err)
	}

}

// Sieve of Eratosthenes
func algo(number int, waitgroup *sync.WaitGroup) {
	max := number
	numbers := make([]bool, max+1)
	// Set values to ture
	for i := range numbers {
		numbers[i] = true
	}

	// main algorithm
	for p := 2; p*p <= max; p++ {
		if numbers[p] {
			for i := p * p; i <= max; i += p {

				numbers[i] = false
			}
		}
	}

	// Print prime numbers
	for p := 2; p <= max; p++ {
		if numbers[p] {
			fmt.Printf("%d ", p)
		}
	}

	waitgroup.Done()
}
