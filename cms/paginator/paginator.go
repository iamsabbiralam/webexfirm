package paginator

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
)

type Page struct {
	Order   int32
	URL     string
	Current bool
}

type Paginator struct {
	Prev         *Page
	Next         *Page
	Total        int32
	PerPage      int32
	TotalShowing int32
	CurrentPage  int32
	ShowingRange string
	Pages        []Page
}

func NewPaginator(currentPage, perPage int32, total int32, req *http.Request) Paginator {
	p := Paginator{}
	p.PerPage = perPage
	p.Total = total
	p.CurrentPage = currentPage
	lastPage := int32(math.Max(math.Ceil(float64(total)/float64(perPage)), 1))
	if lastPage > 1 {
		numberOfItems := currentPage * perPage
		p.TotalShowing = numberOfItems
		if numberOfItems > total {
			p.TotalShowing = total
		}
		p.ShowingRange = fmt.Sprintf("%d-%d", numberOfItems-perPage+1, p.TotalShowing)
	} else {
		p.ShowingRange = strconv.Itoa(int(total))
		p.TotalShowing = total
	}

	for i := int32(1); i <= lastPage; i++ {
		pageURL := *req.URL
		params := pageURL.Query()
		params.Set("page", strconv.Itoa(int(i)))
		pageURL.RawQuery = params.Encode()

		if i == 1 {
			p.Pages = append(p.Pages, Page{
				Order:   i,
				URL:     pageURL.String(),
				Current: i == currentPage,
			})

			if lastPage > 4 && currentPage > 3 {
				p.Pages = append(p.Pages, Page{})
			}
		} else if i == lastPage {
			if lastPage > 4 && currentPage+2 < lastPage {
				p.Pages = append(p.Pages, Page{})
			}
			p.Pages = append(p.Pages, Page{
				Order:   i,
				URL:     pageURL.String(),
				Current: i == currentPage,
			})
		} else if currentPage < 4 && i < 4 {
			p.Pages = append(p.Pages, Page{
				Order:   i,
				URL:     pageURL.String(),
				Current: i == currentPage,
			})
		} else if i+2 >= lastPage && currentPage+2 >= lastPage {
			p.Pages = append(p.Pages, Page{
				Order:   i,
				URL:     pageURL.String(),
				Current: i == currentPage,
			})
		} else if currentPage > 3 && currentPage+2 < lastPage && (currentPage == i+1 || currentPage == i || currentPage == i-1) {
			p.Pages = append(p.Pages, Page{
				Order:   i,
				URL:     pageURL.String(),
				Current: i == currentPage,
			})
		}

	}

	if currentPage != 1 {
		pageURL := *req.URL
		params := pageURL.Query()
		params.Set("page", strconv.Itoa(int(currentPage)-1))
		pageURL.RawQuery = params.Encode()
		p.Prev = &Page{URL: pageURL.String()}
	}

	if currentPage != lastPage {
		pageURL := *req.URL
		params := pageURL.Query()
		params.Set("page", strconv.Itoa(int(currentPage)+1))
		pageURL.RawQuery = params.Encode()
		p.Next = &Page{URL: pageURL.String()}
	}

	return p
}
