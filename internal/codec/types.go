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
