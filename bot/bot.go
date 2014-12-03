package bot

import (
    "os"
    "log"
    "net"
    "time"

    "../config"

    "github.com/thoj/go-ircevent"
)

type Bot struct {
    client      *irc.Connection
    cfg         *config.Config
    Log         log.Logger
    Quitted     chan bool 
}

func (bot *Bot) Quit() {
    bot.Quitted <- true
}

func (bot *Bot) Connect() error {
    bot.Init()
    if (bot.cfg.IRC.Port != "") {
        return bot.client.Connect(net.JoinHostPort(bot.cfg.IRC.Host, bot.cfg.IRC.Port))
    }
    return bot.client.Connect(bot.cfg.IRC.Host)
}

func (bot *Bot) Init() {
    bot.client.AddCallback("001", bot.JoinChannels)
}

func New(cfg *config.Config) (*Bot, error) {
    client := irc.IRC(cfg.IRC.Nick, "go-ircbot")
    client.UseTLS = cfg.IRC.UseSSL
    client.Log = log.New(os.Stdout, "[irc     ] ", log.LstdFlags)

    if cfg.IRC.Password != "" {
        client.Password = cfg.IRC.Password
    }

    if (cfg.IRC.Debug) {
        client.Log.Println("enabling debug mode")
        client.Debug = true
    }

    bot := &Bot{
        cfg:        cfg,
        client:     client,
        Quitted:    make(chan bool, 1),
    }

    return bot, nil
}


// default callbacks
func (bot *Bot) JoinChannels(event *irc.Event) {
    time.Sleep(time.Second)    

    bot.client.Log.Printf("%+v\n", event)

    // only join adminchannel for now
    bot.client.Join(bot.cfg.IRC.AdminChannel)
}

