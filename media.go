/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2019
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE.md
*/

package media

import (
	"strings"

	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type MetadataKey uint32
type MediaType int64
type MediaEventType uint

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

type Media interface {
	gopi.Driver

	// Open and close media files
	Open(filename string) (MediaFile, error)
	Destroy(MediaFile) error
}

type MediaItem interface {

	// Return title for the media item, based on the metadata
	// or the filename
	Title() string

	// Return type for the media item
	Type() MediaType

	// Return all keys included in the metadata
	Keys() []MetadataKey

	// Return additional metadata for the media item and a boolean
	// value which indicates if the metadata value exists
	StringForKey(MetadataKey) (string, bool)
}

type MediaFile interface {
	MediaItem

	// Return filename for the media file
	Filename() string

	// Probe the file and enumerate the streams
	Streams() []MediaStream

	// Return artwork associated with the file, both data
	// and detected format for the artwork data as a mimetype
	ArtworkData() ([]byte, string)
}

type MediaStream interface {
	// Return type for the media stream
	Type() MediaType
}

type MediaLibrary interface {
	gopi.Driver
	gopi.Publisher

	// Scan a path for media files
	AddPath(string) error

	// Query library for media items
	Query(MediaQuery) []MediaItem
}

type MediaEvent interface {
	gopi.Event

	Type() MediaEventType
	Item() MediaItem
	Path() string
	Error() error
}

type MediaQuery interface {
	// MEDIA_TYPE_ALBUM returns all albums
	// MEDIA_TYPE_MUSIC returns all songs within an album
	// MEDIA_TYPE_TVSHOW returns all TV shows
	// MEDIA_TYPE_TVSEASON returns all TV seasons (requires TVSHOW to be set)
	// MEDIA_TYPE_TVEPISODE returns all TV episodes (requires TVSHOW to be set)
	// MEDIA_TYPE_MOVIE returns all movies
	Type() MediaType
	Limit() uint
	Offset() uint
	WhereString(MetadataKey, string)
	WhereBool(MetadataKey, bool)
	WhereUint(MetadataKey, uint)
	WhereYear(MetadataKey, uint)
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	MEDIA_TYPE_FILE MediaType = (1 << iota)
	MEDIA_TYPE_AUDIO
	MEDIA_TYPE_VIDEO
	MEDIA_TYPE_IMAGE
	MEDIA_TYPE_SUBTITLE
	MEDIA_TYPE_DATA
	MEDIA_TYPE_ATTACHMENT
	MEDIA_TYPE_MUSIC
	MEDIA_TYPE_ALBUM
	MEDIA_TYPE_COMPILATION
	MEDIA_TYPE_TVSHOW
	MEDIA_TYPE_TVSEASON
	MEDIA_TYPE_TVEPISODE
	MEDIA_TYPE_AUDIOBOOK
	MEDIA_TYPE_MUSICVIDEO
	MEDIA_TYPE_MOVIE
	MEDIA_TYPE_BOOKLET
	MEDIA_TYPE_RINGTONE
	MEDIA_TYPE_ARTWORK
	MEDIA_TYPE_CAPTIONS
	MEDIA_TYPE_NONE MediaType = 0
	MEDIA_TYPE_MIN  MediaType = MEDIA_TYPE_FILE
	MEDIA_TYPE_MAX  MediaType = MEDIA_TYPE_CAPTIONS
)

var (
	// Invalid key
	METADATA_KEY_NONE = METADATA_KEY(0, 0, 0, 0)

	// File attributes
	METADATA_KEY_FILENAME  = METADATA_KEY('f', 'n', 'a', 'm') // string
	METADATA_KEY_EXTENSION = METADATA_KEY('f', 'e', 'x', 't') // string
	METADATA_KEY_FILESIZE  = METADATA_KEY('f', 's', 'i', 'z') // uint

	// Other strings
	METADATA_KEY_TITLE         = METADATA_KEY('t', 'i', 't', 'x') // string
	METADATA_KEY_TITLE_SORT    = METADATA_KEY('s', 'i', 't', 'x') // string
	METADATA_KEY_COMMENT       = METADATA_KEY('c', 'm', 't', 'x') // string
	METADATA_KEY_DESCRIPTION   = METADATA_KEY('d', 'e', 't', 'x') // string
	METADATA_KEY_SYNOPSIS      = METADATA_KEY('s', 'y', 't', 'x') // string
	METADATA_KEY_GROUPING      = METADATA_KEY('g', 'r', 't', 'x') // string
	METADATA_KEY_COPYRIGHT     = METADATA_KEY('c', 'p', 't', 'x') // string
	METADATA_KEY_LANGUAGE      = METADATA_KEY('l', 'a', 't', 'x') // string
	METADATA_KEY_VERSION_MINOR = METADATA_KEY('m', 'i', 'v', 'e') // uint
	METADATA_KEY_VERSION_MAJOR = METADATA_KEY('m', 'a', 'v', 'e') // uint
	METADATA_KEY_ACCOUNT_ID    = METADATA_KEY('u', 's', 't', 'x') // string

	// Dates
	METADATA_KEY_CREATED   = METADATA_KEY('c', 't', 'i', 'm') // iso date/time
	METADATA_KEY_MODIFIED  = METADATA_KEY('m', 't', 'i', 'm') // iso date/time
	METADATA_KEY_YEAR      = METADATA_KEY('y', 't', 'i', 'm') // iso date/time
	METADATA_KEY_PURCHASED = METADATA_KEY('p', 't', 'i', 'm') // iso date/time

	// Type strings
	METADATA_KEY_BRAND_MAJOR      = METADATA_KEY('m', 'a', 'b', 'r') // string
	METADATA_KEY_BRAND_COMPATIBLE = METADATA_KEY('m', 'i', 'b', 'r') // string
	METADATA_KEY_MEDIA_TYPE       = METADATA_KEY('t', 'y', 'p', 'e') // uint

	// Encoding strings
	METADATA_KEY_ENCODER    = METADATA_KEY('c', 'o', 't', 'x') // string
	METADATA_KEY_ENCODED_BY = METADATA_KEY('e', 'n', 't', 'x') // string

	// Track, disc
	METADATA_KEY_TRACK = METADATA_KEY('t', 'i', 'n', 't') // uint
	METADATA_KEY_DISC  = METADATA_KEY('d', 'i', 'n', 't') // uint

	// Music Item specific
	METADATA_KEY_ALBUM            = METADATA_KEY('a', 'l', 't', 'x') // string
	METADATA_KEY_ALBUM_SORT       = METADATA_KEY('s', 'l', 't', 'x') // string
	METADATA_KEY_ALBUM_ARTIST     = METADATA_KEY('a', 'a', 't', 'x') // string
	METADATA_KEY_ARTIST           = METADATA_KEY('a', 'r', 't', 'x') // string
	METADATA_KEY_ARTIST_SORT      = METADATA_KEY('s', 'r', 't', 'x') // string
	METADATA_KEY_COMPOSER         = METADATA_KEY('c', 'o', 't', 'x') // string
	METADATA_KEY_PERFORMER        = METADATA_KEY('p', 'e', 't', 'x') // string
	METADATA_KEY_PUBLISHER        = METADATA_KEY('p', 'u', 't', 'x') // string
	METADATA_KEY_GENRE            = METADATA_KEY('g', 'e', 't', 'x') // string
	METADATA_KEY_COMPILATION      = METADATA_KEY('c', 'b', 'o', 'l') // bool
	METADATA_KEY_GAPLESS_PLAYBACK = METADATA_KEY('g', 'b', 'o', 'l') // bool

	// TV Item specific
	METADATA_KEY_SHOW         = METADATA_KEY('s', 'h', 't', 'x')
	METADATA_KEY_SEASON       = METADATA_KEY('s', 'i', 'n', 't') // uint
	METADATA_KEY_EPISODE_ID   = METADATA_KEY('e', 'i', 'n', 't') // uint
	METADATA_KEY_EPISODE_SORT = METADATA_KEY('f', 'i', 'n', 't') // uint

	// Broadcasting strings
	METADATA_KEY_SERVICE_NAME     = METADATA_KEY('s', 'n', 't', 'x')
	METADATA_KEY_SERVICE_PROVIDER = METADATA_KEY('s', 'p', 't', 'x')
)

const (
	MEDIA_EVENT_FILE_ADDED MediaEventType = iota
	MEDIA_EVENT_SCAN_START
	MEDIA_EVENT_SCAN_END
	MEDIA_EVENT_ERROR
	MEDIA_EVENT_NONE MediaEventType = 0
)

////////////////////////////////////////////////////////////////////////////////
// METHODS

// METADATA_KEY returns a uint32 version of four bytes
func METADATA_KEY(a, b, c, d byte) MetadataKey {
	return MetadataKey(uint32(a)<<24 | uint32(b)<<16 | uint32(c)<<8 | uint32(d))
}

func (k MetadataKey) String() string {
	switch k {
	case METADATA_KEY_NONE:
		return "METADATA_KEY_NONE"
	case METADATA_KEY_FILENAME:
		return "METADATA_KEY_FILENAME"
	case METADATA_KEY_EXTENSION:
		return "METADATA_KEY_EXTENSION"
	case METADATA_KEY_FILESIZE:
		return "METADATA_KEY_FILESIZE"
	case METADATA_KEY_TITLE:
		return "METADATA_KEY_TITLE"
	case METADATA_KEY_TITLE_SORT:
		return "METADATA_KEY_TITLE_SORT"
	case METADATA_KEY_COMMENT:
		return "METADATA_KEY_COMMENT"
	case METADATA_KEY_DESCRIPTION:
		return "METADATA_KEY_DESCRIPTION"
	case METADATA_KEY_SYNOPSIS:
		return "METADATA_KEY_SYNOPSIS"
	case METADATA_KEY_COPYRIGHT:
		return "METADATA_KEY_COPYRIGHT"
	case METADATA_KEY_LANGUAGE:
		return "METADATA_KEY_LANGUAGE"
	case METADATA_KEY_VERSION_MINOR:
		return "METADATA_KEY_VERSION_MINOR"
	case METADATA_KEY_VERSION_MAJOR:
		return "METADATA_KEY_VERSION_MAJOR"
	case METADATA_KEY_ACCOUNT_ID:
		return "METADATA_KEY_ACCOUNT_ID"
	case METADATA_KEY_CREATED:
		return "METADATA_KEY_CREATED"
	case METADATA_KEY_MODIFIED:
		return "METADATA_KEY_MODIFIED"
	case METADATA_KEY_YEAR:
		return "METADATA_KEY_YEAR"
	case METADATA_KEY_PURCHASED:
		return "METADATA_KEY_PURCHASED"
	case METADATA_KEY_BRAND_MAJOR:
		return "METADATA_KEY_BRAND_MAJOR"
	case METADATA_KEY_BRAND_COMPATIBLE:
		return "METADATA_KEY_BRAND_COMPATIBLE"
	case METADATA_KEY_MEDIA_TYPE:
		return "METADATA_KEY_MEDIA_TYPE"
	case METADATA_KEY_ENCODER:
		return "METADATA_KEY_ENCODER"
	case METADATA_KEY_ENCODED_BY:
		return "METADATA_KEY_ENCODED_BY"
	case METADATA_KEY_TRACK:
		return "METADATA_KEY_TRACK"
	case METADATA_KEY_DISC:
		return "METADATA_KEY_DISC"
	case METADATA_KEY_ALBUM:
		return "METADATA_KEY_ALBUM"
	case METADATA_KEY_ALBUM_SORT:
		return "METADATA_KEY_ALBUM_SORT"
	case METADATA_KEY_ALBUM_ARTIST:
		return "METADATA_KEY_ALBUM_ARTIST"
	case METADATA_KEY_ARTIST:
		return "METADATA_KEY_ARTIST"
	case METADATA_KEY_ARTIST_SORT:
		return "METADATA_KEY_ARTIST_SORT"
	case METADATA_KEY_COMPOSER:
		return "METADATA_KEY_COMPOSER"
	case METADATA_KEY_PERFORMER:
		return "METADATA_KEY_PERFORMER"
	case METADATA_KEY_PUBLISHER:
		return "METADATA_KEY_PUBLISHER"
	case METADATA_KEY_GENRE:
		return "METADATA_KEY_GENRE"
	case METADATA_KEY_COMPILATION:
		return "METADATA_KEY_COMPILATION"
	case METADATA_KEY_GAPLESS_PLAYBACK:
		return "METADATA_KEY_GAPLESS_PLAYBACK"
	case METADATA_KEY_SHOW:
		return "METADATA_KEY_SHOW"
	case METADATA_KEY_SEASON:
		return "METADATA_KEY_SEASON"
	case METADATA_KEY_EPISODE_ID:
		return "METADATA_KEY_EPISODE_ID"
	case METADATA_KEY_EPISODE_SORT:
		return "METADATA_KEY_EPISODE_SORT"
	case METADATA_KEY_SERVICE_NAME:
		return "METADATA_KEY_SERVICE_NAME"
	case METADATA_KEY_SERVICE_PROVIDER:
		return "METADATA_KEY_SERVICE_PROVIDER"
	case METADATA_KEY_GROUPING:
		return "METADATA_KEY_GROUPING"
	default:
		return "[?? Invalid MetadataKey]"
	}
}

func (t MediaType) String() string {
	if t == MEDIA_TYPE_NONE {
		return "MEDIA_TYPE_NONE"
	}
	parts := ""
	for flag := MEDIA_TYPE_MIN; flag <= MEDIA_TYPE_MAX; flag <<= 1 {
		if t&flag == 0 {
			continue
		}
		switch flag {
		case MEDIA_TYPE_FILE:
			parts += "|" + "MEDIA_TYPE_FILE"
		case MEDIA_TYPE_AUDIO:
			parts += "|" + "MEDIA_TYPE_AUDIO"
		case MEDIA_TYPE_VIDEO:
			parts += "|" + "MEDIA_TYPE_VIDEO"
		case MEDIA_TYPE_IMAGE:
			parts += "|" + "MEDIA_TYPE_IMAGE"
		case MEDIA_TYPE_SUBTITLE:
			parts += "|" + "MEDIA_TYPE_SUBTITLE"
		case MEDIA_TYPE_DATA:
			parts += "|" + "MEDIA_TYPE_DATA"
		case MEDIA_TYPE_ATTACHMENT:
			parts += "|" + "MEDIA_TYPE_ATTACHMENT"
		case MEDIA_TYPE_MUSIC:
			parts += "|" + "MEDIA_TYPE_MUSIC"
		case MEDIA_TYPE_ALBUM:
			parts += "|" + "MEDIA_TYPE_ALBUM"
		case MEDIA_TYPE_COMPILATION:
			parts += "|" + "MEDIA_TYPE_COMPILATION"
		case MEDIA_TYPE_TVSHOW:
			parts += "|" + "MEDIA_TYPE_TVSHOW"
		case MEDIA_TYPE_TVSEASON:
			parts += "|" + "MEDIA_TYPE_TVSEASON"
		case MEDIA_TYPE_TVEPISODE:
			parts += "|" + "MEDIA_TYPE_TVEPISODE"
		case MEDIA_TYPE_AUDIOBOOK:
			parts += "|" + "MEDIA_TYPE_AUDIOBOOK"
		case MEDIA_TYPE_MUSICVIDEO:
			parts += "|" + "MEDIA_TYPE_MUSICVIDEO"
		case MEDIA_TYPE_MOVIE:
			parts += "|" + "MEDIA_TYPE_MOVIE"
		case MEDIA_TYPE_BOOKLET:
			parts += "|" + "MEDIA_TYPE_BOOKLET"
		case MEDIA_TYPE_RINGTONE:
			parts += "|" + "MEDIA_TYPE_RINGTONE"
		case MEDIA_TYPE_ARTWORK:
			parts += "|" + "MEDIA_TYPE_ARTWORK"
		case MEDIA_TYPE_CAPTIONS:
			parts += "|" + "MEDIA_TYPE_CAPTIONS"
		default:
			parts += "|" + "[?? Invalid MediaType value]"
		}
	}
	return strings.Trim(parts, "|")
}

func (t MediaEventType) String() string {
	switch t {
	case MEDIA_EVENT_FILE_ADDED:
		return "MEDIA_EVENT_FILE_ADDED"
	case MEDIA_EVENT_SCAN_START:
		return "MEDIA_EVENT_SCAN_START"
	case MEDIA_EVENT_SCAN_END:
		return "MEDIA_EVENT_SCAN_END"
	case MEDIA_EVENT_ERROR:
		return "MEDIA_EVENT_ERROR"
	default:
		return "[?? Invalid MediaEventType value]"
	}
}
