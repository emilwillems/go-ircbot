package main

import (
    "os"
    "log"

    "./config"
    "./bot"
)

func main() {
    logger := log.New(os.Stdout, "[main    ] ", log.LstdFlags)
    logger.Println("loading configuration")
    config := config.LoadFromFile("bot.conf")

    logger.Println("spawning new bot")
    bot, err := bot.New(config)
    if err != nil {
        logger.Printf("FATAL ERROR - unable to init IRC bot: %s\n", err)
        return
    }

    err = bot.Connect()
    if err != nil {
        logger.Printf("FATAL ERROR - unable to connect to IRC: %s\n", err)
        return
    }

    quit := make(chan bool)

    go func() {
        logger.Println("waiting for bot to quit")
        select {
            case <- bot.Quitted:
                logger.Println("bot decided to quit")
                quit <- true
        }
    }()

    <- quit
    logger.Println("exiting gracefully")
    os.Exit(0)
}
