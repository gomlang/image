from pathlib import Path
from struct import pack
from zlib import compress, crc32

root = Path(__file__).parent
signature = bytes.fromhex('89504e470d0a1a0a')


def chunk(name, payload):
    return pack('>I', len(payload)) + name + payload + pack('>I', crc32(name + payload))


def png(name, width, height, depth, color, rows, palette=b'', transparency=b'', interlace=0):
    header = pack('>IIBBBBB', width, height, depth, color, 0, 0, interlace)
    payload = signature + chunk(b'IHDR', header)
    if palette:
        payload += chunk(b'PLTE', palette)
    if transparency:
        payload += chunk(b'tRNS', transparency)
    payload += chunk(b'IDAT', compress(rows)) + chunk(b'IEND', b'')
    (root / name).write_bytes(payload)


def packed(values, depth):
    output = bytearray((len(values) * depth + 7) // 8)
    for index, value in enumerate(values):
        bit = index * depth
        output[bit // 8] |= value << (8 - depth - bit % 8)
    return bytes(output)


for depth, values in [(1, [0, 1]), (2, [1, 2]), (4, [2, 13]), (8, [77, 190]), (16, [0x0101, 0x0102])]:
    row = b''.join(pack('>H', value) for value in values) if depth == 16 else packed(values, depth)
    transparency = pack('>H', values[0]) if depth in (8, 16) else b''
    png(f'gray{depth}.png', 2, 1, depth, 0, b'\0' + row, transparency=transparency)

png('rgb8.png', 2, 1, 8, 2, b'\0' + bytes([10, 20, 30, 200, 100, 50]), transparency=pack('>HHH', 10, 20, 30))
png('rgb16.png', 1, 1, 16, 2, b'\0' + pack('>HHH', 0x1234, 0x5678, 0xabcd))
for depth, count in [(1, 2), (2, 4), (4, 16), (8, 17)]:
    palette = b''.join(bytes([i, 255 - i, i // 2]) for i in range(count))
    png(f'indexed{depth}.png', 2, 1, depth, 3, b'\0' + packed([0, count - 1], depth), palette=palette, transparency=bytes([0, 128]))
png('gray_alpha8.png', 1, 1, 8, 4, bytes([0, 200, 128]))
png('gray_alpha16.png', 1, 1, 16, 4, b'\0' + pack('>HH', 0x8000, 0x8000))
png('rgba8.png', 1, 1, 8, 6, bytes([0, 255, 0, 0, 128]))
png('rgba16.png', 1, 1, 16, 6, b'\0' + pack('>HHHH', 0xffff, 0, 0, 0x8000))

rows = [bytes([255, 0, 0, 255, 0, 255, 0, 255]), bytes([0, 0, 255, 255, 128, 0, 0, 128])]


def paeth(a, b, c):
    p = a + b - c
    distances = [abs(p - a), abs(p - b), abs(p - c)]
    return (a, b, c)[distances.index(min(distances))]


for mode in range(5):
    filtered = bytearray()
    for y, row in enumerate(rows):
        previous = rows[y - 1] if y else bytes(len(row))
        filtered.append(mode)
        for x, value in enumerate(row):
            left = row[x - 4] if x >= 4 else 0
            above = previous[x]
            upper = previous[x - 4] if x >= 4 else 0
            predictor = (0, left, above, (left + above) // 2, paeth(left, above, upper))[mode]
            filtered.append((value - predictor) & 255)
    png(f'filter{mode}.png', 2, 2, 8, 6, bytes(filtered))

passes = [(0, 0, 8, 8), (4, 0, 8, 8), (0, 4, 4, 8), (2, 0, 4, 4), (0, 2, 2, 4), (1, 0, 2, 2), (0, 1, 1, 2)]
interlaced = bytearray()
for sx, sy, dx, dy in passes:
    for y in range(sy, 8, dy):
        interlaced.append(0)
        for x in range(sx, 8, dx):
            interlaced.extend([x * 32, y * 32, 0, 255])
png('adam7.png', 8, 8, 8, 6, bytes(interlaced), interlace=1)

header = pack('>IIBBBBB', 1, 1, 8, 6, 0, 0, 0)
encoded = compress(bytes([0, 255, 0, 0, 255]))
(root / 'split_idat.png').write_bytes(signature + chunk(b'IHDR', header) + chunk(b'tEXt', b'key\0value') + chunk(b'IDAT', encoded[:3]) + chunk(b'IDAT', encoded[3:]) + chunk(b'tEXt', b'after\0data') + chunk(b'IEND', b''))
(root / 'duplicate_ihdr.png').write_bytes(signature + chunk(b'IHDR', header) + chunk(b'IHDR', header) + chunk(b'IDAT', encoded) + chunk(b'IEND', b''))
(root / 'unknown_critical.png').write_bytes(signature + chunk(b'IHDR', header) + chunk(b'ABCD', b'') + chunk(b'IDAT', encoded) + chunk(b'IEND', b''))
(root / 'noncontiguous_idat.png').write_bytes(signature + chunk(b'IHDR', header) + chunk(b'IDAT', encoded[:3]) + chunk(b'tEXt', b'x\0y') + chunk(b'IDAT', encoded[3:]) + chunk(b'IEND', b''))
(root / 'bad_filter.png').write_bytes(signature + chunk(b'IHDR', header) + chunk(b'IDAT', compress(bytes([5, 255, 0, 0, 255]))) + chunk(b'IEND', b''))
indexed_header = pack('>IIBBBBB', 1, 1, 1, 3, 0, 0, 0)
(root / 'bad_palette_index.png').write_bytes(signature + chunk(b'IHDR', indexed_header) + chunk(b'PLTE', bytes([255, 0, 0])) + chunk(b'IDAT', compress(bytes([0, 128]))) + chunk(b'IEND', b''))
