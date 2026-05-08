# Cache + Roble Implementation

## Summary

A MariaDB SQL cache layer sits in front of the Roble HTTP API to reduce latency and load during high-traffic scenarios (exams). Reads hit cache first; writes go to cache and sync to Roble asynchronously via RabbitMQ. The cache is a local service shared by all API instances on the same host.

## Technologies

| Component | Technology | Role |
|---|---|---|
| Cache DB | MariaDB | Persistent SQL cache with full MVCC concurrency |
| Message broker | RabbitMQ (existing) | Durable queue for async Roble sync signals |
| Backend | Go (api_v2) | Cache read/write logic, async consumer |

## Operation Rules

**Read:** Check cache → if miss, fetch from Roble → save to cache.

**Create / Update:** Write to cache with `synced = false` → publish sync message to RabbitMQ → consumer writes to Roble → mark `synced = true`.

**Delete:** Delete from Roble first → delete from cache.

**Eviction:** Each table has a configurable TTL. An async process evicts records past their TTL every 30 minutes. Critical tables (susceptible to external changes) use short TTLs (seconds).

## Business Logic Decisions

- All API instances run on the same host and share one MariaDB instance — no cross-instance cache inconsistency.
- Every cache table has a primary key to prevent duplicate inserts on concurrent cache-miss races.
- Every cache record has a `synced` boolean column. On boot, any `synced = false` record is re-queued to RabbitMQ for reconciliation after a crash.
- RabbitMQ absorbs burst write signals — the async consumer processes them at the rate Roble can handle.
- TTL is per-table and resets on access. Tables with external write risk use short TTLs instead of a global invalidation strategy.
- Cache stampede mitigation: use a per-key singleflight so concurrent misses for the same record only fire one Roble request.
