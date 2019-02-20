-- Create Tables
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name text NOT NULL,
    created timestamp with time zone DEFAULT timezone('utc'::text, now()),
    email text NOT NULL
);

CREATE TABLE IF NOT EXISTS artists (
    id SERIAL PRIMARY KEY,
    name text NOT NULL,
    created timestamp with time zone DEFAULT timezone('utc'::text, now())
);

CREATE TABLE IF NOT EXISTS songs (
    id SERIAL PRIMARY KEY,
    name text NOT NULL,
    created timestamp with time zone DEFAULT timezone('utc'::text, now()),
    artist_id integer NOT NULL REFERENCES artists(id) ON DELETE CASCADE,
    duration bigint NOT NULL
);

CREATE TABLE IF NOT EXISTS likes_songs (
    id SERIAL PRIMARY KEY,
    created timestamp with time zone DEFAULT timezone('utc'::text, now()),
    user_id integer NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    song_id integer NOT NULL REFERENCES songs(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS likes_artists (
    id SERIAL PRIMARY KEY,
    created timestamp with time zone DEFAULT timezone('utc'::text, now()),
    user_id integer NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    artist_id integer NOT NULL REFERENCES artists(id) ON DELETE CASCADE
);

-- Seed data
INSERT INTO users("name", "email") VALUES('Scott', 'scott@test.net');
INSERT INTO users("name", "email") VALUES('Bob', 'bob@test.net');
INSERT INTO users("name", "email") VALUES('Jack', 'jack@test.net');

INSERT INTO artists("name") VALUES('Calvin Harris');
INSERT INTO artists("name") VALUES('Metallica');
INSERT INTO artists("name") VALUES('Sublime');
INSERT INTO artists("name") VALUES('Sugar Ray');
INSERT INTO artists("name") VALUES('Deadmau5');
INSERT INTO artists("name") VALUES('Justin Bieber');
INSERT INTO artists("name") VALUES('Beyonce');
INSERT INTO artists("name") VALUES('Fall Out Boy');
INSERT INTO artists("name") VALUES('Radiohead');
INSERT INTO artists("name") VALUES('Beck');

INSERT INTO songs("name", "artist_id", "duration") VALUES('This is what you came For', 1, 240000000000);
INSERT INTO songs("name", "artist_id", "duration") VALUES('My Way', 1, 270000000000);
INSERT INTO songs("name", "artist_id", "duration") VALUES('How Deep Is Your Love', 1, 300000000000);
INSERT INTO songs("name", "artist_id", "duration") VALUES('Summer', 1, 220000000000);
INSERT INTO songs("name", "artist_id", "duration") VALUES('The Weekend - Funk Wav Remix', 1, 230000000000);

INSERT INTO songs("name", "artist_id", "duration") VALUES('Nothing Else Matters', 2, 240000000000);
INSERT INTO songs("name", "artist_id", "duration") VALUES('Enter Sandman', 2, 270000000000);
INSERT INTO songs("name", "artist_id", "duration") VALUES('Master of Puppets (Remastered)', 2, 300000000000);
INSERT INTO songs("name", "artist_id", "duration") VALUES('The Unforgiven', 2, 220000000000);
INSERT INTO songs("name", "artist_id", "duration") VALUES('One (Remastered)', 2, 230000000000);

INSERT INTO songs("name", "artist_id", "duration") VALUES('Santeria', 3, 240000000000);
INSERT INTO songs("name", "artist_id", "duration") VALUES('What I Got', 3, 270000000000);
INSERT INTO songs("name", "artist_id", "duration") VALUES('Wrong Way', 3, 300000000000);
INSERT INTO songs("name", "artist_id", "duration") VALUES('Caress Me Down', 3, 220000000000);
INSERT INTO songs("name", "artist_id", "duration") VALUES('Smoke Two Joints', 3, 230000000000);

INSERT INTO songs("name", "artist_id", "duration") VALUES('Every Morning', 4, 240000000000);
INSERT INTO songs("name", "artist_id", "duration") VALUES('Fly', 4, 270000000000);
INSERT INTO songs("name", "artist_id", "duration") VALUES('Someday', 4, 300000000000);
INSERT INTO songs("name", "artist_id", "duration") VALUES(E'When It\'s Over - David Kahne Main', 4, 220000000000);
INSERT INTO songs("name", "artist_id", "duration") VALUES('Into Yesterday', 4, 230000000000);
--' bad syntax highlighting
INSERT INTO songs("name", "artist_id", "duration") VALUES('Monophobia', 5, 240000000000);
INSERT INTO songs("name", "artist_id", "duration") VALUES('Drama Free', 5, 270000000000);
INSERT INTO songs("name", "artist_id", "duration") VALUES('Polyphobia', 5, 300000000000);
INSERT INTO songs("name", "artist_id", "duration") VALUES('4ware', 5, 220000000000);
INSERT INTO songs("name", "artist_id", "duration") VALUES('Nyquist', 5, 230000000000);

INSERT INTO songs("name", "artist_id", "duration") VALUES('No Brainer', 6, 240000000000);
INSERT INTO songs("name", "artist_id", "duration") VALUES('Love Yourself', 6, 270000000000);
INSERT INTO songs("name", "artist_id", "duration") VALUES('Sorry', 6, 300000000000);
INSERT INTO songs("name", "artist_id", "duration") VALUES('Friends(with BloodPop)', 6, 220000000000);
INSERT INTO songs("name", "artist_id", "duration") VALUES('What Do You Mean?', 6, 230000000000);

INSERT INTO songs("name", "artist_id", "duration") VALUES('Halo', 7, 240000000000);
INSERT INTO songs("name", "artist_id", "duration") VALUES('Crazy In Love', 7, 270000000000);
INSERT INTO songs("name", "artist_id", "duration") VALUES('Run the World (Girls)', 7, 300000000000);
INSERT INTO songs("name", "artist_id", "duration") VALUES('Love On Top', 7, 220000000000);
INSERT INTO songs("name", "artist_id", "duration") VALUES('Irreplaceable', 7, 230000000000);

INSERT INTO songs("name", "artist_id", "duration") VALUES('Centuries', 8, 240000000000);
INSERT INTO songs("name", "artist_id", "duration") VALUES(E'Sugar, We\'re Goin Down', 8, 270000000000);
--'bad syntax highlighting
INSERT INTO songs("name", "artist_id", "duration") VALUES('Thanks fr the Mmrs', 8, 300000000000);
INSERT INTO songs("name", "artist_id", "duration") VALUES('Dance, Dance', 8, 220000000000);
INSERT INTO songs("name", "artist_id", "duration") VALUES('Immortals', 8, 230000000000);

INSERT INTO songs("name", "artist_id", "duration") VALUES('Creep', 9, 240000000000);
INSERT INTO songs("name", "artist_id", "duration") VALUES('Karma Police', 9, 270000000000);
INSERT INTO songs("name", "artist_id", "duration") VALUES('High And Dry', 9, 300000000000);
INSERT INTO songs("name", "artist_id", "duration") VALUES('Fake Plastic Trees', 9, 220000000000);
INSERT INTO songs("name", "artist_id", "duration") VALUES('No Surprises', 9, 230000000000);

INSERT INTO songs("name", "artist_id", "duration") VALUES('Loser', 10, 240000000000);
INSERT INTO songs("name", "artist_id", "duration") VALUES('Up All Night', 10, 270000000000);
INSERT INTO songs("name", "artist_id", "duration") VALUES('Tarantula', 10, 300000000000);
INSERT INTO songs("name", "artist_id", "duration") VALUES('Dreams', 10, 220000000000);
INSERT INTO songs("name", "artist_id", "duration") VALUES('Colors', 10, 230000000000);