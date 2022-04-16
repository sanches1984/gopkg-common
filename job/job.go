package job

import (
	"context"
	"fmt"
	"github.com/gocraft/work"
	"github.com/sanches1984/gopkg-errors"
	"reflect"
	"time"
)

const nameProperty = "jobName"

type WorkRecord struct {
	Job      interface{}
	Fn       interface{}
	Schedule string
}

func getJobName(ctx context.Context, job interface{}) (string, error) {
	st := reflect.TypeOf(job)
	field, ok := st.FieldByName(nameProperty)
	if !ok {
		return "", errors.Internal.Err(ctx, "Not fount property '"+nameProperty+"' in Job container")
	}
	jobName := field.Tag.Get("job")
	if jobName == "" {
		return "", errors.Internal.Err(ctx, "Empty property '"+nameProperty+"' in Job container")
	}
	return jobName, nil
}

func NewJobWithArgs(ctx context.Context, v interface{}) (*work.Job, error) {
	args, err := packArguments(ctx, v)
	if err != nil {
		return nil, err
	}
	return &work.Job{Args: args}, nil
}

func packArguments(ctx context.Context, job interface{}) (map[string]interface{}, error) {
	st := reflect.ValueOf(job)
	if st.Kind() == reflect.Ptr {
		st = st.Elem()
	}
	numFields := st.NumField()
	ret := make(map[string]interface{}, numFields)
	for i := 0; i < numFields; i++ {
		name := st.Type().Field(i).Name
		if name == nameProperty {
			continue
		}
		kind := st.Type().Field(i).Type.Kind()
		switch kind {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			ret[name] = float64(st.Field(i).Int())
		case reflect.String, reflect.Bool:
			ret[name] = st.Field(i).Interface()
		case reflect.Slice:
			sliceKind := st.Type().Field(i).Type.Elem().Kind()
			switch sliceKind {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				sliceNumField := st.Field(i).Len()
				slice := make([]interface{}, 0, sliceNumField)
				for n := 0; n < sliceNumField; n++ {
					slice = append(slice, float64(st.Field(i).Index(n).Int()))
				}
				ret[name] = slice
			default:
				return nil, errors.Internal.Err(ctx, "Unsupported slice argument type").
					WithLogKV("arg", name, "kind", sliceKind)
			}
		case reflect.Struct:
			switch v := st.Field(i).Interface().(type) {
			case time.Time:
				ret[name] = v.Format(time.RFC3339)
			default:
				return nil, errors.Internal.Err(ctx, "Unsupported struct argument type").
					WithLogKV("arg", name, "type", fmt.Sprintf("%T", v))
			}
		default:
			return nil, errors.Internal.Err(ctx, "Unsupported argument type").
				WithLogKV("arg", name, "kind", kind)
		}
	}
	return ret, nil
}

func UnpackArguments(ctx context.Context, job interface{}, workJob *work.Job) error {
	st := reflect.ValueOf(job)
	if st.Kind() != reflect.Ptr {
		return errors.Internal.Err(ctx, "Job should be pointer")
	}
	st = st.Elem()
	numFields := st.NumField()
	for i := 0; i < numFields; i++ {
		name := st.Type().Field(i).Name
		if name == nameProperty {
			continue
		}
		value, ok := workJob.Args[name]
		if !ok {
			continue
		}
		var castOk bool
		fieldType := st.Type().Field(i).Type
		fieldKind := fieldType.Kind()
		switch fieldKind {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if v, ok := value.(float64); ok {
				st.Field(i).SetInt(int64(v))
				castOk = true
			}
		case reflect.String:
			if v, ok := value.(string); ok {
				st.Field(i).SetString(v)
				castOk = true
			}
		case reflect.Bool:
			if v, ok := value.(bool); ok {
				st.Field(i).SetBool(v)
				castOk = true
			}
		case reflect.Slice:
			sliceKind := st.Type().Field(i).Type.Elem().Kind()
			switch sliceKind {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if v, ok := value.([]interface{}); ok {
					var slice reflect.Value
					sliceNumField := len(v)
					slice = reflect.MakeSlice(fieldType, sliceNumField, sliceNumField)
					for n, vInterface := range v {
						vFloat64, ok := vInterface.(float64)
						if !ok {
							return errors.Internal.Err(ctx, "Cant cast Job argument: not cast slice item").
								WithLogKV("name", name, "kind", sliceKind, "index", n, "value", vInterface)
						}
						item := slice.Index(n)
						switch sliceKind {
						case reflect.Int:
							item.Set(reflect.ValueOf(int(vFloat64)))
						case reflect.Int8:
							item.Set(reflect.ValueOf(int8(vFloat64)))
						case reflect.Int16:
							item.Set(reflect.ValueOf(int16(vFloat64)))
						case reflect.Int32:
							item.Set(reflect.ValueOf(int32(vFloat64)))
						case reflect.Int64:
							item.Set(reflect.ValueOf(int64(vFloat64)))
						default:
							return errors.Internal.Err(ctx, "Can't cast Job argument: not slice implemented").
								WithLogKV("name", name, "kind", sliceKind, "index", n, "value", workJob.Args[name])
						}
					}
					st.Field(i).Set(slice)
					castOk = true
				}
			default:
				return errors.Internal.Err(ctx, "Can't cast Job argument: not slice implemented").
					WithLogKV("name", name, "kind", sliceKind, "value", workJob.Args[name])
			}
		case reflect.Struct:
			switch v := st.Field(i).Interface().(type) {
			case time.Time:
				if valueStr, ok := value.(string); ok {
					dt, err := time.Parse(time.RFC3339, valueStr)
					if err != nil {
						return errors.Internal.Err(ctx, "Cant cast Job argument: not parse struct time.Time item").
							WithLogKV("name", name, "value", v.String())
					}
					st.Field(i).Set(reflect.ValueOf(dt))
					castOk = true
				}
			default:
				return errors.Internal.Err(ctx, "Can't cast Job argument: not struct implemented").
					WithLogKV("arg", name, "type", fmt.Sprintf("%T", v))
			}
		default:
			return errors.Internal.Err(ctx, "Can't cast Job argument: not implemented").
				WithLogKV("name", name, "kind", fieldKind, "value", workJob.Args[name])
		}
		if !castOk {
			return errors.Internal.Err(ctx, "Cant cast Job argument").
				WithLogKV("name", name, "kind", fieldKind, "value", workJob.Args[name])
		}
	}
	return nil
}
