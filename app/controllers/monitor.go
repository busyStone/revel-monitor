// Monitor revel app status

// Base on beego
package controller

import (
  "bytes"
  "fmt"
  "github.com/astaxie/beego/toolbox"
  "github.com/revel/revel"
  "net/http"
  "strings"
  "time"
)

var TimeLayout = "2006-01-02 15:04:05.999999999 -0700 MST"

type Monitor struct {
  *revel.Controller
}

type QPSStatistics struct {
  Url      string
  Method   string
  Times    string
  UsedTime string
  MaxTime  string
  MinTime  string
  AvgTime  string
}

// QpsIndex displays qbs
// url /@qbs
func (c *Monitor) QpsIndex() revel.Result {

  buf := bytes.NewBuffer([]byte{})

  toolbox.StatisticsMap.GetMap(buf)

  qpsStr := buf.String()

  var qps []QPSStatistics

  qpsStrs := strings.Split(qpsStr, "\n")
  if len(qpsStrs) > 1 {
    qpsStrs = qpsStrs[1:]
    for _, s := range qpsStrs {
      s = strings.TrimSpace(s)
      s = strings.TrimPrefix(s, "|")
      s = strings.TrimSuffix(s, "|")
      s = strings.TrimSpace(s)
      if len(s) == 0 {
        continue
      }

      staticStrs := strings.Split(s, "|")
      if len(staticStrs) != 7 {
        revel.INFO.Println("QpsIndex Error: QPSStatistics field error.", len(staticStrs))
        continue
      }

      qps = append(qps, QPSStatistics{
        Url:      strings.TrimSpace(staticStrs[0]),
        Method:   strings.TrimSpace(staticStrs[1]),
        Times:    strings.TrimSpace(staticStrs[2]),
        UsedTime: strings.TrimSpace(staticStrs[3]),
        MaxTime:  strings.TrimSpace(staticStrs[4]),
        MinTime:  strings.TrimSpace(staticStrs[5]),
        AvgTime:  strings.TrimSpace(staticStrs[6]),
      })
    }
  }

  return c.Render(qps)
}

// ProfIndex show profile command or with command to show infos
// url /@prof
func (c *Monitor) ProfIndex(command string) revel.Result {
  var profs []string

  if command != "" {
    buf := bytes.NewBuffer([]byte{})

    toolbox.ProcessInput(command, buf)

    prof := buf.String()
    profs = strings.Split(prof, "\n")
  }

  return c.Render(command, profs)
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

func qpsBegin(c *revel.Controller) revel.Result {
  c.Flash.Data["qpsStartTime"] = time.Now().Format(TimeLayout)

  return nil
}

func qpsEnd(c *revel.Controller) revel.Result {

  startTimeStr, ok := c.Flash.Data["qpsStartTime"]
  if !ok {
    revel.ERROR.Println("Can't find qpsStartTime in flash.")

    return nil
  }

  startTime, err := time.Parse(TimeLayout, startTimeStr)
  if err != nil {
    revel.ERROR.Println(err)

    return nil
  }

  toolbox.StatisticsMap.AddStatistics(c.Request.Method, c.Request.URL.Path, c.Name, time.Since(startTime))

  return nil
}

func init() {
  revel.InterceptFunc(qpsBegin, revel.BEFORE, revel.ALL_CONTROLLERS)
  revel.InterceptFunc(qpsEnd, revel.AFTER, revel.ALL_CONTROLLERS)
}
