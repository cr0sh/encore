package binary

import (
	"bytes"
	"errors"
	"reflect"
	"strings"
	"unsafe"
)

// Marshaler is a interface implemented by types that can marshal themselves into binary stream data.
type Marshaler interface {
	MarshalStream(*bytes.Buffer) error
}

// Unmarshaler is a interface implemented by types that can unmarshal binary stream data to themselves.
type Unmarshaler interface {
	UnmarshalStream(*bytes.Buffer) error
}

// Marshal encodes struct v into buf.
// If v is not a struct, Marshal returns an InvalidMarshalTypeError.
func Marshal(v interface{}, buf *bytes.Buffer) error {
	vv := reflect.ValueOf(v)
	t := vv.Type()
	if t.Kind() != reflect.Struct {
		return InvalidMarshalTypeError{t}
	}

	for i := 0; i < vv.NumField(); i++ {
		field := vv.Field(i)
		tag := strings.ToLower(t.Field(i).Tag.Get("stream"))
		if tag == "pass" || !field.CanInterface() {
			continue
		}

		sp := strings.Split(tag, ",")
		option := make(map[string]struct{})
		for _, v := range sp {
			option[v] = struct{}{}
		}

		fv := field.Interface()
		if marshaler, ok := fv.(Marshaler); ok {
			if err := marshaler.MarshalStream(buf); err != nil {
				return err
			}
			continue
		}

		var endian ByteOrder = BigEndian
		if _, ok := option["little"]; ok {
			endian = LittleEndian
		}

		if field.Kind() == reflect.Slice && i+1 != vv.NumField() {
			return MarshalError{
				errors.New("Exported slice element should be last while marshaling"),
				t.Field(i),
			}
		}

		if err := pack(field, buf, endian); err != nil {
			return MarshalError{err, t.Field(i)}
		}

	}
	return nil
}

func pack(v reflect.Value, buf *bytes.Buffer, endian ByteOrder) error {
	kind := v.Kind()
	fv := v.Interface()

	switch kind {
	case reflect.Int, reflect.Uint:
		break // explicit break expression to clarify separated case with "bool"

	case reflect.Bool:
		fv, ok := fv.(bool)
		if !ok {
			return TypeAssertionError{"bool"}
		}

		var err error
		if fv {
			err = buf.WriteByte(1)
		} else {
			err = buf.WriteByte(0)
		}

		if err != nil {
			return err
		}

	case reflect.Int8:
		fv, ok := fv.(int8)
		if !ok {
			return TypeAssertionError{"int8"}
		}
		if err := buf.WriteByte(*(*byte)(unsafe.Pointer(&fv))); err != nil {
			return err
		}

	case reflect.Uint8:
		fv, ok := fv.(byte)
		if !ok {
			return TypeAssertionError{"byte"}
		}
		if err := buf.WriteByte(fv); err != nil {
			return err
		}

	case reflect.Int16:
		fv, ok := fv.(int16)
		if !ok {
			return TypeAssertionError{"int16"}
		}
		b := make([]byte, 2)
		endian.PutInt16(b, fv)
		if _, err := buf.Write(b); err != nil {
			return err
		}

	case reflect.Uint16:
		fv, ok := fv.(uint16)
		if !ok {
			return TypeAssertionError{"uint16"}
		}
		b := make([]byte, 2)
		endian.PutUint16(b, fv)
		if _, err := buf.Write(b); err != nil {
			return err
		}

	case reflect.Int32:
		fv, ok := fv.(int32)
		if !ok {
			return TypeAssertionError{"int32"}
		}
		b := make([]byte, 4)
		endian.PutInt32(b, fv)
		if _, err := buf.Write(b); err != nil {
			return err
		}

	case reflect.Uint32:
		fv, ok := fv.(uint32)
		if !ok {
			return TypeAssertionError{"uint32"}
		}
		b := make([]byte, 4)
		endian.PutUint32(b, fv)
		if _, err := buf.Write(b); err != nil {
			return err
		}

	case reflect.Int64:
		fv, ok := fv.(int64)
		if !ok {
			return TypeAssertionError{"int64"}
		}
		b := make([]byte, 8)
		endian.PutInt64(b, fv)
		if _, err := buf.Write(b); err != nil {
			return err
		}

	case reflect.Uint64:
		fv, ok := fv.(uint64)
		if !ok {
			return TypeAssertionError{"uint64"}
		}
		b := make([]byte, 8)
		endian.PutUint64(b, fv)
		if _, err := buf.Write(b); err != nil {
			return err
		}

	case reflect.Float32:
		fv, ok := fv.(float32)
		if !ok {
			return TypeAssertionError{"float32"}
		}
		b := make([]byte, 4)
		endian.PutFloat32(b, fv)
		if _, err := buf.Write(b); err != nil {
			return err
		}

	case reflect.Float64:
		fv, ok := fv.(float64)
		if !ok {
			return TypeAssertionError{"float64"}
		}
		b := make([]byte, 8)
		endian.PutFloat64(b, fv)
		if _, err := buf.Write(b); err != nil {
			return err
		}

	case reflect.Struct:
		if err := Marshal(fv, buf); err != nil {
			return err
		}

	case reflect.Array, reflect.Slice:
		length := v.Len()
		if v.Type().Elem().Kind() == reflect.Uint8 {
			if _, err := buf.Write(v.Slice(0, length).Interface().([]byte)); err != nil {
				return err
			}
			break
		}

		for j := 0; j < length; j++ {
			if err := pack(v.Index(j), buf, endian); err != nil {
				return err
			}
		}

	}
	return nil
}

// Unmarshal parses binary stream data from buf and stores the result in the struct pointed by v.
func Unmarshal(v interface{}, buf *bytes.Buffer) error {
	vv := reflect.ValueOf(v).Elem()
	t := vv.Type()
	for i := 0; i < t.NumField(); i++ {
		field := vv.Field(i)
		tag := strings.ToLower(t.Field(i).Tag.Get("stream"))
		if tag == "pass" || !field.CanSet() {
			continue
		}

		sp := strings.Split(tag, ",")
		option := make(map[string]struct{})
		for _, v := range sp {
			option[v] = struct{}{}
		}

		fv := field.Addr().Interface()
		if unmarshaler, ok := fv.(Unmarshaler); ok {
			if err := unmarshaler.UnmarshalStream(buf); err != nil {
				return err
			}
			continue
		}

		var endian ByteOrder = BigEndian
		if _, ok := option["little"]; ok {
			endian = LittleEndian
		}

		if field.Kind() == reflect.Slice && i+1 != vv.NumField() {
			return UnmarshalError{
				errors.New("Exported slice element should be last while unmarshaling"),
				t.Field(i),
			}
		}

		if err := unpack(field, buf, endian); err != nil {
			return UnmarshalError{err, t.Field(i)}
		}
	}
	return nil
}

func unpack(v reflect.Value, buf *bytes.Buffer, endian ByteOrder) error {
	kind := v.Kind()

	switch kind {
	case reflect.Int, reflect.Uint:
		break // explicit break expression to clarify separated case with "bool"

	case reflect.Bool:
		x, err := buf.ReadByte()
		if err != nil {
			return err
		}
		if x != 0 {
			v.SetBool(true)
		}

	case reflect.Int8:
		x, err := buf.ReadByte()
		if err != nil {
			return err
		}
		v.SetInt(int64(*(*int8)(unsafe.Pointer(&x))))

	case reflect.Uint8:
		x, err := buf.ReadByte()
		if err != nil {
			return err
		}
		v.SetUint(uint64(x))

	case reflect.Int16:
		b := make([]byte, 2)
		if _, err := buf.Read(b); err != nil {
			return err
		}
		v.SetInt(int64(endian.Int16(b)))

	case reflect.Uint16:
		b := make([]byte, 2)
		if _, err := buf.Read(b); err != nil {
			return err
		}
		v.SetUint(uint64(endian.Uint16(b)))

	case reflect.Int32:
		b := make([]byte, 4)
		if _, err := buf.Read(b); err != nil {
			return err
		}
		v.SetInt(int64(endian.Int32(b)))

	case reflect.Uint32:
		b := make([]byte, 4)
		if _, err := buf.Read(b); err != nil {
			return err
		}
		v.SetUint(uint64(endian.Uint32(b)))

	case reflect.Int64:
		b := make([]byte, 8)
		if _, err := buf.Read(b); err != nil {
			return err
		}
		v.SetInt(endian.Int64(b))

	case reflect.Uint64:
		b := make([]byte, 8)
		if _, err := buf.Read(b); err != nil {
			return err
		}
		v.SetUint(endian.Uint64(b))

	case reflect.Float32:
		b := make([]byte, 4)
		if _, err := buf.Read(b); err != nil {
			return err
		}
		v.SetFloat(float64(endian.Float32(b)))

	case reflect.Float64:
		b := make([]byte, 8)
		if _, err := buf.Read(b); err != nil {
			return err
		}
		v.SetFloat(endian.Float64(b))

	case reflect.Struct:
		if err := Unmarshal(v.Addr().Interface(), buf); err != nil {
			return err
		}

	case reflect.Array, reflect.Slice:
		length := v.Len()
		if v.Type().Elem().Kind() == reflect.Uint8 {
			if _, err := buf.Read(v.Slice(0, length).Interface().([]byte)); err != nil {
				return err
			}
			break
		}

		for j := 0; j < length; j++ {
			if err := unpack(v.Index(j), buf, endian); err != nil {
				return err
			}
		}
	}
	return nil
}

type InvalidMarshalTypeError struct {
	T reflect.Type
}

func (err InvalidMarshalTypeError) Error() string {
	return "The second argument of Marshal should be struct; given " + err.T.Kind().String()
}

type InvalidUnmarshalTypeError struct {
	T reflect.Type
}

func (err InvalidUnmarshalTypeError) Error() string {
	return "The second argument of Marshal should be struct pointer; given " + err.T.Kind().String()
}

type TypeAssertionError struct {
	totype string
}

func (err TypeAssertionError) Error() string {
	return "Type assertion failed to " + err.totype
}

type MarshalError struct {
	E     error
	field reflect.StructField
}

func (err MarshalError) Error() string {
	return "Error while marshaling " + err.field.Name + ": " + err.E.Error()
}

type UnmarshalError struct {
	E     error
	field reflect.StructField
}

func (err UnmarshalError) Error() string {
	return "Error while marshaling " + err.field.Name + ": " + err.E.Error()
}
