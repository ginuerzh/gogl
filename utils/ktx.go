// ktx
package utils

import (
	"bufio"
	"bytes"
	"encoding/binary"
	gl "github.com/chsc/gogl/gl42"
	"io/ioutil"
	"log"
	"os"
	"unsafe"
)

type header struct {
	Identifier           [12]byte
	Endianness           uint32
	Gltype               uint32
	Gltypesize           uint32
	Glformat             uint32
	Glinternalformat     uint32
	Glbaseinternalformat uint32
	Pixelwidth           uint32
	Pixelheight          uint32
	Pixeldepth           uint32
	Arrayelements        uint32
	Faces                uint32
	Miplevels            uint32
	Keypairbytes         uint32
}

var (
	identifier = []byte{0xAB, 0x4B, 0x54, 0x58, 0x20, 0x31, 0x31, 0xBB, 0x0D, 0x0A, 0x1A, 0x0A}
)

func swap32(u32 uint32) uint32 {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, u32)
	if err != nil {
		log.Fatal("binary.Write failed:", err)
	}

	if buf.Len() != 4 {
		log.Fatal("length error")
	}

	u8 := buf.Bytes()
	u8[0], u8[1], u8[2], u8[3] = u8[3], u8[2], u8[1], u8[0]

	err = binary.Read(buf, binary.LittleEndian, &u32)
	if err != nil {
		log.Fatal("binary.Read failed:", err)
	}
	return u32
}

func swap16(u16 uint16) uint16 {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, u16)
	if err != nil {
		log.Fatalln("binary.Write failed:", err)
	}
	if buf.Len() != 2 {
		log.Fatal("length error")
	}

	u8 := buf.Bytes()
	u8[0], u8[1] = u8[1], u8[0]

	err = binary.Read(buf, binary.LittleEndian, &u16)
	if err != nil {
		log.Fatal("binary.Read failed:", err)
	}
	return u16
}

func calcStride(h *header, width, pad uint32) uint32 {
	var channels uint32 = 0

	switch h.Glbaseinternalformat {
	case gl.RED:
		channels = 1
	case gl.RG:
		channels = 2
	case gl.BGR, gl.RGB:
		channels = 3
	case gl.BGRA, gl.RGBA:
		channels = 4
	}

	var stride uint32 = h.Gltypesize * channels * width
	stride = (stride + (pad - 1)) &^ (pad - 1)

	return stride
}

func calcFaceSize(h *header) uint32 {
	stride := calcStride(h, h.Pixelwidth, 4)
	return stride * h.Pixelheight
}

func LoadKtx(filename string, tex gl.Uint) gl.Uint {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	h := new(header)
	r := bufio.NewReader(file)

	log.Println("header size:", unsafe.Sizeof(*h))
	d, err := r.Peek(int(unsafe.Sizeof(*h)))
	if err != nil {
		log.Fatal(err)
	}

	if bytes.Compare(d[:12], identifier) != 0 {
		log.Fatal("invalid file header")
	}

	if err := binary.Read(bytes.NewBuffer(d), binary.LittleEndian, h); err != nil {
		log.Fatal(err)
	}
	log.Printf("%#v\n", h)

	if h.Endianness == 0x04030201 {
		// No swap needed
	} else if h.Endianness == 0x01020304 {
		// Swap needed
		h.Endianness = swap32(h.Endianness)
		h.Gltype = swap32(h.Gltype)
		h.Gltypesize = swap32(h.Gltypesize)
		h.Glformat = swap32(h.Glformat)
		h.Glinternalformat = swap32(h.Glinternalformat)
		h.Glbaseinternalformat = swap32(h.Glbaseinternalformat)
		h.Pixelwidth = swap32(h.Pixelwidth)
		h.Pixelheight = swap32(h.Pixelheight)
		h.Pixeldepth = swap32(h.Pixeldepth)
		h.Arrayelements = swap32(h.Arrayelements)
		h.Faces = swap32(h.Faces)
		h.Miplevels = swap32(h.Miplevels)
		h.Keypairbytes = swap32(h.Keypairbytes)
	} else {
		log.Fatal("invalid header field endianness:", h.Endianness)
	}

	var target gl.Enum = gl.NONE
	// Guess target (texture type)
	if h.Pixelheight == 0 {
		if h.Arrayelements == 0 {
			target = gl.TEXTURE_1D
		} else {
			target = gl.TEXTURE_1D_ARRAY
		}
	} else if h.Pixeldepth == 0 {
		if h.Arrayelements == 0 {
			if h.Faces == 0 {
				target = gl.TEXTURE_2D
			} else {
				target = gl.TEXTURE_CUBE_MAP
			}
		} else {
			if h.Faces == 0 {
				target = gl.TEXTURE_2D_ARRAY
			} else {
				target = gl.TEXTURE_CUBE_MAP_ARRAY
			}
		}
	} else {
		target = gl.TEXTURE_3D
	}

	// Check for insanity...
	if target == gl.NONE || // Couldn't figure out target
		(h.Pixelwidth == 0) || // Texture has no width???
		(h.Pixelheight == 0 && h.Pixeldepth != 0) { // Texture has depth but no height???
		log.Fatal("invalid dimension")
	}

	if tex == 0 {
		gl.GenTextures(1, &tex)
	}
	gl.BindTexture(target, tex)

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(filename, len(data))
	data = data[unsafe.Sizeof(h):]

	if h.Miplevels == 0 {
		h.Miplevels = 1
	}

	switch target {
	case gl.TEXTURE_1D:
		gl.TexStorage1D(target, gl.Sizei(h.Miplevels), gl.Enum(h.Glinternalformat), gl.Sizei(h.Pixelwidth))
		gl.TexSubImage1D(target, 0, 0, gl.Sizei(h.Pixelwidth), gl.Enum(h.Glformat), gl.Enum(h.Glinternalformat), gl.Pointer(&data[0]))

	case gl.TEXTURE_2D:
		gl.TexStorage2D(target, gl.Sizei(h.Miplevels), gl.Enum(h.Glinternalformat), gl.Sizei(h.Pixelwidth), gl.Sizei(h.Pixelheight))

		height := h.Pixelheight
		width := h.Pixelwidth
		gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)
		for i := 0; i < int(h.Miplevels); i++ {
			gl.TexSubImage2D(target, gl.Int(i), 0, 0, gl.Sizei(width), gl.Sizei(height),
				gl.Enum(h.Glformat), gl.Enum(h.Gltype), gl.Pointer(&data[0]))

			data = data[height*calcStride(h, width, 1):]
			//log.Println(height*calcStride(h, width, 1), len(data))
			height >>= 1
			width >>= 1
			if height == 0 {
				height = 1
			}
			if width == 0 {
				width = 1
			}
		}

	case gl.TEXTURE_3D:
		gl.TexStorage3D(target, gl.Sizei(h.Miplevels), gl.Enum(h.Glinternalformat),
			gl.Sizei(h.Pixelwidth), gl.Sizei(h.Pixelheight), gl.Sizei(h.Pixeldepth))
		gl.TexSubImage3D(target, 0, 0, 0, 0, gl.Sizei(h.Pixelwidth), gl.Sizei(h.Pixelheight),
			gl.Sizei(h.Pixeldepth), gl.Enum(h.Glformat), gl.Enum(h.Gltype), gl.Pointer(&data[0]))

	case gl.TEXTURE_1D_ARRAY:
		gl.TexStorage2D(target, gl.Sizei(h.Miplevels), gl.Enum(h.Glinternalformat), gl.Sizei(h.Pixelwidth), gl.Sizei(h.Arrayelements))
		gl.TexSubImage2D(target, 0, 0, 0, gl.Sizei(h.Pixelwidth), gl.Sizei(h.Arrayelements),
			gl.Enum(h.Glformat), gl.Enum(h.Gltype), gl.Pointer(&data[0]))

	case gl.TEXTURE_2D_ARRAY:
		gl.TexStorage3D(target, gl.Sizei(h.Miplevels), gl.Enum(h.Glinternalformat),
			gl.Sizei(h.Pixelwidth), gl.Sizei(h.Pixelheight), gl.Sizei(h.Arrayelements))
		gl.TexSubImage3D(target, 0, 0, 0, 0, gl.Sizei(h.Pixelwidth),
			gl.Sizei(h.Pixelheight), gl.Sizei(h.Arrayelements), gl.Enum(h.Glformat), gl.Enum(h.Gltype), gl.Pointer(&data[0]))

	case gl.TEXTURE_CUBE_MAP:
		gl.TexStorage2D(target, gl.Sizei(h.Miplevels), gl.Enum(h.Glinternalformat), gl.Sizei(h.Pixelwidth), gl.Sizei(h.Pixelheight))

		face_size := calcFaceSize(h)
		for i := 0; i < int(h.Faces); i++ {
			d := data[int(face_size)*i:]
			gl.TexSubImage2D(gl.Enum(gl.TEXTURE_CUBE_MAP_POSITIVE_X+i), 0, 0, 0, gl.Sizei(h.Pixelwidth), gl.Sizei(h.Pixelheight),
				gl.Enum(h.Glformat), gl.Enum(h.Gltype), gl.Pointer(&d[0]))
		}

	case gl.TEXTURE_CUBE_MAP_ARRAY:
		gl.TexStorage3D(target, gl.Sizei(h.Miplevels), gl.Enum(h.Glinternalformat),
			gl.Sizei(h.Pixelwidth), gl.Sizei(h.Pixelheight), gl.Sizei(h.Arrayelements))
		gl.TexSubImage3D(target, 0, 0, 0, 0, gl.Sizei(h.Pixelwidth), gl.Sizei(h.Pixelheight), gl.Sizei(h.Faces*h.Arrayelements),
			gl.Enum(h.Glformat), gl.Enum(h.Gltype), gl.Pointer(&data[0]))

	default: // Should never happen
		log.Fatal("invalid target", target)
	}

	if h.Miplevels == 1 {
		gl.GenerateMipmap(target)
	}

	return tex
}

/*
func SaveKtx(filename string, target gl.Uint, gl.Uint tex) bool {

}
*/
