CREATE TABLE player_journeys (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    player_profile_id UUID NOT NULL REFERENCES player_profiles(id) ON DELETE CASCADE,
    journey_id UUID NOT NULL REFERENCES journeys(id) ON DELETE CASCADE,
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    completed_at TIMESTAMP WITH TIME ZONE,
    status VARCHAR(50) DEFAULT 'active',

    UNIQUE(player_profile_id, journey_id)
);

CREATE INDEX idx_player_journeys_player_profile_id ON player_journeys(player_profile_id);
CREATE INDEX idx_player_journeys_journey_id ON player_journeys(journey_id);
CREATE INDEX idx_player_journeys_status ON player_journeys(status);