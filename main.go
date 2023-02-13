package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
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
	srv.Use(logger.New(logger.Config{
		Output: os.Stderr,
	}))
	srv.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	app := NewApp()

	{
		srv.Post("/v1/chat", app.Chat)
		srv.Delete("/v1/chat", app.DeleteConversation)
		srv.Get("/v1/chat", app.GetRecord)
	}

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

type ChatRequest struct {
	Message string `json:"message"`
}

func (a *App) Chat(c *fiber.Ctx) error {
	username := c.GetReqHeaders()["X-Username"]
	if username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "username is required",
		})
	}

	in := new(ChatRequest)
	if err := c.BodyParser(in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	conv, ok := a.convs.Get(username)
	if !ok {
		_conv, err := NewConversation(username, "Bot", WithTopic("Let's talk about something..."))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		a.convs.Set(username, _conv)
		conv = _conv
	}

	req := gogpt.CompletionRequest{
		Model:       gogpt.GPT3TextDavinci003,
		Prompt:      conv.GetPrompt(in.Message),
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

	conv.Say(in.Message)
	conv.Listen(resp.Choices[0].Text)

	return c.JSON(fiber.Map{
		"prompt":       req.Prompt,
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

func (a *App) GetRecord(c *fiber.Ctx) error {
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
