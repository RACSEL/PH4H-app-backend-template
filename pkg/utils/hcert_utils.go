package utils

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/fxamacker/cbor/v2"
	"github.com/veraison/go-cose"
)

const base45Alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ $%*+-./:"

var base45reverse map[byte]int

func init() {
	base45reverse = make(map[byte]int)
	for i, c := range []byte(base45Alphabet) {
		base45reverse[c] = i
	}
}

func convertInterfaceMap(m map[interface{}]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range m {
		var keyStr string
		switch kType := k.(type) {
		case string:
			keyStr = kType
		case float64:
			keyStr = strconv.FormatFloat(kType, 'f', -1, 64)
		case int:
			keyStr = strconv.Itoa(kType)
		case uint64:
			keyStr = strconv.FormatUint(kType, 10)
		case int64:
			keyStr = strconv.FormatInt(kType, 10)
		default:
			fmt.Printf("Encountered unexpected key type: %T with value: %v\n", kType, k)
			continue
		}

		// Check if the value is a nested map
		if nestedMap, isMap := v.(map[interface{}]interface{}); isMap {
			result[keyStr] = convertInterfaceMap(nestedMap)
		} else if nestedSlice, isSlice := v.([]interface{}); isSlice {
			result[keyStr] = convertInterfaceSlice(nestedSlice)
		} else {
			result[keyStr] = v
		}
	}
	return result
}

func convertInterfaceSlice(s []interface{}) []interface{} {
	result := make([]interface{}, len(s))
	for i, v := range s {
		if nestedMap, isMap := v.(map[interface{}]interface{}); isMap {
			result[i] = convertInterfaceMap(nestedMap)
		} else if nestedSlice, isSlice := v.([]interface{}); isSlice {
			result[i] = convertInterfaceSlice(nestedSlice)
		} else {
			result[i] = v
		}
	}
	return result
}

func decodeBase45(s string) ([]byte, error) {
	var out []byte
	for i := 0; i < len(s); {
		if len(s)-i < 2 {
			return nil, fmt.Errorf("invalid base45 string length")
		}
		if len(s)-i < 3 {
			// 2-character case
			c1, ok1 := base45reverse[s[i]]
			c2, ok2 := base45reverse[s[i+1]]
			if !ok1 || !ok2 {
				return nil, fmt.Errorf("invalid character in base45 string")
			}
			val := c1 + c2*45
			if val > 255 {
				return nil, fmt.Errorf("invalid 2-character encoding")
			}
			out = append(out, byte(val))
			i += 2
		} else {
			// 3-character case
			c1, ok1 := base45reverse[s[i]]
			c2, ok2 := base45reverse[s[i+1]]
			c3, ok3 := base45reverse[s[i+2]]
			if !ok1 || !ok2 || !ok3 {
				return nil, fmt.Errorf("invalid character in base45 string")
			}
			val := c1 + c2*45 + c3*45*45
			if val > 65535 {
				return nil, fmt.Errorf("invalid 3-character encoding")
			}
			out = append(out, byte(val>>8), byte(val&0xFF))
			i += 3
		}
	}
	return out, nil
}

// DecodeHCert decodes a base45 string encoded using Zlib/COSE/CBOR pipeline.
// The string is an HCERT so it starts with "HC1:".
func DecodeHCert(hcert string) (map[string]interface{}, error) {
	if !strings.HasPrefix(hcert, "HC1:") {
		return nil, fmt.Errorf("invalid HCERT prefix")
	}

	// 1. Remove "HC1:" prefix
	base45Encoded := strings.TrimPrefix(hcert, "HC1:")

	// 2. Base45 decode
	compressedCose, err := decodeBase45(base45Encoded)
	if err != nil {
		return nil, fmt.Errorf("base45 decoding failed: %w", err)
	}

	// 3. Zlib decompress
	var cosePayload []byte
	r, err := zlib.NewReader(bytes.NewReader(compressedCose))
	if err != nil {
		// Not zlib compressed, use as is
		cosePayload = compressedCose
	} else {
		decompressed, err := io.ReadAll(r)
		if err != nil {
			return nil, fmt.Errorf("zlib decompression failed: %w", err)
		}
		cosePayload = decompressed
		defer r.Close()
	}

	// 4. COSE decode
	var payload []byte
	var msg cose.Sign1Message
	if err := msg.UnmarshalCBOR(cosePayload); err == nil {
		payload = msg.Payload
	} else {
		var signMessage cose.SignMessage
		if err2 := signMessage.UnmarshalCBOR(cosePayload); err2 == nil {
			payload = signMessage.Payload
		} else {
			// Both failed, try manual parsing
			var rawCoseMessage interface{}
			if err3 := cbor.Unmarshal(cosePayload, &rawCoseMessage); err3 != nil {
				return nil, fmt.Errorf("cose unmarshalling failed for both Sign1Message and SignMessage and raw: %v, %v, %v", err, err2, err3)
			}

			coseArray, ok := rawCoseMessage.([]interface{})
			if !ok {
				if tagged, ok := rawCoseMessage.(cbor.Tag); ok {
					if content, ok := tagged.Content.([]interface{}); ok {
						coseArray = content
					}
				}
			}

			if coseArray == nil || len(coseArray) < 3 {
				return nil, fmt.Errorf("cose unmarshalling failed for both Sign1Message and SignMessage: %v, %v", err, err2)
			}

			p, ok := coseArray[2].([]byte)
			if !ok {
				return nil, fmt.Errorf("failed to extract payload from cose message: payload is not a byte string")
			}
			payload = p
		}
	}

	// 5. CBOR decode the payload from COSE message
	var cborData map[interface{}]interface{}
	if err := cbor.Unmarshal(payload, &cborData); err != nil {
		return nil, fmt.Errorf("cbor unmarshalling of payload failed: %w", err)
	}

	cborDataConverted := convertInterfaceMap(cborData)
	return cborDataConverted, nil
}
