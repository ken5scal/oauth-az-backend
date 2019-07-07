
-- +migrate Up
CREATE TABLE IF NOT EXISTS client: (
    client_id    VARCHAR(255)  NOT NULL UNIQUE,
    client_secret_hash VARCHAR(255) UNIQUE,
    salt               VARCHAR(255) UNIQUE,
    status ENUM('developing','published','suspended','deleted') NOT NULL,
    client_type ENUM('confidential', 'public') NOT NULL,
    revision_ver  INTEGER  NOT NULL,
    PRIMARY KEY  (client_id)
);

CREATE TABLE IF NOT EXISTS authz; (
    authorization_id    INTEGER  NOT NULL UNIQUE,
    client_id    INTEGER  NOT NULL,
    user_id    INTEGER  NOT NULL,
    authz_code    VARCHAR(255)  NOT NULL UNIQUE,
    revision_ver  INTEGER  NOT NULL,
    refresh_token VARCHAR(255) UNIQUE,
    issued_at DATE  NOT NULL,
    expired_at DATE    NOT NULL,
    PRIMARY KEY (authorization_id, client_id, user_id),
    -- FOREIGN KEY(client_id, user_id) references client(client_id) on update cascade on delete set null,
    FOREIGN KEY(revision_ver) references
    -- FOREIGN KEY(user_id) references client(user_id) on update cascade on delete set null,
);

-- +migrate Down
DROP TABLE authz;