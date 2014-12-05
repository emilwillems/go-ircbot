package bot

import (
    "os"
    "log"
    "net"
    "time"
    "strings"
    "fmt"

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
    bot.client.AddCallback("PRIVMSG", bot.RunDefaultCommands)
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

func (bot *Bot) Say(channel, message string) {
    bot.client.Privmsg(channel, message)
}

func (bot *Bot) Action(channel, action string) {
    bot.client.Privmsg(channel, fmt.Sprintf("\001ACTION %s\001", action))
}

func (bot *Bot) JoinChannels(event *irc.Event) {
    time.Sleep(time.Second)    
    bot.client.Join(bot.cfg.IRC.AdminChannel)
}

func (bot *Bot) RunDefaultCommands(event *irc.Event) {
    //bot.client.Log.Printf("%+v\n", event)

    // split message on space, grab first element, process args
    args := strings.Split(strings.TrimSpace(event.Message()), " ")
    command := args[0]
    if len(args) > 1 {
        args = args[1:]
    } else {
        args = []string{}
    }

    // only trigger on messages starting with our configured command character
    if !strings.HasPrefix(command, bot.cfg.IRC.CommandChar) {
        return 
    }

    command = strings.TrimPrefix(command, "?")
    channel := event.Arguments[0]
    switch command {
        case "quit": 
            bot.Say(channel, "aye aye cap'n!")
            bot.Quit()
    }
}
