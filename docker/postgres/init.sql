CREATE USER docker;
CREATE DATABASE amigo OWNER docker;
\connect amigo;
CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    value text,
    ip_address inet,
    date_created TIMESTAMP DEFAULT NOW()
);
GRANT ALL PRIVILEGES ON TABLE messages TO docker;
GRANT USAGE, SELECT ON SEQUENCE messages_id_seq TO docker;