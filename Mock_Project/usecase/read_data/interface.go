package read_data

type IFile interface {
	ReadFileProcess() ([]string, error)
	InsertCurrentFiles(rows *[]string) error
}
