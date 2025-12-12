# Gator Project - Essential Improvements

## Phase 2: Structural Changes (Affects All Imports)

- [ ] **Restructure `internal/` into domain packages**
  - **Why now:** Changes imports everywhere, so do before other refactoring
  - **Side effect:** Updates required in `main.go` and all internal imports

  **New structure:**
  ```
  internal/
  ├── cli/
  │   ├── dispatcher.go     (from: command.go)
  │   ├── middleware.go     (from: middleware_login.go)
  │   ├── auth.go           (from: login.go, register.go, reset.go, users.go)
  │   ├── feeds.go          (from: addFeed.go, feeds.go, follow.go, following.go, unfollow.go)
  │   └── posts.go          (from: browse.go)
  ├── aggregator/
  │   └── aggregator.go     (from: agg.go)
  ├── rss/
  │   └── client.go         (from: rss.go)
  ├── config/               (unchanged)
  └── database/             (unchanged)
  ```

  **Steps:**
  1. Create directories: `internal/cli/`, `internal/aggregator/`, `internal/rss/`
  2. Move and consolidate files (merge related handlers into single files)
  3. Update package declarations in moved files
  4. Update all imports in `main.go`
  5. Test: `go build` should succeed

---

## Phase 3: Input Validation

- [ ] **Add argument length checks to all handlers**
  - **Files affected:** `internal/cli/auth.go`, `internal/cli/feeds.go`, `internal/cli/posts.go`, `internal/aggregator/aggregator.go`
  - **Side effect:** None (makes code safer)

  **Example:**
  ```go
  // In login handler
  if len(cmd.Args) < 1 {
      return fmt.Errorf("usage: login <username>")
  }

  // In register handler
  if len(cmd.Args) < 1 {
      return fmt.Errorf("usage: register <username>")
  }

  // In addfeed handler
  if len(cmd.Args) < 2 {
      return fmt.Errorf("usage: addfeed <name> <url>")
  }

  // In follow handler
  if len(cmd.Args) < 1 {
      return fmt.Errorf("usage: follow <url>")
  }

  // In unfollow handler
  if len(cmd.Args) < 1 {
      return fmt.Errorf("usage: unfollow <url>")
  }

  // In agg handler
  if len(cmd.Args) < 1 {
      return fmt.Errorf("usage: agg <time_between_reqs>")
  }
  ```

- [ ] **Validate feed URLs**
  - **Files affected:** `internal/cli/feeds.go` (in addfeed and follow handlers)
  - **Side effect:** Prevents invalid data in database

  **Example:**
  ```go
  url := cmd.Args[1] // or cmd.Args[0] for follow
  if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
      return fmt.Errorf("invalid URL: must start with http:// or https://")
  }
  ```

- [ ] **Validate duration format in aggregator**
  - **Files affected:** `internal/aggregator/aggregator.go`
  - **Side effect:** None

  **Example:**
  ```go
  timeBetweenRequests := cmd.Args[0]
  duration, err := time.ParseDuration(timeBetweenRequests)
  if err != nil {
      return fmt.Errorf("invalid duration format: %w", err)
  }
  ```

- [ ] **Sanitize user inputs (trim whitespace, check empty)**
  - **Files affected:** All handlers that accept string input
  - **Side effect:** Cleaner data in database

  **Example:**
  ```go
  name := strings.TrimSpace(cmd.Args[0])
  if name == "" {
      return fmt.Errorf("username cannot be empty")
  }
  ```

---

## Phase 4: Error Handling

- [ ] **Wrap all errors with context**
  - **Files affected:** All handlers in `internal/cli/`, `internal/aggregator/`, `internal/rss/`
  - **Side effect:** Better error messages (no breaking changes)

  **Examples:**
  ```go
  // In CreateUser calls
  user, err := s.Db.CreateUser(...)
  if err != nil {
      return fmt.Errorf("failed to create user %q: %w", name, err)
  }

  // In feed fetching
  rssFeed, err := rss.FetchFeed(ctx, feedURL)
  if err != nil {
      return fmt.Errorf("failed to fetch feed from %q: %w", feedURL, err)
  }

  // In database operations
  feed, err := s.Db.GetFeedByUrl(ctx, url)
  if err != nil {
      return fmt.Errorf("failed to find feed with URL %q: %w", url, err)
  }
  ```

- [ ] **Replace string error matching with type checking**
  - **Files affected:** `internal/aggregator/aggregator.go`
  - **Side effect:** More reliable error detection

  **Replace:**
  ```go
  // OLD (brittle)
  if strings.Contains(err.Error(), "duplicate") {
      continue
  }

  // NEW (reliable)
  var pqErr *pq.Error
  if errors.As(err, &pqErr) && pqErr.Code == "23505" { // unique_violation
      log.Printf("Post already exists, skipping: %s", item.Link)
      continue
  }
  ```

  **Note:** Need to import `"github.com/lib/pq"` and `"errors"`

---

## Phase 5: Extract Magic Constants

- [ ] **Define constants for default values**
  - **Files affected:** `internal/cli/posts.go`
  - **Side effect:** None (makes code clearer)

  **In `internal/cli/posts.go`:**
  ```go
  const DefaultPostLimit = 2

  // Then in browse handler
  limit := DefaultPostLimit
  if len(cmd.Args) > 0 {
      if userLimit, err := strconv.Atoi(cmd.Args[0]); err == nil && userLimit > 0 {
          limit = userLimit
      }
  }
  ```

---

## Phase 6: Security & Configuration

- [ ] **Add environment variable support for DB URL**
  - **Files affected:** `main.go` and/or `internal/config/config.go`
  - **Side effect:** None (fallback to config file)

  **In `main.go` after reading config:**
  ```go
  cfg, err := config.Read()
  if err != nil {
      log.Fatalf("error reading config: %v", err)
  }

  // Prefer environment variable for security
  dbURL := os.Getenv("DATABASE_URL")
  if dbURL == "" {
      dbURL = cfg.DbUrl
  }

  db, err := sql.Open("postgres", dbURL)
  ```

  **Document in README.md:**
  ```markdown
  ## Configuration

  Database URL can be provided via:
  1. Environment variable: `DATABASE_URL` (recommended for security)
  2. Config file: `~/.gatorconfig.json`

  Example:
  ```bash
  export DATABASE_URL="postgres://user:password@localhost:5432/gator"
  ./gator agg 1m
  ```
  ```

- [ ] **Configure database connection pool**
  - **Files affected:** `main.go`
  - **Side effect:** Better performance, prevents connection exhaustion

  **In `main.go` after `sql.Open`:**
  ```go
  db, err := sql.Open("postgres", dbURL)
  if err != nil {
      log.Fatalf("error opening database: %v", err)
  }
  defer db.Close()

  // Configure connection pool
  db.SetMaxOpenConns(25)
  db.SetMaxIdleConns(5)
  db.SetConnMaxLifetime(5 * time.Minute)
  ```

---

## Phase 7: Runtime Improvements

- [ ] **Add graceful shutdown to aggregator**
  - **Files affected:** `internal/aggregator/aggregator.go`
  - **Side effect:** Clean exit on Ctrl+C

  **Implementation:**
  ```go
  import (
      "os"
      "os/signal"
      "syscall"
  )

  func HandlerAggregator(s *State, cmd Command) error {
      // ... existing validation ...

      ticker := time.NewTicker(duration)
      defer ticker.Stop()

      // Setup signal handling
      sigChan := make(chan os.Signal, 1)
      signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

      fmt.Printf("Collecting feeds every %v. Press Ctrl+C to stop.\n", duration)

      for {
          select {
          case <-ticker.C:
              // ... existing scraping logic ...
          case <-sigChan:
              fmt.Println("\nShutting down gracefully...")
              return nil
          }
      }
  }
  ```

- [ ] **Replace `fmt.Println` with structured logging**
  - **Files affected:** All handlers
  - **Side effect:** Better production logging

  **Setup in `main.go`:**
  ```go
  import "log/slog"

  func main() {
      // Setup structured logger
      logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
      slog.SetDefault(logger)

      // ... rest of main
  }
  ```

  **Usage examples:**
  ```go
  // Instead of: fmt.Println("Created user:", user.Name)
  slog.Info("user created", "name", user.Name, "id", user.ID)

  // Instead of: fmt.Printf("Logged in as %s\n", user.Name)
  slog.Info("user logged in", "name", user.Name)

  // Instead of: fmt.Println("Feed added!")
  slog.Info("feed added", "name", feed.Name, "url", feed.Url, "user", user.Name)

  // For errors (instead of returning silently)
  slog.Error("failed to create post", "error", err, "title", item.Title)
  ```

---

## Phase 8: Performance Optimizations

- [ ] **Fetch multiple feeds concurrently in aggregator**
  - **Files affected:** `internal/aggregator/aggregator.go`
  - **Side effect:** Much faster scraping

  **Current approach:** Fetches 1 feed per tick

  **New approach:** Fetch top 10 unfetched feeds in parallel

  **Implementation:**
  ```go
  // In aggregator loop
  for {
      select {
      case <-ticker.C:
          scrapeFeeds(ctx, s)
      case <-sigChan:
          return nil
      }
  }

  func scrapeFeeds(ctx context.Context, s *State) {
      // Get next 10 feeds to fetch (needs new SQL query)
      feeds, err := s.Db.GetNextFeedsToFetch(ctx, 10)
      if err != nil {
          slog.Error("failed to fetch feeds list", "error", err)
          return
      }

      var wg sync.WaitGroup
      for _, feed := range feeds {
          wg.Add(1)
          go func(f database.Feed) {
              defer wg.Done()
              fetchAndStoreFeed(ctx, s, f)
          }(feed)
      }
      wg.Wait()
  }
  ```

  **Requires new SQL query in `sql/queries/feeds.sql`:**
  ```sql
  -- name: GetNextFeedsToFetch :many
  SELECT * FROM feeds
  ORDER BY last_fetched_at ASC NULLS FIRST
  LIMIT $1;
  ```

  **Then regenerate:** `sqlc generate`

- [ ] **Add timeout to HTTP requests**
  - **Files affected:** `internal/rss/client.go`
  - **Side effect:** Prevents hanging on slow/dead feeds

  **Implementation:**
  ```go
  func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
      // Add 30 second timeout
      ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
      defer cancel()

      req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
      if err != nil {
          return nil, err
      }

      req.Header.Add("User-Agent", "gator")

      // ... rest unchanged
  }
  ```

- [ ] **Add database indexes for performance**
  - **Files affected:** New migration file
  - **Side effect:** Faster queries on large datasets

  **Create:** `sql/schema/006_indexes.sql`
  ```sql
  -- +goose Up
  CREATE INDEX IF NOT EXISTS idx_posts_feed_id ON posts(feed_id);
  CREATE INDEX IF NOT EXISTS idx_posts_published_at ON posts(published_at DESC);
  CREATE INDEX IF NOT EXISTS idx_feeds_last_fetched_at ON feeds(last_fetched_at NULLS FIRST);
  CREATE INDEX IF NOT EXISTS idx_feed_follows_user_id ON feed_follows(user_id);

  -- +goose Down
  DROP INDEX IF EXISTS idx_posts_feed_id;
  DROP INDEX IF EXISTS idx_posts_published_at;
  DROP INDEX IF EXISTS idx_feeds_last_fetched_at;
  DROP INDEX IF EXISTS idx_feed_follows_user_id;
  ```

  **Apply:** `make migrate-up`

---

## Phase 9: Developer Experience

- [ ] **Run `gofmt` on all files**
  - **Command:** `gofmt -w .`
  - **Side effect:** Consistent formatting

- [ ] **Setup `golangci-lint`**
  - **Install:** `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`
  - **Create:** `.golangci.yml`

  ```yaml
  linters:
    enable:
      - gofmt
      - govet
      - errcheck
      - staticcheck
      - unused
      - gosimple
      - ineffassign

  issues:
    exclude-use-default: false
  ```

  - **Run:** `golangci-lint run`
  - **Side effect:** Catches bugs and style issues

- [ ] **Add Makefile development targets**
  - **Files affected:** `Makefile`
  - **Side effect:** Easier development workflow

  **Add to existing Makefile:**
  ```makefile
  .PHONY: build run lint fmt clean install-tools

  build:
  	go build -o bin/gator .

  run:
  	go run main.go

  lint:
  	golangci-lint run

  fmt:
  	gofmt -w .

  clean:
  	rm -rf bin/

  install-tools:
  	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
  	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
  	go install github.com/pressly/goose/v3/cmd/goose@latest
  ```

- [ ] **Add package documentation**
  - **Files affected:** All packages
  - **Side effect:** Shows in godoc, better IDE experience

  **Examples:**
  ```go
  // Package cli provides command-line interface handlers for the Gator RSS reader.
  package cli

  // Package rss provides RSS feed fetching and parsing functionality.
  package rss

  // Package aggregator provides background feed scraping functionality.
  package aggregator
  ```

- [ ] **Add function comments for exported functions**
  - **Files affected:** All packages
  - **Side effect:** Go conventions, IDE tooltips

  **Examples:**
  ```go
  // FetchFeed retrieves and parses an RSS feed from the given URL.
  // It returns the parsed feed or an error if the request fails or the XML is invalid.
  func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {

  // Run executes a registered command by name with the given arguments.
  // It returns an error if the command is not found or if execution fails.
  func (c *Commands) Run(s *State, cmd Command) error {
  ```

---

## Summary

**Total:** 24 essential improvements

**Phases:**
1. Prerequisites (3 tasks) - Fix typos, security basics
2. Structure (1 task) - Package reorganization ⚠️ **Affects all imports**
3. Validation (4 tasks) - Input safety
4. Errors (2 tasks) - Better error messages
5. Constants (1 task) - Remove magic numbers
6. Security (2 tasks) - DB URL, connection pooling
7. Runtime (2 tasks) - Graceful shutdown, logging
8. Performance (3 tasks) - Concurrency, timeouts, indexes
9. Dev Experience (6 tasks) - Linting, docs, formatting

**Critical Path:**
- Phase 1 → Phase 2 (structure) → Everything else (order flexible within phases)

**Learning Notes:**

> **Package Organization**: When you have multiple related files, group them by domain (auth, feeds, posts) not by type. Makes codebases scale.

> **Input Validation**: Always validate at the boundary (CLI args, HTTP requests). Never trust external data.

> **Error Wrapping**: Always add context when returning errors: `fmt.Errorf("what failed: %w", err)`. Makes debugging 10x easier.

> **Graceful Shutdown**: Long-running services need signal handling. Use `select` with signal channel.

> **Structured Logging**: Production code uses `slog` or similar. `fmt.Println` is for learning/debugging only.

> **Concurrency**: Use goroutines for I/O (HTTP, DB). Use `sync.WaitGroup` to wait for completion.
