// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package internal

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
	time "time"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonF59a38b1DecodeFruitingJobParserInternal(in *jlexer.Lexer, out *JobsInfo) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "PositionToParse":
			out.PositionToParse = Name(in.String())
		case "MinSalary":
			out.MinSalary = Salary(in.Int32())
		case "MaxSalary":
			out.MaxSalary = Salary(in.Int32())
		case "MedianSalary":
			out.MedianSalary = Salary(in.Int32())
		case "PopularSkills":
			if in.IsNull() {
				in.Skip()
				out.PopularSkills = nil
			} else {
				in.Delim('[')
				if out.PopularSkills == nil {
					if !in.IsDelim(']') {
						out.PopularSkills = make(skills, 0, 4)
					} else {
						out.PopularSkills = skills{}
					}
				} else {
					out.PopularSkills = (out.PopularSkills)[:0]
				}
				for !in.IsDelim(']') {
					var v1 string
					v1 = string(in.String())
					out.PopularSkills = append(out.PopularSkills, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "Parser":
			out.Parser = Parser(in.String())
		case "Jobs":
			if in.IsNull() {
				in.Skip()
				out.Jobs = nil
			} else {
				in.Delim('[')
				if out.Jobs == nil {
					if !in.IsDelim(']') {
						out.Jobs = make([]*Job, 0, 8)
					} else {
						out.Jobs = []*Job{}
					}
				} else {
					out.Jobs = (out.Jobs)[:0]
				}
				for !in.IsDelim(']') {
					var v2 *Job
					if in.IsNull() {
						in.Skip()
						v2 = nil
					} else {
						if v2 == nil {
							v2 = new(Job)
						}
						easyjsonF59a38b1DecodeFruitingJobParserInternal1(in, v2)
					}
					out.Jobs = append(out.Jobs, v2)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "Time":
			if in.IsNull() {
				in.Skip()
				out.Time = nil
			} else {
				if out.Time == nil {
					out.Time = new(time.Time)
				}
				if data := in.Raw(); in.Ok() {
					in.AddError((*out.Time).UnmarshalJSON(data))
				}
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonF59a38b1EncodeFruitingJobParserInternal(out *jwriter.Writer, in JobsInfo) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"PositionToParse\":"
		out.RawString(prefix[1:])
		out.String(string(in.PositionToParse))
	}
	{
		const prefix string = ",\"MinSalary\":"
		out.RawString(prefix)
		out.Int32(int32(in.MinSalary))
	}
	{
		const prefix string = ",\"MaxSalary\":"
		out.RawString(prefix)
		out.Int32(int32(in.MaxSalary))
	}
	{
		const prefix string = ",\"MedianSalary\":"
		out.RawString(prefix)
		out.Int32(int32(in.MedianSalary))
	}
	{
		const prefix string = ",\"PopularSkills\":"
		out.RawString(prefix)
		if in.PopularSkills == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v3, v4 := range in.PopularSkills {
				if v3 > 0 {
					out.RawByte(',')
				}
				out.String(string(v4))
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"Parser\":"
		out.RawString(prefix)
		out.String(string(in.Parser))
	}
	{
		const prefix string = ",\"Jobs\":"
		out.RawString(prefix)
		if in.Jobs == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v5, v6 := range in.Jobs {
				if v5 > 0 {
					out.RawByte(',')
				}
				if v6 == nil {
					out.RawString("null")
				} else {
					easyjsonF59a38b1EncodeFruitingJobParserInternal1(out, *v6)
				}
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"Time\":"
		out.RawString(prefix)
		if in.Time == nil {
			out.RawString("null")
		} else {
			out.Raw((*in.Time).MarshalJSON())
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v JobsInfo) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonF59a38b1EncodeFruitingJobParserInternal(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v JobsInfo) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonF59a38b1EncodeFruitingJobParserInternal(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *JobsInfo) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonF59a38b1DecodeFruitingJobParserInternal(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *JobsInfo) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonF59a38b1DecodeFruitingJobParserInternal(l, v)
}
func easyjsonF59a38b1DecodeFruitingJobParserInternal1(in *jlexer.Lexer, out *Job) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "PositionName":
			out.PositionName = Name(in.String())
		case "Link":
			out.Link = string(in.String())
		case "Salary":
			out.Salary = Salary(in.Int32())
		case "Skills":
			if in.IsNull() {
				in.Skip()
				out.Skills = nil
			} else {
				in.Delim('[')
				if out.Skills == nil {
					if !in.IsDelim(']') {
						out.Skills = make(skills, 0, 4)
					} else {
						out.Skills = skills{}
					}
				} else {
					out.Skills = (out.Skills)[:0]
				}
				for !in.IsDelim(']') {
					var v7 string
					v7 = string(in.String())
					out.Skills = append(out.Skills, v7)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonF59a38b1EncodeFruitingJobParserInternal1(out *jwriter.Writer, in Job) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"PositionName\":"
		out.RawString(prefix[1:])
		out.String(string(in.PositionName))
	}
	{
		const prefix string = ",\"Link\":"
		out.RawString(prefix)
		out.String(string(in.Link))
	}
	{
		const prefix string = ",\"Salary\":"
		out.RawString(prefix)
		out.Int32(int32(in.Salary))
	}
	{
		const prefix string = ",\"Skills\":"
		out.RawString(prefix)
		if in.Skills == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v8, v9 := range in.Skills {
				if v8 > 0 {
					out.RawByte(',')
				}
				out.String(string(v9))
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}
