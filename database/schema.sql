-- ============================================================
-- TICKETMASTER-STYLE MVP SCHEMA
-- Lean version for microservices practice.
--
-- Logically split into 3 owners (User / Ticketing / Payment),
-- each block below would live in that service's own database in
-- a real deployment. Fields like buyer_user_id, host_user_id,
-- order_id in the Payment block are cross-service references —
-- plain UUIDs, not enforced FKs, since they'd point at tables in
-- a different service's database. Kept in one file here only
-- because it's simpler to work with for a single-repo practice
-- project.
-- ============================================================


-- ============================================================
-- USER SERVICE
-- ============================================================

CREATE TABLE users (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username       VARCHAR NOT NULL UNIQUE,
    email          VARCHAR NOT NULL UNIQUE,
    password_hash  VARCHAR NOT NULL,
    role           VARCHAR NOT NULL DEFAULT 'BUYER',  -- BUYER, HOST
    created_at     TIMESTAMP NOT NULL DEFAULT now()
);


-- ============================================================
-- TICKETING SERVICE
-- ============================================================

-- ---- venue layer: static, reused across events ----

CREATE TABLE venues (
    id       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name     VARCHAR NOT NULL,
    address  VARCHAR NOT NULL,
    city     VARCHAR NOT NULL
);

CREATE TABLE sections (
    id        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    venue_id  UUID NOT NULL REFERENCES venues(id) ON DELETE CASCADE,
    name      VARCHAR NOT NULL,          -- "Floor", "Section A"
    UNIQUE (venue_id, name)
);

CREATE TABLE seats (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    section_id   UUID NOT NULL REFERENCES sections(id) ON DELETE CASCADE,
    row_label    VARCHAR NOT NULL,       -- string, not int: "A", "12"
    seat_number  VARCHAR NOT NULL,       -- string, not int: "14", "14A"
    UNIQUE (section_id, row_label, seat_number)
);

-- ---- event layer: per-event data ----

CREATE TABLE events (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    venue_id         UUID NOT NULL REFERENCES venues(id),
    host_user_id     UUID NOT NULL,              -- cross-service ref -> User Service users.id
    name             VARCHAR NOT NULL,
    start_time       TIMESTAMP NOT NULL,
    seating_mode     VARCHAR NOT NULL,             -- ASSIGNED or GENERAL_ADMISSION
    available_count  INT NOT NULL DEFAULT 0         -- denormalized fast-read counter
);

CREATE TABLE price_tiers (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id       UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    section_id     UUID REFERENCES sections(id),   -- <-- added: ties this tier to a specific section
    name           VARCHAR NOT NULL,
    price_cents    BIGINT NOT NULL,
    currency       VARCHAR(3) NOT NULL DEFAULT 'USD',
    ga_capacity    INT,
    ga_sold_count  INT NOT NULL DEFAULT 0
);

-- one row per (seat, event) — only for ASSIGNED events, bulk-generated
-- from the venue's seats at event creation time
CREATE TABLE seat_inventory (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id         UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    seat_id          UUID NOT NULL REFERENCES seats(id),
    price_tier_id    UUID NOT NULL REFERENCES price_tiers(id),
    status           VARCHAR NOT NULL DEFAULT 'AVAILABLE', -- AVAILABLE, HELD, SOLD
    hold_id          UUID,
    hold_expires_at  TIMESTAMP,
    UNIQUE (event_id, seat_id)
);

CREATE INDEX idx_seat_inventory_event_status ON seat_inventory (event_id, status);

-- ---- order layer ----

CREATE TABLE orders (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    buyer_user_id  UUID NOT NULL,          -- cross-service ref -> User Service users.id
    event_id       UUID NOT NULL REFERENCES events(id),
    status         VARCHAR NOT NULL DEFAULT 'PENDING', -- PENDING, CONFIRMED, FAILED, CANCELLED
    total_cents    BIGINT NOT NULL,
    currency       VARCHAR(3) NOT NULL DEFAULT 'USD',
    payment_id     UUID,                   -- cross-service ref -> Payment Service payments.id
    created_at     TIMESTAMP NOT NULL DEFAULT now()
);

-- snapshots price at purchase time — never re-derive by joining back to price_tiers
CREATE TABLE order_line_items (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id          UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    seat_id           UUID REFERENCES seats(id),  -- null for GENERAL_ADMISSION
    price_tier_id     UUID NOT NULL REFERENCES price_tiers(id),
    price_paid_cents  BIGINT NOT NULL,
    currency          VARCHAR(3) NOT NULL DEFAULT 'USD'
);


-- ============================================================
-- PAYMENT SERVICE
-- ============================================================

CREATE TABLE payments (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id         UUID NOT NULL,          -- cross-service ref -> Ticketing Service orders.id
    buyer_user_id    UUID NOT NULL,          -- cross-service ref -> User Service users.id
    amount_cents     BIGINT NOT NULL,
    currency         VARCHAR(3) NOT NULL DEFAULT 'USD',
    status           VARCHAR NOT NULL DEFAULT 'PENDING', -- PENDING, SUCCEEDED, FAILED, REFUNDED
    idempotency_key  VARCHAR NOT NULL UNIQUE, -- prevents double-charging on client retries
    created_at       TIMESTAMP NOT NULL DEFAULT now()
);
