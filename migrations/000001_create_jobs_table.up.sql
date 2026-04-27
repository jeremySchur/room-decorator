CREATE TYPE job_status AS ENUM ('PENDING', 'PROCESSING', 'SUCCESS', 'FAILED');

CREATE TABLE jobs (
    id          uuid PRIMARY KEY,
    status      job_status NOT NULL,
    payload     text NOT NULL,
    created_at  timestamptz NOT NULL DEFAULT now(),
    updated_at  timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX jobs_status_idx ON jobs (status);
