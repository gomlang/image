from pathlib import Path
import subprocess
import tempfile


root = Path(__file__).resolve().parent
cases = [
    ("gradient", 37, 29),
    ("palette", 31, 25),
    ("noise", 19, 17),
    ("predictor", 65, 51),
]
palette = [(0, 0, 0, 255), (255, 0, 0, 255), (17, 47, 83, 128), (255, 255, 255, 0)]

for kind, width, height in cases:
    straight = bytearray()
    premultiplied = bytearray()
    seed = 12345
    for y in range(height):
        for x in range(width):
            if kind == "gradient":
                r, g, b, a = (
                    (x * 7 + y * 3) % 256,
                    (x * 5 + y * 9) % 256,
                    (x * 13 + y * 11) % 256,
                    (x * 17 + y * 19) % 256,
                )
            elif kind == "palette":
                r, g, b, a = palette[(x // 3 + y // 2) % 4]
            elif kind == "noise":
                channels = []
                for _ in range(4):
                    seed = (seed * 1664525 + 1013904223) & 0xFFFFFFFF
                    channels.append(seed >> 24)
                r, g, b, a = channels
            else:
                r, g, b, a = (
                    (x * 3 + y * 4) % 256,
                    (x * 7 + y * 2) % 256,
                    (x * 2 + y * 7) % 256,
                    255,
                )
            straight.extend((r, g, b, a))
            premultiplied.extend(
                ((r * a + 127) // 255, (g * a + 127) // 255, (b * a + 127) // 255, a)
            )
    with tempfile.TemporaryDirectory() as temporary:
        raw = Path(temporary) / "source.rgba"
        raw.write_bytes(straight)
        subprocess.run(
            [
                "ffmpeg", "-v", "error", "-y",
                "-f", "rawvideo", "-pixel_format", "rgba",
                "-video_size", f"{width}x{height}", "-i", str(raw),
                "-frames:v", "1", "-c:v", "libwebp",
                "-lossless", "1", "-compression_level", "6",
                str(root / f"lossless-{kind}.webp"),
            ],
            check=True,
        )
    (root / f"lossless-{kind}.rgba").write_bytes(premultiplied)
