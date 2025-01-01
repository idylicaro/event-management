ALTER TABLE events ADD COLUMN user_id INT NOT NULL;
ALTER TABLE events ADD CONSTRAINT events_user_id_fkey FOREIGN KEY (user_id) REFERENCES users (id);