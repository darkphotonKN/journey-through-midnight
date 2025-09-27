CREATE TABLE player_profiles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    username VARCHAR(50) UNIQUE NOT NULL,
    reputation INTEGER CHECK (reputation >= 1 AND reputation <= 5) DEFAULT 3,
    total_playtime INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_player_profiles_user_id ON player_profiles(user_id);
CREATE UNIQUE INDEX idx_player_profiles_username ON player_profiles(username);