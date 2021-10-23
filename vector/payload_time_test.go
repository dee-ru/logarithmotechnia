package vector

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	emptyNA := []bool{false, false, false}

	testData := []struct {
		name    string
		data    []string
		na      []bool
		outData []string
		isEmpty bool
	}{
		{
			name:    "normal + false na",
			data:    []string{"2006-01-02T15:04:05+07:00", "2021-01-01T12:30:00+03:00", "1800-06-10T11:00:00Z"},
			na:      []bool{false, false, false},
			outData: []string{"2006-01-02T15:04:05+07:00", "2021-01-01T12:30:00+03:00", "1800-06-10T11:00:00Z"},
			isEmpty: false,
		},
		{
			name:    "normal + empty na",
			data:    []string{"2006-01-02T15:04:05+07:00", "2021-01-01T12:30:00+03:00", "1800-06-10T11:00:00Z"},
			na:      []bool{},
			outData: []string{"2006-01-02T15:04:05+07:00", "2021-01-01T12:30:00+03:00", "1800-06-10T11:00:00Z"},
			isEmpty: false,
		},
		{
			name:    "normal + nil na",
			data:    []string{"2006-01-02T15:04:05+07:00", "2021-01-01T12:30:00+03:00", "1800-06-10T11:00:00Z"},
			na:      nil,
			outData: []string{"2006-01-02T15:04:05+07:00", "2021-01-01T12:30:00+03:00", "1800-06-10T11:00:00Z"},
			isEmpty: false,
		},
		{
			name:    "normal + mixed na",
			data:    []string{"2006-01-02T15:04:05+07:00", "2021-01-01T12:30:00+03:00", "1800-06-10T11:00:00Z"},
			na:      []bool{false, false, true},
			outData: []string{"2006-01-02T15:04:05+07:00", "2021-01-01T12:30:00+03:00", "0001-01-01T00:00:00Z"},
			isEmpty: false,
		},
		{
			name:    "normal + incorrect sized na",
			data:    []string{"2006-01-02T15:04:05+07:00", "2021-01-01T12:30:00+03:00", "1800-06-10T11:00:00Z"},
			na:      []bool{false, false, false, false},
			isEmpty: true,
		},
	}

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			timeData := toTimeData(data.data)
			outTimeData := toTimeData(data.outData)

			v := TimeWithNA(timeData, data.na)

			vv := v.(*vector)

			if data.isEmpty {
				naPayload, ok := vv.payload.(*naPayload)
				if !ok || naPayload.Len() > 0 {
					t.Error("Vector's payload is not empty")
				}
			} else {
				length := len(data.data)
				if vv.length != length {
					t.Error(fmt.Sprintf("Vector length (%d) is not equal to data length (%d)\n", vv.length, length))
				}

				payload, ok := vv.payload.(*timePayload)
				if !ok {
					t.Error("Payload is not floatPayload")
				} else {
					if !reflect.DeepEqual(payload.data, outTimeData) {
						t.Error(fmt.Sprintf("Payload data (%v) is not equal to correct data (%v)\n",
							payload.data[1:], timeData))
					}
				}

				if len(data.na) > 0 && len(data.na) == length {
					if !reflect.DeepEqual(payload.na, data.na) {
						t.Error(fmt.Sprintf("Payload na (%v) is not equal to correct na (%v)\n",
							payload.na[1:], data.na))
					}
				} else if len(data.na) == 0 {
					if !reflect.DeepEqual(payload.na, emptyNA) {
						t.Error(fmt.Sprintf("len(data.na) == 0 : incorrect payload.na (%v)", payload.na))
					}
				} else {
					t.Error("error")
				}
			}
		})
	}
}

func TestTimePayload_Type(t *testing.T) {
	vec := TimeWithNA([]time.Time{}, nil)
	if vec.Type() != "time" {
		t.Error("Type is incorrect.")
	}
}

func TestTimePayload_Len(t *testing.T) {
	testData := []struct {
		in        []string
		outLength int
	}{
		{
			in:        []string{"2006-01-02T15:04:05+07:00", "2021-01-01T12:30:00+03:00", "1800-06-10T11:00:00Z"},
			outLength: 3,
		},
		{
			in:        []string{"2006-01-02T15:04:05+07:00"},
			outLength: 1,
		},
		{
			in:        []string{},
			outLength: 0,
		},
		{
			in:        nil,
			outLength: 0,
		},
	}

	for i, data := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			payload := TimeWithNA(toTimeData(data.in), nil).(*vector).payload
			if payload.Len() != data.outLength {
				t.Error(fmt.Sprintf("Payloads's length (%d) is not equal to out (%d)",
					payload.Len(), data.outLength))
			}
		})
	}
}

func TestTimePayload_Strings(t *testing.T) {
	testData := []struct {
		in    []string
		inNA  []bool
		out   []string
		outNA []bool
	}{
		{
			in:    []string{"2006-01-02T15:04:05+07:00", "2021-01-01T12:30:00+03:00", "1800-06-10T11:00:00Z"},
			inNA:  []bool{false, false, false},
			out:   []string{"2006-01-02T15:04:05+07:00", "2021-01-01T12:30:00+03:00", "1800-06-10T11:00:00Z"},
			outNA: []bool{false, false, false},
		},
		{
			in:    []string{"2006-01-02T15:04:05+07:00", "2021-01-01T12:30:00+03:00", "1800-06-10T11:00:00Z"},
			inNA:  []bool{false, true, false},
			out:   []string{"2006-01-02T15:04:05+07:00", "", "1800-06-10T11:00:00Z"},
			outNA: []bool{false, true, false},
		},
	}

	for i, data := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			vec := TimeWithNA(toTimeData(data.in), data.inNA)
			payload := vec.(*vector).payload.(*timePayload)

			strings, na := payload.Strings()
			if !reflect.DeepEqual(strings, data.out) {
				t.Error(fmt.Sprintf("Strings (%v) are not equal to data.out (%v)\n", strings, data.out))
			}
			if !reflect.DeepEqual(na, data.outNA) {
				t.Error(fmt.Sprintf("IsNA (%v) are not equal to data.outNA (%v)\n", na, data.outNA))
			}
		})
	}
}

func TestTimePayload_Times(t *testing.T) {
	testData := []struct {
		in    []string
		inNA  []bool
		out   []time.Time
		outNA []bool
	}{
		{
			in:    []string{"2006-01-02T15:04:05+07:00", "2021-01-01T12:30:00+03:00", "1800-06-10T11:00:00Z"},
			inNA:  []bool{false, false, false},
			out:   toTimeData([]string{"2006-01-02T15:04:05+07:00", "2021-01-01T12:30:00+03:00", "1800-06-10T11:00:00Z"}),
			outNA: []bool{false, false, false},
		},
		{
			in:    []string{"2006-01-02T15:04:05+07:00", "2021-01-01T12:30:00+03:00", "1800-06-10T11:00:00Z"},
			inNA:  []bool{false, true, false},
			out:   toTimeData([]string{"2006-01-02T15:04:05+07:00", "0001-01-01T00:00:00Z", "1800-06-10T11:00:00Z"}),
			outNA: []bool{false, true, false},
		},
	}

	for i, data := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			timeData := toTimeData(data.in)
			vec := TimeWithNA(timeData, data.inNA)
			payload := vec.(*vector).payload.(*timePayload)

			times, na := payload.Times()
			if !reflect.DeepEqual(times, data.out) {
				t.Error(fmt.Sprintf("Times (%v) are not equal to timeData (%v)\n", times, data.out))
			}
			if !reflect.DeepEqual(na, data.outNA) {
				t.Error(fmt.Sprintf("IsNA (%v) are not equal to data.outNA (%v)\n", na, data.outNA))
			}
		})
	}
}

func TestTimePayload_Interfaces(t *testing.T) {
	testData := []struct {
		in    []string
		inNA  []bool
		out   []interface{}
		outNA []bool
	}{
		{
			in:    []string{"2006-01-02T15:04:05+07:00", "2021-01-01T12:30:00+03:00", "1800-06-10T11:00:00Z"},
			inNA:  []bool{false, false, false},
			out:   toTimeInterfaceData([]string{"2006-01-02T15:04:05+07:00", "2021-01-01T12:30:00+03:00", "1800-06-10T11:00:00Z"}),
			outNA: []bool{false, false, false},
		},
		{
			in:    []string{"2006-01-02T15:04:05+07:00", "2021-01-01T12:30:00+03:00", "1800-06-10T11:00:00Z"},
			inNA:  []bool{false, true, false},
			out:   toTimeInterfaceData([]string{"2006-01-02T15:04:05+07:00", "", "1800-06-10T11:00:00Z"}),
			outNA: []bool{false, true, false},
		},
	}

	for i, data := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			timeData := toTimeData(data.in)
			vec := TimeWithNA(timeData, data.inNA)
			payload := vec.(*vector).payload.(*timePayload)

			interfaces, na := payload.Interfaces()
			if !reflect.DeepEqual(interfaces, data.out) {
				t.Error(fmt.Sprintf("Interfaces (%v) are not equal to timeData (%v)\n", interfaces, data.out))
			}
			if !reflect.DeepEqual(na, data.outNA) {
				t.Error(fmt.Sprintf("IsNA (%v) are not equal to data.outNA (%v)\n", na, data.outNA))
			}
		})
	}
}

func TestTimePayload_ByIndices(t *testing.T) {
	vec := TimeWithNA(toTimeData([]string{"2006-01-02T15:04:05+07:00", "2021-01-01T12:30:00+03:00", "1800-06-10T11:00:00Z"}),
		[]bool{false, false, true},
	)
	testData := []struct {
		name    string
		indices []int
		out     []time.Time
		outNA   []bool
	}{
		{
			name:    "all",
			indices: []int{1, 2, 3},
			out:     toTimeData([]string{"2006-01-02T15:04:05+07:00", "2021-01-01T12:30:00+03:00", "0001-01-01T00:00:00Z"}),
			outNA:   []bool{false, false, true},
		},
		{
			name:    "all reverse",
			indices: []int{3, 2, 1},
			out:     toTimeData([]string{"0001-01-01T00:00:00Z", "2021-01-01T12:30:00+03:00", "2006-01-02T15:04:05+07:00"}),
			outNA:   []bool{true, false, false},
		},
		{
			name:    "some",
			indices: []int{3, 1},
			out:     toTimeData([]string{"0001-01-01T00:00:00Z", "2006-01-02T15:04:05+07:00"}),
			outNA:   []bool{true, false},
		},
	}

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			payload := vec.ByIndices(data.indices).(*vector).payload.(*timePayload)
			if !reflect.DeepEqual(payload.data, data.out) {
				t.Error(fmt.Sprintf("payload.data (%v) is not equal to data.out (%v)", payload.data, data.out))
			}
			if !reflect.DeepEqual(payload.na, data.outNA) {
				t.Error(fmt.Sprintf("payload.na (%v) is not equal to data.out (%v)", payload.na, data.out))
			}
		})
	}
}

func TestTimePayload_SupportsWhicher(t *testing.T) {
	testData := []struct {
		name        string
		filter      interface{}
		isSupported bool
	}{
		{
			name:        "func(int, time.Time, bool) bool",
			filter:      func(int, time.Time, bool) bool { return true },
			isSupported: true,
		},
		{
			name:        "func(time.Time, bool) bool",
			filter:      func(time.Time, bool) bool { return true },
			isSupported: true,
		},
		{
			name:        "func(int, int, bool) bool",
			filter:      func(int, int, bool) bool { return true },
			isSupported: false,
		},
	}

	payload := TimeWithNA([]time.Time{}, nil).(*vector).payload.(Whichable)
	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			if payload.SupportsWhicher(data.filter) != data.isSupported {
				t.Error("Selector's support is incorrect.")
			}
		})
	}
}

func TestTimePayload_Whicher(t *testing.T) {
	testData := []struct {
		name   string
		filter interface{}
		out    []bool
	}{
		{
			name:   "func(int, time.Time, bool) bool",
			filter: func(idx int, val time.Time, na bool) bool { return idx == 1 || na == true },
			out:    []bool{true, false, true},
		},
		{
			name:   "func(int, time.Time, bool) bool",
			filter: func(_ time.Time, na bool) bool { return !na },
			out:    []bool{true, true, false},
		},
		{
			name:   "func() bool",
			filter: func() bool { return true },
			out:    []bool{false, false, false},
		},
	}

	payload := TimeWithNA(toTimeData([]string{"2006-01-02T15:04:05+07:00", "2021-01-01T12:30:00+03:00", "1800-06-10T11:00:00Z"}),
		[]bool{false, false, true}).(*vector).payload.(Whichable)

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			filtered := payload.Which(data.filter)
			if !reflect.DeepEqual(payload.Which(data.filter), data.out) {
				t.Error(fmt.Sprintf("payload.Which() (%v) is not equal to out value (%v)",
					filtered, data.out))
			}
		})
	}
}

func TestTimePayload_SupportsApplier(t *testing.T) {
	testData := []struct {
		name        string
		applier     interface{}
		isSupported bool
	}{
		{
			name:        "func(int, time.Time, bool) (time.Time, bool)",
			applier:     func(int, time.Time, bool) (time.Time, bool) { return time.Time{}, true },
			isSupported: true,
		},
		{
			name:        "func(time.Time, bool) (time.Time, bool)",
			applier:     func(time.Time, bool) (time.Time, bool) { return time.Time{}, true },
			isSupported: true,
		},
		{
			name:        "func(int, time.Time, bool) bool",
			applier:     func(int, time.Time, bool) bool { return true },
			isSupported: false,
		},
	}

	payload := TimeWithNA([]time.Time{}, nil).(*vector).payload.(Appliable)
	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			if payload.SupportsApplier(data.applier) != data.isSupported {
				t.Error("Applier's support is incorrect.")
			}
		})
	}
}

func TestTimePayload_Apply(t *testing.T) {
	testData := []struct {
		name        string
		applier     interface{}
		dataIn      []time.Time
		naIn        []bool
		dataOut     []time.Time
		naOut       []bool
		isNAPayload bool
	}{
		{
			name: "regular",
			applier: func(_ int, val time.Time, na bool) (time.Time, bool) {
				return val.Add(24 * time.Hour), na
			},
			dataIn:      toTimeData([]string{"2006-01-02T15:04:05+07:00", "2021-01-01T12:30:00+03:00", "1800-06-10T11:00:00Z"}),
			naIn:        []bool{false, true, false},
			dataOut:     toTimeData([]string{"2006-01-03T15:04:05+07:00", "0001-01-01T00:00:00Z", "1800-06-11T11:00:00Z"}),
			naOut:       []bool{false, true, false},
			isNAPayload: false,
		},
		{
			name: "regular compact",
			applier: func(val time.Time, na bool) (time.Time, bool) {
				return val.Add(24 * time.Hour), na
			},
			dataIn:      toTimeData([]string{"2006-01-02T15:04:05+07:00", "2021-01-01T12:30:00+03:00", "1800-06-10T11:00:00Z"}),
			naIn:        []bool{false, true, false},
			dataOut:     toTimeData([]string{"2006-01-03T15:04:05+07:00", "0001-01-01T00:00:00Z", "1800-06-11T11:00:00Z"}),
			naOut:       []bool{false, true, false},
			isNAPayload: false,
		},
		{
			name: "manipulate na",
			applier: func(idx int, val time.Time, na bool) (time.Time, bool) {
				if idx == 3 {
					return time.Time{}, true
				}
				return val, na
			},
			dataIn:      toTimeData([]string{"2006-01-02T15:04:05+07:00", "2021-01-01T12:30:00+03:00", "1800-06-10T11:00:00Z"}),
			naIn:        []bool{false, false, false},
			dataOut:     toTimeData([]string{"2006-01-02T15:04:05+07:00", "2021-01-01T12:30:00+03:00", "0001-01-01T00:00:00Z"}),
			naOut:       []bool{false, false, true},
			isNAPayload: false,
		},
		{
			name:        "incorrect applier",
			applier:     func(int, string, bool) bool { return true },
			dataIn:      toTimeData([]string{"2006-01-02T15:04:05+07:00", "2021-01-01T12:30:00+03:00", "1800-06-10T11:00:00Z"}),
			naIn:        []bool{false, false, false},
			dataOut:     toTimeData([]string{"0001-01-01T00:00:00Z", "0001-01-01T00:00:00Z", "0001-01-01T00:00:00Z"}),
			naOut:       []bool{true, true, true},
			isNAPayload: true,
		},
	}

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			payload := TimeWithNA(data.dataIn, data.naIn).(*vector).payload.(Appliable).Apply(data.applier)

			if !data.isNAPayload {
				payloadOut := payload.(*timePayload)
				if !reflect.DeepEqual(data.dataOut, payloadOut.data) {
					t.Error(fmt.Sprintf("Output data (%v) does not match expected (%v)",
						payloadOut.data, data.dataOut))
				}
				if !reflect.DeepEqual(data.naOut, payloadOut.na) {
					t.Error(fmt.Sprintf("Output NA (%v) does not match expected (%v)",
						payloadOut.na, data.naOut))
				}
			} else {
				_, ok := payload.(*naPayload)
				if !ok {
					t.Error("Payload is not NA")
				}
			}
		})
	}
}

func toTimeData(times []string) []time.Time {
	timeData := make([]time.Time, len(times))

	for i := 0; i < len(times); i++ {
		t, err := time.Parse(time.RFC3339, times[i])
		if err != nil {
			fmt.Println(err)
		}
		timeData[i] = t
	}

	return timeData
}

func toTimeInterfaceData(times []string) []interface{} {
	timeData := make([]interface{}, len(times))

	for i := 0; i < len(times); i++ {
		if times[i] == "" {
			timeData[i] = nil
		} else {
			t, err := time.Parse(time.RFC3339, times[i])
			if err != nil {
				fmt.Println(err)
			}
			timeData[i] = t
		}
	}

	return timeData
}

func TestTimePayload_SupportsSummarizer(t *testing.T) {
	testData := []struct {
		name        string
		summarizer  interface{}
		isSupported bool
	}{
		{
			name:        "valid",
			summarizer:  func(int, time.Time, time.Time, bool) (time.Time, bool) { return time.Time{}, false },
			isSupported: true,
		},
		{
			name:        "invalid",
			summarizer:  func(int, int, bool) bool { return true },
			isSupported: false,
		},
	}

	payload := TimeWithNA([]time.Time{}, nil).(*vector).payload.(Summarizable)
	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			if payload.SupportsSummarizer(data.summarizer) != data.isSupported {
				t.Error("Summarizer's support is incorrect.")
			}
		})
	}
}

func TestTimePayload_Summarize(t *testing.T) {
	summarizer := func(idx int, prev time.Time, cur time.Time, na bool) (time.Time, bool) {
		if cur.Unix() > prev.Unix() {
			return cur, na
		}

		return prev, na
	}

	testData := []struct {
		name        string
		summarizer  interface{}
		dataIn      []time.Time
		naIn        []bool
		dataOut     []time.Time
		naOut       []bool
		isNAPayload bool
	}{
		{
			name:        "true",
			summarizer:  summarizer,
			dataIn:      toTimeData([]string{"2006-01-02T15:04:05+07:00", "2021-01-01T12:30:00+03:00", "1800-06-10T11:00:00Z"}),
			naIn:        []bool{false, false, false},
			dataOut:     toTimeData([]string{"2021-01-01T12:30:00+03:00"}),
			naOut:       []bool{false},
			isNAPayload: false,
		},
		{
			name:        "NA",
			summarizer:  summarizer,
			dataIn:      toTimeData([]string{"2006-01-02T15:04:05+07:00", "2021-01-01T12:30:00+03:00", "1800-06-10T11:00:00Z"}),
			naIn:        []bool{false, false, true},
			isNAPayload: true,
		},
		{
			name:        "incorrect summarizer",
			summarizer:  func(int, int, bool) bool { return true },
			dataIn:      toTimeData([]string{"2006-01-02T15:04:05+07:00", "2021-01-01T12:30:00+03:00", "1800-06-10T11:00:00Z"}),
			naIn:        []bool{false, false, false},
			isNAPayload: true,
		},
	}

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			payload := TimeWithNA(data.dataIn, data.naIn).(*vector).payload.(Summarizable).Summarize(data.summarizer)

			if !data.isNAPayload {
				payloadOut := payload.(*timePayload)
				if !reflect.DeepEqual(data.dataOut, payloadOut.data) {
					t.Error(fmt.Sprintf("Output data (%v) does not match expected (%v)",
						data.dataOut, payloadOut.data))
				}
				if !reflect.DeepEqual(data.naOut, payloadOut.na) {
					t.Error(fmt.Sprintf("Output NA (%v) does not match expected (%v)",
						data.naOut, payloadOut.na))
				}
			} else {
				naPayload, ok := payload.(*naPayload)
				if ok {
					if naPayload.length != 1 {
						t.Error("Incorrect length of NA payload (not 1)")
					}
				} else {
					t.Error("Payload is not NA")
				}
			}
		})
	}
}

func TestTimePayload_Append(t *testing.T) {
	payload := TimePayload(toTimeData([]string{"2006-01-02T15:04:05+07:00", "2021-01-01T12:30:00+03:00", "1800-06-10T11:00:00Z"}),
		nil)

	testData := []struct {
		name    string
		vec     Vector
		outData []time.Time
		outNA   []bool
	}{
		{
			name: "times",
			vec:  TimeWithNA(toTimeData([]string{"2026-01-02T15:04:05+07:00", "2023-01-01T12:30:00+03:00"}), []bool{true, false}),
			outData: toTimeData([]string{"2006-01-02T15:04:05+07:00", "2021-01-01T12:30:00+03:00", "1800-06-10T11:00:00Z",
				"0001-01-01T00:00:00Z", "2023-01-01T12:30:00+03:00"}),
			outNA: []bool{false, false, false, true, false},
		},
		{
			name: "integer",
			vec:  IntegerWithNA([]int{4, 5}, []bool{true, false}),
			outData: toTimeData([]string{"2006-01-02T15:04:05+07:00", "2021-01-01T12:30:00+03:00", "1800-06-10T11:00:00Z",
				"0001-01-01T00:00:00Z", "0001-01-01T00:00:00Z"}),
			outNA: []bool{false, false, false, true, true},
		},
		{
			name: "na",
			vec:  NA(2),
			outData: toTimeData([]string{"2006-01-02T15:04:05+07:00", "2021-01-01T12:30:00+03:00", "1800-06-10T11:00:00Z",
				"0001-01-01T00:00:00Z", "0001-01-01T00:00:00Z"}),
			outNA: []bool{false, false, false, true, true},
		},
	}

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			outPayload := payload.Append(data.vec.Payload()).(*timePayload)

			if !reflect.DeepEqual(data.outData, outPayload.data) {
				t.Error(fmt.Sprintf("Output data (%v) does not match expected (%v)",
					outPayload.data, data.outData))
			}
			if !reflect.DeepEqual(data.outNA, outPayload.na) {
				t.Error(fmt.Sprintf("Output NA (%v) does not match expected (%v)",
					outPayload.na, data.outNA))
			}
		})
	}
}

func TestTimePayload_Adjust(t *testing.T) {
	payload5 := TimePayload(toTimeData([]string{"2006-01-02T15:04:05+07:00", "2021-01-01T12:30:00+03:00", "1800-06-10T11:00:00Z",
		"0001-01-01T00:00:00Z", "0001-01-01T00:00:00Z"}), nil).(*timePayload)
	payload3 := TimePayload(toTimeData([]string{"2006-01-02T15:04:05+07:00", "2021-01-01T12:30:00+03:00", "1800-06-10T11:00:00Z"}),
		[]bool{false, false, true}).(*timePayload)

	testData := []struct {
		name       string
		inPayload  *timePayload
		size       int
		outPaylout *timePayload
	}{
		{
			inPayload: payload5,
			name:      "same",
			size:      5,
			outPaylout: TimePayload(toTimeData([]string{"2006-01-02T15:04:05+07:00", "2021-01-01T12:30:00+03:00", "1800-06-10T11:00:00Z",
				"0001-01-01T00:00:00Z", "0001-01-01T00:00:00Z"}), nil).(*timePayload),
		},
		{
			inPayload: payload5,
			name:      "lesser",
			size:      3,
			outPaylout: TimePayload(toTimeData([]string{"2006-01-02T15:04:05+07:00", "2021-01-01T12:30:00+03:00", "1800-06-10T11:00:00Z"}),
				nil).(*timePayload),
		},
		{
			inPayload: payload3,
			name:      "bigger",
			size:      10,
			outPaylout: TimePayload(toTimeData([]string{"2006-01-02T15:04:05+07:00", "2021-01-01T12:30:00+03:00", "0001-01-01T00:00:00Z",
				"2006-01-02T15:04:05+07:00", "2021-01-01T12:30:00+03:00", "0001-01-01T00:00:00Z",
				"2006-01-02T15:04:05+07:00", "2021-01-01T12:30:00+03:00", "0001-01-01T00:00:00Z",
				"2006-01-02T15:04:05+07:00"}),
				[]bool{false, false, true, false, false, true, false, false, true, false}).(*timePayload),
		},
	}

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			outPayload := data.inPayload.Adjust(data.size).(*timePayload)

			if !reflect.DeepEqual(outPayload.data, data.outPaylout.data) {
				t.Error(fmt.Sprintf("Output data (%v) does not match expected (%v)",
					outPayload.data, data.outPaylout.data))
			}
			if !reflect.DeepEqual(outPayload.na, data.outPaylout.na) {
				t.Error(fmt.Sprintf("Output NA (%v) does not match expected (%v)",
					outPayload.na, data.outPaylout.na))
			}
		})
	}
}

func TestTimePayload_Find(t *testing.T) {
	payload := TimePayload(toTimeData(
		[]string{"2006-01-02T15:04:05+07:00", "2021-01-01T12:30:00+03:00", "2006-01-02T15:04:05+07:00"},
	), nil).(*timePayload)

	existent, _ := time.Parse(time.RFC3339, "2021-01-01T12:30:00+03:00")
	nonExistent, _ := time.Parse(time.RFC3339, "2026-05-12T10:00:00+00:00")

	testData := []struct {
		name   string
		needle interface{}
		pos    int
	}{
		{"existent", existent, 2},
		{"non-existent", nonExistent, 0},
		{"incorrect type", true, 0},
	}

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			pos := payload.Find(data.needle)

			if pos != data.pos {
				t.Error(fmt.Sprintf("Position (%v) does not match expected (%v)",
					pos, data.pos))
			}
		})
	}
}

func TestTimePayload_FindAll(t *testing.T) {
	payload := TimePayload(toTimeData(
		[]string{"2006-01-02T15:04:05+07:00", "2021-01-01T12:30:00+03:00", "2006-01-02T15:04:05+07:00"},
	), nil).(*timePayload)

	existent, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05+07:00")
	nonExistent, _ := time.Parse(time.RFC3339, "2026-05-12T10:00:00+00:00")

	testData := []struct {
		name   string
		needle interface{}
		pos    []int
	}{
		{"existent", existent, []int{1, 3}},
		{"non-existent", nonExistent, []int{}},
		{"incorrect type", false, []int{}},
	}

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			pos := payload.FindAll(data.needle)

			if !reflect.DeepEqual(pos, data.pos) {
				t.Error(fmt.Sprintf("Positions (%v) does not match expected (%v)",
					pos, data.pos))
			}
		})
	}
}

func TestTimePayload_Eq(t *testing.T) {
	payload := TimePayload(toTimeData([]string{"2006-01-02T15:04:05+07:00", "2021-01-01T12:30:00+03:00",
		"2020-01-01T12:30:00+03:00", "2020-01-01T12:30:00+03:00"}), []bool{false, false, true, false}).(*timePayload)
	date, _ := time.Parse(time.RFC3339, "2020-01-01T12:30:00+03:00")

	testData := []struct {
		eq  interface{}
		cmp []bool
	}{
		{date, []bool{false, false, false, true}},

		{true, []bool{false, false, false, false}},
	}

	for i, data := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			cmp := payload.Eq(data.eq)

			if !reflect.DeepEqual(cmp, data.cmp) {
				t.Error(fmt.Sprintf("Comparator results (%v) do not match expected (%v)",
					cmp, data.cmp))
			}
		})
	}
}

func TestTimePayload_Neq(t *testing.T) {
	payload := TimePayload(toTimeData([]string{"2006-01-02T15:04:05+07:00", "2021-01-01T12:30:00+03:00",
		"2020-01-01T12:30:00+03:00", "2020-01-01T12:30:00+03:00"}), []bool{false, false, true, false}).(*timePayload)
	date, _ := time.Parse(time.RFC3339, "2020-01-01T12:30:00+03:00")

	testData := []struct {
		eq  interface{}
		cmp []bool
	}{
		{date, []bool{true, true, true, false}},

		{true, []bool{true, true, true, true}},
	}

	for i, data := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			cmp := payload.Neq(data.eq)

			if !reflect.DeepEqual(cmp, data.cmp) {
				t.Error(fmt.Sprintf("Comparator results (%v) do not match expected (%v)",
					cmp, data.cmp))
			}
		})
	}
}

func TestTimePayload_Gt(t *testing.T) {
	payload := TimePayload(toTimeData([]string{"2006-01-02T15:04:05+07:00", "2021-01-01T12:30:00+03:00",
		"2006-01-02T15:04:05+07:00", "2020-01-01T12:30:00+03:00"}), []bool{false, false, true, false}).(*timePayload)
	date, _ := time.Parse(time.RFC3339, "2020-01-01T12:30:00+03:00")

	testData := []struct {
		eq  interface{}
		cmp []bool
	}{
		{date, []bool{false, true, false, false}},

		{true, []bool{false, false, false, false}},
	}

	for i, data := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			cmp := payload.Gt(data.eq)

			if !reflect.DeepEqual(cmp, data.cmp) {
				t.Error(fmt.Sprintf("Comparator results (%v) do not match expected (%v)",
					cmp, data.cmp))
			}
		})
	}
}

func TestTimePayload_Lt(t *testing.T) {
	payload := TimePayload(toTimeData([]string{"2006-01-02T15:04:05+07:00", "2021-01-01T12:30:00+03:00",
		"2021-01-01T12:30:00+03:00", "2020-01-01T12:30:00+03:00"}), []bool{false, false, true, false}).(*timePayload)
	date, _ := time.Parse(time.RFC3339, "2020-01-01T12:30:00+03:00")

	testData := []struct {
		eq  interface{}
		cmp []bool
	}{
		{date, []bool{true, false, false, false}},

		{true, []bool{false, false, false, false}},
	}

	for i, data := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			cmp := payload.Lt(data.eq)

			if !reflect.DeepEqual(cmp, data.cmp) {
				t.Error(fmt.Sprintf("Comparator results (%v) do not match expected (%v)",
					cmp, data.cmp))
			}
		})
	}
}

func TestTimePayload_Gte(t *testing.T) {
	payload := TimePayload(toTimeData([]string{"2006-01-02T15:04:05+07:00", "2021-01-01T12:30:00+03:00",
		"2020-01-01T12:30:00+03:00", "2020-01-01T12:30:00+03:00"}), []bool{false, false, true, false}).(*timePayload)
	date, _ := time.Parse(time.RFC3339, "2020-01-01T12:30:00+03:00")

	testData := []struct {
		eq  interface{}
		cmp []bool
	}{
		{date, []bool{false, true, false, true}},

		{true, []bool{false, false, false, false}},
	}

	for i, data := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			cmp := payload.Gte(data.eq)

			if !reflect.DeepEqual(cmp, data.cmp) {
				t.Error(fmt.Sprintf("Comparator results (%v) do not match expected (%v)",
					cmp, data.cmp))
			}
		})
	}
}

func TestTimePayload_Lte(t *testing.T) {
	payload := TimePayload(toTimeData([]string{"2006-01-02T15:04:05+07:00", "2021-01-01T12:30:00+03:00",
		"2020-01-01T12:30:00+03:00", "2020-01-01T12:30:00+03:00"}), []bool{false, false, true, false}).(*timePayload)
	date, _ := time.Parse(time.RFC3339, "2020-01-01T12:30:00+03:00")

	testData := []struct {
		eq  interface{}
		cmp []bool
	}{
		{date, []bool{true, false, false, true}},

		{true, []bool{false, false, false, false}},
	}

	for i, data := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			cmp := payload.Lte(data.eq)

			if !reflect.DeepEqual(cmp, data.cmp) {
				t.Error(fmt.Sprintf("Comparator results (%v) do not match expected (%v)",
					cmp, data.cmp))
			}
		})
	}
}

func TestTimePayload_Groups(t *testing.T) {
	testData := []struct {
		name    string
		payload Payload
		groups  [][]int
		values  []interface{}
	}{
		{
			name: "normal",
			payload: TimePayload(toTimeData([]string{"2006-01-02T15:04:05+07:00", "2021-01-01T12:30:00+03:00",
				"2020-01-01T12:30:00+03:00", "2020-01-01T12:30:00+03:00"}), nil),
			groups: [][]int{{1}, {2}, {3, 4}},
			values: []interface{}{
				toTimeData([]string{"2006-01-02T15:04:05+07:00"})[0],
				toTimeData([]string{"2021-01-01T12:30:00+03:00"})[0],
				toTimeData([]string{"2020-01-01T12:30:00+03:00"})[0],
			},
		},
		{
			name: "with NA",
			payload: TimePayload(toTimeData([]string{"2006-01-02T15:04:05+07:00", "2021-01-01T12:30:00+03:00",
				"2020-01-01T12:30:00+03:00", "2020-01-01T12:30:00+03:00"}), []bool{false, true, false, false}),
			groups: [][]int{{1}, {3, 4}, {2}},
			values: []interface{}{
				toTimeData([]string{"2006-01-02T15:04:05+07:00"})[0],
				toTimeData([]string{"2020-01-01T12:30:00+03:00"})[0],
				nil,
			},
		},
	}

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			groups, _ := data.payload.(*timePayload).Groups()

			if !reflect.DeepEqual(groups, data.groups) {
				t.Error(fmt.Sprintf("Groups (%v) do not match expected (%v)",
					groups, data.groups))
			}
		})
	}
}
