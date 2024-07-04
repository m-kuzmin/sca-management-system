CREATE TABLE missions (
    id UUID PRIMARY KEY,
    assigned_cat UUID REFERENCES cats(id),
    complete BOOLEAN NOT NULL DEFAULT false
);

CREATE TABLE targets (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    country CHAR(3) NOT NULL,
    notes TEXT,
    complete BOOLEAN NOT NULL DEFAULT false
);

CREATE TABLE mission_targets (
    mission_id UUID REFERENCES missions(id),
    target_id UUID REFERENCES targets(id),
    PRIMARY KEY (mission_id, target_id)
);
