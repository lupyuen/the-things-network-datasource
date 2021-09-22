package plugin

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/fxamacker/cbor/v2"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/mqtt-datasource/pkg/mqtt"
)

func ToFrame(topic string, messages []mqtt.Message) *data.Frame {
	log.DefaultLogger.Debug(fmt.Sprintf("ToFrame: topic=%s", topic))

	count := len(messages)
	if count > 0 {
		first := messages[0].Value
		if strings.HasPrefix(first, "{") {
			return jsonMessagesToFrame(topic, messages)
		}
	}

	// Fall through to expecting values
	timeField := data.NewFieldFromFieldType(data.FieldTypeTime, count)
	timeField.Name = "Time"
	valueField := data.NewFieldFromFieldType(data.FieldTypeFloat64, count)
	valueField.Name = "Value"

	for idx, m := range messages {
		if value, err := strconv.ParseFloat(m.Value, 64); err == nil {
			timeField.Set(idx, m.Timestamp)
			valueField.Set(idx, value)
		}
	}

	return data.NewFrame(topic, timeField, valueField)
}

//  Transform the array of MQTT Messages (JSON encoded) into a Grafana Data Frame
func jsonMessagesToFrame(topic string, messages []mqtt.Message) *data.Frame {
	log.DefaultLogger.Debug(fmt.Sprintf("jsonMessagesToFrame: topic=%s", topic))

	count := len(messages)
	if count == 0 {
		log.DefaultLogger.Debug(fmt.Sprintf("jsonMessagesToFrame: No msgs for topic=%s", topic))
		return nil
	}
	log.DefaultLogger.Debug(fmt.Sprintf("jsonMessagesToFrame: msg=%s", messages[0].Value))

	{
		//  Handle the first message. See sample messages: https://github.com/lupyuen/the-things-network-datasource#mqtt-log
		//  Deserialise the message body to a map of String -> interface{}
		frame := data.NewFrame(topic)
		var body map[string]interface{}
		err := json.Unmarshal([]byte(messages[0].Value), &body)
		if err != nil {
			s := fmt.Sprintf("error unmarshalling json message: %s", err.Error())
			frame.AppendNotices(data.Notice{Severity: data.NoticeSeverityError, Text: s})
			log.DefaultLogger.Debug(s)
			return frame
		}

		//  Get the Uplink Message
		if body["uplink_message"] == nil {
			s := "uplink_message missing"
			frame.AppendNotices(data.Notice{Severity: data.NoticeSeverityError, Text: s})
			log.DefaultLogger.Debug(s)
			return frame
		}
		uplink_message := body["uplink_message"].(map[string]interface{})

		//  Get the Payload
		if uplink_message["frm_payload"] == nil {
			s := "frm_payload missing"
			frame.AppendNotices(data.Notice{Severity: data.NoticeSeverityError, Text: s})
			log.DefaultLogger.Debug(s)
			return frame
		}
		frm_payload := uplink_message["frm_payload"].(string)

		//  Base64 decode the Payload
		payload, err := base64.StdEncoding.DecodeString(frm_payload)
		if err != nil {
			s := fmt.Sprintf("Base64 decode failed: %s", err.Error())
			frame.AppendNotices(data.Notice{Severity: data.NoticeSeverityError, Text: s})
			log.DefaultLogger.Debug(s)
			return frame
		}
		log.DefaultLogger.Debug(fmt.Sprintf("payload: %v", payload))

		//  TODO: Decode the payload with CBOR
		//  TODO: Add the decoded fields to the frame
	}

	{
		//  CBOR encoding for {"t": 1234}.  See http://cbor.me/
		payload := []byte{0xA1, 0x61, 0x74, 0x19, 0x04, 0xD2}
		frame := data.NewFrame(topic)

		//  Decode CBOR payload to a map of String -> interface{}
		var body map[string]interface{}
		err := cbor.Unmarshal(payload, &body)
		if err != nil {
			s := fmt.Sprintf("CBOR decode failed: %s", err.Error())
			frame.AppendNotices(data.Notice{Severity: data.NoticeSeverityError, Text: s})
			log.DefaultLogger.Debug(s)
			return frame
		}
		//  Shows: map[t:1234]
		log.DefaultLogger.Debug(fmt.Sprintf("CBOR decoded: %v", body))
	}

	//  Sample body
	body := map[string]float64{}
	body["t"] = 123.0

	timeField := data.NewFieldFromFieldType(data.FieldTypeTime, count)
	timeField.Name = "Time"
	timeField.SetConcrete(0, messages[0].Timestamp)

	// Create a field for each key and set the first value
	keys := make([]string, 0, len(body))
	fields := make(map[string]*data.Field, len(body))
	for key, val := range body {
		field := data.NewFieldFromFieldType(data.FieldTypeNullableFloat64, count)
		field.Name = key
		field.SetConcrete(0, val)
		fields[key] = field
		keys = append(keys, key)
	}
	sort.Strings(keys) // keys stable field order.

	//  TODO: Handle the messages after the first one
	//  We might not need this because The Things Network only supports low-volume messaging
	for row, m := range messages {
		if row == 0 {
			continue
		}

		/*
			//  Deserialise the message body to a map of String -> Float
			err := json.Unmarshal([]byte(m.Value), &body)
			if err != nil {
				log.DefaultLogger.Debug(fmt.Sprintf("error unmarshalling json message: %s", err.Error()))
				continue // bad row?
			}
		*/

		//  TODO: Sample body to indicate that message was dropped
		body = map[string]float64{}
		body["dropped"] = 1.0

		timeField.SetConcrete(row, m.Timestamp)
		for key, val := range body {
			field, ok := fields[key]
			if ok {
				field.SetConcrete(row, val)
			}
		}
	}

	frame := data.NewFrame(topic, timeField)
	for _, key := range keys {
		frame.Fields = append(frame.Fields, fields[key])
	}
	log.DefaultLogger.Debug(fmt.Sprintf("jsonMessagesToFrame: New Frame for topic=%s", topic))
	return frame
}
