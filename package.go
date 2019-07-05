package main

type Package struct {
	Name        string `json:"name"`
	HTMLURL     string `json:"html_url"`
	CloneURL    string `json:"clone_url"`
	Description string `json:"description"`
}
