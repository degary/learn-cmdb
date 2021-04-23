package filters

import (
	"github.com/astaxie/beego/context"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

//每个url请求时间 Histogram 带可变label
var (
	//总请求次数: counter
	totalRequest = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "cmdb_request_total",
		Help: "",
	})
	//每个URL请求次数 counter 带可变label
	urlRequest = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "cmdb_url_request_total",
		Help: "",
	}, []string{"url"})
	//状态码统计: counter 带可变label
	statusCode = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "cmdb_status_code_total",
		Help: "",
	}, []string{"code"})
	//总请求时间
	elapsedTime = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "cmdb_request_url_elapsed_time",
		Help: "",
	}, []string{"url"})
)

func BeforeExec(ctx *context.Context) {
	totalRequest.Inc()
	urlRequest.WithLabelValues(ctx.Input.URL()).Inc()
	ctx.Input.SetData("stime", time.Now())
}

func AfterExec(ctx *context.Context) {

	//fmt.Println(ctx.Request.URL, ctx.ResponseWriter.Status) //todo ctx.ResponseWrite.Status 返回值为0

	statusCode.WithLabelValues(strconv.Itoa(ctx.ResponseWriter.Status)).Inc()
	stime := ctx.Input.GetData("stime")
	if stime != nil {
		if t, ok := stime.(time.Time); ok {
			etime := time.Now().Sub(t)
			elapsedTime.WithLabelValues(ctx.Input.URL()).Observe(float64(etime))
		}
	}
}

func init() {
	prometheus.MustRegister(totalRequest, statusCode, elapsedTime, urlRequest)

}
