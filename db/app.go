package db

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type GetSurahAndAyah struct {
	NameSurah string
	Verses    struct {
		Number        int
		Text          string
		TranslationEn string
		TranslationID string
	}
}

type GetAudioSurah struct {
	NameSurahIND string
	NameSurahENG string
	NameSurahAR  string
	NumberOfAyah int
	Place        string
	Type         string
	NameReciter  string
	Audio        string
}

func Connect() *sql.DB {
	db, err := sql.Open("mysql", os.Getenv("DSN"))
	if err != nil {
		panic(err.Error())
	}
	return db
}

func FindSurahAndAyahByNumber(numberSurah, numberAyah int64) (GetSurahAndAyah, error) {
	db := Connect()
	defer db.Close()

	var surahAndAyah GetSurahAndAyah

	err := db.QueryRow("SELECT surahs.name, surah_verses.number, surah_verses.text, surah_verses.translation_en, surah_verses.translation_id FROM surahs JOIN surah_verses ON surah_verses.surah_id = surahs.id WHERE surahs.number_of_surah = ? AND surah_verses.number = ?", numberSurah, numberAyah).Scan(&surahAndAyah.NameSurah, &surahAndAyah.Verses.Number, &surahAndAyah.Verses.Text, &surahAndAyah.Verses.TranslationEn, &surahAndAyah.Verses.TranslationID)
	if err != nil {
		return surahAndAyah, err
	}

	return surahAndAyah, nil
}

func FindAudioSurahNumber(numberSurah int64) (GetAudioSurah, error) {
	db := Connect()
	defer db.Close()

	var audio GetAudioSurah

	err := db.QueryRow("SELECT surahs.name, surah_translations.en, surah_translations.ind, surahs.number_of_ayah, surahs.place, surahs.type, surah_recitations.name, surah_recitations.audio_url FROM surahs JOIN surah_translations ON surah_translations.surah_id = surahs.id JOIN surah_recitations ON surah_recitations.surah_id = surahs.id WHERE surahs.number_of_surah = ?", numberSurah).
		Scan(&audio.NameSurahAR,
			&audio.NameSurahENG,
			&audio.NameSurahIND,
			&audio.NumberOfAyah,
			&audio.Place,
			&audio.Type,
			&audio.NameReciter,
			&audio.Audio)
	if err != nil {
		return audio, err
	}

	return audio, nil
}
