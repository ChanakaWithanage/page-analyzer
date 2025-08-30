package contract

type AnalyzeParams struct {
	URL                 string
	FetchTimeoutSeconds int
}

type AnalyzeResult struct {
	URL               string            `json:"url"`
	HTMLVersion       string            `json:"html_version"`
	Title             string            `json:"title"`
	Headings          map[string]int    `json:"headings"`
	LinksInternal     int               `json:"links_internal"`
	LinksExternal     int               `json:"links_external"`
	LinksInaccessible int               `json:"links_inaccessible"`
	LoginFormPresent  bool              `json:"login_form_present"`
	Warnings          []string          `json:"warnings,omitempty"`
	Errors            []string          `json:"errors,omitempty"`
}
