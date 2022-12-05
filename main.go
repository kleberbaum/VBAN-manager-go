package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func getConfig(cpath string) (argfiles map[string]string, nextId int) {
	nextId = 0
	argfiles = make(map[string]string)
	re := regexp.MustCompile("[0-9]+")
	f, err := os.Open(cpath)

	if err != nil {
		fmt.Println(err)
	}

	files, err := f.Readdir(0)

	if err != nil {
		fmt.Println(err)
	}

	for _, v := range files {
		fname := v.Name()

		if fname[:1] == "a" {
			id := re.FindString(fname)
			i, _ := strconv.Atoi(id)
			argfiles[id] = cpath + v.Name()
			if nextId <= i {
				nextId = i + 1
			}
		}
	}

	return argfiles, nextId
}

func main() {
	script := "./script/"
	script_name := "vban.sh"
	script_sh := script + script_name
	args := script + "args-"
	args_sub := 14
	plugins_folder := "./plugins/"
	plugins_sub := 10

	// Set the router as the default one provided by Gin
	r := gin.Default()

	// Allow all Origins, fix for cors error
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
	}))

	// Set static path
	r.Static("/assets", "./assets")

	// Process the templates at the start so that they don't have to be loaded
	// from the disk again. This makes serving HTML pages very fast.
	r.LoadHTMLGlob("templates/**/*.html")

	// Print active config
	log.Printf("Configuration")
	log.Printf("Script directory %s", script)
	log.Printf("Script name %s", script_sh)
	log.Printf("Interface config path %s", args)
	log.Printf("No idea yet %d", args_sub)
	log.Printf("Plugin directory %s", plugins_folder)
	log.Printf("F no idea %d", plugins_sub)

	// Define the route for the index page and display the index.html template
	// To start with, we'll use an inline route handler. Later on, we'll create
	// standalone functions that will be used as route handlers.
	r.GET("/", func(c *gin.Context) {
		message := c.Query("message")
		configFiles, nextId := getConfig(script)

		// Call the HTML method of the Context to render a template
		c.HTML(
			// Set the HTTP status to 200 (OK)
			http.StatusOK,
			// Use the index.html template
			"home/index.html",
			// Pass the data that the page uses (in this case, 'title')

			gin.H{
				"page":           "welcome",
				"title":          "Welcome Page",
				"message":        message,
				"script":         script,
				"script_sh":      script_sh,
				"args":           args,
				"args_sub":       args_sub,
				"plugins_folder": plugins_folder,
				"plugins_sub":    plugins_sub,
				"server":         configFiles,
				"nextId":         nextId,
			},
		)

	})

	// Define the route for the server page and display the server.html template
	// To start with, we'll use an inline route handler. Later on, we'll create
	// standalone functions that will be used as route handlers.
	r.GET("/server", func(c *gin.Context) {
		id := c.Query("id")
		message := c.Query("message")
		newi, _ := strconv.ParseBool(c.Query("new"))
		configFiles, nextId := getConfig(script)
		args := make(map[string]string)
		argscontent, err := make([]byte, 5), *new(error)

		if newi {
			args["type"] = "receptor"
		} else {
			argscontent, err = os.ReadFile(configFiles[id])

			if string(argscontent) != "" {
				for i, argsar := range strings.Split(string(argscontent), " -") {
					if i == 0 {
						args["type"] = argsar
					}

					args[argsar[:1]] = argsar[1:]
				}

				if err != nil {
					log.Fatal(err)
				}
			}
		}

		log.Printf("Interface config %s", argscontent)

		// Call the HTML method of the Context to render a template
		c.HTML(
			// Set the HTTP status to 200 (OK)
			http.StatusOK,

			// Use the server.html template
			"server/server.html",

			// Pass the data that the page uses (in this case, 'title')
			gin.H{
				"page":           "server",
				"title":          "Server Page",
				"id":             id,
				"nextId":         nextId,
				"message":        message,
				"script":         script,
				"script_sh":      script_sh,
				"args":           args,
				"args_sub":       args_sub,
				"plugins_folder": plugins_folder,
				"plugins_sub":    plugins_sub,
				"server":         configFiles,
			},
		)

	})

	// Define the route for the server page and display the server.html template
	// To start with, we'll use an inline route handler. Later on, we'll create
	// standalone functions that will be used as route handlers.
	r.POST("/modify_args", func(c *gin.Context) {
		c.Request.ParseForm()
		args := ""

		for key, value := range c.Request.PostForm {
			if len(key) == 1 && value[0] != "" {
				args += " -" + string(key) + string(value[0])
			}
		}

		cmd := exec.Command("./"+script_name, "args", c.PostForm("nb"), c.PostForm("type")+args)
		cmd.Dir = script
		stdout, err := cmd.Output()
		out := string(stdout)

		if err != nil {
			out = err.Error()
		}

		c.Redirect(http.StatusMovedPermanently, "/server?id="+string(c.PostForm("nb"))+"&message=Arguments+changed!"+out)
	})

	// Define the route for the server page and display the server.html template
	// To start with, we'll use an inline route handler. Later on, we'll create
	// standalone functions that will be used as route handlers.
	r.GET("/log", func(c *gin.Context) {
		id := c.Query("id")
		cmd := exec.Command(script_sh, "status", id)
		stdout, err := cmd.Output()
		out := string(stdout)

		if err != nil {
			out = err.Error()
		}

		// Call the HTML method of the Context to render a template
		c.HTML(
			// Set the HTTP status to 200 (OK)
			http.StatusOK,
			// Use the index.html template
			"server/log.html",
			// Pass the data that the page uses (in this case, 'title')

			gin.H{
				"log": out,
				"err": err,
			},
		)

	})

	// Define the route for the server page and display the server.html template
	// To start with, we'll use an inline route handler. Later on, we'll create
	// standalone functions that will be used as route handlers.
	r.GET("/status", func(c *gin.Context) {
		id := c.Query("id")
		cmd := exec.Command(script_sh, "is-active", id)
		stdout, err := cmd.Output()
		out := string(stdout)

		if err != nil {
			out = err.Error()
		}

		c.JSON(200, gin.H{
			"status": out,
		})

	})

	// Define the route for the server page and display the server.html template
	// To start with, we'll use an inline route handler. Later on, we'll create
	// standalone functions that will be used as route handlers.
	r.GET("/action", func(c *gin.Context) {
		id := c.Query("id")
		atype := c.Query("type")
		cmd := exec.Command("./"+script_name, atype, id)
		cmd.Dir = script
		stdout, err := cmd.Output()
		out := string(stdout)

		if err != nil {
			out = err.Error()
		}

		log.Printf("Authorized on account %s", cmd.String())
		log.Printf("Authorized on account %s", out)

		c.JSON(200, gin.H{
			"status": out,
		})

	})

	// Test if server is available
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
			"query":   c.Request.URL.Query(),
		})
	})

	// //show user template
	// r.GET("/users", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "users/users.tmpl", gin.H{
	// 		"title": "Users Page",
	// 		"query": c.Request.URL.Query(),
	// 	})
	// })

	// Start serving the application
	// Listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	r.Run()
}
