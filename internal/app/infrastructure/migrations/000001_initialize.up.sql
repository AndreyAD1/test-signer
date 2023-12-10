BEGIN;

CREATE TABLE signatures(
    id uuid PRIMARY KEY,
    request_id varchar CONSTRAINT "request" NOT NULL CHECK (request_id <> ''),
    user_id varchar CONSTRAINT "user" NOT NULL CHECK (user_id <> ''),
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT request_user_id UNIQUE (request_id, user_id)
);

CREATE TABLE test_details(
    id SERIAL PRIMARY KEY,
    signature_id uuid,
    question varchar,
    answer varchar
);

ALTER TABLE test_details ADD CONSTRAINT "signature_id"
FOREIGN KEY (signature_id) REFERENCES signatures (id);

COMMIT;