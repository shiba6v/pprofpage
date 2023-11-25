# PProfPage

PProfPage は、fgprof や pprof の結果に対して URL を発行し、ブラウザから見られるようにするツールです。

まず、以下のようなコードをアプリケーションに追加します。
以下のコードでは /flush_pprof を叩いたときに /home/isucon/cpu.pprof, /home/isucon/cpu.fgprof が吐かれるようにしていますが、この処理を最後に叩くようにすれば Web 以外にも使えます。

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/pprof"

	"github.com/felixge/fgprof"
	"github.com/labstack/echo/v4"
)

func main(){
	e := echo.New()
    // 計測したいAPI定義など
    setupAllMetrics(e, true)
    e.Start(addr)
}

func setupAllMetrics(e *echo.Echo, enable bool) {
	if enable {
		flushPProf := setupPProf()
		flushFGProf := setupFGProf()
		e.GET("/flush_pprof", func(c echo.Context) error {
			flushPProf()
			flushFGProf()
			return c.String(http.StatusOK, "ok")
		})
	}
}

func setupPProf() func() {
	f, err := os.Create("/home/isucon/tmp/cpu.pprof")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	return func() {
		if f == nil {
			return
		}
		pprof.StopCPUProfile()
		f.Close()
		f = nil
	}
}

func setupFGProf() func() {
	f, err := os.Create("/home/isucon/tmp/cpu.fgprof")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	cleanup := fgprof.Start(f, "pprof")
	if err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	return func() {
		if f == nil {
			return
		}
		if err := cleanup(); err != nil {
			fmt.Printf("SetupFGProf cleanup failed, %v", err)
		}
		f.Close()
		f = nil
	}
}
```

以下は、使用時のサンプルコードです。
以下で使用しているエンドポイントは、試用のためにご自由に使ってくださって大丈夫です。
(ただし、動作の保証はいたしません。)

```bash
    curl -k http://127.0.0.1:8080/flush_pprof
    PPROF_PATH=$(curl -X POST -F file=@/home/isucon/tmp/cpu.pprof https://emaaxnj6hk.execute-api.ap-northeast-1.amazonaws.com/prod/pprof/register)
    echo https://emaaxnj6hk.execute-api.ap-northeast-1.amazonaws.com/prod${PPROF_PATH} | post_slack
    PPROF_PATH=$(curl -X POST -F file=@/home/isucon/tmp/cpu.fgprof https://emaaxnj6hk.execute-api.ap-northeast-1.amazonaws.com/prod/pprof/register)
    echo https://emaaxnj6hk.execute-api.ap-northeast-1.amazonaws.com/prod${PPROF_PATH} | post_slack
```

FlagSet 周りのコードは kaz/pprotein からお借りしています。
