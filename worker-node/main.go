package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
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
			"/index",
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
		EndpointsCall: []string{
			"",
			"asd__basf",
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

type Response struct {
	ServiceName      string    // Name of the service
	Host             string    // Name of the host who answer
	ConfigOK         bool      // Results of parameter check
	CalledEnpoint    string    // Name of the called endpoint
	CPU              int32     // CPU usage for the endpoint
	Delay            int32     // Delay time of the endpoint
	CalloutParameter string    // Commandline parameter given in start
	Callouts         []string  //[]Response // Responses from callouts
	ActualDelay      int32     // Actual delay in response
	Time             time.Time // Current time in response
	RequestMethod    string    // Method of request
	RequestURL       *url.URL  // Full URL of request
	RequestAddress   string    // Remote address from request
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
					response := Response{}

					fmt.Printf("[REQUEST]\t%s\n", endp)
					fmt.Fprintf(w, "<h1>Hello from: %s!</h1><hl>\n", cfg.Name)
					response.ServiceName = cfg.Name
					fmt.Fprintln(w, "<h3>Config</h3><ul>")
					fmt.Fprintf(w, "<li>Config values: %t</li>\n", cfg.check())
					response.ConfigOK = cfg.check()
					fmt.Fprintf(w, "<li>Endpoint: %s</li>\n", endp)
					response.CalledEnpoint = endp
					fmt.Fprintf(w, "<li>CPU usage: %d</li>\n", cfg.EndpointsCPU[k])
					response.CPU = int32(cfg.EndpointsCPU[k])
					fmt.Fprintf(w, "<li>Delay time: %d</li>\n", cfg.EndpointsDelay[k])
					response.Delay = int32(cfg.EndpointsDelay[k])
					fmt.Fprintf(w, "<li>Call out: %s</li>\n", cfg.EndpointsCall[k])
					response.CalloutParameter = cfg.EndpointsCall[k]

					// Remove ' character if get from command line
					call := strings.ReplaceAll(cfg.EndpointsCall[k], "'", "")
					// Check for empty/only whitespace string
					length := len(strings.TrimSpace(call))

					if length != 0 {
						fmt.Printf("Callout: '%s' - '%s' --> len: %d", endpoint, cfg.EndpointsCall[k], length)
						fmt.Fprintf(w, "<ul>")

						// Split callout parameter by separate character: '__'
						calloutStringArray := strings.Split(cfg.EndpointsCall[k], "__")

						// Create array to collect callout responses
						calloutResponses := make([]string, len(calloutStringArray))

						for i, callOut := range calloutStringArray {
							callOut = strings.ReplaceAll(callOut, "'", "")
							fmt.Printf("[CALL_OUT]\t#no%d --> %s\n", i, callOut)
							url := "http://" + callOut
							resp, err := http.Get(url)

							if err != nil {
								calloutResponses = append(calloutResponses, "Oops, calling out failed")
							} else {
								// Convert response body to string
								buf := new(strings.Builder)
								_, err := io.Copy(buf, resp.Body)
								if err != nil {
									// Convertion failed
									calloutResponses = append(calloutResponses, "Oops, failed to convert response to string")
								} else {
									// Convertion was successfull
									calloutResponses = append(calloutResponses, string(buf.String()))
								}

								fmt.Fprintf(w, "<li>%d: <b>%s</b>: %s</li>\n", i, callOut, resp.Status)
							}

						}
						fmt.Fprintf(w, "</ul>")

						response.Callouts = calloutResponses
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
					response.Time = time.Now()
					fmt.Fprintf(w, "<li>Method: %s</li>\n", r.Method)
					response.RequestMethod = r.Method
					fmt.Fprintf(w, "<li>URL: %s</li>\n", r.URL)
					response.RequestURL = r.URL
					fmt.Fprintf(w, "<li>RemoteAddr: %s</li>\n", r.RemoteAddr)
					response.RequestAddress = r.RemoteAddr
					fmt.Fprintf(w, "<li>Host: %s</li>\n", r.Host)
					response.Host = r.Host
					fmt.Fprintln(w, "</ul>")

					responseJson, err := json.Marshal(response)

					if err != nil {
						fmt.Println(err)
					}
					fmt.Fprintf(w, "json: %s", string(responseJson))
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
