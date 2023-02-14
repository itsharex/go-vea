package constant

type BusinessType int64

const (
	BUS_OTHER BusinessType = iota
	BUS_INSERT
	BUS_UPDATE
	BUS_DELETE
	BUS_GRANT
	BUS_EXPORT
	BUS_IMPORT
	BUS_FORCE
	BUS_GENCODE
	BUS_CLEAN
)
