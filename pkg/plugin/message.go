package plugin

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

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

//  Transform the array of MQTT Messages (JSON encoded) into a Grafana Data Frame.
//  See sample messages: https://github.com/lupyuen/the-things-network-datasource#mqtt-log
func jsonMessagesToFrame(topic string, messages []mqtt.Message) *data.Frame {
	//  Quit if no messages to transform
	count := len(messages)
	if count == 0 {
		log.DefaultLogger.Debug(fmt.Sprintf("jsonMessagesToFrame: No msgs for topic=%s", topic))
		return nil
	}

	//  Transform the first message
	msg := messages[0]
	log.DefaultLogger.Debug(fmt.Sprintf("jsonMessagesToFrame: topic=%s, msg=%s", topic, msg.Value))

	//  Constuct the Timestamp field
	timeField := data.NewFieldFromFieldType(data.FieldTypeTime, count)
	timeField.Name = "Time"
	timeField.SetConcrete(0, msg.Timestamp)

	//  Construct the Data Frame
	frame := data.NewFrame(topic, timeField)

	//  Decode the CBOR payload
	body, err := decodeCborPayload(msg.Value)
	if err != nil {
		return set_error(frame, err)
	}

	// Create a field for each key and set the first value
	keys := make([]string, 0, len(body))
	fields := make(map[string]*data.Field, len(body))

	//  Compose the fields for the Data Frame
	for key, val := range body {
		//  Get the Data Frame Type for the field
		typ := get_type(val)

		//  Create the field
		field := data.NewFieldFromFieldType(typ, count)
		field.Name = key
		field.SetConcrete(0, val)
		fields[key] = field
		keys = append(keys, key)
	}
	sort.Strings(keys) // keys stable field order.

	//  TODO: Transform the messages after the first one
	//  We might not need this because The Things Network only supports low-volume messaging
	for row, m := range messages {
		if row == 0 {
			continue
		}
		log.DefaultLogger.Debug(fmt.Sprintf("jsonMessagesToFrame: Dropped msg=%s", topic, m.Value))
	}

	//  Append the fields to the Data Frame
	for _, key := range keys {
		frame.Fields = append(frame.Fields, fields[key])
	}

	//  Dump the Data Frame
	log.DefaultLogger.Debug(fmt.Sprintf("jsonMessagesToFrame: Frame=%+v", frame))
	for _, field := range frame.Fields {
		log.DefaultLogger.Debug(fmt.Sprintf("  field=%+v", field))
	}
	return frame
}

//  Decode the CBOR payload in the JSON message.
//  See sample messages: https://github.com/lupyuen/the-things-network-datasource#mqtt-log
func decodeCborPayload(msg string) (map[string]interface{}, error) {
	//  Deserialise the message doc to a map of String -> interface{}
	var doc map[string]interface{}
	err := json.Unmarshal([]byte(msg), &doc)
	if err != nil {
		return nil, err
	}

	//  Get the Uplink Message
	if doc["uplink_message"] == nil {
		return nil, errors.New("uplink_message missing")
	}
	uplink_message := doc["uplink_message"].(map[string]interface{})

	//  Get the Payload
	if uplink_message["frm_payload"] == nil {
		return nil, errors.New("frm_payload missing")
	}
	frm_payload := uplink_message["frm_payload"].(string)

	//  Base64 decode the Payload
	payload, err := base64.StdEncoding.DecodeString(frm_payload)
	if err != nil {
		return nil, err
	}
	log.DefaultLogger.Debug(fmt.Sprintf("payload: %v", payload))

	//  TODO: Testing CBOR encoding for {"t": 1234}.  See http://cbor.me/
	if payload[0] == 0 {
		payload = []byte{0xA1, 0x61, 0x74, 0x19, 0x04, 0xD2}
	}

	//  Decode CBOR payload to a map of String -> interface{}
	var body map[string]interface{}
	err = cbor.Unmarshal(payload, &body)
	if err != nil {
		return nil, err
	}

	//  Add the Device ID to the body: end_device_ids -> device_id
	if doc["end_device_ids"] != nil {
		end_device_ids := doc["end_device_ids"].(map[string]interface{})
		if end_device_ids["device_id"] != nil {
			body["device_id"] = end_device_ids["device_id"].(string)
		}
	}

	//  TODO: Test various field types
	body["f64"] = float64(1234)
	body["u64"] = uint64(1234)
	body["str"] = "Test"

	//  Shows: map[device_id:eui-70b3d57ed0045669 t:1234]
	log.DefaultLogger.Debug(fmt.Sprintf("CBOR decoded: %v", body))
	return body, nil
}

//  Return the Data Frame Type for the CBOR decoded value
func get_type(val interface{}) data.FieldType {
	//  Based on https://github.com/fxamacker/cbor/blob/master/decode.go#L43-L53
	switch v := val.(type) {
	//  CBOR booleans decode to bool.
	case bool:
		return data.FieldTypeBool

	//  CBOR positive integers decode to uint64.
	case uint64:
		return data.FieldTypeNullableUint64

	//  CBOR negative integers decode to int64 (big.Int if value overflows).
	case int64:
		return data.FieldTypeInt64

	//  CBOR floating points decode to float64.
	case float64:
		return data.FieldTypeNullableFloat64

	//  CBOR text strings decode to string.
	case string:
		return data.FieldTypeNullableString

	//  CBOR times (tag 0 and 1) decode to time.Time.
	case time.Time:
		return data.FieldTypeNullableTime

	//  TODO: CBOR byte strings decode to []byte.
	//  TODO: CBOR arrays decode to []interface{}.
	//  TODO: CBOR maps decode to map[interface{}]interface{}.
	//  TODO: CBOR null and undefined values decode to nil.
	//  TODO: CBOR bignums (tag 2 and 3) decode to big.Int.
	default:
		log.DefaultLogger.Debug(fmt.Sprintf("Unknown type %T for %v", v, val))
		return data.FieldTypeUnknown
	}
}

//  Return the Data Frame set to the given error
func set_error(frame *data.Frame, err error) *data.Frame {
	frame.AppendNotices(data.Notice{
		Severity: data.NoticeSeverityError,
		Text:     err.Error(),
	})
	log.DefaultLogger.Debug(err.Error())
	return frame
}
