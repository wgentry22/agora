package api_test

import (
	"encoding/json"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/wgentry22/agora/modules/api"
	"github.com/wgentry22/agora/types/config"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
)

var _ = Describe("Router", func() {
	var (
		router api.Router
		ts     *httptest.Server
	)

	BeforeEach(func() {
		router = api.NewRouter(config.API{
			Port:       8123,
			PathPrefix: "/api",
		})
	})

	AfterEach(func() {
		ts.Close()
	})

	Context("without any registered controllers", func() {
		It("should serve the /info endpoint", func() {
			ts = httptest.NewServer(router.Handler())

			res, err := http.Get(fmt.Sprintf("%s/api/info/version", ts.URL))
			Expect(err).To(BeNil())
			Expect(res.StatusCode).To(Equal(http.StatusOK))

			data, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())

			var versionAsJSON struct {
				Version string `json:"version"`
			}

			err = json.Unmarshal(data, &versionAsJSON)
			Expect(err).To(BeNil())

			version, err := config.ParseSemanticVersion(versionAsJSON.Version)
			Expect(err).To(BeNil())

			Expect(version).ToNot(BeNil())

			Expect(*version).To(Equal(config.SemanticVersion{
				Major: 0,
				Minor: 0,
				Patch: 1,
			}))
		})
	})

	Context("with registered controllers", func() {
		It("should serve all endpoints", func() {
			router.Register(controller)
			ts = httptest.NewServer(router.Handler())

			url := func(uri string) string {
				return fmt.Sprintf("%s/api/test%s", ts.URL, uri)
			}

			By("serving GET all")
			getAllRes := makeRequest(url("/"), http.MethodGet)
			Expect(getAllRes.StatusCode).To(Equal(http.StatusOK))
			Expect(resData(getAllRes)).To(Equal("get all"))

			By("serving GET one")
			getOneRes := makeRequest(url("/1"), http.MethodGet)
			Expect(getOneRes.StatusCode).To(Equal(http.StatusOK))
			Expect(resData(getOneRes)).To(Equal("get 1"))

			By("serving POST")
			postRes := makeRequest(url("/"), http.MethodPost)
			Expect(postRes.StatusCode).To(Equal(http.StatusCreated))
			Expect(resData(postRes)).To(BeEmpty())

			By("serving PUT")
			putRes := makeRequest(url("/1"), http.MethodPut)
			Expect(putRes.StatusCode).To(Equal(http.StatusOK))
			Expect(resData(putRes)).To(Equal("put 1"))

			By("serving PATCH")
			patchRes := makeRequest(url("/1"), http.MethodPatch)
			Expect(patchRes.StatusCode).To(Equal(http.StatusOK))
			Expect(resData(patchRes)).To(Equal("patch 1"))

			By("serving DELETE")
			deleteRes := makeRequest(url("/1"), http.MethodDelete)
			Expect(deleteRes.StatusCode).To(Equal(http.StatusNoContent))
			Expect(resData(deleteRes)).To(BeEmpty())
		})
	})
})

func makeRequest(url, method string) *http.Response {
	req, err := http.NewRequest(method, url, strings.NewReader(""))
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	return res
}

func resData(response *http.Response) string {
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	return string(data)
}
