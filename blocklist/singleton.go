package blocklist

var writer Writer

const persistanceFileName = "blocked_summoners.txt"

func init() {
	writer = NewFileWriter(persistanceFileName)
}

func Append(str string) error {
	return writer.Append(str)
}

func Contains(str string) (bool, error) {
	return writer.Contains(str)
}
