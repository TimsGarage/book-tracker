package bookconnector

import (
	"bookend/internal/models"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Intermediate XML structs for DNB MARC21-xml response
type dnbSRUResponse struct {
	XMLName         xml.Name `xml:"searchRetrieveResponse"`
	NumberOfRecords int      `xml:"numberOfRecords"`
	Records         []struct {
		Data dnbRecord `xml:"recordData>record"`
	} `xml:"records>record"`
}

type dnbRecord struct {
	DataFields []dnbDataField `xml:"datafield"`
}

type dnbDataField struct {
	Tag       string        `xml:"tag,attr"`
	Subfields []dnbSubfield `xml:"subfield"`
}

type dnbSubfield struct {
	Code  string `xml:"code,attr"`
	Value string `xml:",chardata"`
}

func (r *dnbRecord) getSubfield(tag, code string) string {
	for _, df := range r.DataFields {
		if df.Tag == tag {
			for _, sf := range df.Subfields {
				if sf.Code == code {
					return strings.TrimSpace(sf.Value)
				}
			}
		}
	}
	return ""
}

// Helper to extract an integer page count from bibliographic text (e.g. "320 S.", "xiv, 250 p.")
func parsePageCount(raw string) int {
	re := regexp.MustCompile(`\d+`)
	match := re.FindString(raw)
	if match == "" {
		return 0
	}
	pages, _ := strconv.Atoi(match)
	return pages
}

func DnbLookupIsbn(isbn string) BookLookupResponse {
	cleanISBN := strings.ReplaceAll(strings.ReplaceAll(isbn, "-", ""), " ", "")

	endpoint := "https://services.dnb.de/sru/dnb"
	params := url.Values{}
	params.Add("version", "1.1")
	params.Add("operation", "searchRetrieve")
	params.Add("query", fmt.Sprintf("num=%s", cleanISBN))
	params.Add("recordSchema", "MARC21-xml")
	params.Add("maximumRecords", "1")

	reqURL := fmt.Sprintf("%s?%s", endpoint, params.Encode())

	reqClient := http.Client{
		Timeout: time.Second * 15,
	}

	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return BookLookupResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to create request",
			Book:       models.LookupBook{},
			Success:    false,
		}
	}

	// DNB recommends identifying your client
	req.Header.Set("User-Agent", "BookendApp/1.0 (contact: info@example.com)")

	res, getErr := reqClient.Do(req)
	if getErr != nil {
		log.Printf("Error requesting DNB SRU API: %v", getErr)
		return BookLookupResponse{
			StatusCode: http.StatusBadGateway,
			Message:    "Failed to contact Deutsche Nationalbibliothek API",
			Book:       models.LookupBook{},
			Success:    false,
		}
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Printf("DNB API returned status code: %d", res.StatusCode)
		return BookLookupResponse{
			StatusCode: res.StatusCode,
			Message:    fmt.Sprintf("DNB API returned HTTP %d", res.StatusCode),
			Book:       models.LookupBook{},
			Success:    false,
		}
	}

	var sru dnbSRUResponse
	if err := xml.NewDecoder(res.Body).Decode(&sru); err != nil {
		log.Printf("Error decoding XML: %v", err)
		return BookLookupResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to parse API response",
			Book:       models.LookupBook{},
			Success:    false,
		}
	}

	if sru.NumberOfRecords == 0 || len(sru.Records) == 0 {
		return BookLookupResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Book not found",
			Book:       models.LookupBook{},
			Success:    false,
		}
	}

	rec := sru.Records[0].Data

	// Title assembly (Main Title 245$a + Subtitle 245$b)
	title := rec.getSubfield("245", "a")
	if sub := rec.getSubfield("245", "b"); sub != "" {
		title = strings.TrimRight(title, " /:") + ": " + sub
	}
	title = strings.TrimRight(title, " /:")

	// Author fallback: Primary creator (100$a) -> Added entry (700$a)
	author := rec.getSubfield("100", "a")
	if author == "" {
		author = rec.getSubfield("700", "a")
	}
	author = strings.TrimRight(author, ", /")

	// Publisher & Release year fallback: RDA (264) -> AACR2 (260)
	publisher := rec.getSubfield("264", "b")
	if publisher == "" {
		publisher = rec.getSubfield("260", "b")
	}
	publisher = strings.TrimRight(publisher, ", ;/")

	releaseYear := rec.getSubfield("264", "c")
	if releaseYear == "" {
		releaseYear = rec.getSubfield("260", "c")
	}
	releaseYear = strings.Trim(releaseYear, " .[]c")

	// Physical description (e.g., "320 S. : Ill. ; 21 cm") -> integer pages
	pagesRaw := rec.getSubfield("300", "a")
	pageCount := parsePageCount(pagesRaw)

	// Since DNB does not provide cover images via MARC21 XML,
	// resolve directly from Open Library's deterministic static image endpoint.
	thumbnailURL := fmt.Sprintf("https://covers.openlibrary.org/b/isbn/%s-L.jpg", cleanISBN)

	book := models.LookupBook{
		Isbn:          cleanISBN,
		Title:         title,
		Author:        author,
		Release:       releaseYear,
		Publisher:     publisher,
		ThumbnailLink: thumbnailURL,
		Pages:         pageCount,
	}

	return BookLookupResponse{
		StatusCode: http.StatusOK,
		Message:    "",
		Book:       book,
		Success:    true,
	}
}

// DnbSearchBooks queries DNB using a free-text search across all catalog fields.
// limit controls how many records to return (DNB max per page is 100).
func DnbSearchBooks(searchTerm string, limit int) BookSearchResponse {
	trimmedTerm := strings.TrimSpace(searchTerm)
	if trimmedTerm == "" {
		return BookSearchResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Search term cannot be empty",
			Books:      []models.LookupBook{},
			TotalFound: 0,
			Success:    false,
		}
	}

	if limit <= 0 {
		limit = 10
	} else if limit > 100 {
		limit = 100
	}

	// Escape any double quotes in user input to keep CQL valid
	safeTerm := strings.ReplaceAll(trimmedTerm, `"`, `\"`)
	// 'woe' searches across all catalog fields (words everywhere)
	cqlQuery := fmt.Sprintf(`woe="%s"`, safeTerm)

	endpoint := "https://services.dnb.de/sru/dnb"
	params := url.Values{}
	params.Add("version", "1.1")
	params.Add("operation", "searchRetrieve")
	params.Add("query", cqlQuery)
	params.Add("recordSchema", "MARC21-xml")
	params.Add("maximumRecords", strconv.Itoa(limit))

	reqURL := fmt.Sprintf("%s?%s", endpoint, params.Encode())

	reqClient := http.Client{
		Timeout: time.Second * 15,
	}

	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return BookSearchResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to create request",
			Books:      []models.LookupBook{},
			Success:    false,
		}
	}

	req.Header.Set("User-Agent", "BookendApp/1.0 (contact: info@example.com)")

	res, getErr := reqClient.Do(req)
	if getErr != nil {
		log.Printf("Error requesting DNB SRU API: %v", getErr)
		return BookSearchResponse{
			StatusCode: http.StatusBadGateway,
			Message:    "Failed to contact Deutsche Nationalbibliothek API",
			Books:      []models.LookupBook{},
			Success:    false,
		}
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Printf("DNB API returned status code: %d", res.StatusCode)
		return BookSearchResponse{
			StatusCode: res.StatusCode,
			Message:    fmt.Sprintf("DNB API returned HTTP %d", res.StatusCode),
			Books:      []models.LookupBook{},
			Success:    false,
		}
	}

	var sru dnbSRUResponse
	if err := xml.NewDecoder(res.Body).Decode(&sru); err != nil {
		log.Printf("Error decoding XML: %v", err)
		return BookSearchResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to parse API response",
			Books:      []models.LookupBook{},
			Success:    false,
		}
	}

	if sru.NumberOfRecords == 0 || len(sru.Records) == 0 {
		return BookSearchResponse{
			StatusCode: http.StatusOK,
			Message:    "No books found matching search criteria",
			Books:      []models.LookupBook{},
			TotalFound: 0,
			Success:    true,
		}
	}

	books := make([]models.LookupBook, 0, len(sru.Records))

	for _, item := range sru.Records {
		rec := item.Data

		// Title assembly (Main Title 245$a + Subtitle 245$b)
		title := rec.getSubfield("245", "a")
		if sub := rec.getSubfield("245", "b"); sub != "" {
			title = strings.TrimRight(title, " /:") + ": " + sub
		}
		title = strings.TrimRight(title, " /:")

		// Author fallback: Primary creator (100$a) -> Added entry (700$a)
		author := rec.getSubfield("100", "a")
		if author == "" {
			author = rec.getSubfield("700", "a")
		}
		author = strings.TrimRight(author, ", /")

		// Publisher & Release year fallback: RDA (264) -> AACR2 (260)
		publisher := rec.getSubfield("264", "b")
		if publisher == "" {
			publisher = rec.getSubfield("260", "b")
		}
		publisher = strings.TrimRight(publisher, ", ;/")

		releaseYear := rec.getSubfield("264", "c")
		if releaseYear == "" {
			releaseYear = rec.getSubfield("260", "c")
		}
		releaseYear = strings.Trim(releaseYear, " .[]c")

		// Extract ISBN from 020$a (may contain suffixes like "(hardcover)")
		isbnRaw := rec.getSubfield("020", "a")
		cleanISBN := ""
		if parts := strings.Fields(isbnRaw); len(parts) > 0 {
			cleanISBN = strings.ReplaceAll(parts[0], "-", "")
		}

		// Pages
		pagesRaw := rec.getSubfield("300", "a")
		pageCount := parsePageCount(pagesRaw)

		// Cover image link (only generated if an ISBN exists)
		thumbnailURL := ""
		if cleanISBN != "" {
			thumbnailURL = fmt.Sprintf("https://covers.openlibrary.org/b/isbn/%s-L.jpg", cleanISBN)
		}

		books = append(books, models.LookupBook{
			Isbn:          cleanISBN,
			Title:         title,
			Author:        author,
			Release:       releaseYear,
			Publisher:     publisher,
			ThumbnailLink: thumbnailURL,
			Pages:         pageCount,
		})
	}

	return BookSearchResponse{
		StatusCode: http.StatusOK,
		Message:    "",
		Books:      books,
		TotalFound: sru.NumberOfRecords,
		Success:    true,
	}
}
