package github


import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// SearchIssues queries the GitHub issue tracker.
func SearchIssues(terms []string) (*IssuesSearchResult, error) {
	q := url.QueryEscape(strings.Join(terms, " "))
	resp, err := http.Get(IssuesURL + "?q=" + q)
	if err != nil {
		return nil, err
	}

	// We must close resp.Body on all execution paths.
	// (Chapter 5 presents 'defer', which makes this simpler.)
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("search query failed: %s", resp.Status)
	}

	var result IssuesSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		resp.Body.Close()
		return nil, err
	}
	resp.Body.Close()

	classified(&result)
	
	return &result, nil
}

//练习 4.10： 修改issues程序，根据问题的时间进行分类，比如不到一个月的、不到一年的、超过一年
func classified(result *IssuesSearchResult) {

	result.Dict = make(map[string][]*Issue)

	result.Dict[Ltm] = []*Issue{}
	result.Dict[Lty] = []*Issue{}
	result.Dict[Mty] = []*Issue{}

	if len(result.Items) <= 0 {
		return
	}

	k := time.Now()
	mm := k.Add((-30 * 24 * time.Hour))
	yy := k.Add((-365 * 24 * time.Hour))

	arr1 := []*Issue{}//不超过一个月
	arr2 := []*Issue{}//不超过一年
	arr3 := []*Issue{}//一年以上

	for _, item := range result.Items {

		if mm.Before(item.CreatedAt) {
			arr1 = append(arr1, item)
		}else if yy.After(item.CreatedAt) {
			arr3 = append(arr3, item)
		}else {
			arr2 = append(arr2, item)
		}
	}

	result.Dict[Ltm] = arr1
	result.Dict[Lty] = arr2
	result.Dict[Mty] = arr3
}