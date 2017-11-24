package main

import (
	"fmt"
	smpp "github.com/esazykin/smpp34"
)

func main() {
	// connect and bind
	transmitter := newTransmitter()

	// for latin and cyrillic symbols
	messageIds, err := sendSMS(transmitter, "7913...", "test")

	// Pdu gen errors
	if err != nil {
		fmt.Println("SubmitSm err:", err)
	}

	// Should save this to match with message_id
	fmt.Println(messageIds)

	for {
		pdu, err := transmitter.Read() // This is blocking
		if err != nil {
			fmt.Println("Read Err:", err)
			break
		}

		// EnquireLinks are auto handles
		switch pdu.GetHeader().Id {
		case smpp.SUBMIT_SM_RESP:
			// message_id should match this with seq message
			fmt.Println("MSG ID:", pdu.GetField("message_id").Value())
		default:
			// ignore all other PDUs or do what you link with them
			fmt.Println("PDU ID:", pdu.GetHeader().Id)
		}
	}

	fmt.Println("ending...")
}

func newTransmitter() *smpp.Transmitter {
	transmitter, err := smpp.NewTransmitter(
		"host",
		1234,
		5,
		smpp.Params{
			smpp.SYSTEM_TYPE: "",
			smpp.SYSTEM_ID:   "",
			smpp.PASSWORD:    "",
		},
	)
	if err != nil {
		panic(err)
	}

	return transmitter
}

func sendSMS(transmitter *smpp.Transmitter, phone, message string) ([]uint32, error) {
	pduParams := smpp.Params{
		smpp.DATA_CODING: smpp.ENCODING_ISO10646,
	}

	messageIds, err := transmitter.SubmitSm(
		"...",
		phone,
		message,
		&pduParams,
	)

	if err != nil {
		return []uint32{}, err
	}

	return messageIds, nil
}
