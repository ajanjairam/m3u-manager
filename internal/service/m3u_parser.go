package service

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/ajanjairam/m3u-manager/internal/models"
	"github.com/samber/lo"
)

var extinfTagRe = regexp.MustCompile(`([a-zA-Z0-9-]+?)="([^"]+)"`)
var urlRegex = regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)

var (
	ErrReadInvalidM3UUrl    = errors.New("m3u: invalid url")
	ErrReadInvalidM3UFile   = errors.New("m3u: invalid file")
	ErrReadInvalidM3UHeader = errors.New("m3u: missing #EXTM3U header")
	ErrReadInvalidEXTINF    = errors.New("m3u: invalid EXTINF metadata")
	ErrWriteChannelsEmpty   = errors.New("m3u: invalid EXTINF metadata")
	ErrWriteM3UInternal     = errors.New("m3u: error during write")
	ErrURIPrecedesChannels  = errors.New("m3u: URI before any track definition")
)

func ParseM3U(ctx context.Context, m3uSource string) (models.Playlist, error) {
	reader, err := parseM3USource(ctx, m3uSource)
	if err != nil {
		return models.Playlist{}, err
	}
	defer reader.Close()
	input := bufio.NewScanner(reader)
	var channels []models.Channel
	first := true

	for input.Scan() {
		if err := ctx.Err(); err != nil {
			return models.Playlist{}, err
		}
		line := input.Text()

		if first {
			first = false
			if !strings.HasPrefix(line, "#EXTM3U") {
				return models.Playlist{}, ErrReadInvalidM3UHeader
			}
			continue
		}

		switch {
		case strings.HasPrefix(line, "#EXTINF:"):
			ch, err := parseEXTINF(line)
			if err != nil {
				return models.Playlist{}, err
			}
			channels = append(channels, ch)

		case strings.HasPrefix(line, "#"), strings.HasPrefix(line, "--"), line == "":

		case len(channels) == 0:
			return models.Playlist{}, fmt.Errorf("%w: %q", ErrURIPrecedesChannels, line)

		case channels[len(channels)-1].Uri == "":
			channels[len(channels)-1].Uri = strings.TrimSpace(line)
		}
	}

	if err := input.Err(); err != nil {
		return models.Playlist{}, fmt.Errorf("m3u scan: %w", err)
	}
	return models.Playlist{Channels: channels}, nil
}

func parseM3USource(ctx context.Context, m3uSource string) (io.ReadCloser, error) {
	var reader io.ReadCloser
	if urlRegex.MatchString(m3uSource) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, m3uSource, nil)
		if err != nil {
			return nil, ErrReadInvalidM3UUrl
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, ErrReadInvalidM3UUrl
		}

		if resp.StatusCode != http.StatusOK {
			return nil, ErrReadInvalidM3UUrl
		}
		reader = resp.Body
	} else {
		file, err := os.Open(m3uSource)
		if err != nil {
			return nil, ErrReadInvalidM3UFile
		}
		reader = file
	}
	return reader, nil
}

func parseEXTINF(line string) (models.Channel, error) {
	line = strings.TrimPrefix(line, "#EXTINF:")

	meta, name, ok := strings.Cut(line, ",")
	if !ok {
		return models.Channel{}, ErrReadInvalidEXTINF
	}

	rawLen, tags, _ := strings.Cut(meta, " ")
	length, err := strconv.ParseInt(rawLen, 10, 64)
	if err != nil {
		return models.Channel{}, ErrReadInvalidEXTINF
	}

	channel := models.Channel{
		Name:   strings.TrimSpace(name),
		Length: length,
	}
	channel = mapTags(tags, channel)
	return channel, nil
}

func mapTags(line string, ch models.Channel) models.Channel {
	matches := extinfTagRe.FindAllStringSubmatch(line, -1)
	lo.ForEach(matches, func(item []string, _ int) {
		switch item[1] {
		case "tvg-id":
			ch.TvgID = item[2]
		case "tvg-name":
			ch.TvgName = item[2]
		case "tvg-logo":
			ch.TvgLogo = item[2]
		case "group-title":
			ch.GroupTitle = item[2]
		}
	})
	return ch
}

func MarshalM3U(channels []models.Channel) (string, error) {
	if channels == nil || len(channels) == 0 {
		return "", ErrWriteChannelsEmpty
	}
	var writer strings.Builder
	_, err := writer.WriteString("#EXTM3U\n")
	if err != nil {
		return "", ErrWriteM3UInternal
	}

	for i := range channels {
		_, err := writer.WriteString(fmt.Sprintf("#EXTINF:%d,", channels[i].Length))
		if err != nil {
			return "", ErrWriteM3UInternal
		}

		err = writeTagAttrs(&writer, channels[i])
		if err != nil {
			return "", ErrWriteM3UInternal
		}
		_, err = writer.WriteString(fmt.Sprintf(",%s\n%s\n", channels[i].Name, channels[i].Uri))
		if err != nil {
			return "", ErrWriteM3UInternal
		}
	}
	return writer.String(), nil
}

func writeTagAttrs(w *strings.Builder, ch models.Channel) error {
	for _, kv := range [][2]string{
		{"tvg-id", ch.TvgID},
		{"tvg-name", ch.TvgName},
		{"tvg-logo", ch.TvgLogo},
		{"group-title", ch.GroupTitle},
	} {
		if kv[1] == "" {
			continue
		}
		_, err := w.WriteString(fmt.Sprintf(` %s="%s"`, kv[0], kv[1]))
		if err != nil {
			return ErrWriteM3UInternal
		}
	}
	return nil
}
