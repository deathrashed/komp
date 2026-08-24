package codec

type Kind int

const (
	KindArchive Kind = iota // zip, 7z, rar… containers
	KindStream              // gzip, xz… single-file
	KindTarPipe             // folder→tar→stream (tar.gz)
)

type Codec struct {
	Name        string
	Extensions  []string // longest-first matching matters
	Kind        Kind
	Bin         string
	BrewFormula string
	CreateArgs  []string // {in} {out} {level} substituted
	AddArgs     []string // empty ⇒ add unsupported
	DeleteArgs  []string // empty ⇒ clean-in-place unsupported
	TestArgs    []string
	ListArgs    []string
	ExtractArgs []string
	ExtractBin  string   // override for extraction; defaults to Bin
	Fast        string // level token, e.g. "-1"
	Normal      string
	Max         string
	WrapExt     string // ".tar" prefix for TarPipe outputs, e.g. ".tar.gz"
}

func (c Codec) CanAdd() bool     { return len(c.AddArgs) > 0 }
func (c Codec) CanList() bool    { return len(c.ListArgs) > 0 }
func (c Codec) CanTest() bool    { return len(c.TestArgs) > 0 }
func (c Codec) CanExtract() bool { return len(c.ExtractArgs) > 0 || c.ExtractBin != "" }
func (c Codec) CanClean() bool {
	return len(c.DeleteArgs) > 0 || c.Name == "tar" || c.Kind == KindTarPipe
}
