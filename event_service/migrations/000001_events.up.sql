CREATE TABLE IF NOT EXISTS public.events
(
    id serial PRIMARY KEY,
    creator_id INT NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    time TIMESTAMP NOT NULL,
    place TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS public.topics
(
    id serial PRIMARY KEY,
    event_id INT NOT NULL,
    topic TEXT,
    CONSTRAINT fk_event FOREIGN KEY (event_id) REFERENCES public.events(id)
);

CREATE TABLE IF NOT EXISTS public.participants
(
    user_id INT PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS public.registrations
(
    event_id INT,
    participant_id INT,
    ready_to_chat BOOL,
    PRIMARY KEY (event_id, participant_id),
    CONSTRAINT fk_event FOREIGN KEY (event_id) REFERENCES public.events(id),
    CONSTRAINT fk_participant FOREIGN KEY (participant_id) REFERENCES public.participants(user_id)
);

CREATE TABLE IF NOT EXISTS public.interests
(
    id serial PRIMARY KEY,
    participant_id INT NOT NULL,
    interest TEXT,
    CONSTRAINT fk_participant FOREIGN KEY (participant_id) REFERENCES public.participants(user_id)
);
