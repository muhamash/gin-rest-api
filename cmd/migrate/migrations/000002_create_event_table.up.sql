CREATE TABLE IF NOT EXISTS events (
    id SERIAL PRIMARY KEY,
    owner_name VARCHAR(50) NOT NULL UNIQUE,
    owner_id INT NOT NULL, 
    description TEXT NOT NULL,
    date DATE NOT NULL,
    location TEXT NOT NULL, 
    FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE
);