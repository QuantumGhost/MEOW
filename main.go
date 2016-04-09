package main

import (
	"fmt"
	"os"
	"runtime"
	"sync"
)

func main() {
	// Parse flags after load config to allow override options in config
	cmdLineConfig := parseCmdLineConfig()
	if cmdLineConfig.PrintVer {
		printVersion()
		os.Exit(0)
	}

	fmt.Printf(`
       /\
   )  ( ')     MEOW Proxy %s
  (  /  )      http://renzhn.github.io/MEOW/
   \(__)|      
	`, version)
	fmt.Println()

	parseConfig(cmdLineConfig.RcFile, cmdLineConfig)
	if config.AuthDB != "" {
		fmt.Printf("Using database %s as auth info storage.", config.AuthDB)
		auth.storage = NewSQLiteStorage(config.AuthDB)
	} else {
		auth.storage = NewMemoryStorage()
		fmt.Println("Reading auth info from config file...")
		addUserPasswd(config.UserPasswd)
		loadUserPasswdFile(config.UserPasswdFile)
	}

	initSelfListenAddr()
	initLog()
	initAuth()
	initStat()

	initParentPool()

	if config.JudgeByIP {
		initCNIPData()
	}

	if config.Core > 0 {
		runtime.GOMAXPROCS(config.Core)
	}

	go runSSH()

	var wg sync.WaitGroup
	wg.Add(len(listenProxy))
	for _, proxy := range listenProxy {
		go proxy.Serve(&wg)
	}
	wg.Wait()
}
