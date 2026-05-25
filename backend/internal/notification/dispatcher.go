package notification

import (
	"encoding/json"
	"fmt"
	"log"

	"upp/internal/model"
	"upp/internal/repository"
	"upp/internal/scheduler"
)

type Dispatcher struct {
	notificationRepo *repository.NotificationRepository
}

func NewDispatcher(repo *repository.NotificationRepository) *Dispatcher {
	return &Dispatcher{notificationRepo: repo}
}

func (d *Dispatcher) Start(ch chan scheduler.NotificationEvent) {
	go func() {
		for event := range ch {
			d.dispatch(event)
		}
	}()
}

func (d *Dispatcher) dispatch(event scheduler.NotificationEvent) {
	notifications, err := d.notificationRepo.FindEnabledByUserID(event.Monitor.UserID)
	if err != nil {
		log.Printf("[Notification] Error fetching notifications: %v", err)
		return
	}

	var title, message string
	if event.NewStatus == "down" {
		title = fmt.Sprintf("[DOWN] %s", event.Monitor.Name)
		message = fmt.Sprintf("Monitor %s is DOWN\nTarget: %s\nMessage: %s", event.Monitor.Name, event.Monitor.Target, event.Result.Message)
	} else {
		title = fmt.Sprintf("[UP] %s", event.Monitor.Name)
		message = fmt.Sprintf("Monitor %s is UP\nTarget: %s\nLatency: %.0fms", event.Monitor.Name, event.Monitor.Target, event.Result.Latency)
	}

	for _, n := range notifications {
		notifier, err := d.createNotifier(n)
		if err != nil {
			log.Printf("[Notification] Error creating notifier %d: %v", n.ID, err)
			continue
		}
		if err := notifier.Send(title, message); err != nil {
			log.Printf("[Notification] Error sending to %s notifier %d: %v", n.Type, n.ID, err)
		}
	}
}

func (d *Dispatcher) createNotifier(n model.Notification) (Notifier, error) {
	switch n.Type {
	case "lark":
		var config LarkNotifier
		if err := json.Unmarshal([]byte(n.Config), &config); err != nil {
			return nil, err
		}
		return &config, nil
	case "dingtalk":
		var config DingTalkNotifier
		if err := json.Unmarshal([]byte(n.Config), &config); err != nil {
			return nil, err
		}
		return &config, nil
	case "webhook":
		var config WebHookNotifier
		if err := json.Unmarshal([]byte(n.Config), &config); err != nil {
			return nil, err
		}
		if config.Method == "" {
			config.Method = "POST"
		}
		return &config, nil
	case "email":
		var config EmailNotifier
		if err := json.Unmarshal([]byte(n.Config), &config); err != nil {
			return nil, err
		}
		return &config, nil
	default:
		return nil, fmt.Errorf("unknown notification type: %s", n.Type)
	}
}

func (d *Dispatcher) SendTest(n *model.Notification) error {
	notifier, err := d.createNotifier(*n)
	if err != nil {
		return err
	}
	return notifier.Send("[TEST] UPP Notification Test", "This is a test notification from UPP monitoring system.")
}
