package test

import (
	"encoding/json"
	"ips-lacpass-backend/pkg/utils"
	"testing"
)

func TestHCertICVPDecode(t *testing.T) {
	hcert := "HC1:6BFOXN%TSMAHN-HJM80DOO8W%TG34UE726*2OC9Y.TW1ANU9SCE7JM:UC*ELIQ5B264IM:/42JO2 7V35U:7+V4YC5/HQ6EOHCRBK81EPFJM5C9YCBJ%GBVCL+9-0G2PBUDBACARDAEI97KE*LHXQM.FDBIK4LD JM3.K/HLNOI3.KH+G7IKSH9NOIEJK5+K6IASD9YHI1KKK3MYII3IKEIAM0G6JK%86%X49/SQN4:U45ALD-4$XKHBTQ1LTA3$73HRJFRJ9STE-4/-KFU4-EF:57MUBMTF*MCXJL  RGBFH*RK%4U7U*+RDQJHY23QPX4MQ2S1$U4ST236MDNW*PGNETTU4DK/$TJ7PS4JLDV%0K1GDMDP $A*EK/JP:T3%.4OYB"

	expectedJSON := `{
  "-260": {
    "-6": {
      "dob": "1905-08-23",
      "n": "Aulo Agerio",
      "ndt": "NI",
      "nid": "16337361-9",
      "s": "male",
      "v": {
        "bo": "123123123",
        "dt": "2017-12-11",
        "vls": "2017-12-11",
        "vp": "YellowFeverProductd2c75a15ed309658b3968519ddb31690"
      }
    }
  },
  "1": "XCL",
  "6": 1757187943
}`

	decoded, err := utils.DecodeHCert(hcert)
	if err != nil {
		t.Fatalf("Failed to decode HCert: %v", err)
	}

	decodedPretty, err := json.MarshalIndent(decoded, "", "  ")
	if err != nil {
		t.Fatalf("Failed to decode JSON: %v", err)
	}

	if string(expectedJSON) != string(decodedPretty) {
		t.Errorf("Decoded value is not correct. Decoded: %s, Expected: %s", decodedPretty, expectedJSON)
	}
}

func TestHCertVHLDecode(t *testing.T) {
	hcert := "HC1:6BFOXNMG2N9HZBPYHQ3D69SO5D6%9L60JO DJS4L:P:R8LCDO%08JJG.NSOEV 9OG6%6Q4TJ7AJENS:NK7VCECM:MQ0FE%JC5Y479D/*8G.CV3NV3OVLD86J:KE2HF86GX2BTLHA9A86GNY8XOIROBZQMQOB9MEBED:KE87B MH:8DZYK%KNU9O%UL75E2*KH42$T8CRJ.V89:GF-K8JV Y8GJNKY8%97JR8ZV0:JVIP46+8KD35T8/Z8ZIV-YKAUVH40DQL2I6AI8LZP9WHHK5.SMIY9TO6YN6MJE2I6DF5P.P%OE-M6U JGKETW7YP6GUMY.HBNMAP50TBIM5GUMSWPVYB5RH+PEGKE5SG7UT4L5%K82OO-L8+$RTNKCZUN.DSB1971PFU%0F$5MH6QTMUEO1HB5*%L4NH7KEK%56VEUS17%E2F14LETP5 9VZ*MTJR.*U6.CH8795KTD8B836B4X/9+JIQT24GA-+DVE9B2K9FDJ4N172IM2%-2SFL -UNNF0GJG0AG16%$V%*C9:A8+I2QOHUQDVJ7VF +AU61$8IE0U4NOKIS1RE0BBSEWUVKI9K4/TQQP5U974CI9JQI10DEG30QUKL1"

	expectedJSON := `{
  "-260": {
    "5": [
      {
        "u": "shlink://eyJ1cmwiOiJodHRwOi8vbGFjcGFzcy5jcmVhdGUuY2w6ODE4Mi92Mi9tYW5pZmVzdHMvYjEzYzA0Y2QtMDc1Yy00YjY4LTgyOTQtMzJhZTMwN2YxYjA5IiwiZmxhZyI6IlAiLCJleHAiOjE3NjAxNDAyMjEwMDAsImtleSI6IkxTTnVaTXFHZEo1cmdQLUpJSEoySllLaWtuYzJXZDcwaG1VMFBSZFAwSHM9IiwibGFiZWwiOiJHREhDTiBWYWxpZGF0b3IifQ=="
      }
    ]
  },
  "1": "XJ",
  "4": 1760140221,
  "6": 1757271643804
}`

	decoded, err := utils.DecodeHCert(hcert)
	if err != nil {
		t.Fatalf("Failed to decode HCert: %v", err)
	}

	decodedPretty, err := json.MarshalIndent(decoded, "", "  ")
	if err != nil {
		t.Fatalf("Failed to decode JSON: %v", err)
	}

	if string(expectedJSON) != string(decodedPretty) {
		t.Errorf("Decoded value is not correct. Decoded: %s, Expected: %s", decodedPretty, expectedJSON)
	}
}