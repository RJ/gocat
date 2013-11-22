package main

import (
    "syscall"
    "log"
    "io"
    "fmt"
    xmlx "github.com/jteeuwen/go-pkg-xmlx"
    gcfg "code.google.com/p/gcfg"
)

type Settings struct {
    nick string
    
    connection_host string
    connection_ssl bool
    connection_port int16
    connection_password string

    postconnectcmds string

    cat_port int16
    cat_ip string
    cat_messagedelay int16

    script_path string
    script_maxlines int

    // name of channel, membership of means you are authed with the bot
    auth_channel string

    // #channel -> password
    channels map[string]string
}

func charsetIdentityFun(charset string, input io.Reader) (io.Reader,error) {
    return input, nil
}

func fromINI(path string) Settings {
    s := Settings{}
    cfg := struct {
        Section struct {
            Name string
        }
    }{}
    err := gcfg.ReadFileInto(&cfg, path)
    if err != nil {
        log.Println("Error reading ini file: ", err)
        syscall.Exit(1)
    }
    return s
}

func fromXML(path string) Settings {
    s := Settings{}

    doc := xmlx.New()
    if err := doc.LoadFile(path, charsetIdentityFun); err != nil {
        fmt.Println("Error reading xml config file: ", err)
        syscall.Exit(1)
    }
    node := doc.SelectNode("", "server")
    s.connection_host = node.S("","address")
    s.connection_port = node.I16("","port")
    s.connection_password = node.S("","password")

    node = doc.SelectNode("", "bot")
    s.nick = node.S("","nick")
    s.cat_messagedelay = node.I16("","messagedelay")

    node = doc.SelectNode("", "cat")
    s.cat_port = node.I16("", "port")
    s.cat_ip = node.S("", "ip")

    node = doc.SelectNode("", "script")
    s.script_path = node.S("", "cmdhandler")
    s.script_maxlines = node.I("", "maxresponselines")


    s.channels = make(map[string]string)
    node = doc.SelectNode("", "channels")
    channodes := node.SelectNodes("", "channel")
    for i, n := range channodes {
        chan_name := n.S("", "name")
        chan_pass := n.S("", "password")
        if i == 0 {
            s.auth_channel = chan_name
        }
        s.channels[chan_name] = chan_pass
    }

    log.Println("Conf: ", s)

    syscall.Exit(0)

    return s
}
