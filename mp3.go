package hammer

import (
	"io"
	"os"
)

type Mp3 struct {
	Header *ID3V2_3Header
	Body   []byte
}

func (mp3 *Mp3) Tag(name, value string) {
	//修改frame
	//修改header size
	frame := mp3.Header.frames[0]
	if frame == nil {
		return
	}
	frame.Content(value)
	tags := mp3.Header.Frames()
	mp3.Header.size = reSize(len(tags))
}

func (mp3 *Mp3) Byte() []byte {
	bytes := mp3.Header.Byte()
	bytes = append(bytes, mp3.Body...)
	return bytes
}

func openMp3(f *os.File) ([]byte, error) {
	var chunk []byte
	buf := make([]byte, 1024)
	for {
		n, err := f.Read(buf)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if n == 0 {
			break
		}
		chunk = append(chunk, buf[:n]...)
	}
	return chunk, nil
}

type ID3V2_3Header struct {
	id3v2    [3]byte
	major    byte
	revision byte
	flags    byte
	size     [4]byte
	frames   map[int]*ID3V2_3Frame
}

type ID3V2_3Frame struct {
	frameId [4]byte
	size    [4]byte
	flags   [2]byte
	content []byte
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
	l := len(h.frames)
	for i := 0; i < l; i++ {
		tags = append(tags, h.frames[i].Byte()...)
	}
	return tags
}

func (h *ID3V2_3Header) Length() int {
	return size(h.size) + 10
}

func (h *ID3V2_3Header) ContentSize() int {
	return size(h.size)
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

func Mp3_ID3V2_3(bytes []byte) (mp3 *Mp3) {
	mp3 = new(Mp3)
	mp3.Header = readID3V2_3Header(bytes)
	mp3.Body = bytes[mp3.Header.Length():]
	frames := bytes[10 : mp3.Header.ContentSize()+10]
	m := make(map[int]*ID3V2_3Frame)
	key := 0
	for {
		if len(frames) <= 0 {
			break
		}
		frame := readFrame(frames)
		m[key] = frame
		key++
		i := frame.Length()
		if len(frames) < i {
			break
		}
		frames = frames[i:]
	}
	mp3.Header.frames = m
	return
}

//func Mp3_ID3V2(mp3 []byte) {
//	open, _ := os.Open(`C:\down\NeteaseCloudMusic\test.mp3`)
//	defer open.Close()
//	open01, _ := os.Create(`C:\down\NeteaseCloudMusic\test01.mp3`)
//	defer open01.Close()
//	bts, err := openMp3(open)
//	header := readID3V2_3Header(bts)
//	body := bts[header.Length():]
//	frames := header.content
//	m := make(map[int]*ID3V2_3Frame)
//	var key = 0
//	count := 0
//	for {
//		if len(frames) <= 0 {
//			break
//		}
//		frame := readFrame(frames)
//		count += frame.Length()
//		fmt.Println(frame.String())
//		m[key] = frame
//		key++
//		i := frame.Length()
//		if len(frames) < i {
//			break
//		}
//		frames = frames[i:]
//	}
//	all := make([]byte, 0)
//	for i := 0; i < key; i++ {
//		v2Frame := m[i]
//		if string(v2Frame.frameId[:]) == "TIT2" {
//			fmt.Println(len(v2Frame.content))
//			v2Frame.Content("dsadsa")
//			fmt.Println(len(v2Frame.content))
//		}
//		all = append(all, v2Frame.Byte()...)
//	}
//	header.Content(all)
//	all = header.Byte()
//	all = append(all, body...)
//	write, err := open01.Write(all)
//	fmt.Println(write, err)
//}
