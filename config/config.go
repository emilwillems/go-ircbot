package config

import (
    "os"
    "log"
    "io/ioutil"
    "encoding/json"
)

type Config_IRC struct {
    Host            string
    Port            string
    Nick            string
    Password        string
    UseSSL          bool

    AdminChannel    string
    CommandChar     string

    Debug           bool
}

type Config struct {
    IRC         Config_IRC    
}

// Loads configuration from given filename, this function is critical
// so it will os.exit(1) using log.Fatal when something bad[tm] happens.
func LoadFromFile(filename string) *Config {
    var conf Config

    log := log.New(os.Stdout, "[config  ] ", log.LstdFlags)
    log.Printf("reading configuration from %s\n", filename)

    file, err := ioutil.ReadFile(filename)
    if err != nil {
        log.Fatalf("could not read configuration: %s\n", err)
    }

    err = json.Unmarshal(file, &conf)
    if err != nil {
        log.Fatalf("could not parse configuration: %s\n", err)
    }

    log.Println("configuration loaded")
    return &conf 
}
