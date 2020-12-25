package common

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"reflect"
)

func decode(buf *bytes.Buffer, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Bool:
		var i int8
		err := binary.Read(buf, binary.LittleEndian, &i)
		if err != nil {
			return err
		}
		if i == 1 {
			v.SetBool(true)
		}
	case reflect.Int:
		var i int32
		err := binary.Read(buf, binary.LittleEndian, &i)
		if err != nil {
			return err
		}
		px := v.Addr().Interface().(*int)
		*px = int(i)
	case reflect.Int8:
		var i int8
		err := binary.Read(buf, binary.LittleEndian, &i)
		if err != nil {
			return err
		}
		px := v.Addr().Interface().(*int8)
		*px = i
	case reflect.Int16:
		var i int16
		err := binary.Read(buf, binary.LittleEndian, &i)
		if err != nil {
			return err
		}
		px := v.Addr().Interface().(*int16)
		*px = i
	case reflect.Int32:
		var i int32
		err := binary.Read(buf, binary.LittleEndian, &i)
		if err != nil {
			return err
		}
		px := v.Addr().Interface().(*int32)
		*px = i
	case reflect.Int64:
		var i int64
		err := binary.Read(buf, binary.LittleEndian, &i)
		if err != nil {
			return err
		}
		px := v.Addr().Interface().(*int64)
		*px = i
	case reflect.Uint:
		var i uint32
		err := binary.Read(buf, binary.LittleEndian, &i)
		if err != nil {
			return err
		}
		px := v.Addr().Interface().(*uint)
		*px = uint(i)
	case reflect.Uint8:
		var i uint8
		err := binary.Read(buf, binary.LittleEndian, &i)
		if err != nil {
			return err
		}
		px := v.Addr().Interface().(*uint8)
		*px = i
	case reflect.Uint16:
		var i uint16
		err := binary.Read(buf, binary.LittleEndian, &i)
		if err != nil {
			return err
		}
		px := v.Addr().Interface().(*uint16)
		*px = i
	case reflect.Uint32:
		var i uint32
		err := binary.Read(buf, binary.LittleEndian, &i)
		if err != nil {
			return err
		}
		px := v.Addr().Interface().(*uint32)
		*px = i
	case reflect.Uint64:
		var i uint64
		err := binary.Read(buf, binary.LittleEndian, &i)
		if err != nil {
			return err
		}
		px := v.Addr().Interface().(*uint64)
		*px = i
	case reflect.Float32:
		var f32 float32
		err := binary.Read(buf, binary.LittleEndian, &f32)
		if err != nil {
			return err
		}
		px := v.Addr().Interface().(*float32)
		*px = f32

	case reflect.Float64:
		var f64 float64
		err := binary.Read(buf, binary.LittleEndian, &f64)
		if err != nil {
			return err
		}
		px := v.Addr().Interface().(*float64)
		*px = f64

	case reflect.String:
		var s string
		b, err := buf.ReadBytes(0)
		if err != nil {
			return err
		}
		s = string(b[:len(b)-1])
		px := v.Addr().Interface().(*string)
		*px = s
	case reflect.Ptr:
		return decode(buf, v.Elem())
	case reflect.Array:
		for i := 0; i < v.Len(); i++ {
			if err := decode(buf, v.Index(i)); err != nil {
				return err
			}
		}
	case reflect.Slice:
		var l uint32
		err := binary.Read(buf, binary.LittleEndian, &l)
		if err != nil {
			return err
		}
		for i := uint32(0); i < l; i++ {
			item := reflect.New(v.Type().Elem()).Elem()
			if err := decode(buf, item); err != nil {
				return err
			}
			v.Set(reflect.Append(v, item))
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			if err := decode(buf, v.Field(i)); err != nil {
				return err
			}
		}
	case reflect.Map:
		var l uint32
		err := binary.Read(buf, binary.LittleEndian, &l)
		if err != nil {
			return err
		}
		v.Set(reflect.MakeMap(v.Type()))
		for i := uint32(0); i < l; i++ {
			key := reflect.New(v.Type().Key()).Elem()
			err = decode(buf, key)
			if err != nil {
				return err
			}
			value := reflect.New(v.Type().Elem()).Elem()
			err = decode(buf, value)
			if err != nil {
				return err
			}
			v.SetMapIndex(key, value)
		}
	default:
		panic(fmt.Sprintf("%v is not support", v.Type()))
	}
	return nil
}
func encode(buf *bytes.Buffer, v reflect.Value) {
	switch v.Kind() {
	case reflect.Bool:
		if v.Bool() {
			binary.Write(buf, binary.LittleEndian, int8(1))
		} else {
			binary.Write(buf, binary.LittleEndian, int8(0))
		}
	case reflect.Int:
		binary.Write(buf, binary.LittleEndian, int32(v.Int()))
	case reflect.Int8:
		binary.Write(buf, binary.LittleEndian, int8(v.Int()))
	case reflect.Int16:
		binary.Write(buf, binary.LittleEndian, int16(v.Int()))
	case reflect.Int32:
		binary.Write(buf, binary.LittleEndian, int32(v.Int()))
	case reflect.Int64:
		binary.Write(buf, binary.LittleEndian, int64(v.Int()))
	case reflect.Uint:
		binary.Write(buf, binary.LittleEndian, uint32(v.Uint()))
	case reflect.Uint8:
		binary.Write(buf, binary.LittleEndian, uint8(v.Uint()))
	case reflect.Uint16:
		binary.Write(buf, binary.LittleEndian, uint16(v.Uint()))
	case reflect.Uint32:
		binary.Write(buf, binary.LittleEndian, uint32(v.Uint()))
	case reflect.Uint64:
		binary.Write(buf, binary.LittleEndian, uint64(v.Uint()))
	case reflect.Float32:
		f := math.Float32bits(float32(v.Float()))
		binary.Write(buf, binary.LittleEndian, f)
	case reflect.Float64:
		f := math.Float64bits(v.Float())
		binary.Write(buf, binary.LittleEndian, f)
	case reflect.String:
		b := make([]byte, v.Len()+1)
		copy(b, v.String())
		buf.Write(b)
	case reflect.Ptr:
		encode(buf, v.Elem())
	case reflect.Array:
		for i := 0; i < v.Len(); i++ {
			encode(buf, v.Index(i))
		}
	case reflect.Slice:
		binary.Write(buf, binary.LittleEndian, uint32(v.Len()))
		for i := 0; i < v.Len(); i++ {
			encode(buf, v.Index(i))
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			encode(buf, v.Field(i))
		}
	case reflect.Map:
		binary.Write(buf, binary.LittleEndian, uint32(v.Len()))
		for _, key := range v.MapKeys() {
			encode(buf, key)
			encode(buf, v.MapIndex(key))
		}
	default:
		panic(fmt.Sprintf("%v is not support", v.Type()))
	}
}

//Packet 消息包处理类
type Packet struct {
	buffer *bytes.Buffer
	cmdId  uint16
}

//Initialize 初始化操作数
func (p *Packet) Initialize(cmdId uint16) {
	p.buffer = new(bytes.Buffer)
	p.cmdId = cmdId
}

//GetOpCode 操作数
func (p *Packet) GetCmd() uint16 {
	return p.cmdId
}
func (p *Packet) WriteBytes(b []byte) {
	p.buffer.Write(b)
}

func (p *Packet) WriteData(data interface{}) {
	v := reflect.Indirect(reflect.ValueOf(data))
	encode(p.buffer, v)
}

func (p *Packet) ReadData(data interface{}) error {
	v := reflect.ValueOf(data)
	return decode(p.buffer, v)
}

func (p *Packet) Len() int {
	return p.buffer.Len()
}

func (p *Packet) Bytes() []byte {
	return p.buffer.Bytes()
}

func (p *Packet) GetBuffer() *bytes.Buffer {
	return p.buffer
}

func (p *Packet) SetBuffer(b *bytes.Buffer) {
	p.buffer = b
}
