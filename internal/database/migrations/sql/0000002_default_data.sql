-- Default Profiles
INSERT INTO profiles (
    name, container, extension, pass_thru_common_metadata, 
    flipping, rotation, cropping, limit_value, anamorphic, fill, 
    color, detelecine, interlace_detection, deinterlace, 
    deinterlace_preset, deblock, deblock_tune, denoise, 
    denoise_preset, denoise_tune, chroma_smooth, chroma_smooth_tune,
    sharpen, sharpen_preset, sharpen_tune, colorspace, grayscale, codec, encoder,
    framerate, framerate_type, quality_type, constant_quality,
    average_bitrate, multipass_encoding, preset, tune, profile,
    level, fast_decode, map_untagged_audio_tracks, map_untagged_subtitle_tracks
) VALUES 
-- Any Profile
('Any', 'matroska', 'mkv', 1, 0, 0, 'off', 'none', 'off', 'none', 'black',
'off', 'off', 'off', 'default', 'off', 'medium', 'off', 'light', 'none', 'off',
'none', 'off', 'medium', 'none', 'off', 0, 'Any', '', 'same as source', 'peak Framerate',
'constant quality', 22, 15000, 0, 'medium', 'none', 'auto', 'auto', 0, 1, 1),

-- h264 Profile
('h264', 'matroska', 'mkv', 1, 0, 0, 'off', 'none', 'off', 'none', 'black',
'off', 'off', 'off', 'default', 'off', 'medium', 'off', 'light', 'none', 'off',
'none', 'off', 'medium', 'none', 'off', 0, 'h264', 'libx264', 'same as source', 'peak Framerate',
'constant quality', 22, 15000, 0, 'medium', 'none', 'auto', 'auto', 0, 1, 1),

-- hevc Profile
('hevc', 'matroska', 'mkv', 1, 0, 0, 'off', 'none', 'off', 'none', 'black',
'off', 'off', 'off', 'default', 'off', 'medium', 'off', 'light', 'none', 'off',
'none', 'off', 'medium', 'none', 'off', 0, 'hevc', 'libx265', 'same as source', 'peak Framerate',
'constant quality', 22, 15000, 0, 'medium', 'none', 'auto', 'auto', 0, 1, 1),

-- mpeg4 Profile
('mpeg4', 'matroska', 'mkv', 1, 0, 0, 'off', 'none', 'off', 'none', 'black',
'off', 'off', 'off', 'default', 'off', 'medium', 'off', 'light', 'none', 'off',
'none', 'off', 'medium', 'none', 'off', 0, 'mpeg4', 'mpeg4', 'same as source', 'peak Framerate',
'constant quality', 22, 15000, 0, '', '', '', '', 0, 1, 1),

-- vp8 Profile
('vp8', 'matroska', 'mkv', 1, 0, 0, 'off', 'none', 'off', 'none', 'black',
'off', 'off', 'off', 'default', 'off', 'medium', 'off', 'light', 'none', 'off',
'none', 'off', 'medium', 'none', 'off', 0, 'vp8', 'libvpx', 'same as source', 'peak Framerate',
'constant quality', 22, 15000, 0, 'medium', 'none', 'auto', 'auto', 0, 1, 1),

-- vp9 Profile
('vp9', 'matroska', 'mkv', 1, 0, 0, 'off', 'none', 'off', 'none', 'black',
'off', 'off', 'off', 'default', 'off', 'medium', 'off', 'light', 'none', 'off',
'none', 'off', 'medium', 'none', 'off', 0, 'vp9', 'libvpx-vp9', 'same as source', 'peak Framerate',
'constant quality', 22, 15000, 0, 'medium', 'none', 'auto', 'auto', 0, 1, 1),

-- av1 Profile
('av1', 'matroska', 'mkv', 1, 0, 0, 'off', 'none', 'off', 'none', 'black',
'off', 'off', 'off', 'default', 'off', 'medium', 'off', 'light', 'none', 'off',
'none', 'off', 'medium', 'none', 'off', 0, 'av1', 'libaom-av1', 'same as source', 'peak Framerate',
'constant quality', 22, 15000, 0, '7', 'none', 'auto', 'auto', 0, 1, 1);

-- Default Settings
INSERT INTO settings (id, value) VALUES
('theme', 'auto'),
('defaultProfile', '1'),
('queueStatus', 'active'),
('queueStartupState', 'previous'),
('logLevel', 'info'),
('mediaView', 'posters'),
('mediaSort', 'title'),
('massEditorSort', 'title'),
('massEditorSortDirection', 'ascending'),
('massEditorFilter', 'all'),
('mediaSortDirection', 'ascending'),
('mediaFilter', 'all'),
('TMDB', 'ZXlKaGJHY2lPaUpJVXpJMU5pSjkuZXlKaGRXUWlPaUprT1RCalpqQmhaREEyT0dJd01XVXpNVFkxTWpjNVltWXpPRE0xWmpRNU9TSXNJbk4xWWlJNklqWTFOR0UxWVRReE5qZGlOakV6TURFeFpqUXdaV0ZpWVNJc0luTmpiM0JsY3lJNld5SmhjR2xmY21WaFpDSmRMQ0oyWlhKemFXOXVJam94ZlEuNU1LVjViaXV0RmZvQkRuMk14aFMxQU1wbV9DTmE4QTh4WE5XTkFKUVNnTQ=='),
('mediaPosterPosterSize', 'medium'),
('mediaPosterDetailedProgressBar', 'false'),
('mediaPosterShowTitle', 'true'),
('mediaPosterShowMonitored', 'true'),
('mediaPosterShowProfile', 'true'),
('mediaTableShowNetwork', 'false'),
('mediaTableShowProfile', 'true'),
('mediaTableShowSeasons', 'true'),
('mediaTableShowEpisodes', 'true'),
('mediaTableShowEpisodeCount', 'false'),
('mediaTableShowYear', 'true'),
('mediaTableShowType', 'true'),
('mediaTableShowSizeOnDisk', 'true'),
('mediaTableShowSizeSaved', 'true'),
('mediaTableShowGenre', 'false'),
('mediaOverviewPosterSize', 'medium'),
('mediaOverviewDetailedProgressBar', 'false'),
('mediaOverviewShowMonitored', 'true'),
('mediaOverviewShowNetwork', 'true'),
('mediaOverviewShowProfile', 'true'),
('mediaOverviewShowSeasonCount', 'true'),
('mediaOverviewShowPath', 'false'),
('mediaOverviewShowSizeOnDisk', 'true'),
('queueFilter', 'all'),
('queuePageSize', '12'),
('historyFilter', 'all'),
('historyPageSize', '15'),
('eventsFilter', 'all'),
('eventsPageSize', '15'),
('port', '7889');

-- Default User with generated secret
WITH generated_secret AS (
    SELECT lower(hex(randomblob(32))) as secret
)
INSERT INTO users (username, password, secret, created_at, updated_at)
SELECT '', '', secret, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
FROM generated_secret;