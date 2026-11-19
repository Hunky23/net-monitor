CREATE TABLE IF NOT EXISTS processes (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    pid          INTEGER NOT NULL,
    process_name TEXT    NOT NULL,
    display_name TEXT,
    path         TEXT,
    status       TEXT    NOT NULL,
    type         TEXT    NOT NULL,
    icon         TEXT,
    created_at   TEXT    NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_processes_pid        ON processes (pid);
CREATE INDEX IF NOT EXISTS idx_process_name         ON processes (process_name);
CREATE INDEX IF NOT EXISTS idx_process_display_name ON processes (display_name);
CREATE INDEX IF NOT EXISTS idx_processes_status     ON processes (status);

CREATE TABLE IF NOT EXISTS connections (
    id                   INTEGER PRIMARY KEY AUTOINCREMENT,
    process_id           INTEGER NOT NULL,
    source_ip            TEXT    NOT NULL,
    source_port          INTEGER NOT NULL,
    destination_ip       TEXT    NOT NULL,
    destination_port     INTEGER NOT NULL,
    domain               TEXT,
    bytes_received       INTEGER NOT NULL DEFAULT 0,
    bytes_sent           INTEGER NOT NULL DEFAULT 0,
    transport_protocol   TEXT    NOT NULL,
    application_protocol TEXT,
    status               TEXT    NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_connections_process_id           ON connections (process_id);
CREATE INDEX IF NOT EXISTS idx_connections_source_ip            ON connections (source_ip);
CREATE INDEX IF NOT EXISTS idx_connections_source_port          ON connections (source_port);
CREATE INDEX IF NOT EXISTS idx_connections_destination_ip       ON connections (destination_ip);
CREATE INDEX IF NOT EXISTS idx_connections_destination_port     ON connections (destination_port);
CREATE INDEX IF NOT EXISTS idx_connections_transport_protocol   ON connections (transport_protocol);
CREATE INDEX IF NOT EXISTS idx_connections_application_protocol ON connections (application_protocol);
CREATE INDEX IF NOT EXISTS idx_connections_status               ON connections (status);

CREATE TABLE IF NOT EXISTS dns (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    process_id INTEGER NOT NULL,
    domain     TEXT    NOT NULL,
    type       TEXT    NOT NULL,
    value      TEXT    NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_dns_process_id ON dns (process_id);
CREATE INDEX IF NOT EXISTS idx_dns_domain     ON dns (domain);