package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lemonade-command/lemonade/pkg/utils"
)

type MessageData struct {
	Message string `json:"message"`
}

func (m *MessageData) toByte() []byte {
	data, _ := json.Marshal(m)
	return data
}

var currentMessage MessageData
var mlock sync.Mutex

var bindIP = "127.0.0.1:1789"
var clipboardDISPLAY = ""

func main() {
	if len(os.Args) < 2 {
		return
	}

	Command := os.Args[1]
	switch Command {
	case "copy":
		var m MessageData

		if len(os.Args) < 3 {

			bData, err := io.ReadAll(os.Stdin)
			if err != nil {
				panic(err)
			}
			m.Message = string(bData)
		} else {
			m.Message = os.Args[2]
		}

		if _, exists := os.LookupEnv("TERMUX_VERSION"); exists {
			cmd := exec.Command("termux-clipboard-set", m.Message)
			cmd.Stdout = os.Stdout
			cmd.Stdin = os.Stdin
			cmd.Stderr = os.Stderr
			cmd.WaitDelay = 2 * time.Second
			cmd.Run()
			return
		}
		err := utils.Post("http://"+bindIP+"/copy", m.toByte())
		if err != nil {
			panic(err)
		}
	case "paste":
		if _, exists := os.LookupEnv("TERMUX_VERSION"); exists {
			cmd := exec.Command("termux-clipboard-get")
			cmd.Stdout = os.Stdout
			cmd.Stdin = os.Stdin
			cmd.Stderr = os.Stderr
			cmd.WaitDelay = 2 * time.Second
			cmd.Run()
			return
		}
		// else if _, exists := os.LookupEnv("DISPLAY"); exists {
		// 	err := clipboard.Init()
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// 	mess := clipboard.Read(clipboard.FmtText)
		// 	fmt.Printf("%s", string(mess))
		// 	return
		// }
		result, err := utils.Get("http://" + bindIP + "/paste")
		if err != nil {
			panic(err)
		}
		var m MessageData
		err = json.Unmarshal(result, &m)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s", m.Message)
	case "serve":
		handleServe()
	default:
	}
}
func handleServe() {
	r := gin.Default()
	r.GET("/paste", func(c *gin.Context) {
		mlock.Lock()
		defer mlock.Unlock()
		c.JSON(http.StatusOK, currentMessage)
	})
	r.POST("/copy", func(ctx *gin.Context) {
		mlock.Lock()
		defer mlock.Unlock()
		err := ctx.BindJSON(&currentMessage)
		if err != nil {
			log.Println(err)
			return
		}
		cmd := exec.Command("termux-clipboard-set", currentMessage.Message)
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		cmd.WaitDelay = 2 * time.Second
		cmd.Run()

	})
	r.Run(bindIP) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
