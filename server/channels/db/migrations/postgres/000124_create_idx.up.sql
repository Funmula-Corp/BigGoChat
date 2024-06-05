-- morph:nontransactional
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_channelblockusers_channelid ON channelblockusers(channelid);
-- I cannot merge this statement right into the migration of create table
