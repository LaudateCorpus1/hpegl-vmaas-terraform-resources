// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package utils

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	ErrInvalidType   = "error : Invalid Type"
	ErrKeyNotDefined = "error : Key is not defined"
	ErrSet           = "error : Failed to set"
	NAN              = -999999
)

type Data struct {
	d *schema.ResourceData
	// errors will hold list of errors for each attrib
	// attrib name will be the key
	errors map[string][]string
}

// NewData returns new Data instance
func NewData(d *schema.ResourceData) *Data {
	return &Data{
		d:      d,
		errors: make(map[string][]string),
	}
}

func (d *Data) Error() error {
	if d.hasError() {
		errStr := ""
		for k, v := range d.errors {
			errStr += k + " : " + strings.Join(v, ",") + "\n"
		}

		return errors.New(errStr)
	}

	return nil
}

func (d *Data) hasError() bool {
	return len(d.errors) > 0
}

func (d *Data) err(key, msg string) {
	d.errors[key] = append(d.errors[key], msg)
}

// GetListMap take key as parameter and returns []map[string]interfac{}.
// This function can be used for retrieving list of map or list of set
func (d *Data) GetListMap(key string) []map[string]interface{} {
	src := d.get(key)
	if src == nil {
		return nil
	}
	list, ok := src.([]interface{})
	if !ok {
		d.err(key, ErrInvalidType)

		return nil
	}
	dst := make([]map[string]interface{}, 0, len(list))
	for _, s := range list {
		ds, ok := s.(map[string]interface{})
		if ok {
			dst = append(dst, ds)
		}
	}

	return dst
}

func (d *Data) get(key string) interface{} {
	return d.d.Get(key)
}

func (d *Data) GetID() int64 {
	id, err := ParseInt(d.d.Id())
	if err != nil {
		d.err("id", ErrInvalidType)

		return NAN
	}

	return id
}

func (d *Data) GetIDString() string {
	return d.d.Id()
}

func (d *Data) SetID(v string) {
	d.d.SetId(v)
}

func (d *Data) set(key string, value interface{}) error {
	return d.d.Set(key, value)
}

func (d *Data) GetStringList(key string) []string {
	src := d.get(key)
	list, ok := src.([]interface{})
	if !ok {
		return nil
	}
	dst := make([]string, 0, len(list))
	for _, s := range list {
		ds, ok := s.(string)
		if ok {
			dst = append(dst, ds)
		}
	}

	return dst
}

func (d *Data) GetInt(key string) int64 {
	valString, ok := d.d.GetOk(key)
	if !ok {
		d.err(key, ErrKeyNotDefined)

		return NAN
	}
	val, err := ParseInt(valString.(string))
	if err != nil {
		d.err(key, ErrInvalidType)
	}

	return val
}

// GetSMap for get map for a Set
func (d *Data) GetSMap(key string) map[string]interface{} {
	src, ok := d.d.GetOk(key)
	if !ok {
		return nil
	}
	set, ok := src.(*schema.Set)
	if !ok {
		return nil
	}
	if set == nil {
		return nil
	}
	list := set.List()
	if len(list) == 0 {
		return nil
	}

	return list[0].(map[string]interface{})
}

func (d *Data) GetMap(key string) map[string]interface{} {
	src, ok := d.d.GetOk(key)
	if !ok {
		return nil
	}
	dst, ok := src.(map[string]interface{})
	if !ok {
		return nil
	}

	return dst
}

func (d *Data) GetString(key string) string {
	val := d.get(key)
	if val != nil {
		return val.(string)
	}

	return ""
}

func (d *Data) GetJSONNumber(key string) json.Number {
	in := d.get(key)

	return json.Number(in.(string))
}

func (d *Data) SetString(key string, value string) {
	if err := d.set(key, value); err != nil {
		d.err(key, ErrSet+" : "+err.Error())
	}
}

func (d *Data) ListToIntSlice(key string) []int {
	src := d.get(key)
	list, ok := src.([]interface{})
	if !ok {
		return nil
	}

	dst := make([]int, 0, len(list))
	for i, s := range list {
		ds, ok := s.(int)
		if !ok {
			d.err(key, ErrInvalidType+" at index "+strconv.Itoa(i))
		} else {
			dst = append(dst, ds)
		}
	}

	return dst
}