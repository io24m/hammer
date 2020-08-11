package util

import "io"

func ReadBytes(reader io.Reader) ([]byte, error) {
	var chunk []byte
	buf := make([]byte, 1024)
	for {
		n, err := reader.Read(buf)
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
