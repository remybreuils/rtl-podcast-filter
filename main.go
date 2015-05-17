package main

import (
	"fmt"
	"net/http"
	"encoding/xml"
	"io/ioutil"
	"regexp"
	"log"
)

func main() {
	http.HandleFunc("/rtl/podcast/integral-les-grosses-tetes.xml", getRtlLgti)
	http.HandleFunc("/", getRtlLgti)
    http.ListenAndServe(":8080", nil)
	
	fmt.Println("Launch server on 8080")
}


type rssFeed struct {
	XMLName xml.Name   `xml:"rss"` 
	XmlnsAtom		string			`xml:"xmlns:atom,attr"`  
	XmlnsItunes		string          `xml:"xmlns:itunes,attr,omitempty"`  
	XmlnsMedia		string          `xml:"xmlns:media,attr,omitempty"`  
	XmlnsDcTerms	string			`xml:"xmlns:dcterms,attr,omitempty"`                                                                                    
	XmlVersion 		string			`xml:"version,attr"`
	
	Title			string			`xml:"channel>title"`
	Description		string			`xml:"channel>description"`
	ImageList		rssFeedImage	`xml:"channel>image,omitempty"`
	Language		string			`xml:"channel>language,omitempty"`
	Link 			string			`xml:"channel>link"`
	Copyright		string			`xml:"channel>copyright,omitempty"`
	Ttl				int32			`xml:"channel>ttl,omitempty"`
	ItemList 		[]*rssFeedItem	`xml:"channel>item,omitempty"`
}

type rssFeedItem struct {
	Author		string		`xml:"author"`
	Guid			string		`xml:"guid,omitempty"`
	Link	 		string		`xml:"link"`
	PubDate		string		`xml:"pubDate,omitempty"`
	Title 		string 		`xml:"title"`
	Enclosure 	struct {
		Url 		string   	`xml:"url,attr"`
		Length 	int32 	 	`xml:"length,attr"`
		Type 	string	 	`xml:"type,attr"`
	} `xml:"enclosure"`
}

type rssFeedImage struct {
	Url string `xml:"url"`
	Title string `xml:"title"`
	Link string `xml:"link"`
}

func getRtlLgti(w http.ResponseWriter, r *http.Request) {
	link := "http://www.rtl.fr/podcast/les-grosses-tetes.xml"
	rss := fetchRss(link)
	
	var items []*rssFeedItem
	for _, item := range rss.ItemList {
		match, _ := regexp.MatchString("int(e|é)gral(e.?).?", item.Title)
		if (match) {
			items = append(items, item)
		}
	}
	
	result := rssFeed{}
	result.XmlVersion = "2.0"
	result.XmlnsAtom = "http://www.w3.org/2005/Atom"
	result.XmlnsDcTerms = "http://purl.org/dc/terms/"
	result.XmlnsItunes = "http://www.itunes.com/dtds/podcast-1.0.dtd"
	result.XmlnsMedia = "http://search.yahoo.com/mrss/"
	result.Title = rss.Title
	result.Description = rss.Description
	result.Link = link
	result.Ttl = rss.Ttl
	result.Language = rss.Language
	result.ImageList = rss.ImageList
	result.ItemList = items
	
 	x, err := xml.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
    		return
  	}

  	w.Header().Set("Content-Type", "application/xml")
  	w.Write(x)
}

func fetchRss(uri string) rssFeed {
	resp, err := http.Get(uri)
	defer resp.Body.Close()
	
	if err != nil {
		log.Println(err)
	}
	res := rssFeed{}
	data, err := ioutil.ReadAll(resp.Body)
	xml.Unmarshal(data, &res)
	
	return res
}