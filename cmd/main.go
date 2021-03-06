package main

import (
	"log"
	"os"

	"github.com/alebsys/telegram-article-bot/internal/devto/article"
	"github.com/alebsys/telegram-article-bot/internal/devto/podcast"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	// descp        = "`Request example:\n/article go 10 5\n* go - topic (tag);\n* 10 - search period in days;\n* 5 - number of posts.`"
	descpArticle = "`Request example:\n/article go 10 5\n* go - topic (tag);\n* 10 - search period in days;\n* 5 - number of posts.`"
	descpPodcast = "`Request example:\n/podcast gotime\n* gotime - topic (tag).`"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Panic("getting TELEGRAM_APITOKEN: ", err)
	}
	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.EditedMessage != nil {
			continue
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		msg.ParseMode = "markdown"
		msg.DisableWebPagePreview = true

		input := update.Message.Text

		log.Printf("[%s] %s", update.Message.From.UserName, input)

		switch update.Message.Command() {
		case "help":
			msg.Text = "`Hello! I can find articles and podcasts of interest to you on DEV.TO\n\n`" + descpArticle + "\n\n`or`\n\n" + descpPodcast
		case "article":
			note := "`Enter the correct command!\n\n`" + descpArticle

			b, err := article.ValidateInput(input)
			if err != nil {
				log.Print(err)
			}
			if !b {
				msg.Text = note
				break
			}

			query, err := article.ParseInput(input)
			if err != nil {
				log.Print(err)
				continue
			}
			articles, err := article.GetArticles(query.Tag, query.Freshness)
			if err != nil {
				log.Print(err)
				continue
			}

			msg.Text = articles.WriteArticles(query.Limit)
		case "podcast":
			note := "`Enter the correct command!\n\n`" + descpPodcast

			b, err := podcast.ValidateInput(input)
			if err != nil {
				log.Print(err)
			}
			if !b {
				msg.Text = note
				break
			}

			query := podcast.ParseInput(input)

			podcasts, err := podcast.GetPodcasts(query.Tag)
			if err != nil {
				log.Print(err)
				continue
			}

			msg.Text = podcasts.WritePodcasts()
		default:
			msg.Text = "`I don't know this command. Enter /help`"
		}

		bot.Send(msg)
	}

}
