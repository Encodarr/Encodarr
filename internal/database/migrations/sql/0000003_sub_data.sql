-- Insert audio languages for each profile
INSERT INTO profile_audio_languages (profile_id, language) 
SELECT 1, 'all'
UNION ALL SELECT 2, 'all'
UNION ALL SELECT 3, 'all'
UNION ALL SELECT 4, 'all'
UNION ALL SELECT 5, 'all'
UNION ALL SELECT 6, 'all'
UNION ALL SELECT 7, 'all';

-- Insert codecs for each profile
INSERT INTO profile_codecs (profile_id, codec_id)
SELECT 1, 'Any'
UNION ALL SELECT 2, 'Any'
UNION ALL SELECT 3, 'Any'
UNION ALL SELECT 4, 'Any'
UNION ALL SELECT 5, 'Any'
UNION ALL SELECT 6, 'Any'
UNION ALL SELECT 7, 'Any';

-- Insert subtitle languages for each profile
INSERT INTO profile_subtitle_languages (profile_id, language)
SELECT 1, 'all'
UNION ALL SELECT 2, 'all'
UNION ALL SELECT 3, 'all'
UNION ALL SELECT 4, 'all'
UNION ALL SELECT 5, 'all'
UNION ALL SELECT 6, 'all'
UNION ALL SELECT 7, 'all';