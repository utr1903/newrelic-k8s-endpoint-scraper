package scraper

type Parser interface {
	Run(data []byte) map[string]string
}
