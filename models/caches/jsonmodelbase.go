package caches

import (
	"log"
	"reflect"

	"github.com/bitly/go-simplejson"
)

type TypeData struct {
	typ        reflect.Kind
	fieldIndex int
}

var (
	typeInfoMap = make(map[string](map[string]*TypeData))
)

func registerJsonTypeInfo(i interface{}) {
	typ := reflect.TypeOf(i)
	if typ.Kind() != reflect.Ptr {
		log.Fatalf("Failed to register json model type info: %s, not ptr type", typ.Name())
	}

	val := reflect.Indirect(reflect.ValueOf(i))
	typ = val.Type()
	if typ.Kind() != reflect.Struct {
		log.Fatalf("Failed to register json model type info: %s, not struct type", typ.Name())
	}

	fieldNum := typ.NumField()
	typeDataMap := make(map[string]*TypeData)
	for ix := 0; ix < fieldNum; ix++ {
		field := typ.Field(ix)
		if field.Tag.Get("json") == "" {
			continue
		}

		if !val.Field(ix).CanSet() {
			log.Fatalf("Failed to register json model type info: %dth field(%s) of %s can not be setted", ix, field.Name, typ.Name())
		}

		typeDataMap[field.Tag.Get("json")] = &TypeData{typ: field.Type.Kind(), fieldIndex: ix}
	}

	typeInfoMap[typ.Name()] = typeDataMap
}

func fillJsonModelInfo(i interface{}, data *simplejson.Json) {
	val := reflect.ValueOf(i)
	if val.Kind() != reflect.Ptr {
		return
	}

	val = reflect.Indirect(val)
	jsonTypeInfo := typeInfoMap[val.Type().Name()]
	if jsonTypeInfo == nil {
		return
	}

	for jsonKey, typeData := range jsonTypeInfo {
		switch typeData.typ {
		case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
			val.Field(typeData.fieldIndex).SetInt(data.Get(jsonKey).MustInt64())
		case reflect.Float64, reflect.Float32:
			val.Field(typeData.fieldIndex).SetFloat(data.Get(jsonKey).MustFloat64())
		case reflect.String:
			val.Field(typeData.fieldIndex).SetString(data.Get(jsonKey).MustString())
		case reflect.Bool:
			val.Field(typeData.fieldIndex).SetBool(data.Get(jsonKey).MustBool())
		default:
			//
		}
	}
}

func modelTransfer2JsonModel(i interface{}) (json *simplejson.Json) {
	val := reflect.ValueOf(i)
	if val.Kind() == reflect.Ptr {
		val = reflect.Indirect(val)
	}

	jsonTypeInfo := typeInfoMap[val.Type().Name()]
	if jsonTypeInfo == nil {
		return
	}

	json = simplejson.New()
	for jsonKey, typeData := range jsonTypeInfo {
		switch typeData.typ {
		case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
			json.Set(jsonKey, val.Field(typeData.fieldIndex).Int())
		case reflect.Float64, reflect.Float32:
			json.Set(jsonKey, val.Field(typeData.fieldIndex).Float())
		case reflect.String:
			json.Set(jsonKey, val.Field(typeData.fieldIndex).String())
		case reflect.Bool:
			json.Set(jsonKey, val.Field(typeData.fieldIndex).Bool())
		default:
			//
		}
	}

	return
}
