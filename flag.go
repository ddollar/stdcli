package stdcli

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Flag struct {
	Default     interface{}
	Description string
	Name        string
	Kind        reflect.Kind
	Short       string
	Value       interface{}
}

// func StringFlag(name, short, description string, def interface{}) Flag {
//   return Flag{Name: name, Short: short, Description: description, Default: def, Kind: "string"}
// }

func (f *Flag) Set(v string) error {
	switch f.Kind {
	case reflect.Bool:
		f.Value = (v == "true")
	case reflect.Int:
		fmt.Printf("v = %+v\n", v)
		i, err := strconv.Atoi(v)
		if err != nil {
			return err
		}
		f.Value = i
	case reflect.String:
		f.Value = v
	default:
		return fmt.Errorf("unknown flag type: %s", f.Kind)
	}

	return nil
}

func (f *Flag) String() string {
	return fmt.Sprintf("%s: %v", f.Name, f.Value)
}

func (f *Flag) Type() string {
	switch f.Kind {
	case reflect.Bool:
		return "bool"
	case reflect.Int:
		return "int"
	case reflect.String:
		return "string"
	default:
		panic(fmt.Sprintf("unknown flag type: %s", f.Kind))
	}
}

func (f *Flag) Usage(v string) string {
	switch f.Kind {
	case reflect.Bool:
		return v
	case reflect.Int, reflect.String:
		return fmt.Sprintf("%s <u><info>%s</info></u>", v, f.Name)
	default:
		fmt.Printf("f = %+v\n", f)
		panic(fmt.Sprintf("unknown flag type: %s", f.Kind))
	}
}

func (f *Flag) UsageLong() string {
	return f.Usage(fmt.Sprintf("--%s", f.Name))
}

func (f *Flag) UsageShort() string {
	if f.Short == "" {
		return ""
	}

	return f.Usage(fmt.Sprintf("-%s", f.Short))
}

func OptionFlags(opts interface{}) []Flag {
	flags := []Flag{}

	v := reflect.ValueOf(opts)
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		if n := f.Tag.Get("flag"); n != "" {
			parts := strings.Split(n, ",")
			flag := Flag{
				Default:     f.Tag.Get("default"),
				Description: f.Tag.Get("desc"),
				Kind:        f.Type.Elem().Kind(),
				Name:        parts[0],
			}
			if len(parts) > 1 {
				flag.Short = parts[1]
			}
			flags = append(flags, flag)
		}
	}

	return flags
}

func StringFlag(name, short, description string) Flag {
	return Flag{
		Description: description,
		Kind:        reflect.String,
		Name:        name,
		Short:       short,
	}
}
