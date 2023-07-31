package read_data

type IFile interface {
	ReadFileProcess(path string) ([]string, error)
}
