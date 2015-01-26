package spotify

import (
	"net/http"
	"testing"
)

const albumsResponse = `
{
	"href" : "https://api.spotify.com/v1/artists/1vCWHaC5f2uS3yhpwWbIA6/albums?offset=0&limit=2&album_type=single",
	"items" : [ {
		"album_type" : "single",
		"available_markets" : [ "AD", "AR", "AT", "AU", "BE", "BG", "BO", "BR", "CH", "CL", "CO", "CR", "CY", "CZ", "DE", "DK", "DO", "EC", "EE", "ES", "FI", "FR", "GR", "GT", "HK", "HN", "HU", "IE", "IS", "IT", "LI", "LT", "LU", "LV", "MC", "MT", "MY", "NI", "NL", "NO", "NZ", "PA", "PE", "PH", "PL", "PT", "PY", "RO", "SE", "SG", "SI", "SK", "SV", "TR", "TW", "UY" ],
		"external_urls" : {
			"spotify" : "https://open.spotify.com/album/3ckwyt0bTOcDbXovWbweMp"
			},
		"href" : "https://api.spotify.com/v1/albums/3ckwyt0bTOcDbXovWbweMp",
		"id" : "3ckwyt0bTOcDbXovWbweMp",
		"images" : [ {
			"height" : 640,
			"url" : "https://i.scdn.co/image/144ac57ad073741e99b5243c59abebe1500ada0a",
			"width" : 640
			}, {
			"height" : 300,
			"url" : "https://i.scdn.co/image/4680e5f3af02219fd9e79ce432c1b18f97af6426",
			"width" : 300
			}, {
			"height" : 64,
			"url" : "https://i.scdn.co/image/8c803d6cb612b6f2b37a7276deb2ff05f5a77097",
			"width" : 64
			} ],
		"name" : "The Days / Nights",
		"type" : "album",
		"uri" : "spotify:album:3ckwyt0bTOcDbXovWbweMp"
		}, {
			"album_type" : "single",
			"available_markets" : [ "CA", "MX", "US" ],
			"external_urls" : {
				"spotify" : "https://open.spotify.com/album/1WXM7DYQRT7QX8AKBJRfK9"
			},
		"href" : "https://api.spotify.com/v1/albums/1WXM7DYQRT7QX8AKBJRfK9",
		"id" : "1WXM7DYQRT7QX8AKBJRfK9",
		"images" : [ {
			"height" : 640,
			"url" : "https://i.scdn.co/image/590dbe5504d2898c120b942bee2b699404783896",
			"width" : 640
			}, {
			"height" : 300,
			"url" : "https://i.scdn.co/image/9a4db24b1930e8683b4dfd19c7bd2a40672c6718",
			"width" : 300
			}, {
			"height" : 64,
			"url" : "https://i.scdn.co/image/d5cfc167e03ed328ae7dfa9b56d3628d81b6831b",
			"width" : 64
			} ],
			"name" : "The Days / Nights",
			"type" : "album",
			"uri" : "spotify:album:1WXM7DYQRT7QX8AKBJRfK9"
			} ],
	"limit" : 2,
	"next" : "https://api.spotify.com/v1/artists/1vCWHaC5f2uS3yhpwWbIA6/albums?offset=2&limit=2&album_type=single",
	"offset" : 0,
	"previous" : null,
	"total" : 157
}`

func TestFindArtist(t *testing.T) {
	client := testClientFile(http.StatusOK, "test_data/find_artist.txt")
	artist, err := client.FindArtist(ID("0TnOYISbd1XYRBk9myaseg"))
	if err != nil {
		t.Error(err)
		return
	}
	if followers := artist.Followers.Count; followers != 2265279 {
		t.Errorf("Got %d followers, want 2265279\n", followers)
		return
	}
	if artist.Name != "Pitbull" {
		t.Error("Got ", artist.Name, ", wanted Pitbull")
	}
}

func TestArtistTopTracks(t *testing.T) {
	client := testClientFile(http.StatusOK, "test_data/artist_top_tracks.txt")
	tracks, err := client.ArtistsTopTracks(ID("43ZHCT0cAZBISjO8DG9PnE"), "SE")
	if err != nil {
		t.Error(err)
		return
	}
	l := len(tracks)
	if l != 10 {
		t.Errorf("Got %d tracks, expected 10\n", l)
	}
	track := tracks[9]
	if track.Name != "(You're The) Devil in Disguise" {
		t.Error("Incorrect track name")
	}
	if track.TrackNumber != 24 {
		t.Errorf("Track number was %d, expected 24\n", track.TrackNumber)
	}
}

func TestRelatedArtists(t *testing.T) {
	client := testClientFile(http.StatusOK, "test_data/related_artists.txt")
	artists, err := client.FindRelatedArtists(ID("43ZHCT0cAZBISjO8DG9PnE"))
	if err != nil {
		t.Error(err)
		return
	}
	if count := len(artists); count != 20 {
		t.Errorf("Got %d artists, wanted 20\n", count)
		return
	}
	a2 := artists[2]
	if a2.Name != "Carl Perkins" {
		t.Error("Expected Carl Perkins, got ", a2.Name)
	}
	if a2.Popularity != 54 {
		t.Errorf("Expected popularity 54, got %d\n", a2.Popularity)
	}
}

func TestArtistAlbumsFiltered(t *testing.T) {
	client := testClientString(http.StatusOK, albumsResponse)

	l := 2
	var typ AlbumType = AlbumTypeSingle

	options := Options{}
	options.Limit = &l

	albums, err := client.ArtistAlbumsOpt(ID("1vCWHaC5f2uS3yhpwWbIA6"), &options, &typ)
	if err != nil {
		t.Error(err)
		return
	}
	if albums == nil {
		t.Error("Result is nil")
		return
	}
	// since we didn't specify a country, we got duplicate albums
	// (the album has a different ID in different regions)
	if l := len(albums.Albums); l != 2 {
		t.Errorf("Expected 2 albums, got %d\n", l)
		return
	}
	if albums.Albums[0].Name != "The Days / Nights" {
		t.Error("Expected 'The Days / Nights', got ", albums.Albums[0].Name)
	}
}