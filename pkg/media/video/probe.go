package video

import (
	"errors"
	"io"
	"math"
	"os"
	"time"

	"github.com/abema/go-mp4"

	"github.com/photoprism/photoprism/pkg/fs"
	"github.com/photoprism/photoprism/pkg/media"
	"github.com/photoprism/photoprism/pkg/media/http/header"
)

// ProbeFile returns information for the given filename.
func ProbeFile(fileName string) (info Info, err error) {
	if fileName == "" {
		return info, errors.New("filename missing")
	}

	var stat os.FileInfo
	var file *os.File

	// Ensure the file exists and is not a directory.
	if stat, err = os.Stat(fileName); err != nil {
		return info, err
	} else if stat.IsDir() {
		return info, errors.New("invalid filename")
	}

	// Open the file for reading.
	if file, err = os.Open(fileName); err != nil {
		return info, err
	}

	defer file.Close()

	// Get video information.
	info, err = Probe(file)

	// Add file name.
	info.FileName = fileName
	info.FileSize = stat.Size()

	// Set file type based on filename?
	if info.FileType == fs.TypeUnknown {
		info.FileType = fs.FileType(fileName)
	}

	// Set media type based on filename?
	if info.MediaType == media.Unknown {
		info.MediaType = media.Formats[info.FileType]
	}

	return info, err
}

// Probe returns information on the provided video file.
// see https://pkg.go.dev/github.com/abema/go-mp4#ProbeInfo
func Probe(file io.ReadSeeker) (info Info, err error) {
	// Set defaults.
	info = NewInfo()

	if file == nil {
		return info, errors.New("file is nil")
	}

	// Probe file from the beginning.
	if _, err = file.Seek(0, io.SeekStart); err != nil {
		return info, err
	}

	// Find file type start offset.
	if offset, findErr := CompatibleBrands.FileTypeOffset(file); findErr != nil {
		return info, findErr
	} else if offset < 0 {
		return info, nil
	} else {
		info.VideoOffset = int64(offset)
	}

	// Ignore any data before the video offset.
	videoFile := NewReadSeeker(file, info.VideoOffset)

	var video *mp4.ProbeInfo

	// Detect Mp4 video metadata.
	if video, err = mp4.Probe(videoFile); err != nil {
		return info, err
	}

	// Check compatibility.
	if CompatibleBrands.ContainsAny(video.CompatibleBrands) {
		info.Compatible = true
		info.VideoType = Mp4
		info.VideoMimeType = header.ContentTypeMp4
		info.FPS = 30.0 // TODO: Detect actual frames per second!

		if info.VideoOffset > 0 {
			info.MediaType = media.Live
			info.ThumbOffset = 0
		} else {
			info.MediaType = media.Video
			info.FileType = fs.VideoMp4
		}
	}

	// Check major brand.
	switch video.MajorBrand {
	case ChunkQT.Get():
		info.VideoType = Mov
		info.VideoMimeType = header.ContentTypeMov
		if info.MediaType == media.Video {
			info.FileType = fs.VideoMov
		}
	}

	// Additional properties.
	info.FastStart = video.FastStart
	info.Tracks = len(video.Tracks)

	// Check tracks for codec, resolution and encryption.
	for _, track := range video.Tracks {
		if track.Encrypted && !info.Encrypted {
			info.Encrypted = track.Encrypted
		}

		switch track.Codec {
		case mp4.CodecAVC1:
			info.VideoCodec = CodecAvc1
		}

		if avc := track.AVC; avc != nil {
			if info.VideoMimeType == "" {
				info.VideoMimeType = header.ContentTypeMp4
			}
			if w := int(avc.Width); w > info.VideoWidth {
				info.VideoWidth = w
			}
			if h := int(avc.Height); h > info.VideoHeight {
				info.VideoHeight = h
			}
		}
	}

	// If no AVC video was found, search the video data for High Efficiency Video Coding (HEVC) chunks,
	// see https://stackoverflow.com/questions/63468587/what-hevc-codec-tag-to-use-with-fmp4-hvc1-or-hev1.
	if info.VideoCodec == "" {
		// To improve performance, only search for "hvc1" as that is the most common HEVC video identifier.
		if fileOffset, fileErr := ChunkHVC1.DataOffset(file, -1); fileOffset > 0 && fileErr == nil {
			info.VideoCodec = CodecHvc1
		}
	}

	// Calculate video duration in seconds.
	if video.Duration > 0 {
		s := float64(video.Duration) / float64(video.Timescale)

		info.Duration = time.Duration(s * float64(time.Second))

		if info.FPS > 0 {
			info.Frames = int(math.Round(info.FPS * s))
		}
	}

	return info, err
}
