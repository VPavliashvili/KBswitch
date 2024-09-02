-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS switches (
    id SERIAL PRIMARY KEY,
	lifespan         INT,
    operatingForce   INT,
	activationTravel real,
	totalTravel      real,
	image            BYTEA,
	manufacturer     VARCHAR(255),
	model            VARCHAR(255),
	actuationType    VARCHAR(255),
	soundProfile     VARCHAR(255),
	triggerMethod    VARCHAR(255),
	profile          VARCHAR(255)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS switches;
-- +goose StatementEnd
