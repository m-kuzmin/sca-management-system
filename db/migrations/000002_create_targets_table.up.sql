BEGIN;

CREATE TABLE missions (
    id UUID PRIMARY KEY,
    assigned_cat UUID REFERENCES cats(id),
    complete BOOLEAN NOT NULL DEFAULT false
);

CREATE TABLE targets (
    id UUID PRIMARY KEY,
    mission_id UUID NOT NULL REFERENCES missions(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    country CHAR(3) NOT NULL,
    notes TEXT,
    complete BOOLEAN NOT NULL DEFAULT false
);

COMMIT;
