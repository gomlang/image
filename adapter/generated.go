package adapter

func GomlBind_decode_jpeg(p0 string, p1 int, p2 int, p3 int, p4 int) (string, int, int, string) {
    return DecodeJPEG(p0, p1, p2, p3, p4)
}

func GomlBind_decode_webp(p0 string, p1 int, p2 int, p3 int, p4 int) (string, int, int, string) {
    return DecodeWebP(p0, p1, p2, p3, p4)
}

func GomlBind_encode_jpeg(p0 string, p1 int, p2 int, p3 int, p4 int, p5 int, p6 int, p7 int, p8 uint32) (string, string) {
    return EncodeJPEG(p0, p1, p2, p3, p4, p5, p6, p7, p8)
}

func GomlBind_encode_webp(p0 string, p1 int, p2 int, p3 int, p4 int, p5 int, p6 int, p7 int) (string, string) {
    return EncodeWebP(p0, p1, p2, p3, p4, p5, p6, p7)
}
