package filepath

import (
	"regexp"
	"strings"

	"github.com/patrickhuber/go-cross/os"
	"github.com/patrickhuber/go-cross/platform"
)

type Provider interface {
	Separator() PathSeparator
	Comparison() Comparison
	VolumeName(path string) string
	Abs(path string) (string, error)
	Join(elements ...string) string
	Rel(sourcepath string, targetpath string) (string, error)
	Root(path string) string
	Clean(path string) string
	Dir(path string) string
	Ext(path string) string
	Base(path string) string
	Normalize(path string) (string, error)
	Parse(path string) (FilePath, error)
	String(fp FilePath) string
}

type provider struct {
	os         os.OS
	separator  PathSeparator
	comparison Comparison
	parser     Parser
}

func NewProvider(o os.OS, parser Parser, separator PathSeparator, cmp Comparison) Provider {
	return &provider{
		os:         o,
		parser:     parser,
		separator:  separator,
		comparison: cmp,
	}
}

func NewProviderFromOS(o os.OS) Provider {
	plat := o.Platform()
	parser := NewParserFromPlatform(plat)

	p := &provider{
		parser: parser,
		os:     o,
	}
	if platform.IsWindows(plat) {
		p.comparison = IgnoreCase
		p.separator = BackwardSlash
	} else {
		p.comparison = CaseSensitive
		p.separator = ForwardSlash
	}
	return p
}

func (p *provider) Abs(path string) (string, error) {
	wd, err := p.os.WorkingDirectory()
	if err != nil {
		return "", err
	}
	return p.abs(wd, path)
}

func (p *provider) abs(wd, rel string) (string, error) {
	fp, err := p.parser.Parse(rel)
	if err != nil {
		return "", err
	}
	if fp.IsAbs() {
		return p.String(fp.Clean()), nil
	}
	wdp, err := p.parser.Parse(wd)
	if err != nil {
		return "", err
	}
	abs := wdp.Join(fp)
	return p.String(abs.Clean()), nil
}

// Join implements Processor
func (p *provider) Join(elements ...string) string {
	if len(elements) == 0 {
		return ""
	}

	var accumulator FilePath
	first := true
	for _, element := range elements {

		// skip empty elements
		if len(element) == 0 {
			continue
		}

		// call parse on the first element
		// set the first element as the accumulator
		if first {
			accumulator, _ = p.parser.Parse(element)
			first = false
			continue
		}

		// call parse on each next element
		next, _ := p.parser.Parse(element)

		// and then join the accumulator to that element
		accumulator = accumulator.Join(next)
	}

	// call clean on the result
	return p.String(accumulator.Clean())
}

// Rel implements Processor
func (p *provider) Rel(sourcepath string, targetpath string) (string, error) {

	source, err := p.parser.Parse(sourcepath)
	if err != nil {
		return "", err
	}

	target, err := p.parser.Parse(targetpath)
	if err != nil {
		return "", err
	}

	result, err := source.Rel(target, p.comparison)
	if err != nil {
		return "", err
	}

	return result.String(p.separator), nil
}

// Clean implements Processor
func (p *provider) Clean(path string) string {
	fp, _ := p.parser.Parse(path)

	// for empty unc paths, normalize the original string
	// (is there a way to do this in the String method?)
	fp = fp.Clean()

	// on the windows platform if the first segment matches the drive pattern
	// the current directory needs to be added in the front
	if platform.IsWindows(p.os.Platform()) && fp.IsRel() && len(fp.Segments) > 0 {
		matched, err := regexp.MatchString(`^[a-zA-Z][:]`, fp.Segments[0])
		if (err == nil) && matched {
			fp.Segments = append([]string{CurrentDirectory}, fp.Segments...)
		}
	}

	cleaned := fp.String(p.separator)
	return cleaned
}

// Root is a helper function to print the root of the filepath
func (p *provider) Root(path string) string {
	fp, _ := p.parser.Parse(path)
	return p.String(fp.Root())
}

// VolumeName behaves similar to filepath.VolumeName in the path/filepath package
func (p *provider) VolumeName(path string) string {
	fp, _ := p.parser.Parse(path)
	return fp.VolumeName(p.separator)
}

func (p *provider) Ext(path string) string {
	fp, _ := p.parser.Parse(path)
	return fp.Ext()
}

func (p *provider) Dir(path string) string {
	fp, _ := p.parser.Parse(path)
	dir := fp.Dir()
	return dir.String(p.separator)
}

// Base returns the last element of path. Trailing path separators are removed before extracting the last element. If the path is empty, Base returns ".". If the path consists entirely of separators, Base returns a single separator.
func (p *provider) Base(path string) string {
	fp, _ := p.parser.Parse(path)
	base := fp.Base()
	return base.String(p.separator)
}

// String returns the string representation of the file path
func (p *provider) String(fp FilePath) string {
	return fp.String(p.separator)
}

func (p *provider) Separator() PathSeparator {
	return p.separator
}

func (p *provider) Comparison() Comparison {
	return p.comparison
}

func (p *provider) Normalize(path string) (string, error) {
	fp, err := p.parser.Parse(path)
	if err != nil {
		return "", err
	}
	path = fp.String(p.separator)
	if p.comparison == IgnoreCase {
		return strings.ToLower(path), nil
	}
	return path, nil
}

func (p *provider) Parse(path string) (FilePath, error) {
	return p.parser.Parse(path)
}
