package go_order

import (
	"fmt"
	"sync"
)

type Order struct {
	funcs map[string]*orderFunc
}
type worker func(res map[string]interface{}) (interface{}, error)

type orderFunc struct {
	sync.Once

	Deps []string
	Ctr  int
	Fn   worker
	C    chan interface{}
}

func (o *orderFunc) Done(r interface{}) {
	for i := 0; i < o.Ctr; i++ {
		o.C <- r
	}
}

func (o *orderFunc) Close() {
	o.Do(func() {
		close(o.C)
	})
}

func (o *orderFunc) Init() {
	o.C = make(chan interface{}, o.Ctr)
}

func New() *Order {
	return &Order{
		funcs: make(map[string]*orderFunc),
	}
}

func (o *Order) Add(name string, d []string, fn worker) *Order {
	o.funcs[name] = &orderFunc{
		Deps: d,
		Fn:   fn,
		Ctr:  1,
	}
	return o
}

func (o *Order) Run() (map[string]interface{}, error) {
	for name, fn := range o.funcs {
		for _, dep := range fn.Deps {
			if dep == name {
				return nil, fmt.Errorf("Error: Function \"%s\" depends of it self!", name)
			}
			if _, has := o.funcs[dep]; !has {
				return nil, fmt.Errorf("Error: Function \"%s\" not exists!", dep)
			}
			o.funcs[dep].Ctr++
		}
	}
	return o.run()
}

func (o *Order) run() (map[string]interface{}, error) {
	var err error
	res := make(map[string]interface{}, len(o.funcs))

	for _, f := range o.funcs {
		f.Init()
	}

	for name, fn := range o.funcs {
		go func(name string, of *orderFunc) {
			defer func() {
				of.Close()
			}()

			results := make(map[string]interface{}, len(of.Deps))

			for _, dep := range of.Deps {
				results[dep] = <-o.funcs[dep].C
			}
			r, fnErr := of.Fn(results)
			if fnErr != nil {
				for _, f := range o.funcs {
					f.Close()
				}
				err = fnErr
				return
			}
			if err != nil {
				return
			}
			of.Done(r)
		}(name, fn)
	}

	//
	for name, fs := range o.funcs {
		res[name] = <-fs.C
	}

	return res, err
}
