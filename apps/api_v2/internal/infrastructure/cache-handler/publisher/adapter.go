package cache_publisher

import (
	rabbitmq_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/publisher/rabbitmq"
)

type CachePublisherAdapter struct {
	publisher *rabbitmq_infrastructure.RabbitMQAdapter
}

func NewCachePublisherAdapter() *CachePublisherAdapter {
	return &CachePublisherAdapter{}
}

func (a *CachePublisherAdapter) Save(message SyncMessage) error {
	return nil
}

func (a *CachePublisherAdapter) Update(message SyncMessage) error {
	return nil
}