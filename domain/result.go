package domain

type AnalysisResult struct {
	URL         string
	Reachable   string
	HTMLVersion string
	Headings    string
	Links       string
	LoginForm   string
}
