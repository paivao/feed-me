CREATE TABLE feeds (
    id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(128) NOT NULL UNIQUE,
    description VARCHAR(255),
    is_public TINYINT(1) NOT NULL,
    type ENUM('ip', 'domain', 'url') NOT NULL
);

CREATE TABLE ip_entries (
    id BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    value VARBINARY(17) NOT NULL,
    enabled TINYINT(1) NOT NULL DEFAULT 1,
    description VARCHAR(255),
    valid_until DATETIME,
    feed_id INT NOT NULL,
    UNIQUE KEY(value, feed_id),
    FOREIGN KEY (feed_id) REFERENCES feeds(id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE domain_entries (
    id BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    value VARCHAR(255) NOT NULL UNIQUE,
    enabled TINYINT(1) NOT NULL DEFAULT 1,
    description VARCHAR(255),
    valid_until DATETIME,
    feed_id INT NOT NULL,
    UNIQUE KEY(value, feed_id),
    FOREIGN KEY (feed_id) REFERENCES feeds(id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE url_entries (
    id BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    value VARCHAR(8192) NOT NULL UNIQUE,
    enabled TINYINT(1) NOT NULL DEFAULT 1,
    description VARCHAR(255),
    valid_until DATETIME,
    feed_id INT NOT NULL,
    UNIQUE KEY(value, feed_id),
    FOREIGN KEY (feed_id) REFERENCES feeds(id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE users (
    id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(128) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL
);

CREATE TABLE permissions (
    id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE `groups` (
    id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE _user_permission_ (
    user_id INT NOT NULL,
    permission_id INT NOT NULL,
    PRIMARY KEY (user_id,permission_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions(id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE _group_permission_ (
    group_id INT NOT NULL,
    permission_id INT NOT NULL,
    PRIMARY KEY (group_id,permission_id),
    FOREIGN KEY (group_id) REFERENCES `groups`(id) ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions(id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE _user_group_ (
    user_id INT NOT NULL,
    group_id INT NOT NULL,
    PRIMARY KEY (user_id,group_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY (group_id) REFERENCES `groups`(id) ON UPDATE CASCADE ON DELETE CASCADE
);

-- admin:gofeed
INSERT INTO users (name, email, password_hash) VALUES ("admin", "admin@feed.me", "$2b$05$gF2CW3YsRxtlc9o1msB7uOwwmRvd14/AKrQPJ3NZXBf/LcZUbvJam")