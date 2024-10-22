CREATE TABLE files
(
    hash        TEXT                                NOT NULL PRIMARY KEY,
    filename    TEXT                                NOT NULL,
    extension   TEXT                                NOT NULL,
    size        INTEGER                             NOT NULL,
    password    TEXT      DEFAULT NULL,

    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    expires_at  TIMESTAMP DEFAULT NULL
);

CREATE TABLE pastes
(
    hash        TEXT                                NOT NULL PRIMARY KEY,
    filename    TEXT                                NOT NULL,
    size        INTEGER                             NOT NULL,
    language    TEXT      DEFAULT NULL,
    password    TEXT      DEFAULT NULL,

    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    expires_at  TIMESTAMP DEFAULT NULL
);

CREATE TABLE tokens
(
    token      TEXT                                NOT NULL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);