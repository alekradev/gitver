package version

type IHistoryEntry interface {
	GetVersion() IVersion
	GetDate() string
}

type HistoryEntry struct {
	Version IVersion `yaml:"Version"`
	Date    string   `yaml:"Date"`
}

func (h *HistoryEntry) GetVersion() IVersion {
	return h.Version
}

func (h *HistoryEntry) GetDate() string {
	return h.Date
}

type HistoryEntryBuilderObj struct {
	obj     HistoryEntry
	version IVersion
	date    string
}

func HistoryEntryBuilder() *HistoryEntryBuilderObj {
	return new(HistoryEntryBuilderObj)
}

func (b *HistoryEntryBuilderObj) WithVersion(version IVersion) *HistoryEntryBuilderObj {
	b.version = version
	return b
}

func (b *HistoryEntryBuilderObj) WithDate(date string) *HistoryEntryBuilderObj {
	b.date = date
	return b
}

func (b *HistoryEntryBuilderObj) Build() IHistoryEntry {
	h := new(HistoryEntry)
	h.Version = b.version
	h.Date = b.date
	return h
}
