package sitemap

import "encoding/xml"

/**
?xml version="1.0" encoding="UTF-8"?>

<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">

   <url>

      <loc>http://www.example.com/</loc>

      <lastmod>2005-01-01</lastmod>

      <changefreq>monthly</changefreq>

      <priority>0.8</priority>

   </url>

</urlset>
*/

type URL struct {
	Loc        string  `xml:"loc"`
	LastMod    string  `xml:"lastmod"`
	ChangeFreq string  `xml:"changefreq"`
	Priority   float32 `xml:"priority"`
}

type URLset struct {
	XMLName xml.Name `xml:"urlset"`
	URLS    []URL    `xml:"url"`
}

type SiteMap struct {
	Entries URLset
}

func NewSiteMap() *SiteMap {
	return &SiteMap{
		Entries: URLset{},
	}
}

func (s *SiteMap) AddPage(path string, options ...func(*URL)) {
	url := URL{
		Loc: path,
	}
	for _, option := range options {
		option(&url)
	}
	s.Entries.URLS = append(s.Entries.URLS, url)
}

func WithLastMod(lastMod string) func(*URL) {
	return func(url *URL) {
		url.LastMod = lastMod
	}
}

func WithChangeFreq(changeFreq string) func(*URL) {
	return func(url *URL) {
		url.ChangeFreq = changeFreq
	}
}

func WithPriority(priority float32) func(*URL) {
	return func(url *URL) {
		url.Priority = priority
	}
}

func (s SiteMap) Generate() string {
	b, _ := xml.MarshalIndent(s.Entries, "", "  ")
	return string(b)
}
