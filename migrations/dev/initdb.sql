
-- +migrate Up
CREATE TYPE client_status as ENUM('developing','published','suspended','deleted');
CREATE TYPE client_type as ENUM('confidential', 'public');
CREATE TABLE IF NOT EXISTS clients (
    client_id    VARCHAR(255)   NOT NULL UNIQUE,
    client_secret_hash VARCHAR(255) UNIQUE,
    salt               VARCHAR(255) UNIQUE,
    status client_status default 'developing' NOT NULL,
    type client_type NOT NULL,
    revision_ver  INTEGER default 1 NOT NULL UNIQUE,
    PRIMARY KEY  (client_id)
);

CREATE INDEX index_client_on_client_id ON clients(client_id);
CREATE INDEX index_client_on_revision_ver ON clients(revision_ver);

CREATE TABLE IF NOT EXISTS authz (
    authorization_id    INTEGER   NOT NULL UNIQUE,
    client_id    VARCHAR(255)   NOT NULL UNIQUE,
    user_id    VARCHAR(255)   NOT NULL,
    authz_code    VARCHAR(255)  NOT NULL,
    revision_ver  INTEGER default 1 NOT NULL,
    refresh_token VARCHAR(255),
    issued_at DATE  NOT NULL,
    expired_at DATE    NOT NULL,
    PRIMARY KEY (authorization_id, client_id, user_id),
    FOREIGN KEY(client_id) references clients(client_id) on update cascade,
    FOREIGN KEY(revision_ver) references clients(revision_ver)  on update cascade on delete restrict
--      FOREIGN KEY(user_id) references client(user_id) on update cascade on delete set null,
);


-- +migrate Down
-- DROP TABLE authz;
-- drop INDEX index_client_on_client_id;
-- drop INDEX index_client_on_revision_ver;
-- DROP TABLE clients;