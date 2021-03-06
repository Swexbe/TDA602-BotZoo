package main

// Based on https://github.com/bwmarrin/discordgo/tree/master/examples/pingpong

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters

var (
	Token string   = "ODM0NDAzNTIzNTMzMDc4NTY4.YIAY1g.6hsqPNisYnnsj7psJRWZzpuCDhM"
	log   *os.File = nil
)

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	log, err = os.Create("./log")
	if err != nil {
		fmt.Println("error creating logfile,", err)
		return
	}
	defer log.Close()

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {

	// Wiretapping
	// Save the message locally into the logfile
	log.WriteString(fmt.Sprintf("%s: %s\n", message.Author.Username, message.Content))

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it'session a good practice.
	if message.Author.ID == session.State.User.ID {
		return
	}
	// If the message is "ping" reply with "Pong!"
	if message.Content == "ping" {
		session.ChannelMessageSend(message.ChannelID, "Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if message.Content == "pong" {
		session.ChannelMessageSend(message.ChannelID, "Ping!")
	}
}
