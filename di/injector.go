// https://github.com/codegangsta/inject
package core

import (
    `errors`
    `reflect`
    `sync`
)

type Injector interface {
    Applicator
    Invoker
    TypeMapper
}

type Applicator interface {
    Apply(interface{}) error
}

type Invoker interface {
    Invoke(interface{}) ([]reflect.Value, error)
}

type TypeMapper interface {
    Map(interface{}) TypeMapper
    MapTo(interface{}, interface{}) TypeMapper
    Set(reflect.Type, reflect.Value) TypeMapper
    Get(reflect.Type) reflect.Value
}

func InterfaceOf(value interface{}) reflect.Type {
    rt := reflect.TypeOf(value)
    for rt.Kind() == reflect.Ptr {
        rt = rt.Elem()
    }
    if rt.Kind() != reflect.Interface {
        panic("Called inject.InterfaceOf with a value that is not a pointer to an interface")
    }
    return rt
}

type injector struct {
    values map[reflect.Type]reflect.Value
    lock   *sync.RWMutex
}

func NewInjector() Injector {
    i := new(injector)
    i.lock = &sync.RWMutex{}
    i.values = make(map[reflect.Type]reflect.Value)
    return i
}

var (
    ErrorParamType     = errors.New("参数类型错误")
    ErrorValueNotFound = errors.New("未找到需要注入的参数")
)

func (injector *injector) Invoke(f interface{}) ([]reflect.Value, error) {
    rt := reflect.TypeOf(f)
    if rt.Kind() != reflect.Func {
        return nil, ErrorParamType
    }
    args := make([]reflect.Value, rt.NumIn())
    for i := 0; i < rt.NumIn(); i++ {
        argType := rt.In(i)
        val := injector.Get(argType)
        if !val.IsValid() {
            return nil, ErrorValueNotFound
        }
        args[i] = val
    }
    return reflect.ValueOf(f).Call(args), nil
}

func (injector *injector) Apply(val interface{}) error {
    rv := reflect.ValueOf(val)
    for rv.Kind() == reflect.Ptr {
        rv = rv.Elem()
    }
    if rv.Kind() != reflect.Struct {
        return ErrorParamType
    }
    rt := rv.Type()
    for i := 0; i < rv.NumField(); i++ {
        fv := rv.Field(i)
        ft := rt.Field(i)
        if fv.CanSet() && (ft.Tag.Get("inject") == "true") {
            _ft := fv.Type()
            _v := injector.Get(_ft)
            if !_v.IsValid() {
                return ErrorValueNotFound
            }
            fv.Set(_v)
        }
    }
    return nil
}

func (injector *injector) Map(val interface{}) TypeMapper {
    injector.Set(reflect.TypeOf(val), reflect.ValueOf(val))
    return injector
}

func (injector *injector) MapTo(val interface{}, iFace interface{}) TypeMapper {
    injector.Set(reflect.TypeOf(iFace), reflect.ValueOf(val))
    return injector
}

func (injector *injector) Set(typ reflect.Type, val reflect.Value) TypeMapper {
    injector.lock.Lock()
    injector.values[typ] = val
    injector.lock.Unlock()
    return injector
}

func (injector *injector) Get(rt reflect.Type) reflect.Value {
    injector.lock.RLock()
    val, ok := injector.values[rt]
    injector.lock.RUnlock()
    if ok && val.IsValid() {
        return val
    }
    if rt.Kind() == reflect.Interface {
        for k, v := range injector.values {
            if k.Implements(rt) {
                val = v
                break
            }
        }
    }
    return val
}
