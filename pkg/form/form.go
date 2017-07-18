package form

import (
	"reflect"
	"github.com/go-macaron/binding"
	"gopkg.in/macaron.v1"
	"github.com/Unknwon/com"
	"strings"
	"regexp"
	"fmt"
	log "gopkg.in/clog.v1"
	"io/ioutil"
	"dev.sigpipe.me/dashie/reel2bits/pkg/tool"
	"mime/multipart"
)

const errAlphaDashDotSlash = "AlphaDashDotSlashError"
const errOnlyAudioFile = "OnlyAudioFileError"

var alphaDashDotSlashPattern = regexp.MustCompile("[^\\d\\w-_\\./]")

func init() {
	binding.SetNameMapper(com.ToSnakeCase)
	binding.AddRule(&binding.Rule{
		IsMatch: func(rule string) bool {
			return rule == "AlphaDashDotSlash"
		},
		IsValid: func(errs binding.Errors, name string, v interface{}) (bool, binding.Errors) {
			if alphaDashDotSlashPattern.MatchString(fmt.Sprintf("%v", v)) {
				errs.Add([]string{name}, errAlphaDashDotSlash, "AlphaDashDotSlash")
				return false, errs
			}
			return true, errs
		},
	})

	binding.AddRule(&binding.Rule{
		IsMatch: func(rule string) bool {
			return rule == "OnlyAudioFile"
		},
		IsValid: func(errs binding.Errors, name string, v interface{}) (bool, binding.Errors) {
			fr, err := v.(*multipart.FileHeader).Open()
			if err != nil {
				log.Error(2, "Error opening temp uploaded file to read content: %s", err)
				errs.Add([]string{name}, errOnlyAudioFile, "OnlyAudioFile")
				return false, errs
			}
			defer fr.Close()

			data, err := ioutil.ReadAll(fr)
			if err != nil {
				log.Error(2, "Error reading temp uploaded file: %s", err)
				errs.Add([]string{name}, errOnlyAudioFile, "OnlyAudioFile")
				return false, errs
			}

			mimetype, err := tool.GetBlobMimeType(data)
			if err != nil {
				log.Error(2, "Error reading temp uploaded file: %s", err)
				errs.Add([]string{name}, errOnlyAudioFile, "OnlyAudioFile")
				return false, errs
			}
			log.Trace("Got mimetype %s for file %s", mimetype, v.(*multipart.FileHeader).Filename)

			if mimetype != "audio/mpeg" && mimetype != "audio/x-wav" && mimetype != "audio/ogg" && mimetype != "audio/x-flac" {
				errs.Add([]string{name}, errOnlyAudioFile, "OnlyAudioFile")
				return false, errs
			}

			return true, errs
		},
	})
}

// Form interface for binding validator
type Form interface {
	binding.Validator
}

// Assign assign form values back to the template data.
func Assign(form interface{}, data map[string]interface{}) {
	typ := reflect.TypeOf(form)
	val := reflect.ValueOf(form)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		fieldName := field.Tag.Get("form")
		// Allow ignored fields in the struct
		if fieldName == "-" {
			continue
		} else if len(fieldName) == 0 {
			fieldName = com.ToSnakeCase(field.Name)
		}

		data[fieldName] = val.Field(i).Interface()
	}
}

func getRuleBody(field reflect.StructField, prefix string) string {
	for _, rule := range strings.Split(field.Tag.Get("binding"), ";") {
		if strings.HasPrefix(rule, prefix) {
			return rule[len(prefix) : len(rule)-1]
		}
	}
	return ""
}

func getSize(field reflect.StructField) string {
	return getRuleBody(field, "Size(")
}

func getMinSize(field reflect.StructField) string {
	return getRuleBody(field, "MinSize(")
}

func getMaxSize(field reflect.StructField) string {
	return getRuleBody(field, "MaxSize(")
}

func getInclude(field reflect.StructField) string {
	return getRuleBody(field, "Include(")
}

func getIn(field reflect.StructField) string {
	return getRuleBody(field, "In(")
}

func validate(errs binding.Errors, data map[string]interface{}, f Form, l macaron.Locale) binding.Errors {
	log.Trace("Validating form")

	if errs.Len() == 0 {
		return errs
	}

	data["HasError"] = true
	Assign(f, data)

	typ := reflect.TypeOf(f)
	val := reflect.ValueOf(f)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		fieldName := field.Tag.Get("form")
		// Allow ignored fields in the struct
		if fieldName == "-" {
			continue
		}

		if errs[0].FieldNames[0] == field.Name {
			data["Err_"+field.Name] = true

			trName := field.Tag.Get("locale")
			if len(trName) == 0 {
				trName = l.Tr("form." + field.Name)
			} else {
				trName = l.Tr(trName)
			}

			switch errs[0].Classification {
			case binding.ERR_REQUIRED:
				data["ErrorMsg"] = trName + l.Tr("form.require_error")
			case binding.ERR_ALPHA_DASH:
				data["ErrorMsg"] = trName + l.Tr("form.alpha_dash_error")
			case binding.ERR_ALPHA_DASH_DOT:
				data["ErrorMsg"] = trName + l.Tr("form.alpha_dash_dot_error")
			case errAlphaDashDotSlash:
				data["ErrorMsg"] = trName + l.Tr("form.alpha_dash_dot_slash_error")
			case errOnlyAudioFile:
				data["ErrorMsg"] = trName + l.Tr("form.audio_not_supported")
			case binding.ERR_SIZE:
				data["ErrorMsg"] = trName + l.Tr("form.size_error", getSize(field))
			case binding.ERR_MIN_SIZE:
				data["ErrorMsg"] = trName + l.Tr("form.min_size_error", getMinSize(field))
			case binding.ERR_MAX_SIZE:
				data["ErrorMsg"] = trName + l.Tr("form.max_size_error", getMaxSize(field))
			case binding.ERR_EMAIL:
				data["ErrorMsg"] = trName + l.Tr("form.email_error")
			case binding.ERR_URL:
				data["ErrorMsg"] = trName + l.Tr("form.url_error")
			case binding.ERR_INCLUDE:
				data["ErrorMsg"] = trName + l.Tr("form.include_error", getInclude(field))
			case binding.ERR_IN:
				data["ErrorMsg"] = trName + l.Tr("form.in_error", getIn(field))
			default:
				data["ErrorMsg"] = l.Tr("form.unknown_error") + " " + errs[0].Classification
			}
			return errs
		}
	}
	return errs
}
