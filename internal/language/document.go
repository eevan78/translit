package language

type Document interface {
	open()
	transliterate()
	getInputFilePath() string
	getOuputFilePath() string
}
