package mp3

import "github.com/pkg/errors"

type Mp3 struct {
	Header *ID3V2_3Header
	Body   []byte
}

func (mp3 *Mp3) Tag(frameID FrameID, value string) {
	frame := mp3.Header.frames[frameID.Name()]
	if frame == nil {
		frame = mp3.Header.NewFrame(frameID)
	}
	frame.Content(value)
	tags := mp3.Header.Frames()
	mp3.Header.size = reSize(len(tags))
}

func (mp3 *Mp3) Tags() (m map[string]string) {
	m = make(map[string]string)
	for k, v := range mp3.Header.frames {
		m[k] = string(v.content)
	}
	return
}

func (mp3 *Mp3) Byte() []byte {
	bytes := mp3.Header.Byte()
	bytes = append(bytes, mp3.Body...)
	var end [10]byte
	bytes = append(bytes, end[:]...)
	return bytes
}

type ID3V2_3Header struct {
	id3v2    [3]byte
	major    byte
	revision byte
	flags    byte
	size     [4]byte
	frames   map[string]*ID3V2_3Frame
}

func (h *ID3V2_3Header) NewFrame(frameID FrameID) *ID3V2_3Frame {
	if f := h.frames[frameID.Name()]; f != nil {
		return f
	}
	f := &ID3V2_3Frame{}
	f.frameId = frameID
	h.frames[frameID.Name()] = f
	return f
}

func (h *ID3V2_3Header) ID3V2_3() bool {
	return string(h.id3v2[:]) == "ID3" && h.major == 3
}

func (h *ID3V2_3Header) Byte() []byte {
	b := make([]byte, 0)
	b = append(b, h.id3v2[:]...)
	b = append(b, h.major)
	b = append(b, h.revision)
	b = append(b, h.flags)
	b = append(b, h.size[:]...)
	b = append(b, h.Frames()...)
	return b
}

func (h *ID3V2_3Header) Frames() []byte {
	tags := make([]byte, 0)
	for _, v := range h.frames {
		tags = append(tags, v.Byte()...)
	}
	return tags
}

func (h *ID3V2_3Header) Length() int {
	return size(h.size) + 10
}

func (h *ID3V2_3Header) ContentSize() int {
	return size(h.size)
}

type ID3V2_3Frame struct {
	frameId FrameID
	size    [4]byte
	flags   [2]byte
	content []byte
}

func (frame *ID3V2_3Frame) Content(content string) {
	c := make([]byte, 0)
	c = append(c, byte(0))
	c = append(c, []byte(content)...)
	c = append(c, byte(0))
	frame.content = c
	l := len(frame.content)
	frame.size = reSize(l)
}

func (frame *ID3V2_3Frame) Byte() []byte {
	b := make([]byte, 0)
	b = append(b, frame.frameId[:]...)
	b = append(b, frame.size[:]...)
	b = append(b, frame.flags[:]...)
	b = append(b, frame.content...)
	return b
}

func (frame *ID3V2_3Frame) Length() int {
	return size(frame.size) + 10
}

func (frame *ID3V2_3Frame) ContentSize() int {
	return size(frame.size)
}

type FrameID [4]byte

func (f *FrameID) Name() string {
	return string(f[:])
}

//Declared ID3v2 frames https://id3.org/id3v2.3.0#sec4
var (
	AENC FrameID = [4]byte{65, 69, 78, 67}
	APIC FrameID = [4]byte{65, 80, 73, 67}
	COMM FrameID = [4]byte{67, 79, 77, 77}
	COMR FrameID = [4]byte{67, 79, 77, 82}
	ENCR FrameID = [4]byte{69, 78, 67, 82}
	EQUA FrameID = [4]byte{69, 81, 85, 65}
	ETCO FrameID = [4]byte{69, 84, 67, 79}
	GEOB FrameID = [4]byte{71, 69, 79, 66}
	GRID FrameID = [4]byte{71, 82, 73, 68}
	IPLS FrameID = [4]byte{73, 80, 76, 83}
	LINK FrameID = [4]byte{76, 73, 78, 75}
	MCDI FrameID = [4]byte{77, 67, 68, 73}
	MLLT FrameID = [4]byte{77, 76, 76, 84}
	OWNE FrameID = [4]byte{79, 87, 78, 69}
	PRIV FrameID = [4]byte{80, 82, 73, 86}
	PCNT FrameID = [4]byte{80, 67, 78, 84}
	POPM FrameID = [4]byte{80, 79, 80, 77}
	POSS FrameID = [4]byte{80, 79, 83, 83}
	RBUF FrameID = [4]byte{82, 66, 85, 70}
	RVAD FrameID = [4]byte{82, 86, 65, 68}
	RVRB FrameID = [4]byte{82, 86, 82, 66}
	SYLT FrameID = [4]byte{83, 89, 76, 84}
	SYTC FrameID = [4]byte{83, 89, 84, 67}
	TALB FrameID = [4]byte{84, 65, 76, 66}
	TBPM FrameID = [4]byte{84, 66, 80, 77}
	TCOM FrameID = [4]byte{84, 67, 79, 77}
	TCON FrameID = [4]byte{84, 67, 79, 78}
	TCOP FrameID = [4]byte{84, 67, 79, 80}
	TDAT FrameID = [4]byte{84, 68, 65, 84}
	TDLY FrameID = [4]byte{84, 68, 76, 89}
	TENC FrameID = [4]byte{84, 69, 78, 67}
	TEXT FrameID = [4]byte{84, 69, 88, 84}
	TFLT FrameID = [4]byte{84, 70, 76, 84}
	TIME FrameID = [4]byte{84, 73, 77, 69}
	TIT1 FrameID = [4]byte{84, 73, 84, 49}
	TIT2 FrameID = [4]byte{84, 73, 84, 50}
	TIT3 FrameID = [4]byte{84, 73, 84, 51}
	TKEY FrameID = [4]byte{84, 75, 69, 89}
	TLAN FrameID = [4]byte{84, 76, 65, 78}
	TLEN FrameID = [4]byte{84, 76, 69, 78}
	TMED FrameID = [4]byte{84, 77, 69, 68}
	TOAL FrameID = [4]byte{84, 79, 65, 76}
	TOFN FrameID = [4]byte{84, 79, 70, 78}
	TOLY FrameID = [4]byte{84, 79, 76, 89}
	TOPE FrameID = [4]byte{84, 79, 80, 69}
	TORY FrameID = [4]byte{84, 79, 82, 89}
	TOWN FrameID = [4]byte{84, 79, 87, 78}
	TPE1 FrameID = [4]byte{84, 80, 69, 49}
	TPE2 FrameID = [4]byte{84, 80, 69, 50}
	TPE3 FrameID = [4]byte{84, 80, 69, 51}
	TPE4 FrameID = [4]byte{84, 80, 69, 52}
	TPOS FrameID = [4]byte{84, 80, 79, 83}
	TPUB FrameID = [4]byte{84, 80, 85, 66}
	TRCK FrameID = [4]byte{84, 82, 67, 75}
	TRDA FrameID = [4]byte{84, 82, 68, 65}
	TRSN FrameID = [4]byte{84, 82, 83, 78}
	TRSO FrameID = [4]byte{84, 82, 83, 79}
	TSIZ FrameID = [4]byte{84, 83, 73, 90}
	TSRC FrameID = [4]byte{84, 83, 82, 67}
	TSSE FrameID = [4]byte{84, 83, 83, 69}
	TYER FrameID = [4]byte{84, 89, 69, 82}
	TXXX FrameID = [4]byte{84, 88, 88, 88}
	UFID FrameID = [4]byte{85, 70, 73, 68}
	USER FrameID = [4]byte{85, 83, 69, 82}
	USLT FrameID = [4]byte{85, 83, 76, 84}
	WCOM FrameID = [4]byte{87, 67, 79, 77}
	WCOP FrameID = [4]byte{87, 67, 79, 80}
	WOAF FrameID = [4]byte{87, 79, 65, 70}
	WOAR FrameID = [4]byte{87, 79, 65, 82}
	WOAS FrameID = [4]byte{87, 79, 65, 83}
	WORS FrameID = [4]byte{87, 79, 82, 83}
	WPAY FrameID = [4]byte{87, 80, 65, 89}
	WPUB FrameID = [4]byte{87, 80, 85, 66}
	WXXX FrameID = [4]byte{87, 88, 88, 88}
)

func size(b [4]byte) int {
	l := int(b[0]) << 21
	l += int(b[1]) << 14
	l += int(b[2]) << 7
	l += int(b[3])
	return l
}

func reSize(l int) [4]byte {
	var b [4]byte
	b[3] = byte(l & 127)
	b[2] = byte(l >> 7 & 127)
	b[1] = byte(l >> 14 & 127)
	b[0] = byte(l >> 21 & 127)
	return b
}

func readID3V2_3Header(bts []byte) (header *ID3V2_3Header) {
	header = &ID3V2_3Header{}
	if len(bts) < 11 {
		return
	}
	header.id3v2[0] = bts[0]
	header.id3v2[1] = bts[1]
	header.id3v2[2] = bts[2]
	header.major = bts[3]
	header.revision = bts[4]
	header.flags = bts[5]
	header.size[0] = bts[6]
	header.size[1] = bts[7]
	header.size[2] = bts[8]
	header.size[3] = bts[9]
	return
}

func readFrame(bts []byte) *ID3V2_3Frame {
	frame := &ID3V2_3Frame{}
	frame.frameId[0] = bts[0]
	frame.frameId[1] = bts[1]
	frame.frameId[2] = bts[2]
	frame.frameId[3] = bts[3]
	frame.size[0] = bts[4]
	frame.size[1] = bts[5]
	frame.size[2] = bts[6]
	frame.size[3] = bts[7]
	frame.flags[0] = bts[8]
	frame.flags[1] = bts[9]
	frame.content = bts[10:frame.Length()]
	return frame
}

func Mp3_ID3V2_3(bytes []byte) (mp3 *Mp3, err error) {
	mp3 = new(Mp3)
	mp3.Header = readID3V2_3Header(bytes)
	if !mp3.Header.ID3V2_3() {
		return nil, errors.New("not ID3V2_3")
	}
	mp3.Body = bytes[mp3.Header.Length():]
	frames := bytes[10 : mp3.Header.ContentSize()+10]
	m := make(map[string]*ID3V2_3Frame)
	for {
		if len(frames) < 10 {
			break
		}
		frame := readFrame(frames)
		if len(frame.content) != 0 {
			m[string(frame.frameId[:])] = frame
		}
		i := frame.Length()
		if len(frames) < i {
			break
		}
		frames = frames[i:]
	}
	mp3.Header.frames = m
	return
}
