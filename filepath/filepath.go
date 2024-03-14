package filepath

import (
	"fmt"
	"os"
	"strings"
)

type PathType string

type FilePath struct {
	Volume   Volume
	Absolute bool
	Segments []string
}

type Volume struct {
	Host  Nullable[string]
	Share Nullable[string]
	Drive Nullable[string]
}

type Nullable[T any] struct {
	HasValue bool
	Value    T
}

type PathSeparator rune
type PathListSeparator rune

const (
	ForwardSlash             PathSeparator     = '/'
	BackwardSlash            PathSeparator     = '\\'
	Colon                    PathListSeparator = ':'
	SemiColon                PathListSeparator = ';'
	DefaultPathListSeparator PathListSeparator = os.PathListSeparator
	DefaultPathSeparator     PathSeparator     = os.PathSeparator

	CurrentDirectory = "."
	ParentDirectory  = ".."
	EmptyDirectory   = ""
)

func (fp FilePath) IsAbs() bool {
	return fp.Absolute
}

func (fp FilePath) IsRel() bool {
	return !fp.Absolute
}

func (fp FilePath) isWindows() bool {
	return fp.Volume.Drive.HasValue
}

func (fp FilePath) isUNC() bool {
	return fp.Volume.Host.HasValue
}

func (fp FilePath) Root() FilePath {
	return FilePath{
		Volume:   fp.Volume,
		Absolute: fp.Absolute,
	}
}

func (fp FilePath) Join(other FilePath) FilePath {
	return FilePath{
		Volume: Volume{
			Host:  fp.Volume.Host,
			Share: fp.Volume.Share,
			Drive: fp.Volume.Drive,
		},
		Absolute: fp.Absolute,
		Segments: append(fp.Segments, other.Segments...),
	}
}

func (fp FilePath) VolumeName(sep PathSeparator) string {
	switch {

	case fp.isWindows():
		// Windows
		var builder strings.Builder
		if fp.Volume.Drive.HasValue {
			builder.WriteString(fp.Volume.Drive.Value)
		}
		return builder.String()

	case fp.isUNC():
		return fp.uncVolumeName(sep)

	case fp.IsRel():
		// Relative
		return ""
	}

	// Unix
	return ""
}

func (fp FilePath) uncVolumeName(sep PathSeparator) string {
	var builder strings.Builder

	// write //
	builder.WriteByte(byte(sep))
	builder.WriteByte(byte(sep))

	// write the hostname
	if fp.Volume.Host.HasValue {
		builder.WriteString(fp.Volume.Host.Value)
	}

	// write separator and share if share exists
	if fp.Volume.Share.HasValue {
		builder.WriteByte(byte(sep))
		builder.WriteString(fp.Volume.Share.Value)
	}

	return builder.String()
}

func (fp FilePath) String(sep PathSeparator) string {
	var builder strings.Builder

	// write the volume name
	builder.WriteString(fp.VolumeName(sep))

	switch {
	// relative paths don't need a separator
	case fp.IsRel():

	// absolute windows and unix paths need a separator
	case !fp.isUNC():
		builder.WriteRune(rune(sep))

	// unc paths with segments need a separator
	case len(fp.Segments) > 0:
		builder.WriteRune(rune(sep))
	}

	// write the segments
	for i, seg := range fp.Segments {
		if i > 0 {
			builder.WriteRune(rune(sep))
		}
		builder.WriteString(seg)
	}
	return builder.String()
}

func (fp FilePath) Clean() FilePath {
	clean := fp.clean()

	// add a current directory indicator for empty relative paths
	if len(clean.Segments) == 0 && clean.IsRel() {
		clean.Segments = append(clean.Segments, CurrentDirectory)
	}
	return clean
}

func (fp FilePath) clean() FilePath {
	// if the path has no segments it is already clean
	if len(fp.Segments) == 0 {
		return fp
	}

	// unc paths with one empty segment are already clean
	if fp.isUNC() && len(fp.Segments) == 1 && fp.Segments[0] == "" {
		return fp
	}

	var s []string
	for _, segment := range fp.Segments {

		switch segment {
		// remove . and empty directories (resulting from // in path)
		case CurrentDirectory:
			continue
		case EmptyDirectory:
			continue
		case ParentDirectory:
			switch {
			case len(s) > 0:
				// if the current segment is parent
				// and the last segment is parent
				// we have already processed a parent so write both to the output
				var item string
				item, s = pop(s)
				if item == ParentDirectory {
					s = push(s, ParentDirectory)
					s = push(s, ParentDirectory)
				}
			case fp.IsRel():
				// if the current segment is parent
				// and there are no elements to process (s.Length() == 0)
				// push the ..
				s = push(s, ParentDirectory)
			}
		default:
			// if the segment matches the pattern \w[:] it is a drive letter in windows
			s = push(s, segment)
		}
	}

	return FilePath{
		Volume:   fp.Volume,
		Absolute: fp.Absolute,
		Segments: s,
	}
}

func (fp FilePath) Rel(other FilePath, cmp Comparison) (FilePath, error) {
	source := fp.Clean()
	target := other.Clean()

	// if paths are equal, return CurrentDirectory string
	if source.Equal(target, cmp) {
		return FilePath{
			Segments: []string{CurrentDirectory},
		}, nil
	}

	// remove any current directory '.' only source paths
	if len(source.Segments) == 1 && source.Segments[0] == CurrentDirectory {
		source.Segments = nil
	}

	// both paths must be either relative or absolute
	// if absolute both paths must match volumes
	if source.Absolute != target.Absolute || !source.Volume.Equal(target.Volume, cmp) {
		return FilePath{}, fmt.Errorf("can't make target relative to source: absolute paths must share a prefix")
	}

	// get the first index where segments differ
	firstDiff := source.firstSegmentDiff(target, cmp)

	if firstDiff < len(source.Segments) && source.Segments[firstDiff] == ".." {
		return FilePath{}, fmt.Errorf("can't make target relative to source")
	}

	var segments []string

	// run the source to the end by adding ..
	for s := firstDiff; s < len(source.Segments); s++ {
		segments = append(segments, ParentDirectory)
	}

	// run the target to the end by adding target[firstDiff:]
	if firstDiff < len(target.Segments) {
		segments = append(segments, target.Segments[firstDiff:]...)
	}

	return FilePath{
		Segments: segments,
		Absolute: false,
	}, nil
}

func (source FilePath) firstSegmentDiff(target FilePath, cmp Comparison) int {

	sourceLen := len(source.Segments)
	targetLen := len(target.Segments)

	segmentLen := sourceLen
	if targetLen < sourceLen {
		segmentLen = targetLen
	}

	// find the first differing element
	for diffPosition := 0; diffPosition < segmentLen; diffPosition++ {
		if !cmp.Equal(source.Segments[diffPosition], target.Segments[diffPosition]) {
			return diffPosition
		}
	}
	return segmentLen
}

func (fp FilePath) Dir() FilePath {
	last := len(fp.Segments) - 1
	if last >= 0 && fp.Segments[last] != "" {
		fp.Segments[last] = ""
	}
	return fp.Clean()
}

func (fp FilePath) Ext() string {
	if len(fp.Segments) == 0 {
		return ""
	}
	last := fp.Segments[len(fp.Segments)-1]
	for i := len(last) - 1; i >= 0; i-- {
		if last[i] == '.' {
			return last[i:]
		}
	}
	return ""
}

// Base returns the last element of path. Trailing path separators are removed before extracting the last element. If the path is empty, Base returns ".". If the path consists entirely of separators, Base returns a single separator.
func (fp FilePath) Base() FilePath {

	// the empty path case
	if fp.IsRel() && len(fp.Segments) == 0 {
		return FilePath{
			Absolute: false,
			Segments: []string{CurrentDirectory},
		}
	}

	// find the last non empty element
	for i := len(fp.Segments) - 1; i >= 0; i-- {
		if fp.Segments[i] == "" {
			continue
		}
		return FilePath{
			Absolute: false,
			Segments: fp.Segments[i : i+1],
		}
	}

	return FilePath{
		Absolute: true,
	}
}

// Equal compares two paths using case sensitive comparison
func (fp FilePath) Equal(other FilePath, cmp Comparison) bool {
	if fp.Absolute != other.Absolute {
		return false
	}
	if !fp.Volume.Equal(other.Volume, cmp) {
		return false
	}
	if len(fp.Segments) != len(other.Segments) {
		return false
	}

	for i := 0; i < len(fp.Segments); i++ {
		if !cmp.Equal(fp.Segments[i], other.Segments[i]) {
			return false
		}
	}
	return true
}

// Equal compares two volumes using case sensetive comparison
func (v Volume) Equal(other Volume, cmp Comparison) bool {
	if !nullableStringEqual(v.Drive, other.Drive) {
		return false
	}
	if !nullableStringEqual(v.Host, other.Host) {
		return false
	}
	return nullableStringEqual(v.Share, other.Share)
}

func nullableStringEqual(first Nullable[string], second Nullable[string]) bool {
	if !first.HasValue && !second.HasValue {
		return true
	}
	if first.HasValue != second.HasValue {
		return false
	}
	return strings.Compare(first.Value, second.Value) == 0
}

func push[T any](list []T, item T) []T {
	return append(list, item)
}

func pop[T any](list []T) (item T, remainder []T) {
	var zero T
	if len(list) == 0 {
		return zero, list
	}
	item = list[len(list)-1]
	remainder = list[:len(list)-1]
	return
}
