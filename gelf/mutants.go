package gelf

import (
	"encoding/json"
	"fmt"
	"strings"
)

type MessageMutant func(*Message) error

type Mutant func(rawMsg []byte) MessageMutant

func DefaultParser(rawMsg []byte) MessageMutant {
	return func(m *Message) error {
		i := make(map[string]interface{}, 16)
		if err := json.Unmarshal(rawMsg, &i); err != nil {
			return err
		}
		for k, v := range i {
			if k[0] == '_' {
				if m.Extra == nil {
					m.Extra = make(map[string]interface{}, 1)
				}
				m.Extra[k] = v
				continue
			}

			ok := true
			switch k {
			case "version":
				m.Version, ok = v.(string)
			case "host":
				m.Host, ok = v.(string)
			case "short_message":
				m.Short, ok = v.(string)
			case "full_message":
				m.Full, ok = v.(string)
			case "timestamp":
				m.TimeUnix, ok = v.(float64)
			case "level":
				var level float64
				level, ok = v.(float64)
				m.Level = int32(level)
			case "facility":
				m.Facility, ok = v.(string)
			}

			if !ok {
				return fmt.Errorf("invalid type for field %s", k)
			}
		}
		return nil
	}
}

func ZeroLogParser(rawMsg []byte) MessageMutant {
	return func(m *Message) error {
		var (
			value string
			list  []string
		)
		i := make(map[string]interface{}, 16)
		if err := json.Unmarshal(rawMsg, &i); err != nil {
			return err
		}
		for k, v := range i {
			ok := true
			switch k {
			case "error":
				m.Extra["error"], ok = v.(string)
			case "caller":
				value, ok = v.(string)
				if ok {
					value = srcFile(value)
					list = strings.Split(value, ":")
					m.Extra["caller"] = list[0]
					if len(list) > 1 {
						m.Extra["caller_line"] = list[1]
					}
				}
			case "time":
				m.Extra["time"], ok = v.(string)
			case "version":
				m.Version, ok = v.(string)
			case "host":
				m.Host, ok = v.(string)
			case "short_message":
				m.Short, ok = v.(string)
			case "message":
				m.Short, ok = v.(string)
			case "full_message":
				m.Full, ok = v.(string)
			case "timestamp":
				m.TimeUnix, ok = v.(float64)
			case "level":
				switch v.(type) {
				case float64:
					{
						var level float64
						if level, ok = v.(float64); ok {
							m.Level = int32(level)
							continue
						}
					}
				case string:
					{
						if value, ok = v.(string); ok {
							switch value {
							case "debug":
								{
									m.Level = LOG_DEBUG
									continue
								}

							case "info":
								{
									m.Level = LOG_INFO
									continue
								}
							case "warn":
								{
									m.Level = LOG_WARNING
									continue
								}
							case "error":
								{
									m.Level = LOG_ERR
									continue
								}
							case "fatal":
								{
									m.Level = LOG_CRIT
									continue
								}
							case "panic":
								{
									m.Level = LOG_CRIT
									continue
								}
							case "notice":
								{
									m.Level = LOG_NOTICE
									continue
								}
							case "trace":
								{
									m.Level = -1
									continue
								}

							}
							m.Level = LOG_EMERG
							continue
						}
					}
				}
			case "facility":
				m.Facility, ok = v.(string)
			default:
				m.Extra[k] = v
			}

			if !ok {
				return fmt.Errorf("invalid type for field %s", k)
			}
		}
		return nil
	}
}

func Caller(rawMsg []byte) MessageMutant {
	const callDepth = 2
	return func(m *Message) error {
		m.Extra["file"], m.Extra["line"] = getCaller(callDepth, "/pkg/log/log.go", "/pkg/io/multi.go")
		return nil
	}
}
