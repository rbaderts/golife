package main

import (
	_ "encoding/json"
	"flag"
	"fmt"
	glog "github.com/ccding/go-logging/logging"
	"github.com/kataras/iris"
	"github.com/kataras/iris/config"
	_ "github.com/kataras/iris/websocket"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

var addr = flag.String("addr", ":8080", "http service address")
var Logger *glog.Logger

var Games map[string]*Game

type GameAPI struct {
	*iris.Context
}

func initialize() {

	var err error
	Logger, err =
		glog.FileLogger("logfile",
			glog.DEBUG,
			glog.BasicFormat,
			glog.DefaultTimeFormat,
			"logfile", true)

	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}
}

// GET /api/game/:id
/*
func (this GameAPI) GetBy(id string) {
	Logger.Debugf("GetBy - id = %v\n", id)
	jsonbytes, err := json.Marshal(Games[id])
	if err != nil {
		panic(err)
	}

	Logger.Debugf("game = %v\n", Games[id])
	Logger.Debugf("game json = %v\n", string(jsonbytes))
	this.JSON(http.StatusOK, Games[id])
}
*/

// POST /api/game
/*
func (this GameAPI) Post() {

	g := NewGame(50, 50)
	g.Seed()
	Games[g.GameId] = g

	Logger.Debugf("Post - game created: %v\n", g)

	loc := this.URI().String() + "/" + g.GameId
	this.SetHeader("Location", loc)
	this.SetStatusCode(http.StatusCreated)
}
*/

func main() {

	initialize()
	Games = make(map[string]*Game)

	config := config.Iris{
		Profile:     true,
		ProfilePath: "",
	}
	_ = config

	iris.Static("/css", "./resources/css", 1)
	iris.Static("/js", "./resources/js", 1)
	iris.Static("/img", "./resources/img", 1)

	iris.Get("/", func(ctx *iris.Context) {
		if err := ctx.Render("index.html", nil); err != nil {
			println(err.Error())
		}
	})

	iris.Get("/api/game/:id", func(ctx *iris.Context) {
		Logger.Debugf("GET /api/game/:id\n")
		id := ctx.Param("id")
		g := Games[id]
		ctx.JSON(http.StatusOK, g)
	})

	iris.Post("/api/game", func(ctx *iris.Context) {

		Logger.Debugf("URLParams = %v\n", ctx.URLParams())
		Logger.Debugf("POST /api/game\n")

		pattern := ctx.URLParam("pattern")
		Logger.Debugf("pattern param = %v\n", pattern)
		patternType := PatternTypeFromString(pattern)

		Logger.Debugf("patterType = %v\n", patternType)
		g := NewGame(50, 50)
		if patternType == 0 {
			g.Seed()
		} else {
			g.SeedPattern(patternType)
		}

		Games[g.GameId] = g

		Logger.Debugf("Post - game created: %v\n", g)

		loc := "/api/game/" + g.GameId
		ctx.SetHeader("Location", loc)
		//ctx.SetStatusCode(http.StatusCreated)
		ctx.JSON(http.StatusCreated, Games[g.GameId])

	})

	iris.Post("/api/game/:id/step", func(ctx *iris.Context) {
		Logger.Debugf("GET /api/game/:id/step\n")
		// Retrieve the parameters fullname and friendID
		id := ctx.Param("id")
		g := Games[id]
		if g == nil {
			ctx.NotFound()
		}

		g.Evolve()
		ctx.JSON(http.StatusOK, g)
	})

	iris.Post("/api/game/:id/field", func(ctx *iris.Context) {

		Logger.Debugf("POST  /api/game/:id/field\n")
		id := ctx.Param("id")
		pattern := ctx.PostFormValue("pattern")
		x := ctx.PostFormValue("xpos")
		y := ctx.PostFormValue("ypos")

		patternType := PatternTypeFromString(pattern)

		xpos, err := strconv.Atoi(x)
		if err != nil {
			xpos = 10
		}
		ypos, err := strconv.Atoi(y)
		if err != nil {
			ypos = 10
		}

		g := Games[id]
		if g != nil {
			g.AddPattern(patternType, xpos, ypos)
		}
		ctx.JSON(http.StatusOK, g)
	})

	/*
		jsonbytes, err := json.Marshal(g)
		if err != nil {
			panic(err)
		}
	*/
	//		ctx.JSON(http.StatusOK, g)

	/*
		iris.Post("/game", func(ctx *iris.Context) {
			game := NewGame(20, 20)
			game.Seed()
			ctx.Text(iris.StatusOK, game.gameId)
		})

	*/
	iris.Config().Render.Template.Layout = "layouts/layout.html"

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		os.Exit(1)
	}()

	iris.API("/api/game", GameAPI{})

	iris.Listen(":8080")
}
