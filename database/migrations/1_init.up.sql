CREATE TABLE IF NOT EXISTS users (
                                     id INTEGER PRIMARY KEY,
                                     name TEXT,
                                     balance INTEGER DEFAULT 100,
                                     numLama INTEGER DEFAULT 100,
                                     numGPT INTEGER DEFAULT 100,
                                     apiKey TEXT,
                                     modelID TEXT
);