// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/sapcc/go-api-declarations/liquid"
	"github.com/sapcc/go-bits/logg"
)

func TestLiquidOptionType(t *testing.T) {
	optionTypes := make(map[reflect.Type]struct{})
	getOptionsTypesRecursively(reflect.ValueOf(liquid.ServiceInfo{}).Type(), optionTypes)
	getOptionsTypesRecursively(reflect.ValueOf(liquid.ServiceUsageReport{}).Type(), optionTypes)
	getOptionsTypesRecursively(reflect.ValueOf(liquid.ServiceCapacityReport{}).Type(), optionTypes)
	for optionType := range optionTypes {
		isHandled := false
		for _, liquidOptionType := range LiquidOptionTypes {
			field, ok := reflect.TypeOf(liquidOptionType).FieldByName("value")
			if !ok {
				t.Errorf("expected type majewsky/gg/option.Option with field 'value' but got %q", optionType)
			}
			if field.Type == optionType {
				isHandled = true
				break
			}
		}
		if !isHandled {
			t.Errorf("compare option missing for type Option[%s]", optionType)
		}
	}
}

func getOptionsTypesRecursively(t reflect.Type, optionTypes map[reflect.Type]struct{}) {
	switch t.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Array, reflect.Map:
		getOptionsTypesRecursively(t.Elem(), optionTypes)
	case reflect.Struct:
		if strings.HasPrefix(t.Name(), "Option") {
			field, ok := t.FieldByName("value")
			if !ok {
				logg.Error(fmt.Sprintf("expected type majewsky/gg/option.Option with field 'value' but got %q", t.Name()))
				return
			}
			optionTypes[field.Type] = struct{}{}
			getOptionsTypesRecursively(field.Type, optionTypes)
		} else {
			for idx := range t.NumField() {
				f := t.Field(idx)
				getOptionsTypesRecursively(f.Type, optionTypes)
			}
		}
	}
}
