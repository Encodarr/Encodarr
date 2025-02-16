-- Core tables
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL,
    password TEXT NOT NULL,
    secret TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS profiles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    container TEXT,
    extension TEXT,
    pass_thru_common_metadata BOOLEAN DEFAULT FALSE,
    flipping BOOLEAN DEFAULT FALSE,
    rotation INTEGER,
    cropping TEXT,
    limit_value TEXT,
    anamorphic TEXT,
    fill TEXT,
    color TEXT,
    detelecine TEXT,
    interlace_detection TEXT,
    deinterlace TEXT,
    deinterlace_preset TEXT,
    deblock TEXT,
    deblock_tune TEXT,
    denoise TEXT,
    denoise_preset TEXT,
    denoise_tune TEXT,
    chroma_smooth TEXT,
    chroma_smooth_tune TEXT,
    sharpen TEXT,
    sharpen_preset TEXT, 
    sharpen_tune TEXT,
    colorspace TEXT,
    grayscale BOOLEAN DEFAULT FALSE,
    codec TEXT,
    encoder TEXT,
    framerate TEXT,
    framerate_type TEXT,
    quality_type TEXT,
    constant_quality INTEGER,
    average_bitrate INTEGER,
    multipass_encoding BOOLEAN DEFAULT FALSE,
    preset TEXT,
    tune TEXT,
    profile TEXT,
    level TEXT,
    fast_decode BOOLEAN DEFAULT FALSE,
    map_untagged_audio_tracks BOOLEAN DEFAULT FALSE,
    map_untagged_subtitle_tracks BOOLEAN DEFAULT FALSE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS profile_audio_languages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    profile_id INTEGER NOT NULL,
    language TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS profile_subtitle_languages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    profile_id INTEGER NOT NULL,
    language TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS profile_codecs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    profile_id INTEGER NOT NULL,
    codec_id TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS series (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    release_date TEXT,
    genre TEXT,
    status TEXT,
    last_air_date TEXT,
    networks TEXT,
    overview TEXT,
    profile_id INTEGER,
    monitored BOOLEAN DEFAULT FALSE,
    episode_count INTEGER DEFAULT 0,
    size INTEGER DEFAULT 0,
    seasons_count INTEGER DEFAULT 0,
    space_saved INTEGER DEFAULT 0,
    missing_episodes INTEGER DEFAULT 0,
    runtime INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (profile_id) REFERENCES profiles(id)
);

CREATE TABLE IF NOT EXISTS seasons (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    season_number INTEGER NOT NULL,
    episode_count INTEGER DEFAULT 0,
    size INTEGER DEFAULT 0,
    series_id TEXT NOT NULL,
    space_saved INTEGER DEFAULT 0,
    missing_episodes INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (series_id) REFERENCES series(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS episodes (
    id TEXT PRIMARY KEY,
    series_id TEXT NOT NULL,
    season_id TEXT NOT NULL,
    episode_number INTEGER NOT NULL,
    season_name TEXT,
    season_number INTEGER,
    filename TEXT,
    episode_name TEXT,
    video_codec TEXT,
    air_date TEXT,
    size INTEGER DEFAULT 0,
    space_saved INTEGER DEFAULT 0,
    original_size INTEGER DEFAULT 0,
    path TEXT,
    missing BOOLEAN DEFAULT FALSE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (series_id) REFERENCES series(id) ON DELETE CASCADE,
    FOREIGN KEY (season_id) REFERENCES seasons(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS movies (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    release_date TEXT,
    genre TEXT,
    status TEXT,
    filename TEXT,
    video_codec TEXT,
    overview TEXT,
    size INTEGER DEFAULT 0,
    space_saved INTEGER DEFAULT 0,
    profile_id INTEGER,
    monitored BOOLEAN DEFAULT FALSE,
    missing BOOLEAN DEFAULT FALSE,
    studio TEXT,
    original_size INTEGER DEFAULT 0,
    path TEXT,
    runtime INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (profile_id) REFERENCES profiles(id)
);

CREATE TABLE IF NOT EXISTS history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    media_id TEXT NOT NULL,
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    season_number INTEGER,
    episode_number INTEGER,
    profile_id INTEGER,
    prev_codec TEXT,
    new_codec TEXT,
    prev_size INTEGER,
    new_size INTEGER,
    date TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (profile_id) REFERENCES profiles(id)
);

CREATE TABLE IF NOT EXISTS events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    timestamp TEXT NOT NULL,
    level TEXT NOT NULL,
    service TEXT NOT NULL,
    message TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS settings (
    id TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS systems (
    id TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS migrations (
    version INTEGER PRIMARY KEY,
    applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_series_profile ON series(profile_id);
CREATE INDEX IF NOT EXISTS idx_movies_profile ON movies(profile_id);
CREATE INDEX IF NOT EXISTS idx_episodes_series ON episodes(series_id);
CREATE INDEX IF NOT EXISTS idx_episodes_season ON episodes(season_id);
CREATE INDEX IF NOT EXISTS idx_seasons_series ON seasons(series_id);
CREATE INDEX IF NOT EXISTS idx_history_profile ON history(profile_id);
CREATE INDEX IF NOT EXISTS idx_profile_audio_lang ON profile_audio_languages(profile_id);
CREATE INDEX IF NOT EXISTS idx_profile_subtitle_lang ON profile_subtitle_languages(profile_id);
CREATE INDEX IF NOT EXISTS idx_profile_codecs ON profile_codecs(profile_id);