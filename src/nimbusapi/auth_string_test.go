package nimbusapi

import (
	"testing"
)

type testEntry struct {
	credentials Credentials
	timestamp   int64
	method      string
	uri         string
	authString  string
}

// test data generated from python library lumberyard
var testData = []testEntry{
	{
		Credentials{
			"motoboto-benchmark-000",
			2,
			[]byte("iKn/OxpggHSXzB0oUAihTMTf+n6Bsyywwm3bXMQfdKo"),
		},
		1344970933,
		"HEAD",
		"/data/xxx",
		"NIMBUS.IO 2:c4d946d1089f11d310ae8bd8f52501212a3eb50cc800b0688d3df7bf65f04ce7",
	},
	{
		Credentials{
			"motoboto-benchmark-001",
			3,
			[]byte("WoCsHcHGvm0gf1wzLAeuArC1VmdFTxo9Sd5S4Ag6SjQ"),
		},
		1344970934,
		"HEAD",
		"/data/xxx",
		"NIMBUS.IO 3:8b26220618f8da55f371d3108221f4fd982f4e61b02d8eb00ee8cc3c50ca2707",
	},
	{
		Credentials{
			"motoboto-benchmark-002",
			4,
			[]byte("FjzbqM8bEuZ3Hw7cA4XaiUZgJxKlQHKluL8ovtpxTbE"),
		},
		1344970935,
		"HEAD",
		"/data/xxx",
		"NIMBUS.IO 4:45e71c4839d68d2e2c67c80b916e1cb3e75abd6bc84c4499c00a770dbaf62785",
	},
	{
		Credentials{
			"motoboto-benchmark-003",
			5,
			[]byte("9dJwBh1YWNlrs31F+Fx2BY2KZUtFYmBH3hmocH+6Ggk"),
		},
		1344970936,
		"POST",
		"/data/xxx",
		"NIMBUS.IO 5:015c53c88a0531ef1f998fcfee01612627de2c9de8cbaa76d1c99b0de0f32d6d",
	},
	{
		Credentials{
			"motoboto-benchmark-004",
			6,
			[]byte("Vo+YJDRVmLNCKzy20nFuxLZQv9ohpq2FfFfZXdPk2kk"),
		},
		1344970937,
		"POST",
		"/data/xxx",
		"NIMBUS.IO 6:ab6b13b21230e142c5ca8de245953001a6877e434c07ff991074ff8ec95f487c",
	},
	{
		Credentials{
			"motoboto-benchmark-005",
			7,
			[]byte("n0OClxSMaDZhqqhPTUfaIxyysQ1PHHj8O4dWPVmEKNk"),
		},
		1344970938,
		"POST",
		"/data/xxx",
		"NIMBUS.IO 7:87949c05b6b566332a2d42d8de66b47b9828d7e8ebedceac1bd128c23b268100",
	},
	{
		Credentials{
			"motoboto-benchmark-006",
			8,
			[]byte("nziQgsJ4mbuxCKfA2HJgFoabdhg2A5b8RHOVkOOfLVo"),
		},
		1344970939,
		"GET",
		"/data/xxx",
		"NIMBUS.IO 8:317af9956227a5f3ec1f86579df48a1c85dca470bb8cc999cf686b7dd030935d",
	},
	{
		Credentials{
			"motoboto-benchmark-007",
			9,
			[]byte("laHX0cXLqAU8b2OVvNlAATagwBeQ/zuFF+tDi7QLfSc"),
		},
		1344970940,
		"GET",
		"/data/xxx",
		"NIMBUS.IO 9:23cfee26b7329e0f672c1b77336f2b3084fc0c0e4cefac88391cd124fd3948d2",
	},
	{
		Credentials{
			"motoboto-benchmark-008",
			10,
			[]byte("uBgJfTKOIdh6mezC4T+6Dk8AKLfDxtphKjqLHXm5MmY"),
		},
		1344970941,
		"POST",
		"/data/xxx",
		"NIMBUS.IO 10:79a8ca622bcba5999e4a63adccb5383f7b850b187cb68b9ba79da927612f8660",
	},
	{
		Credentials{
			"motoboto-benchmark-009",
			11,
			[]byte("CJiKEHRPMdwWaiwDB157bwW/piEyoZCuysXdyKvFPCk"),
		},
		1344970942,
		"GET",
		"/data/xxx",
		"NIMBUS.IO 11:629fc4344699c0c3cd04dcd318d10ce32ba36541539ece6bc6dba116519b81b3",
	},
}

func TestAuthString(t *testing.T) {
	for _, entry := range testData {
		authString, err:= ComputeAuthString(&entry.credentials, entry.method,
			entry.timestamp, entry.uri)
		if err != nil {
			t.Fatalf("%s, %s, %s", authString, entry.authString, err)
		}
		if authString != entry.authString {
			t.Fatalf("%s != %s, %v", authString, entry.authString, entry)
		}
	}
}
