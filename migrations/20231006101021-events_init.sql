
-- +migrate Up

CREATE TABLE events
(
    `id`           char(36)                                       NOT NULL,
    `name`         varchar(100)                                   NOT NULL,
    `date`         datetime                                       NOT NULL,
    `description`  varchar(1000)                                          ,
    `created_at`   datetime                                       NOT NULL,
    `updated_at`   datetime                                       NOT NULL,
    PRIMARY KEY (`id`)
);


CREATE TABLE event_classifiers
(
    `id`           char(36)                                             NOT NULL,
    `type`         enum('language', 'video_quality', 'audio_quality')   NOT NULL,
    `value`        varchar(20)                                          NOT NULL,
    `event_id`     char(36)                                             NOT NULL,
    `created_at`   datetime                                             NOT NULL,
    `updated_at`   datetime                                             NOT NULL,
    PRIMARY KEY (`id`),
    KEY `event_id` (`event_id`),
    CONSTRAINT `fk_event_classifiers_event_id`
        FOREIGN KEY (`event_id`)
            REFERENCES `events` (`id`)
            ON DELETE CASCADE
);


CREATE TABLE event_invitees
(
    `id`           char(36)                                       NOT NULL,
    `invitee`      varchar(100)                                    NOT NULL,
    `event_id`     char(36)                                       NOT NULL,
    `created_at`   datetime                                       NOT NULL,
    `updated_at`   datetime                                       NOT NULL,
    PRIMARY KEY (`id`),
    KEY `event_id` (`event_id`),
    CONSTRAINT `fk_event_invitees_event_id`
        FOREIGN KEY (`event_id`)
            REFERENCES `events` (`id`)
            ON DELETE CASCADE
);

-- +migrate Down

DROP TABLE event_classifiers;
DROP TABLE event_invitees;
DROP TABLE events;
