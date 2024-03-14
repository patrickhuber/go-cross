package filepath

import (
	"regexp"

	"github.com/patrickhuber/go-cross/os"
	"github.com/patrickhuber/go-cross/platform"
	"github.com/patrickhuber/go-types"
	"github.com/patrickhuber/go-types/option"
)

type Provider interface {
	Separator() PathSeparator
	VolumeName(path string) string
	Abs(path string) (string, error)
	Join(elements ...string) string
	Rel(sourcepath string, targetpath string) (string, error)
	Root(path string) string
	Clean(path string) string
	Dir(path string) string
	Ext(path string) string
	Base(path string) string
}

type provider struct {
	os         os.OS
	separator  PathSeparator
	comparison Comparison
	parser     Parser
}

type providerOptions struct {
	os         types.Option[os.OS]
	separator  types.Option[PathSeparator]
	comparison types.Option[Comparison]
}

type ProviderOption func(p *providerOptions)

func WithOS(os os.OS) ProviderOption {
	return func(p *providerOptions) {
		p.os = option.Some(os)
	}
}

func WithSeparator(sep PathSeparator) ProviderOption {
	return func(p *providerOptions) {
		p.separator = option.Some(sep)
	}
}

func WithComparison(cmp Comparison) ProviderOption {
	return func(p *providerOptions) {
		p.comparison = option.Some(cmp)
	}
}

func NewProvider(options ...ProviderOption) Provider {

	pop := &providerOptions{
		os:         option.None[os.OS](),
		separator:  option.None[PathSeparator](),
		comparison: option.None[Comparison](),
	}

	for _, option := range options {
		option(pop)
	}

	p := &provider{}

	// is the OS Set? if not use the default
	if o, ok := pop.os.Deconstruct(); ok {
		p.os = o
	} else {
		p.os = os.New()
	}

	// set the default seperator and comparison operations
	if platform.IsPosix(p.os.Platform()) {
		p.separator = ForwardSlash
		p.comparison = CaseSensitive{}
	} else {
		p.separator = BackwardSlash
		p.comparison = IgnoreCase{}
	}

	// set overrides
	if s, ok := pop.separator.Deconstruct(); ok {
		p.separator = s
	}
	if c, ok := pop.comparison.Deconstruct(); ok {
		p.comparison = c
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
