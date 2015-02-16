package cf

import (
	"net/http"
	"os"
	"testing"

	"github.com/getlantern/cloudflare"

	"github.com/getlantern/fdcount"
	"github.com/getlantern/testify/assert"
)

func TestAll(t *testing.T) {
	_, counter, err := fdcount.Matching("TCP")
	if err != nil {
		t.Fatalf("Unable to get starting fdcount: %v", err)
	}

	u := New("getiantem.org", os.Getenv("CF_USER"), os.Getenv("CF_API_KEY"))
	u.Client.Http.Transport = &http.Transport{
		DisableKeepAlives: true,
	}
	recs, err := u.GetAllRecords()
	if assert.NoError(t, err, "Should be able to get all records") {
		for _, r := range recs {
			log.Tracef("%v : %v", r.Domain, r.Value)
		}
		assert.True(t, len(recs) > 0, "There should be some records")
	}

	assert.NoError(t, counter.AssertDelta(0), "All file descriptors should have been closed")
}

func TestEnsureRegistered(t *testing.T) {
	_, counter, err := fdcount.Matching("TCP")
	if err != nil {
		t.Fatalf("Unable to get starting fdcount: %v", err)
	}

	u := New("getiantem.org", os.Getenv("CF_USER"), os.Getenv("CF_API_KEY"))
	u.Client.Http.Transport = &http.Transport{
		DisableKeepAlives: true,
	}

	// Test with no existing record
	name, ip := "cf-test-entry", "127.0.0.1"
	rec, proxying, err := u.EnsureRegistered(name, ip, nil)
	if assert.NoError(t, err, "Should be able to register with no record") {
		assert.NotNil(t, rec, "A new record should have been returned")
		assert.True(t, proxying, "Proxying (orange cloud) should be on")
	}

	// Test with existing record, but not passing it in
	rec, proxying, err = u.EnsureRegistered(name, ip, nil)
	if assert.NoError(t, err, "Should be able to register with unspecified existing record") {
		assert.NotNil(t, rec, "Existing record should have been returned")
		assert.True(t, proxying, "Proxying (orange cloud) should be on")

		// Test with existing record, passing it in
		rec, proxying, err = u.EnsureRegistered(name, ip, rec)
		if assert.NoError(t, err, "Should be able to register with specified existing record") {
			assert.NotNil(t, rec, "Existing record should have been returned")
			assert.True(t, proxying, "Proxying (orange cloud) should be on")
		}
	}

	if rec != nil {
		err := u.DestroyRecord(rec)
		assert.NoError(t, err, "Should be able to destroy record")
	}

	assert.NoError(t, counter.AssertDelta(0), "All file descriptors should have been closed")
}

func doTestEnsureRegistered(t *testing.T, rec *cloudflare.Record) *cloudflare.Record {

	return rec
}
