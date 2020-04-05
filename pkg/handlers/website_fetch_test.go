package handlers

import "testing"

func TestParseAnnotationData(t *testing.T) {
	exampleAnnotationContent := "joke=curl-a-joke.herokuapp.com"

	fetchReq, err := parseAnnotationData(exampleAnnotationContent)
	if err != nil {
		t.Logf("failed to parse the annotation data: %s\n", err.Error())
		t.FailNow()
	}

	if fetchReq.IntoKey != "joke" {
		t.FailNow()
	}

	if fetchReq.FromSite.Host != "curl-a-joke.herokuapp.com" {
		t.FailNow()
	}
}

func TestFetchSiteData(t *testing.T) {
	fRequest, _ := parseAnnotationData("joke=curl-a-joke.herokuapp.com")
	fResp, err := fetchSiteData(fRequest)
	if err != nil {
		t.Logf("received error from fetchSiteData: %s\n", err.Error())
		t.FailNow()
	}

	if fResp.Key != fRequest.IntoKey {
		t.FailNow()
	}

	if len(fResp.Value) == 0 {
		t.FailNow()
	}
	t.Logf("fResp: %s\n", fResp)
}
