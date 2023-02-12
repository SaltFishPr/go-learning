package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	cmap "github.com/orcaman/concurrent-map/v2"
	gogpt "github.com/sashabaranov/go-gpt3"
)

const (
	MaxTokens = 500
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	srv := fiber.New(fiber.Config{
		CaseSensitive: true,
	})

	app := NewApp()

	srv.Post("/v1/chat", app.Chat)
	srv.Delete("/v1/chat", app.DeleteConversation)
	srv.Get("/v1/chat", app.GetConversation)

	go func() {
		if err := srv.Listen(":3000"); err != nil {
			log.Panic(err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	<-sig
	log.Println("Gracefully shutting down...")
	_ = srv.Shutdown()

	log.Println("Running cleanup tasks...")

	// Your cleanup tasks go here
	// db.Close()
	// redisConn.Close()
	log.Println("Fiber was successful shutdown.")
}

type App struct {
	client *gogpt.Client
	convs  cmap.ConcurrentMap[string, *Conversation]
}

func NewApp() *App {
	return &App{
		client: gogpt.NewClient(os.Getenv("OPENAI_API_KEY")),
		convs:  cmap.New[*Conversation](),
	}
}

func (a *App) Chat(c *fiber.Ctx) error {
	username := c.GetReqHeaders()["X-Username"]
	if username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "username is required",
		})
	}

	conv, ok := a.convs.Get(username)
	if !ok {
		_conv, err := NewConversation(username, "Bot", WithPromptPrefix("Let's talk about something..."))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		a.convs.Set(username, _conv)
		conv = _conv
	}

	conv.Say(utils.CopyString(c.Query("message")))

	req := gogpt.CompletionRequest{
		Model:       gogpt.GPT3TextDavinci003,
		Prompt:      conv.GetPrompt(),
		MaxTokens:   MaxTokens,
		Temperature: 0.4,
		Stop:        []string{conv.User(), conv.Bot()},
	}
	resp, err := a.client.CreateCompletion(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	conv.Listen(resp.Choices[0].Text)

	return c.JSON(fiber.Map{
		"message":      resp.Choices[0].Text,
		"total_tokens": resp.Usage.TotalTokens,
	})
}

func (a *App) DeleteConversation(c *fiber.Ctx) error {
	username := c.GetReqHeaders()["X-Username"]
	if username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "username is required",
		})
	}

	a.convs.Remove(username)
	return c.SendStatus(fiber.StatusNoContent)
}

func (a *App) GetConversation(c *fiber.Ctx) error {
	username := c.GetReqHeaders()["X-Username"]
	if username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "username is required",
		})
	}

	conv, ok := a.convs.Get(username)
	if !ok {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "conversation not found",
		})
	}

	return c.JSON(fiber.Map{
		"record": conv.GetRecord(),
	})
}
