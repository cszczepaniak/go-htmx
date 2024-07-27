package httpwrap

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/shoenig/test"
	"github.com/shoenig/test/must"
)

func TestUnmarshal_Strings(t *testing.T) {
	// Set dummy query data in URL
	req := newRequest(t, "?userName=spongebob")

	// Set dummy data in form
	req.Request.Form = url.Values{}
	req.Request.Form.Set("userID", "abc123")

	// Set dummy path value
	req.Request.SetPathValue("userEmail", "email@abc.com")

	type testData struct {
		UserID    string `req:"form:userID"`
		UserEmail string `req:"path:userEmail"`
		UserName  string `req:"query:userName"`
	}

	var data testData
	err := req.Unmarshal(&data)
	must.NoError(t, err)

	test.Eq(t, "abc123", data.UserID)
	test.Eq(t, "email@abc.com", data.UserEmail)
	test.Eq(t, "spongebob", data.UserName)

	tests := []struct {
		desc      string
		getReq    func(t testing.TB) Request
		expErrMsg string
	}{{
		desc: "missing form data",
		getReq: func(t testing.TB) Request {
			return newRequest(t, "")
		},
		expErrMsg: `field "formData" is required but was empty`,
	}, {
		desc: "missing path data",
		getReq: func(t testing.TB) Request {
			req := newRequest(t, "")
			req.Request.Form = url.Values{}
			req.Request.Form.Set("formData", "a")
			return req
		},
		expErrMsg: `field "pathData" is required but was empty`,
	}, {
		desc: "missing query data",
		getReq: func(t testing.TB) Request {
			req := newRequest(t, "")
			req.Request.Form = url.Values{}
			req.Request.Form.Set("formData", "a")

			req.Request.SetPathValue("pathData", "b")
			return req
		},
		expErrMsg: `field "queryData" is required but was empty`,
	}}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			type validationTestData struct {
				FormData  string `req:"form:formData,required"`
				PathData  string `req:"path:pathData,required"`
				QueryData string `req:"query:queryData,required"`
			}

			req = tc.getReq(t)

			err = req.Unmarshal(new(validationTestData))
			test.ErrorContains(t, err, tc.expErrMsg)
			test.Eq(t, http.StatusBadRequest, StatusCodeForError(err))
		})
	}
}

func TestUnmarshal_Ints(t *testing.T) {
	// Set dummy query data in URL
	req := newRequest(t, "?int3=789")

	// Set dummy data in form
	req.Request.Form = url.Values{}
	req.Request.Form.Set("int1", "123")

	// Set dummy path value
	req.Request.SetPathValue("int2", "456")

	type testData struct {
		Int1 int `req:"form:int1"`
		Int2 int `req:"path:int2"`
		Int3 int `req:"query:int3"`
	}

	var data testData
	err := req.Unmarshal(&data)
	must.NoError(t, err)

	test.Eq(t, 123, data.Int1)
	test.Eq(t, 456, data.Int2)
	test.Eq(t, 789, data.Int3)

	tests := []struct {
		desc      string
		getReq    func(t testing.TB) Request
		expErrMsg string
		expCode   int
	}{{
		desc: "empty int",
		getReq: func(t testing.TB) Request {
			return newRequest(t, "")
		},
		expErrMsg: `strconv.ParseInt: parsing "": invalid syntax`,
		expCode:   http.StatusBadRequest,
	}, {
		desc: "unparseable int",
		getReq: func(t testing.TB) Request {
			req := newRequest(t, "")
			req.Request.Form = url.Values{}
			req.Request.Form.Set("formData", "a")
			return req
		},
		expErrMsg: `strconv.ParseInt: parsing "a": invalid syntax`,
		expCode:   http.StatusBadRequest,
	}, {
		desc: "invalid validation",
		getReq: func(t testing.TB) Request {
			req := newRequest(t, "")
			req.Request.Form = url.Values{}
			req.Request.Form.Set("formData", "123")

			req.Request.SetPathValue("pathData", "123")
			return req
		},
		expErrMsg: `required option not supported for type int`,
		expCode:   http.StatusInternalServerError,
	}}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			type validationTestData struct {
				FormData int `req:"form:formData"`
				PathData int `req:"path:pathData,required"`
			}

			req = tc.getReq(t)

			err = req.Unmarshal(new(validationTestData))
			test.ErrorContains(t, err, tc.expErrMsg)
			test.Eq(t, tc.expCode, StatusCodeForError(err))
		})
	}
}

func newRequest(t testing.TB, url string) Request {
	t.Helper()

	r, err := http.NewRequest(http.MethodGet, url, nil)
	must.NoError(t, err)

	return Request{
		Request: r,
	}
}
