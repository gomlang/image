from pathlib import Path
import struct

ROOT = Path(__file__).parent


def segment(marker, payload):
    return bytes([255, marker]) + struct.pack('>H', len(payload) + 2) + payload


def quantization(value, wide=False):
    entry = struct.pack('>H', value) if wide else bytes([value])
    return segment(219, bytes([16 if wide else 0]) + entry * 64)


def frame(ids, width=8, height=8):
    components = b''.join(bytes([identifier, 17, 0]) for identifier in ids)
    return segment(192, bytes([8]) + struct.pack('>HHB', height, width, len(ids)) + components)


def huffman():
    dc = bytes([0]) + bytes([0, 0, 0, 12] + [0] * 12) + bytes(range(12))
    ac = bytes([16]) + bytes([1] + [0] * 15) + bytes([0])
    return segment(196, dc + ac)


def entropy(values):
    bits = ''
    for value in values:
        size = abs(value).bit_length()
        bits += f'{size:04b}'
        if size:
            encoded = value if value >= 0 else value - 1
            bits += f'{encoded & ((1 << size) - 1):0{size}b}'
        bits += '0'
    bits += '1' * (-len(bits) % 8)
    data = bytes(int(bits[index:index + 8], 2) for index in range(0, len(bits), 8))
    return data.replace(b'\xff', b'\xff\x00')


def scan(ids, values):
    header = bytes([len(ids)]) + b''.join(bytes([identifier, 0]) for identifier in ids) + bytes([0, 63, 0])
    return segment(218, header) + entropy(values)


def adobe(transform):
    return segment(238, b'Adobe' + struct.pack('>HHHB', 100, 0, 0, transform))


for name, ids, values, transform in [
    ('rgb', [82, 71, 66], [32, 96, 160], 0),
    ('cmyk', [1, 2, 3, 4], [40, 80, 160, 200], 0),
    ('ycck', [1, 2, 3, 4], [96, 128, 128, 200], 2),
]:
    data = b'\xff\xd8' + adobe(transform) + quantization(1) + frame(ids) + huffman()
    data += scan(ids, [(value - 128) * 8 for value in values]) + b'\xff\xd9'
    (ROOT / (name + '.jpg')).write_bytes(data)

for wide in [False, True]:
    data = b'\xff\xd8' + frame([1, 2, 3]) + huffman()
    for component, quant, value in [(1, 1, -512), (2, 2, 0), (3, 3, 0)]:
        data += quantization(quant * (256 if wide else 1), wide) + scan([component], [value // 256 if wide else value])
    data += b'\xff\xd9'
    (ROOT / ('quantization-wide.jpg' if wide else 'quantization-change.jpg')).write_bytes(data)

ids = [1]
data = b'\xff\xd8' + quantization(1) + frame(ids, 16, 8) + huffman()
data += segment(221, bytes([0, 1])) + scan(ids, [-512]) + bytes([255, 208]) + entropy([-512]) + b'\xff\xd9'
(ROOT / 'restart.jpg').write_bytes(data)
