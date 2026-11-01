package main

import (
	"fmt"

	"go-course/hw_4/internal/config"
	"go-course/hw_4/internal/handler"
	"go-course/hw_4/internal/notifier"
	"go-course/hw_4/internal/repository"
	"go-course/hw_4/internal/service"
)

func main() {
	cfg := config.Load()
	repo := repository.NewMemoryRepository()
	notifiers := []service.Notifier{notifier.NewEmail(cfg.Email), notifier.NewPush(cfg.PushDevice)}
	svc := service.NewService(repo, notifiers)
	h := handler.New(svc)

	for _, message := range []string{"Заказ готов", ""} {
		n, err := h.Send(handler.SendRequest{Message: message})
		if err != nil {
			fmt.Println("ошибка:", err)
		}
		fmt.Printf("сохранено: %+v\n", n)
	}
}
