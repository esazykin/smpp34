package smpp34

import (
	"github.com/esazykin/smpp34/gsmutil"
	"math"
	"crypto/rand"
)

func splitShortMessage(shortMessage string, params *Params) ([][]byte, error) {
	dataCoding, ok := (*params)[DATA_CODING];
	if !ok {
		dataCoding = ENCODING_DEFAULT
	}

	var octetLimit int
	var message []byte
	switch dataCoding {
	case ENCODING_DEFAULT:
		octetLimit = 160
		message = []byte(shortMessage)
		break

	case ENCODING_ISO10646:
		octetLimit = 140
		message = gsmutil.EncodeUcs2(shortMessage)

	default:
		octetLimit = 254
		message = []byte(shortMessage)
	}

	messageLen := len(message)

	if messageLen > octetLimit {
		totalParts := byte(int(math.Ceil(float64(messageLen) / 134.0)))
		(*params)[ESM_CLASS] = ESM_CLASS_GSMFEAT_UDHI

		uid := make([]byte, 1)
		_, err := rand.Read(uid)
		if err != nil {
			return nil, err
		}

		partNum := 1
		parts := make([][]byte, 0)
		for i := 0; i < messageLen; i += 134 {
			start := i
			end := i + 134
			if end > messageLen {
				end = messageLen
			}
			part := []byte{0x05, 0x00, 0x03, uid[0], totalParts, byte(partNum)}
			part = append(part, message[start:end]...)
			partNum++

			parts = append(parts, part)
		}

		return parts, nil
	}

	return [][]byte{message}, nil
}
