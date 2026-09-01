	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
		username TEXT NOT NULL UNIQUE,
		rating INT NOT NULL DEFAULT 0,
		hashed_password TEXT NOT NULL,
		created_at DATETIME NOT NULL
	);
	CREATE TABLE IF NOT EXISTS friendships (
	    user_id INTEGER NOT NULL,
	    friend_id INTEGER NOT NULL,
	    status TEXT NOT NULL DEFAULT 'accepted', -- 'pending', 'accepted', 'blocked'
	    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	    PRIMARY KEY (user_id, friend_id),
	    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
	    FOREIGN KEY (friend_id) REFERENCES users(id) ON DELETE CASCADE,
	    CHECK (user_id != friend_id)
	);

	-- Create indexing for friend_id to search efficiantly.

	CREATE INDEX IF NOT EXISTS idx_friendships_friend_id ON friendships(friend_id);
