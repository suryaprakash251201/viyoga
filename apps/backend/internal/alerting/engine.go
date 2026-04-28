package alerting

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

// AlertRule defines a threshold-based alert.
type AlertRule struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	MetricType    string  `json:"metric_type"` // "cpu", "memory", "disk"
	Condition     string  `json:"condition"`   // "gt", "lt"
	Threshold     float64 `json:"threshold"`
	NotifyChannel string  `json:"notify_channel"` // "webhook", "telegram"
	NotifyTarget  string  `json:"notify_target"`
	Enabled       bool    `json:"enabled"`
	CooldownMins  int     `json:"cooldown_mins"`
}

// AlertEvent records a triggered alert.
type AlertEvent struct {
	ID          int       `json:"id"`
	RuleID      int       `json:"rule_id"`
	RuleName    string    `json:"rule_name"`
	MetricType  string    `json:"metric_type"`
	Value       float64   `json:"value"`
	Threshold   float64   `json:"threshold"`
	TriggeredAt time.Time `json:"triggered_at"`
	ResolvedAt  *time.Time `json:"resolved_at,omitempty"`
	Acknowledged bool     `json:"acknowledged"`
}

// Engine evaluates metrics against alert rules.
type Engine struct {
	rules      []AlertRule
	lastFired  map[int]time.Time
	events     []AlertEvent
	mu         sync.RWMutex
	notifier   *Notifier
}

// NewEngine creates a new alert engine.
func NewEngine() *Engine {
	return &Engine{
		lastFired: make(map[int]time.Time),
		notifier:  NewNotifier(),
	}
}

// SetRules updates the alert rules.
func (e *Engine) SetRules(rules []AlertRule) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.rules = rules
}

// GetRules returns current alert rules.
func (e *Engine) GetRules() []AlertRule {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return append([]AlertRule{}, e.rules...)
}

// GetEvents returns recent alert events.
func (e *Engine) GetEvents(limit int) []AlertEvent {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if limit <= 0 || limit > len(e.events) {
		limit = len(e.events)
	}
	start := len(e.events) - limit
	if start < 0 {
		start = 0
	}
	return append([]AlertEvent{}, e.events[start:]...)
}

// Evaluate checks metrics against all rules.
func (e *Engine) Evaluate(metricType string, value float64) {
	e.mu.Lock()
	defer e.mu.Unlock()

	now := time.Now()
	for _, rule := range e.rules {
		if !rule.Enabled || rule.MetricType != metricType {
			continue
		}

		triggered := false
		switch rule.Condition {
		case "gt":
			triggered = value > rule.Threshold
		case "lt":
			triggered = value < rule.Threshold
		case "eq":
			triggered = value == rule.Threshold
		}

		if !triggered {
			continue
		}

		// Check cooldown
		if lastTime, ok := e.lastFired[rule.ID]; ok {
			cooldown := time.Duration(rule.CooldownMins) * time.Minute
			if cooldown == 0 {
				cooldown = 5 * time.Minute
			}
			if now.Sub(lastTime) < cooldown {
				continue
			}
		}

		e.lastFired[rule.ID] = now

		event := AlertEvent{
			ID:          len(e.events) + 1,
			RuleID:      rule.ID,
			RuleName:    rule.Name,
			MetricType:  metricType,
			Value:       value,
			Threshold:   rule.Threshold,
			TriggeredAt: now,
		}
		e.events = append(e.events, event)

		// Keep only last 1000 events
		if len(e.events) > 1000 {
			e.events = e.events[len(e.events)-500:]
		}

		log.Warn().
			Str("rule", rule.Name).
			Str("metric", metricType).
			Float64("value", value).
			Float64("threshold", rule.Threshold).
			Msg("Alert triggered")

		// Fire notification in background
		go e.notifier.Send(rule, event)
	}
}

// Notifier sends alert notifications.
type Notifier struct{}

func NewNotifier() *Notifier {
	return &Notifier{}
}

// Send dispatches a notification for an alert event.
func (n *Notifier) Send(rule AlertRule, event AlertEvent) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	switch rule.NotifyChannel {
	case "webhook":
		n.sendWebhook(ctx, rule.NotifyTarget, event)
	case "telegram":
		n.sendTelegram(ctx, rule.NotifyTarget, event)
	default:
		log.Warn().Str("channel", rule.NotifyChannel).Msg("Unknown notification channel")
	}
}

func (n *Notifier) sendWebhook(ctx context.Context, url string, event AlertEvent) {
	if url == "" {
		return
	}

	payload, _ := json.Marshal(map[string]interface{}{
		"alert":     event.RuleName,
		"metric":    event.MetricType,
		"value":     event.Value,
		"threshold": event.Threshold,
		"timestamp": event.TriggeredAt.Format(time.RFC3339),
		"source":    "viyoga",
	})

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(payload))
	if err != nil {
		log.Error().Err(err).Msg("Failed to create webhook request")
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("Webhook notification failed")
		return
	}
	defer resp.Body.Close()

	log.Info().Int("status", resp.StatusCode).Str("url", url).Msg("Webhook notification sent")
}

func (n *Notifier) sendTelegram(ctx context.Context, target string, event AlertEvent) {
	// target format: "bot_token:chat_id"
	// Implementation placeholder — requires Telegram Bot API
	msg := fmt.Sprintf("🚨 *Viyoga Alert*\n\n*%s*\nMetric: `%s`\nValue: `%.1f`\nThreshold: `%.1f`\nTime: %s",
		event.RuleName, event.MetricType, event.Value, event.Threshold,
		event.TriggeredAt.Format("15:04:05"))

	log.Info().Str("message", msg).Msg("Telegram notification (stub)")
}
