package learning

import (
	"context"
	"errors"
	"go/token"
	"reflect"
)

var (
	// Precompute the reflect type for error. Can't use error directly
	// because Typeof takes an empty interface value. This is annoying.
	typeOfError = reflect.TypeOf((*error)(nil)).Elem()
	// Precompute the reflect type for context.Context.
	typeOfContext = reflect.TypeOf((*context.Context)(nil)).Elem()
)

type methodType struct {
	method    reflect.Method
	ArgType   reflect.Type
	ReplyType reflect.Type
}

type service struct {
	name    string                 // name of service
	rcvr    reflect.Value          // receiver of methods for the service
	typ     reflect.Type           // type of the receiver
	methods map[string]*methodType // registered methods
}

// Register publishes in the server the set of methods of the
// receiver value that satisfy the following conditions:
//   - 导出类型的导出方法
//   - 方法有两个参数, context.Context 和 导出类型的指针
//   - 方法有两个返回值, 导出类型的指针 和 error
//
// It returns an error if the receiver is not an exported type or has
// no suitable methods. The client accesses each method using a string of
// the form "Type.Method", where Type is the receiver's concrete type.
func (s *Server) Register(rcvr any) error {
	name := reflect.Indirect(reflect.ValueOf(rcvr)).Type().Name()
	return s.register(rcvr, name)
}

// RegisterName is like Register but uses the provided name for the type
// instead of the receiver's concrete type.
func (s *Server) RegisterName(name string, rcvr any) error {
	return s.register(rcvr, name)
}

func (s *Server) register(rcvr any, name string) error {
	svc := new(service)
	svc.typ = reflect.TypeOf(rcvr)
	svc.rcvr = reflect.ValueOf(rcvr)
	svc.name = name
	if svc.name == "" {
		return errors.New("rpc.Register: no service name for type " + svc.typ.String())
	}
	if !token.IsExported(svc.name) {
		return errors.New("rpc.Register: type " + svc.name + " is not exported")
	}

	// Install the methods
	svc.methods = suitableMethods(svc.typ)

	if len(svc.methods) == 0 {
		var errorStr string

		// To help the user, see if a pointer receiver would work.
		method := suitableMethods(reflect.PointerTo(svc.typ))
		if len(method) != 0 {
			errorStr = "rpc.Register: type " + svc.name + " has no exported methods of suitable type (hint: pass a pointer to value of that type)"
		} else {
			errorStr = "rpc.Register: type " + svc.name + " has no exported methods of suitable type"
		}
		return errors.New(errorStr)
	}

	s.serviceMapMu.Lock()
	defer s.serviceMapMu.Unlock()

	if _, dup := s.serviceMap[svc.name]; dup {
		return errors.New("rpc: service already defined: " + svc.name)
	}
	s.serviceMap[svc.name] = svc

	return nil
}

// suitableMethods returns suitable Rpc methods of typ. It will log
// errors if logErr is true.
func suitableMethods(typ reflect.Type) map[string]*methodType {
	methods := make(map[string]*methodType)
	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)
		mtype := method.Type
		mname := method.Name
		// 必须是导出的方法
		if !method.IsExported() {
			continue
		}
		// 方法必须有三个参数: receiver, context.Context, *Args.
		if mtype.NumIn() != 3 {
			continue
		}
		// 方法的第二个参数必须为 context.Context
		if ctxType := mtype.In(1); ctxType != typeOfContext {
			continue
		}
		// 方法的第三个参数必须是指针, 且为导出类型
		argType := mtype.In(2)
		if argType.Kind() != reflect.Pointer || !isExportedOrBuiltinType(argType) {
			continue
		}
		// 方法必须有两个返回值: *Reply, error
		if mtype.NumOut() != 2 {
			continue
		}
		// 方法的第一个返回值必须是指针, 且为导出类型
		replyType := mtype.Out(0)
		if replyType.Kind() != reflect.Pointer || !isExportedOrBuiltinType(replyType) {
			continue
		}
		// 方法的第二个返回值必须为 error
		if errorType := mtype.Out(1); errorType != typeOfError {
			continue
		}
		methods[mname] = &methodType{
			method:    method,
			ArgType:   argType,
			ReplyType: replyType,
		}
	}
	return methods
}

// Is this type exported or a builtin?
func isExportedOrBuiltinType(t reflect.Type) bool {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	// PkgPath will be non-empty even for an exported type,
	// so we need to check the type name as well.
	return token.IsExported(t.Name()) || t.PkgPath() == ""
}
