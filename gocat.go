package main

import "strings"
import "time"
import "fmt"
import "log"
import "flag"
import irc "github.com/fluffle/goirc/client"

func CatMsgSender(ch chan CatMsg, client *irc.Conn) {
    defaultChan := "#rjtest"
    for {
        cm := <-ch
        if len(cm.To) == 0 {
            client.Privmsg(defaultChan, cm.Msg)
        } else {
            for _, to := range cm.To {
                client.Privmsg(to, cm.Msg)
            }
        }
    }
}

func setupClient(c *irc.Conn, chConnected chan bool) {
    c.AddHandler(irc.CONNECTED, func(conn *irc.Conn, line *irc.Line) {
        chConnected <- true
    })
    c.AddHandler(irc.DISCONNECTED, func(conn *irc.Conn, line *irc.Line) {
        chConnected <- false
    })
    c.AddHandler("PRIVMSG", func(conn *irc.Conn, line *irc.Line) {
        if len(line.Args) < 2 || !strings.HasPrefix(line.Args[1], "!") {
            return
        }
        to := line.Args[0]
        sender := to
        if to == c.Me.Nick {
            // TODO: check if sender is in main chan, else return
            log.Println("Got ! via PM from " + line.Src)
            sender = line.Src // replies go via PM too.
        } else {
            log.Println("Got ! via chan: " + line.Args[0] + " from " + line.Src)
        }
        log.Println(line.Args)
        switch line.Args[0] {
        case "!join":
            if len(line.Args) == 2 {
                c.Join(line.Args[1])
            } else if len(line.Args) == 3 {
                c.Join(line.Args[1] + " " + line.Args[2])
            } else {
                c.Privmsg(sender, "Usage: !join #chan  or  !join #chan key")
            }
        case "!part":
            if len(line.Args) == 2 {
                c.Part(line.Args[1])
            } else {
                c.Privmsg(sender, "Usage: !part #chan")
            }
        default:
            c.Privmsg(sender, "Invalid command: "+strings.Join(line.Args[1:], " "))
            return
        }
    })

}

func main() {
    fromXML("./conf.xml")
    var irchost = flag.String(
        "irchost",
        "localhost",
        "Hostname of IRC server, eg: irc.example.org:6667")
    var ircnick = flag.String(
        "ircnick",
        "gocat",
        "Nickname to use for IRC")
    var ircssl = flag.Bool(
        "ircssl",
        false,
        "Use SSL for IRC connection")
    var catbind = flag.String(
        "catbind",
        ":12345",
        "net.Listen spec, to listen for IRCCat msgs")
    var catfam = flag.String(
        "catfamily",
        "tcp4",
        "net.Listen address family for IRCCat msgs")
    flag.Parse()

    fmt.Println("Go irccat")

    // to block main:
    control := make(chan string)
    // msgs from tcp catport to this chan:
    catmsgs := make(chan CatMsg)
    // channel signaling irc connection status
    chConnected := make(chan bool)
    // setup IRC client:
    c := irc.SimpleClient(*ircnick)
    c.SSL = *ircssl
    // Listen on catport:
    go CatportServer(catmsgs, *catfam, *catbind)
    go CatMsgSender(catmsgs, c)
    // loop on IRC dis/connected events
    setupClient(c, chConnected)
    for {
        log.Println("Connecting to IRC...")
        err := c.Connect(*irchost)
        if err != nil {
            log.Println("Failed to connect to IRCd")
            log.Println(err)
            continue
        }
        for {
            status := <-chConnected
            if status {
                log.Println("Connected to IRC")
                c.Join("#rjtest")
            } else {
                log.Println("Disconnected from IRC")
                break
            }
        }
        time.Sleep(5 * time.Second)
    }

    <-control
}
