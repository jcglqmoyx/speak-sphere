package conf

import "speak-sphere/pkg/server/model"

type ServerConfig struct {
	Port int `yaml:"port"`
}

type JWTConfig struct {
	Secret string `yaml:"secret"`
}

type LogConfig struct {
	Path       string `yaml:"path"`
	MaxAge     int    `yaml:"max_age"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	Compress   bool   `yaml:"compress"`
	LocalTime  bool   `yaml:"local_time"`
}

type DBConfig struct {
	Path            string `yaml:"path"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime"`
}

type VocabularySetConfig struct {
	MaxTitleLength    int    `yaml:"max_title_length"`
	MaxCategoryLength int    `yaml:"max_category_length"`
	MaxFileSize       int64  `yaml:"max_file_size"`
	UploadPath        string `yaml:"upload_path"`
}

type VocabularyConfig struct {
	MaxVocabularyLength int `yaml:"max_vocabulary_length"`
	MaxMeaningLength    int `yaml:"max_meaning_length"`
	MaxNoteLength       int `yaml:"max_note_length"`
	DefaultDateToReview int `yaml:"default_date_to_review"`
}

type DictionaryConfig struct {
	MaxTitleLength  int                 `yaml:"max_title_length"`
	MaxPrefixLength int                 `yaml:"max_prefix_length"`
	MaxSuffixLength int                 `yaml:"max_suffix_length"`
	Dictionaries    []*model.Dictionary `yaml:"dictionaries"`
}

type Config struct {
	Mode          string               `yaml:"mode"`
	Server        *ServerConfig        `yaml:"server"`
	JWT           *JWTConfig           `yaml:"jwt"`
	Log           *LogConfig           `yaml:"log"`
	Sqlite        *DBConfig            `yaml:"sqlite"`
	VocabularySet *VocabularySetConfig `yaml:"vocabulary_set"`
	Vocabulary    *VocabularyConfig    `yaml:"vocabulary"`
	Dictionary    *DictionaryConfig    `yaml:"dictionary"`
}

var Cfg Config
