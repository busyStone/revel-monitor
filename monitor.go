// Monitor
// base on beego
// config revmonitor.httpaddr revmonitor.httpport
package monitor

import (
  "fmt"
  "github.com/astaxie/beego/toolbox"
  "github.com/revel/revel"
  "net/http"
  "time"
)

// QpsIndex is the http.Handler for writing qbs statistics map result info in http.ResponseWriter.
// it's registered with url pattern "/qbs" in admin module.
func qpsIndex(rw http.ResponseWriter, r *http.Request) {
  toolbox.StatisticsMap.GetMap(rw)
}

// ProfIndex is a http.Handler for showing profile command.
// it's in url pattern "/prof" in admin module.
func profIndex(rw http.ResponseWriter, r *http.Request) {
  r.ParseForm()
  command := r.Form.Get("command")
  if command != "" {
    toolbox.ProcessInput(command, rw)
  } else {
    rw.Write([]byte("<html><head><title>beego admin dashboard</title></head><body>"))
    rw.Write([]byte("request url like '/prof?command=lookup goroutine'<br>\n"))
    rw.Write([]byte("the command have below types:<br>\n"))
    rw.Write([]byte("1. <a href='?command=lookup goroutine'>lookup goroutine</a><br>\n"))
    rw.Write([]byte("2. <a href='?command=lookup heap'>lookup heap</a><br>\n"))
    rw.Write([]byte("3. <a href='?command=lookup threadcreate'>lookup threadcreate</a><br>\n"))
    rw.Write([]byte("4. <a href='?command=lookup block'>lookup block</a><br>\n"))
    rw.Write([]byte("5. <a href='?command=start cpuprof'>start cpuprof</a><br>\n"))
    rw.Write([]byte("6. <a href='?command=stop cpuprof'>stop cpuprof</a><br>\n"))
    rw.Write([]byte("7. <a href='?command=get memprof'>get memprof</a><br>\n"))
    rw.Write([]byte("8. <a href='?command=gc summary'>gc summary</a><br>\n"))
    rw.Write([]byte("</body></html>"))
  }
}

// Healthcheck is a http.Handler calling health checking and showing the result.
// it's in "/healthcheck" pattern in admin module.
func healthcheck(rw http.ResponseWriter, req *http.Request) {
  for name, h := range toolbox.AdminCheckList {
    if err := h.Check(); err != nil {
      fmt.Fprintf(rw, "%s : %s\n", name, err.Error())
    } else {
      fmt.Fprintf(rw, "%s : ok\n", name)
    }
  }
}

type Monitor struct {
  *revel.Controller
  StartTime time.Time
}

func (m *Monitor) QpsBegin() revel.Result {
  m.StartTime = time.Now()

  return nil
}

func (m *Monitor) QpsEnd() revel.Result {

  toolbox.StatisticsMap.AddStatistics(m.Request.Method, m.Request.URL.Path, m.Name, time.Since(m.StartTime))

  return nil
}

func init() {
  revel.OnAppStart(func() {
    // default httpAddr = localhost:8088
    httpAddr, found := revel.Config.String("revmonitor.httpaddr")
    if !found {
      httpAddr = "localhost"
    }

    httpPort, found := revel.Config.Int("revmonitor.httpport")
    if !found {
      httpAddr += ":8088"
    } else if httpPort != 0 {
      httpAddr = fmt.Sprintf("%s:%d", httpAddr, httpPort)
    }

    http.HandleFunc("/healthcheck", healthcheck)
    http.HandleFunc("/prof", profIndex)
    http.HandleFunc("/qps", qpsIndex)

    err := http.ListenAndServe(httpAddr, nil)
    if err != nil {
      revel.ERROR.Println("Monitor server start filed!")
    } else {
      revel.INFO.Printf("Monitor lisend on %s", httpAddr)

      revel.InterceptMethod((*Monitor).QpsBegin, revel.BEFORE)
      revel.InterceptMethod((*Monitor).QpsEnd, revel.AFTER)
    }
  })
}
